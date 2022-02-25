package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/evcc-io/evcc/api"
)

// setPhases sets the number of enabled phases without modifying the charger
func (lp *LoadPoint) setPhases(phases int) {
	lp.Lock()
	defer lp.Unlock()

	if lp.Phases != phases {
		lp.Phases = phases
		lp.publish("phases", lp.Phases)
	}
}

// activePhases returns the number of expectedly active phases for the meter.
// If unknown for 1p3p chargers during startup it will assume 1p.
func (lp *LoadPoint) activePhases() int {
	const unknownPhases = 1

	vehicle := lp.vehicleCapablePhases()
	physical := lp.GetPhases()

	if vehicle > 0 && (vehicle <= physical || physical == 0) {
		return vehicle
	}

	// if we don't have a valid value yet assume phase configuration
	if physical > 1 {
		return unknownPhases
	}

	// assume 1p for switchable charger during startup
	return 1
}

func (lp *LoadPoint) vehicleCapablePhases() int {
	if lp.vehicle != nil {
		if phases := lp.vehicle.Phases(); phases > 0 {
			return phases
		}
	}

	// if vehicle is charging >1p then assume that is the
	// number of phases that the vehicle can use
	if lp.measuredPhases > 1 {
		return lp.measuredPhases
	}

	return 0
}

// scalePhasesIfAvailable scales if api.ChargePhases is available
func (lp *LoadPoint) scalePhasesIfAvailable(phases int) error {
	err := lp.scalePhases(phases)
	if errors.Is(err, api.ErrNotAvailable) {
		return nil
	}
	return err
}

// scalePhases adjusts the number of active phases and returns the appropriate charging current.
// Returns api.ErrNotAvailable if api.ChargePhases is not available.
func (lp *LoadPoint) scalePhases(phases int) error {
	if phases != 1 && phases != 3 {
		return fmt.Errorf("invalid number of phases: %d", phases)
	}

	cp, ok := lp.charger.(api.ChargePhases)
	if !ok {
		return api.ErrNotAvailable
	}

	if lp.GetPhases() != phases {
		// disable charger - this will also stop the car charging using the api if available
		if err := lp.setLimit(0, true); err != nil {
			return err
		}

		// switch phases
		if err := cp.Phases1p3p(phases); err != nil {
			return fmt.Errorf("switch phases: %w", err)
		}

		// update setting
		lp.setPhases(phases)

		// disable phase timer
		lp.phaseTimer = time.Time{}

		// allow pv mode to re-enable charger right away
		lp.elapsePVTimer()
	}

	return nil
}

// pvScalePhases switches phases if necessary and returns if switch occurred
func (lp *LoadPoint) pvScalePhases(availablePower, minCurrent, maxCurrent float64) bool {
	phases := lp.GetPhases()

	// observed phase state inconsistency (https://github.com/evcc-io/evcc/issues/1572, https://github.com/evcc-io/evcc/issues/2230)
	if phases > 0 && phases < lp.measuredPhases {
		lp.log.WARN.Printf("ignoring inconsistent phases: %dp < %dp observed active", phases, lp.measuredPhases)
	}

	// this can happen the first time for a 1p3p-capable charger, see https://github.com/evcc-io/evcc/issues/2520
	// if phases == 0 && lp.activePhases == 0 {
	// 	lp.log.DEBUG.Printf("assuming initial phase state: 3p")
	// 	lp.phaseTimer = elapsed
	// 	lp.activePhases = 3
	// }

	var waiting bool
	activePhases := lp.activePhases()
	targetCurrent := powerToCurrent(availablePower, activePhases)

	// scale down phases
	if targetCurrent < minCurrent && activePhases > 1 {
		lp.log.DEBUG.Printf("available power below %dp min threshold of %.0fW", activePhases, float64(activePhases)*Voltage*minCurrent)

		if lp.phaseTimer.IsZero() {
			lp.log.DEBUG.Printf("start phase disable timer: %v", lp.Disable.Delay)
			lp.phaseTimer = lp.clock.Now()
		}

		lp.publishTimer(phaseTimer, lp.Disable.Delay, phaseScale1p)

		elapsed := lp.clock.Since(lp.phaseTimer)
		if elapsed >= lp.Disable.Delay {
			lp.log.DEBUG.Println("phase disable timer elapsed")
			if err := lp.scalePhases(1); err == nil {
				lp.log.DEBUG.Printf("switched phases: 1p @ %.0fW", availablePower)
				return true
			} else {
				lp.log.ERROR.Printf("switch phases: %v", err)
			}
		}

		waiting = true
		lp.log.DEBUG.Printf("phase disable timer remaining: %v", (lp.Disable.Delay - elapsed).Round(time.Second))
	}

	// scale up phases
	if min3pCurrent := powerToCurrent(availablePower, 3); min3pCurrent >= minCurrent && activePhases == 1 {
		lp.log.DEBUG.Printf("available power above 3p min threshold of %.0fW", 3*Voltage*minCurrent)

		if lp.phaseTimer.IsZero() {
			lp.log.DEBUG.Printf("start phase enable timer: %v", lp.Enable.Delay)
			lp.phaseTimer = lp.clock.Now()
		}

		lp.publishTimer(phaseTimer, lp.Disable.Delay, phaseScale3p)

		elapsed := lp.clock.Since(lp.phaseTimer)
		if elapsed >= lp.Disable.Delay {
			lp.log.DEBUG.Println("phase enable timer elapsed")
			if err := lp.scalePhases(3); err == nil {
				lp.log.DEBUG.Printf("switched phases: 3p @ %.0fW", availablePower)
				return true
			} else {
				lp.log.ERROR.Printf("switch phases: %v", err)
			}
		}

		waiting = true
		lp.log.DEBUG.Printf("phase enable timer remaining: %v", (lp.Disable.Delay - elapsed).Round(time.Second))
	}

	// reset timer to disabled state
	if !waiting && !lp.phaseTimer.IsZero() {
		lp.log.DEBUG.Printf("phase timer reset")
		lp.phaseTimer = time.Time{}

		lp.publishTimer(phaseTimer, 0, timerInactive)
	}

	return false
}
