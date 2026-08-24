package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/bakoven/internal/model"
)

func (a *App) WarmupStatus() (ready bool, detail string) {
	snap := a.Snapshot()
	if snap.Burner.ProofStartedAt.IsZero() {
		return false, "proof not started"
	}
	if !a.proofWindow.Ready(snap.Burner.ProofStartedAt) {
		return false, "proof window open"
	}
	if !snap.Burner.IgnitionAt.IsZero() && !a.warmupWindow.Ready(snap.Burner.IgnitionAt) {
		return false, "burner warmup window open"
	}
	if !snap.Tray.LastSwellAt.IsZero() {
		if err := a.tray.RequireSettled(snap.Tray); err != nil {
			return false, "tray swell settling"
		}
	}
	return true, "ready"
}

func (a *App) WaitWarmup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", model.ErrContextDone)
		default:
		}
		ready, _ := a.WarmupStatus()
		if ready {
			return nil
		}
	}
}

func (a *App) ProofRemaining() string {
	snap := a.Snapshot()
	if snap.Burner.ProofStartedAt.IsZero() {
		return "not started"
	}
	if a.proofWindow.Ready(snap.Burner.ProofStartedAt) {
		return "complete"
	}
	return "in progress"
}

func (a *App) BurnerWarmupRemaining() string {
	snap := a.Snapshot()
	if snap.Burner.IgnitionAt.IsZero() {
		return "not ignited"
	}
	if a.warmupWindow.Ready(snap.Burner.IgnitionAt) {
		return "complete"
	}
	return "in progress"
}
