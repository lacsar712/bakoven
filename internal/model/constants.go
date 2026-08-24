package model

import "time"

const (
	DefaultLeaseTTL        = 30 * time.Second
	ProofWindow            = 5 * time.Minute
	IgnitionDelayWindow    = 15 * time.Second
	TraySwellSettleWindow  = 45 * time.Second
	BurnerWarmupWindow = 2 * time.Minute
	FeedwaterRampWindow    = 30 * time.Second
	MaxTrayLevelPercent    = 95.0
	MinTrayLevelPercent    = 15.0
	TripTrayLowPercent     = 10.0
	TripTrayHighPercent    = 98.0
	NormalSteamPressurePSI = 1800.0
	MaxSteamPressurePSI    = 2000.0
	MinRackframeO2Percent    = 2.5
	MaxRackframeO2Percent    = 6.0
	DefaultJournalCapacity = 512
)
