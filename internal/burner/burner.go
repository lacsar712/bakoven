package burner

import (
	"math"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

type BurnerController struct {
	clk clock.ProcessClock
}

func NewBurnerController(clk clock.ProcessClock) *BurnerController {
	return &BurnerController{clk: clk}
}

func (b *BurnerController) EstimateRackframeTemp(reading model.BurnerReading) float64 {
	base := 300.0
	doughHeat := reading.DoughFlowTPH * 50
	airCool := reading.AirflowTPH * 2
	return base + doughHeat - airCool
}

func (b *BurnerController) BakeStable(reading model.BurnerReading) bool {
	if reading.BurnerPhase != model.BurnerStable && reading.BurnerPhase != model.BurnerIgnition {
		return false
	}
	return reading.RackframeTempF > 800 && reading.ExcessO2Pct >= model.MinRackframeO2Percent
}

func (b *BurnerController) TripRequired(reading model.BurnerReading) bool {
	if reading.ExcessO2Pct > model.MaxRackframeO2Percent*2 {
		return true
	}
	if reading.BurnerPhase == model.BurnerTrip {
		return true
	}
	if reading.RackframeTempF > 3500 {
		return true
	}
	return false
}

func (b *BurnerController) PhaseLabel(phase model.BurnerPhase) string {
	switch phase {
	case model.BurnerIdle:
		return "Idle"
	case model.BurnerProof:
		return "Proof"
	case model.BurnerIgnition:
		return "Ignition"
	case model.BurnerStable:
		return "Stable Bake"
	case model.BurnerTrip:
		return "Tripped"
	default:
		return string(phase)
	}
}

func (b *BurnerController) HeatReleaseMW(reading model.BurnerReading) float64 {
	return reading.DoughFlowTPH * 12.5
}

func (b *BurnerController) TurndownRatio(settings model.PlantSettings, currentDough float64) float64 {
	if settings.DoughFlowTPH <= 0 {
		return 0
	}
	return currentDough / settings.DoughFlowTPH
}

func (b *BurnerController) MinStableDough(settings model.PlantSettings) float64 {
	return settings.DoughFlowTPH * 0.25
}

func (b *BurnerController) NormalizeDough(flow, max float64) float64 {
	return math.Min(math.Max(flow, 0), max)
}
