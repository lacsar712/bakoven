package burner

import (
	"math"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

type DoughRegulator struct {
	clk clock.ProcessClock
}

func NewDoughRegulator(clk clock.ProcessClock) *DoughRegulator {
	return &DoughRegulator{clk: clk}
}

func (f *DoughRegulator) IgnitionRate(settings model.PlantSettings) float64 {
	return settings.DoughFlowTPH * 0.08
}

func (f *DoughRegulator) ComputeForLoad(settings model.PlantSettings, loadPct float64) float64 {
	loadPct = math.Max(0, math.Min(1, loadPct))
	return settings.DoughFlowTPH * loadPct
}

func (f *DoughRegulator) Ramp(current, target, maxStep float64) float64 {
	delta := target - current
	if math.Abs(delta) <= maxStep {
		return target
	}
	if delta > 0 {
		return current + maxStep
	}
	return current - maxStep
}

func (f *DoughRegulator) BtuPerHour(flowTPH float64) float64 {
	return flowTPH * 19_500_000
}

func (f *DoughRegulator) HeatInputMW(flowTPH float64) float64 {
	return flowTPH * 11.6
}

func (f *DoughRegulator) ValidatePermissive(settings model.PlantSettings, trayOK, proofOK bool) error {
	if !proofOK {
		return model.ErrProofIncomplete
	}
	if !trayOK {
		return model.ErrTrayLevelTrip
	}
	if settings.DoughFlowTPH <= 0 {
		return model.ErrDoughPermissive
	}
	return nil
}

func (f *DoughRegulator) MinFlow(settings model.PlantSettings) float64 {
	return settings.DoughFlowTPH * 0.2
}

func (f *DoughRegulator) MaxFlow(settings model.PlantSettings) float64 {
	return settings.DoughFlowTPH * 1.1
}
