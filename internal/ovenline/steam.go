package ovenline

import (
	"math"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

type SteamPath struct {
	clk clock.ProcessClock
}

func NewSteamPath(clk clock.ProcessClock) *SteamPath {
	return &SteamPath{clk: clk}
}

func (s *SteamPath) ComputeFlow(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	loadFactor := snap.Ovenline.OutputMW / math.Max(snap.Settings.TargetMW, 1)
	pressureFactor := snap.Ovenline.SteamPressurePSI / math.Max(snap.Settings.TargetSteamPSI, 1)
	return snap.Settings.FeedwaterFlowTPH * loadFactor * math.Min(pressureFactor, 1.1)
}

func (s *SteamPath) ComputeMW(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	flow := s.ComputeFlow(snap, firing)
	enthalpy := s.enthalpyFromPressure(snap.Ovenline.SteamPressurePSI)
	return flow * enthalpy * 0.001
}

func (s *SteamPath) ComputeTemp(pressurePSI float64) float64 {
	if pressurePSI <= 0 {
		return 70
	}
	return 212 + math.Sqrt(pressurePSI)*8.5
}

func (s *SteamPath) enthalpyFromPressure(pressurePSI float64) float64 {
	return 1000 + pressurePSI*0.3
}

func (s *SteamPath) EnthalpyBalance(snap model.PlantSnapshot) float64 {
	inEnergy := snap.Tray.FeedwaterTPH * 200
	outEnergy := snap.Tray.SteamFlowTPH * s.enthalpyFromPressure(snap.Ovenline.SteamPressurePSI)
	return inEnergy - outEnergy
}

func (s *SteamPath) LoadPercent(snap model.PlantSnapshot) float64 {
	if snap.Settings.TargetMW <= 0 {
		return 0
	}
	return (snap.Ovenline.OutputMW / snap.Settings.TargetMW) * 100
}
