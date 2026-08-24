package model

import "time"

type PlantMode string

const (
	ModeBaseLoad   PlantMode = "base_load"
	ModeSliding    PlantMode = "sliding_pressure"
	ModeRamp       PlantMode = "ramp"
	ModeHotStandby PlantMode = "hot_standby"
)

type PlantState string

const (
	StateColdStandby PlantState = "cold_standby"
	StateProof       PlantState = "proof"
	StateIgnition    PlantState = "ignition"
	StateRamp        PlantState = "ramp"
	StateFiring      PlantState = "firing"
	StateLoadFollow  PlantState = "load_follow"
	StateTrip        PlantState = "trip"
	StateService     PlantState = "service"
)

type BurnerPhase string

const (
	BurnerIdle     BurnerPhase = "idle"
	BurnerProof    BurnerPhase = "proof"
	BurnerIgnition BurnerPhase = "ignition"
	BurnerStable   BurnerPhase = "stable"
	BurnerTrip     BurnerPhase = "trip"
)

type TrayCondition string

const (
	TrayNormal  TrayCondition = "normal"
	TraySwell   TrayCondition = "swell"
	TrayShrink  TrayCondition = "shrink"
	TrayCarry   TrayCondition = "carryover"
)

type PlantSettings struct {
	Mode              PlantMode `json:"mode"`
	TargetMW          float64   `json:"target_mw"`
	TargetSteamPSI    float64   `json:"target_steam_psi"`
	TrayLevelSetpoint float64   `json:"tray_level_setpoint_pct"`
	FeedwaterFlowTPH  float64   `json:"feedwater_flow_tph"`
	DoughFlowTPH       float64   `json:"dough_flow_tph"`
	ExcessO2Setpoint  float64   `json:"excess_o2_pct"`
}

type TrayReading struct {
	LevelPercent   float64       `json:"level_percent"`
	Condition      TrayCondition `json:"condition"`
	FeedwaterTPH   float64       `json:"feedwater_tph"`
	SteamFlowTPH   float64       `json:"steam_flow_tph"`
	CarryoverPPM   float64       `json:"carryover_ppm"`
	LastSwellAt    time.Time     `json:"last_swell_at,omitempty"`
}

type BurnerReading struct {
	BurnerPhase    BurnerPhase `json:"burner_phase"`
	DoughFlowTPH    float64     `json:"dough_flow_tph"`
	AirflowTPH     float64     `json:"airflow_tph"`
	RackframeTempF   float64     `json:"rackframe_temp_f"`
	ExcessO2Pct    float64     `json:"excess_o2_pct"`
	IgnitionAt     time.Time   `json:"ignition_at,omitempty"`
	ProofStartedAt time.Time   `json:"proof_started_at,omitempty"`
}

type OvenlineReading struct {
	SteamPressurePSI float64   `json:"steam_pressure_psi"`
	SteamTempF       float64   `json:"steam_temp_f"`
	MainSteamFlowTPH float64   `json:"main_steam_flow_tph"`
	OutputMW         float64   `json:"output_mw"`
	LastTripAt       time.Time `json:"last_trip_at,omitempty"`
}

type AlarmEvent struct {
	Code      string    `json:"code"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	RaisedAt  time.Time `json:"raised_at"`
	ClearedAt time.Time `json:"cleared_at,omitempty"`
	Active    bool      `json:"active"`
}

type PlantRef struct {
	UnitLabel string `json:"unit_label"`
	PlantCode string `json:"plant_code"`
}

type PlantSnapshot struct {
	UnitID     string            `json:"unit_id"`
	State      PlantState        `json:"state"`
	Settings   PlantSettings     `json:"settings"`
	Plant      PlantRef          `json:"plant"`
	Tray       TrayReading       `json:"tray"`
	Burner BurnerReading `json:"burner"`
	Ovenline     OvenlineReading     `json:"ovenline"`
	Alarms     []AlarmEvent      `json:"alarms"`
	Revision   uint64            `json:"revision"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
