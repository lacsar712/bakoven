package model

import "time"

func CloneSnapshot(s PlantSnapshot) PlantSnapshot {
	out := s
	out.Alarms = append([]AlarmEvent(nil), s.Alarms...)
	return out
}

func DefaultSnapshot(unitID string) PlantSnapshot {
	now := time.Now()
	return PlantSnapshot{
		UnitID: unitID,
		State:  StateColdStandby,
		Settings: PlantSettings{
			Mode:              ModeBaseLoad,
			TargetMW:          150,
			TargetSteamPSI:    NormalSteamPressurePSI,
			TrayLevelSetpoint: 55,
			FeedwaterFlowTPH:  400,
			DoughFlowTPH:       35,
			ExcessO2Setpoint:  3.5,
		},
		Plant: PlantRef{UnitLabel: unitID, PlantCode: "STEAM-PLT"},
		Tray: TrayReading{
			LevelPercent: 50,
			Condition:    TrayNormal,
			FeedwaterTPH: 0,
			SteamFlowTPH: 0,
		},
		Burner: BurnerReading{
			BurnerPhase: BurnerIdle,
		},
		Ovenline: OvenlineReading{
			SteamPressurePSI: 0,
			SteamTempF:       70,
		},
		UpdatedAt: now,
	}
}

func (s PlantSnapshot) IsFiring() bool {
	return s.State == StateFiring || s.State == StateLoadFollow || s.State == StateRamp
}

func (s PlantSnapshot) TrayWithinLimits() bool {
	return s.Tray.LevelPercent >= MinTrayLevelPercent && s.Tray.LevelPercent <= MaxTrayLevelPercent
}

func (s PlantSnapshot) PressureWithinLimits() bool {
	if !s.IsFiring() {
		return true
	}
	return s.Ovenline.SteamPressurePSI <= MaxSteamPressurePSI
}
