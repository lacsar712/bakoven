package model

import "errors"

var (
	ErrContextDone      = errors.New("operation cancelled")
	ErrPlantNotFound    = errors.New("plant unit not found")
	ErrLeaseHeld        = errors.New("interlock lease held by another operator")
	ErrLeaseMissing     = errors.New("interlock lease missing or expired")
	ErrGateBlocked      = errors.New("safety gate blocked")
	ErrDoughPermissive   = errors.New("dough permissive not satisfied")
	ErrIgnitionBlocked  = errors.New("ignition sequence blocked")
	ErrTrayLevelTrip    = errors.New("tray level trip condition")
	ErrPressureTrip     = errors.New("steam pressure trip condition")
	ErrBurnerTrip   = errors.New("burner trip condition")
	ErrIllegalState     = errors.New("illegal plant state transition")
	ErrSnapshotStale    = errors.New("snapshot revision stale")
	ErrWindowOpen       = errors.New("timing window still open")
	ErrProofIncomplete  = errors.New("rackframe proof incomplete")
	ErrCoordinationLock = errors.New("coordination lock held")
	ErrTrayLevelLow     = errors.New("tray level below low limit")
	ErrBakeLoss        = errors.New("rackframe bake lost")
	ErrSteamrelLimit    = errors.New("steamrel valve at limit")
)
