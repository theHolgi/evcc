package core

import (
	"reflect"
	"testing"
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/mock"
	"github.com/evcc-io/evcc/util"
	"github.com/golang/mock/gomock"
)

func TestEffectivePhases(t *testing.T) {
	const unknownPhases = 1

	clock := clock.NewMock()
	ctrl := gomock.NewController(t)

	tcs := []struct {
		// capable=0 signals 1p3p as set during loadpoint init
		// physical/vehicle=0 signals unknown
		// previousActive<>0 signals previous measurement
		capable, physical, vehicle, previousActive, expected int
	}{
		// 1p
		{1, 1, 0, 0, 1},
		{1, 1, 0, 1, 1},
		{1, 1, 1, 0, 1},
		{1, 1, 2, 0, 1},
		{1, 1, 3, 0, 1},
		// 3p
		{3, 3, 0, 0, unknownPhases},
		{3, 3, 0, 1, 1},
		{3, 3, 0, 2, 2},
		{3, 3, 0, 3, 3},
		{3, 3, 1, 0, 1},
		{3, 3, 2, 0, 2},
		{3, 3, 3, 0, 3},
		// 1p3p initial
		{0, 0, 0, 0, unknownPhases}, // TODO gelbe Markierung
		{0, 0, 0, 1, 1},             // TODO gelbe Markierung
		{0, 0, 0, 2, 2},             // TODO gelbe Markierung
		{0, 0, 0, 3, 3},             // TODO gelbe Markierung
		{0, 0, 1, 0, 1},
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 0, unknownPhases},
		{0, 1, 0, 1, 1},
		// {0, 1, 0, 2, 2}, // 2p active > 1p configured must not happen
		// {0, 1, 0, 3, 3}, // 3p active > 1p configured must not happen
		{0, 1, 1, 0, 1},
		{0, 1, 2, 0, 1},
		{0, 1, 3, 0, 1},
		// 1p3p, 3 currently active
		{0, 3, 0, 0, unknownPhases},
		{0, 3, 0, 1, 1},
		{0, 3, 0, 2, 2},
		{0, 3, 0, 3, 3},
		{0, 3, 1, 0, 1},
		{0, 3, 2, 0, 2},
		{0, 3, 3, 0, 3},
	}

	scaleDown := []struct{ capable, physical, vehicle, previousActive, expected int }{
		// 1p3p initial
		{0, 0, 0, 2, 2}, // TODO gelbe Markierung
		{0, 0, 0, 3, 3}, // TODO gelbe Markierung
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 2, 2},
		{0, 1, 0, 3, 3},
		// 1p3p, 3 currently active
		{0, 3, 0, 2, 2},
		{0, 3, 0, 3, 3},
		{0, 3, 2, 0, 2},
		{0, 3, 3, 0, 3},
	}

	scaleUp := []struct{ capable, physical, vehicle, previousActive, expected int }{
		// 1p3p initial
		{0, 0, 0, 1, 1},
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 0, unknownPhases},
		{0, 1, 2, 0, 1},
		{0, 1, 3, 0, 1},
		// 1p3p, 3 currently active
		{0, 3, 0, 2, 2},
		{0, 3, 0, 3, 3},
		{0, 3, 2, 0, 2},
		{0, 3, 3, 0, 3},
	}

	for _, tc := range tcs {
		t.Log(tc)

		var charger struct {
			*mock.MockCharger
			*mock.MockChargePhases
		}

		charger.MockCharger = mock.NewMockCharger(ctrl)
		charger.MockCharger.EXPECT().Enabled().Return(true, nil)
		charger.MockCharger.EXPECT().MaxCurrent(int64(minA)).Return(nil) // MaxCurrentEx not implemented

		// 1p3p
		if tc.capable == 0 {
			charger.MockChargePhases = mock.NewMockChargePhases(ctrl)
		}

		vehicle := mock.NewMockVehicle(ctrl)
		vehicle.EXPECT().Phases().Return(tc.vehicle).MinTimes(1)

		lp := &LoadPoint{
			log:         util.NewLogger("foo"),
			bus:         evbus.New(),
			clock:       clock,
			chargeMeter: &Null{},            // silence nil panics
			chargeRater: &Null{},            // silence nil panics
			chargeTimer: &Null{},            // silence nil panics
			progress:    NewProgress(0, 10), // silence nil panics
			wakeUpTimer: NewTimer(),         // silence nil panics
			Mode:        api.ModeNow,
			MinCurrent:  minA,
			MaxCurrent:  maxA,
			charger:     charger,
			vehicle:     vehicle,
			Phases:      tc.physical,
		}

		attachListeners(t, lp)

		// TODO reset activePhases when vehicle disconnects
		lp.measuredPhases = tc.previousActive
		if tc.previousActive > 0 && tc.vehicle > 0 {
			t.Fatalf("%v invalid test case", tc)
		}

		if lp.Phases != tc.physical {
			t.Error("wrong phases", lp.Phases, tc.physical)
		}

		if phs := lp.activePhases(); phs != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, phs)
		}

		// scaling
		if charger.MockChargePhases != nil {
			// scale down
			for _, tc2 := range scaleDown {
				if reflect.DeepEqual(tc, tc2) {
					charger.MockCharger.EXPECT().Enable(false).Return(nil)
					charger.MockChargePhases.EXPECT().Phases1p3p(1).Return(nil)

					if !lp.pvScalePhases(1*minA*Voltage, minA, maxA) {
						t.Errorf("%v missing scale down", tc)
					}

					break
				}
			}

			// scale up
			lp.phaseTimer = time.Time{}
			for _, tc2 := range scaleUp {
				if reflect.DeepEqual(tc, tc2) {
					charger.MockCharger.EXPECT().Enable(false).Return(nil)
					charger.MockChargePhases.EXPECT().Phases1p3p(3).Return(nil)

					if !lp.pvScalePhases(3*minA*Voltage, minA, maxA) {
						t.Errorf("%v missing scale up", tc)
					}

					break
				}
			}
		}

		ctrl.Finish()
	}
}
