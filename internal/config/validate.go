package config

import (
	"fmt"
	"strings"

	"github.com/lacsar712/bakoven/internal/model"
)

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.UnitID) == "" {
		return fmt.Errorf("unit_id required")
	}
	if cfg.ListenAddr == "" {
		return fmt.Errorf("listen_addr required")
	}
	if err := validateSettings(cfg.Settings); err != nil {
		return err
	}
	return nil
}

func validateSettings(s model.PlantSettings) error {
	if s.TargetMW < 0 {
		return fmt.Errorf("target_mw cannot be negative")
	}
	if s.TargetSteamPSI <= 0 {
		return fmt.Errorf("target_steam_psi must be positive")
	}
	if s.TrayLevelSetpoint < model.MinTrayLevelPercent || s.TrayLevelSetpoint > model.MaxTrayLevelPercent {
		return fmt.Errorf("tray level setpoint out of range")
	}
	if s.DoughFlowTPH <= 0 {
		return fmt.Errorf("dough_flow_tph must be positive")
	}
	if s.ExcessO2Setpoint < model.MinRackframeO2Percent || s.ExcessO2Setpoint > model.MaxRackframeO2Percent {
		return fmt.Errorf("excess_o2 setpoint out of range")
	}
	return nil
}
