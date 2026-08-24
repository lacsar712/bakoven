package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/bakoven/internal/ovenline"
	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/burner"
	"github.com/lacsar712/bakoven/internal/config"
	"github.com/lacsar712/bakoven/internal/tray"
	"github.com/lacsar712/bakoven/internal/fsm"
	"github.com/lacsar712/bakoven/internal/interlock"
	"github.com/lacsar712/bakoven/internal/model"
	"github.com/lacsar712/bakoven/internal/store"
)

type App struct {
	cfg           config.Config
	clk           clock.ProcessClock
	store         *store.PlantStore
	journal       *store.Journal
	fsm           *fsm.OvenlineFSM
	ovenline        *ovenline.Controller
	burner    *burner.Coordinator
	tray          *tray.Coordinator
	interlock     *interlock.Interlock
	permissives   *interlock.PermissiveSet
	coordLock     *interlock.CoordinationLock
	scheduler     *clock.Scheduler
	proofWindow   *clock.ProofWindow
	warmupWindow  *clock.BurnerWarmupWindow
	telemetry     *Telemetry
	tickCancels    map[string]context.CancelFunc
	doughLoopCancels map[string]context.CancelFunc
	mu             sync.RWMutex
}

func New(cfg config.Config, clk clock.ProcessClock) *App {
	return &App{
		cfg:          cfg,
		clk:          clk,
		store:        store.NewPlantStore(),
		journal:      store.NewJournal(cfg.JournalPath, cfg.JournalCapacity),
		fsm:          fsm.NewOvenlineFSM(cfg.UnitID),
		ovenline:       ovenline.NewController(clk),
		burner:   burner.NewCoordinator(clk),
		tray:         tray.NewCoordinator(clk),
		interlock:    interlock.NewInterlock(cfg.LeaseTTL),
		permissives:  interlock.NewPermissiveSet(),
		coordLock:    interlock.NewCoordinationLock(),
		scheduler:    clock.NewScheduler(clk),
		proofWindow:  clock.NewProofWindow(clk),
		warmupWindow: clock.NewBurnerWarmupWindow(clk),
		telemetry:    NewTelemetry(cfg.UnitID),
		tickCancels:     make(map[string]context.CancelFunc),
		doughLoopCancels: make(map[string]context.CancelFunc),
	}
}

func (a *App) Snapshot() model.PlantSnapshot {
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return model.DefaultSnapshot(a.cfg.UnitID)
	}
	return snap
}

func (a *App) Config() config.Config              { return a.cfg }
func (a *App) Clock() clock.ProcessClock          { return a.clk }
func (a *App) FSM() *fsm.OvenlineFSM                { return a.fsm }
func (a *App) UnitID() string                     { return a.cfg.UnitID }
func (a *App) Store() *store.PlantStore           { return a.store }
func (a *App) Interlock() *interlock.Interlock    { return a.interlock }
func (a *App) Telemetry() TelemetrySnapshot       { return a.telemetry.Snapshot() }
func (a *App) Journal() *store.Journal            { return a.journal }

func (a *App) journalEvent(ev, payload string) {
	_, _ = a.journal.Append(a.cfg.UnitID, ev, payload)
}

func (a *App) syncState(state model.PlantState) {
	_ = a.store.UpdateState(a.cfg.UnitID, state)
}

func (a *App) isFiring(state model.PlantState) bool {
	return state == model.StateFiring || state == model.StateLoadFollow || state == model.StateRamp
}

func (a *App) refreshPermissives(snap model.PlantSnapshot) {
	a.permissives.SetTray(a.tray.Level().WithinLimits(snap.Tray.LevelPercent))
	a.permissives.SetPressure(a.ovenline.Pressure().WithinTripLimits(snap.Ovenline.SteamPressurePSI, a.isFiring(snap.State)))
	a.permissives.SetBurner(a.burner.Burner().BakeStable(snap.Burner))
	a.permissives.SetDough(snap.Burner.DoughFlowTPH > 0 || snap.State == model.StateProof)
	a.permissives.SetIgnition(snap.Burner.BurnerPhase == model.BurnerStable || snap.Burner.BurnerPhase == model.BurnerIgnition)
	a.fsm.SetDoughPermissive(a.permissives.DoughOK())
	a.fsm.SetProofComplete(a.proofWindow.Ready(snap.Burner.ProofStartedAt))
}

func (a *App) tickLabel() string {
	return fmt.Sprintf("%s-tick", a.cfg.UnitID)
}
