package app

import (
	"context"
	"fmt"
	"time"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

func (a *App) advanceClock(d time.Duration) {
	if mc, ok := a.clk.(*clock.ManualClock); ok {
		mc.Advance(d)
		time.Sleep(time.Millisecond)
	} else {
		time.Sleep(d)
	}
}

func (a *App) bindDoughLoop(holder string, ctx context.Context) context.Context {
	a.mu.Lock()
	if cancel, ok := a.doughLoopCancels[holder]; ok {
		cancel()
	}
	child, cancel := context.WithCancel(ctx)
	a.doughLoopCancels[holder] = cancel
	a.mu.Unlock()
	return child
}

func (a *App) cancelDoughLoop(holder string) {
	a.mu.Lock()
	if cancel, ok := a.doughLoopCancels[holder]; ok {
		cancel()
		delete(a.doughLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) cancelAllDoughLoops() {
	a.mu.Lock()
	for holder, cancel := range a.doughLoopCancels {
		cancel()
		delete(a.doughLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) CoalFeedTPH() float64 {
	return a.Snapshot().Burner.DoughFlowTPH
}

func (a *App) RunDoughRamp(ctx context.Context, holder string, targetTPH float64) error {
	loopCtx := a.bindDoughLoop(holder, ctx)
	defer a.cancelDoughLoop(holder)
	for {
		if err := loopCtx.Err(); err != nil {
			return fmt.Errorf("%w", model.ErrContextDone)
		}
		snap := a.Snapshot()
		current := snap.Burner.DoughFlowTPH
		if current >= targetTPH {
			return nil
		}
		comb := snap.Burner
		comb.DoughFlowTPH = current + 1.0
		_ = a.store.UpdateBurner(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.DoughFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
}

func (a *App) RunCoalFeed(ctx context.Context, holder string, steps int) error {
	loopCtx := a.bindDoughLoop(holder, ctx)
	defer a.cancelDoughLoop(holder)
	for i := 0; steps <= 0 || i < steps; i++ {
		if err := loopCtx.Err(); err != nil {
			return fmt.Errorf("%w", model.ErrContextDone)
		}
		snap := a.Snapshot()
		comb := snap.Burner
		comb.DoughFlowTPH += 0.5
		_ = a.store.UpdateBurner(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.DoughFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
	return nil
}
