package tray

import (
	"math"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

type LevelController struct {
	clk clock.ProcessClock
}

func NewLevelController(clk clock.ProcessClock) *LevelController {
	return &LevelController{clk: clk}
}

func (l *LevelController) Compute(snap model.PlantSnapshot, firing bool) (float64, model.TrayCondition) {
	level := snap.Tray.LevelPercent
	if !firing {
		return level, model.TrayNormal
	}
	balance := snap.Tray.FeedwaterTPH - snap.Tray.SteamFlowTPH
	level += balance * 0.01
	level = math.Max(model.MinTrayLevelPercent, math.Min(model.MaxTrayLevelPercent, level))
	cond := l.classify(level, snap)
	return level, cond
}

func (l *LevelController) classify(level float64, snap model.PlantSnapshot) model.TrayCondition {
	setpoint := snap.Settings.TrayLevelSetpoint
	if level > setpoint+15 {
		return model.TraySwell
	}
	if level < setpoint-15 {
		return model.TrayShrink
	}
	if snap.Ovenline.SteamPressurePSI > snap.Settings.TargetSteamPSI*0.9 && level > setpoint+5 {
		return model.TrayCarry
	}
	return model.TrayNormal
}

func (l *LevelController) RecommendFeedwater(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	err := snap.Settings.TrayLevelSetpoint - snap.Tray.LevelPercent
	return snap.Settings.FeedwaterFlowTPH + err*3
}

func (l *LevelController) WithinLimits(level float64) bool {
	return level >= model.MinTrayLevelPercent && level <= model.MaxTrayLevelPercent
}

func (l *LevelController) TripLow(level float64) bool  { return level < model.TripTrayLowPercent }
func (l *LevelController) TripHigh(level float64) bool { return level > model.TripTrayHighPercent }

func (l *LevelController) LevelError(snap model.PlantSnapshot) float64 {
	return snap.Tray.LevelPercent - snap.Settings.TrayLevelSetpoint
}

func (l *LevelController) ThreeElementBias(snap model.PlantSnapshot) float64 {
	steam := snap.Tray.SteamFlowTPH
	feed := snap.Tray.FeedwaterTPH
	levelErr := l.LevelError(snap)
	return feed + (steam-feed)*0.5 + levelErr*2
}
