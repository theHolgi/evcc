package charger

// Code generated by github.com/evcc-io/evcc/cmd/tools/decorate.go. DO NOT EDIT.

import (
	"github.com/evcc-io/evcc/api"
)

func decoratePhoenixEVEth(phoenixEVEth *PhoenixEVEth, meter func() (float64, error), meterEnergy func() (float64, error), meterCurrent func() (float64, float64, float64, error)) api.Charger {
	switch {
	case meter == nil && meterCurrent == nil && meterEnergy == nil:
		return &struct {
			*PhoenixEVEth
		}{
			PhoenixEVEth: phoenixEVEth,
		}

	case meter != nil && meterCurrent == nil && meterEnergy == nil:
		return &struct {
			*PhoenixEVEth
			api.Meter
		}{
			PhoenixEVEth: phoenixEVEth,
			Meter: &decoratePhoenixEVEthMeterImpl{
				meter: meter,
			},
		}

	case meter == nil && meterCurrent == nil && meterEnergy != nil:
		return &struct {
			*PhoenixEVEth
			api.MeterEnergy
		}{
			PhoenixEVEth: phoenixEVEth,
			MeterEnergy: &decoratePhoenixEVEthMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}

	case meter != nil && meterCurrent == nil && meterEnergy != nil:
		return &struct {
			*PhoenixEVEth
			api.Meter
			api.MeterEnergy
		}{
			PhoenixEVEth: phoenixEVEth,
			Meter: &decoratePhoenixEVEthMeterImpl{
				meter: meter,
			},
			MeterEnergy: &decoratePhoenixEVEthMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}

	case meter == nil && meterCurrent != nil && meterEnergy == nil:
		return &struct {
			*PhoenixEVEth
			api.MeterCurrent
		}{
			PhoenixEVEth: phoenixEVEth,
			MeterCurrent: &decoratePhoenixEVEthMeterCurrentImpl{
				meterCurrent: meterCurrent,
			},
		}

	case meter != nil && meterCurrent != nil && meterEnergy == nil:
		return &struct {
			*PhoenixEVEth
			api.Meter
			api.MeterCurrent
		}{
			PhoenixEVEth: phoenixEVEth,
			Meter: &decoratePhoenixEVEthMeterImpl{
				meter: meter,
			},
			MeterCurrent: &decoratePhoenixEVEthMeterCurrentImpl{
				meterCurrent: meterCurrent,
			},
		}

	case meter == nil && meterCurrent != nil && meterEnergy != nil:
		return &struct {
			*PhoenixEVEth
			api.MeterCurrent
			api.MeterEnergy
		}{
			PhoenixEVEth: phoenixEVEth,
			MeterCurrent: &decoratePhoenixEVEthMeterCurrentImpl{
				meterCurrent: meterCurrent,
			},
			MeterEnergy: &decoratePhoenixEVEthMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}

	case meter != nil && meterCurrent != nil && meterEnergy != nil:
		return &struct {
			*PhoenixEVEth
			api.Meter
			api.MeterCurrent
			api.MeterEnergy
		}{
			PhoenixEVEth: phoenixEVEth,
			Meter: &decoratePhoenixEVEthMeterImpl{
				meter: meter,
			},
			MeterCurrent: &decoratePhoenixEVEthMeterCurrentImpl{
				meterCurrent: meterCurrent,
			},
			MeterEnergy: &decoratePhoenixEVEthMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}
	}

	return nil
}

type decoratePhoenixEVEthMeterImpl struct {
	meter func() (float64, error)
}

func (impl *decoratePhoenixEVEthMeterImpl) CurrentPower() (float64, error) {
	return impl.meter()
}

type decoratePhoenixEVEthMeterCurrentImpl struct {
	meterCurrent func() (float64, float64, float64, error)
}

func (impl *decoratePhoenixEVEthMeterCurrentImpl) Currents() (float64, float64, float64, error) {
	return impl.meterCurrent()
}

type decoratePhoenixEVEthMeterEnergyImpl struct {
	meterEnergy func() (float64, error)
}

func (impl *decoratePhoenixEVEthMeterEnergyImpl) TotalEnergy() (float64, error) {
	return impl.meterEnergy()
}
