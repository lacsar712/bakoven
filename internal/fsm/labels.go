package fsm

import "github.com/lacsar712/bakoven/internal/model"

type PlantEvent string

const (
	EvStartProof     PlantEvent = "start_proof"
	EvProofComplete  PlantEvent = "proof_complete"
	EvIgnite         PlantEvent = "ignite"
	EvIgnitionStable PlantEvent = "ignition_stable"
	EvRampLoad       PlantEvent = "ramp_load"
	EvReachFiring    PlantEvent = "reach_firing"
	EvLoadFollow     PlantEvent = "load_follow"
	EvTrip           PlantEvent = "trip"
	EvResetTrip      PlantEvent = "reset_trip"
	EvEnterService   PlantEvent = "enter_service"
	EvLeaveService   PlantEvent = "leave_service"
	EvShutdown       PlantEvent = "shutdown"
)

func StateLabel(s model.PlantState) string {
	switch s {
	case model.StateColdStandby:
		return "Cold Standby"
	case model.StateProof:
		return "Rackframe Proof"
	case model.StateIgnition:
		return "Ignition"
	case model.StateRamp:
		return "Ramp"
	case model.StateFiring:
		return "Firing"
	case model.StateLoadFollow:
		return "Load Following"
	case model.StateTrip:
		return "Trip"
	case model.StateService:
		return "Service"
	default:
		return string(s)
	}
}

func EventLabel(e PlantEvent) string {
	switch e {
	case EvStartProof:
		return "Start Proof"
	case EvProofComplete:
		return "Proof Complete"
	case EvIgnite:
		return "Ignite"
	case EvIgnitionStable:
		return "Ignition Stable"
	case EvRampLoad:
		return "Ramp Load"
	case EvReachFiring:
		return "Reach Firing"
	case EvLoadFollow:
		return "Load Follow"
	case EvTrip:
		return "Trip"
	case EvResetTrip:
		return "Reset Trip"
	case EvEnterService:
		return "Enter Service"
	case EvLeaveService:
		return "Leave Service"
	case EvShutdown:
		return "Shutdown"
	default:
		return string(e)
	}
}
