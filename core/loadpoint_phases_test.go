package core

import (
	"testing"

	evbus "github.com/asaskevich/EventBus"
	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/mock"
	"github.com/evcc-io/evcc/util"
	"github.com/golang/mock/gomock"
)

func TestEffectivePhases(t *testing.T) {
	clock := clock.NewMock()
	ctrl := gomock.NewController(t)

	// wrap vehicle with estimator
	// vehicle.EXPECT().Capacity().Return(int64(10))
	// socEstimator := soc.NewEstimator(util.NewLogger("foo"), charger, vehicle, false)

	tcs := []struct {
		capable, physical, vehicle, expected int
	}{
		{1, 1, 0, 1},
		{1, 1, 1, 1},
		{1, 1, 2, 1},
		{1, 1, 3, 1},
		{3, 3, 0, 1}, // Annahme
		{3, 3, 1, 1},
		{3, 3, 2, 2},
		{3, 3, 3, 3},
		{0, 0, 0, 1}, // TODO Annahme gelbe Markierung
		{0, 0, 1, 1},
		{0, 0, 2, 2},
		{0, 0, 3, 3},
		{0, 1, 0, 1}, // Annahme
		{0, 1, 1, 1},
		{0, 1, 2, 1},
		{0, 1, 3, 1},
		{0, 3, 0, 1}, // Annahme
		{0, 3, 1, 1},
		{0, 3, 2, 2},
		{0, 3, 3, 3},
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
		vehicle.EXPECT().Phases().Return(tc.vehicle)

		lp := &LoadPoint{
			log:         util.NewLogger("foo"),
			bus:         evbus.New(),
			clock:       clock,
			charger:     charger,
			chargeMeter: &Null{},            // silence nil panics
			chargeRater: &Null{},            // silence nil panics
			chargeTimer: &Null{},            // silence nil panics
			progress:    NewProgress(0, 10), // silence nil panics
			wakeUpTimer: NewTimer(),         // silence nil panics
			MinCurrent:  minA,
			MaxCurrent:  maxA,
			vehicle:     vehicle, // needed for targetSoC check
			// socEstimator: socEstimator, // instead of vehicle: vehicle,
			Mode:   api.ModeNow,
			Phases: tc.physical,
		}

		attachListeners(t, lp)

		if lp.Phases != tc.physical {
			t.Error("wrong phases", lp.Phases, tc.physical)
		}

		if phs := lp.effectivePhases(); phs != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, phs)
		}
	}
}
