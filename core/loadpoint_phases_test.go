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

type testCase struct {
	// capable=0 signals 1p3p as set during loadpoint init
	// physical/vehicle=0 signals unknown
	// measuredPhases<>0 signals previous measurement
	capable, physical, vehicle, measuredPhases, expected int
}

var (
	phaseTests = []testCase{
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
		{0, 0, 0, 0, unknownPhases},
		{0, 0, 0, 1, 1},
		{0, 0, 0, 2, 2},
		{0, 0, 0, 3, 3},
		{0, 0, 1, 0, 1},
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 0, 1},
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

	scaleDown = []testCase{
		// 1p3p initial
		{0, 0, 0, 0, unknownPhases},
		{0, 0, 0, 2, 2},
		{0, 0, 0, 3, 3},
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 2, 2},
		{0, 1, 0, 3, 3},
		// 1p3p, 3 currently active
		{0, 3, 0, 0, unknownPhases},
		{0, 3, 0, 2, 2},
		{0, 3, 0, 3, 3},
		{0, 3, 2, 0, 2},
		{0, 3, 3, 0, 3},
	}

	scaleUp = []testCase{
		// 1p3p initial
		{0, 0, 0, 0, unknownPhases},
		{0, 0, 0, 0, 1},
		{0, 0, 0, 1, 1},
		{0, 0, 0, 2, 2},
		{0, 0, 0, 3, 3},
		{0, 0, 2, 0, 2},
		{0, 0, 3, 0, 3},
		// 1p3p, 1 currently active
		{0, 1, 0, 0, 1},
		{0, 1, 0, 1, 1},
		{0, 1, 2, 0, 1},
		{0, 1, 3, 0, 1},
		// 1p3p, 3 currently active
		{0, 3, 0, 0, unknownPhases}, // TODO remove: Fehler bei scaleUp
		{0, 3, 0, 1, 1},             // TODO remove: Fehler bei scaleUp
		{0, 3, 0, 2, 2},
		{0, 3, 0, 3, 3},
		{0, 3, 2, 0, 2},
		{0, 3, 3, 0, 3},
	}
)

func caseMatches(tc testCase, cases []testCase) bool {
	for _, tc2 := range cases {
		if reflect.DeepEqual(tc, tc2) {
			return true
		}
	}
	return false
}

func testScale(t *testing.T, lp *LoadPoint, power float64, direction string, tc testCase, cases []testCase) {
	scaled := lp.pvScalePhases(power, minA, maxA)

	if caseMatches(tc, cases) {
		if !scaled {
			t.Errorf("%v missing scale %s", tc, direction)
		}
	} else if scaled {
		t.Errorf("%v unexpected scale %s", tc, direction)
	}
}

func TestPhaseHandling(t *testing.T) {
	clock := clock.NewMock()
	ctrl := gomock.NewController(t)

	for _, tc := range phaseTests {
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

		lp.measuredPhases = tc.measuredPhases
		if tc.measuredPhases > 0 && tc.vehicle > 0 {
			t.Fatalf("%v invalid test case", tc)
		}

		if lp.Phases != tc.physical {
			t.Error("wrong phases", lp.Phases, tc.physical)
		}

		if phs := lp.activePhases(); phs != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, phs)
		}
		ctrl.Finish()

		// scaling
		if charger.MockChargePhases != nil {
			// scale down
			min1p := 1 * minA * Voltage
			lp.phaseTimer = time.Time{}

			charger.MockCharger.EXPECT().Enable(false).Return(nil).MaxTimes(1)
			charger.MockChargePhases.EXPECT().Phases1p3p(1).Return(nil).MaxTimes(1)

			testScale(t, lp, min1p, "down", tc, scaleDown)
			ctrl.Finish()

			// scale up
			min3p := 3 * minA * Voltage
			lp.phaseTimer = time.Time{}

			charger.MockCharger.EXPECT().Enable(false).Return(nil).MaxTimes(1)
			charger.MockChargePhases.EXPECT().Phases1p3p(3).Return(nil).MaxTimes(1)

			testScale(t, lp, min3p, "up", tc, scaleUp)
			ctrl.Finish()
		}
	}
}
