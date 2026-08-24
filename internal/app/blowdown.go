package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/bakoven/internal/model"
)

const maxSteamrelOpeningPct = 100.0

func (a *App) OpenSteamrel(ctx context.Context, holder string, openingPct float64) error {
	_ = holder
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if openingPct >= maxSteamrelOpeningPct {
		return fmt.Errorf("steamrel: %w", model.ErrSteamrelLimit)
	}
	return nil
}

func (a *App) SteamrelAfterShutdown(ctx context.Context, openingPct float64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	snap := a.Snapshot()
	if snap.State != model.StateTrip && snap.State != model.StateColdStandby {
		return fmt.Errorf("plant not shut down")
	}
	if openingPct >= maxSteamrelOpeningPct {
		return fmt.Errorf("steamrel: %w", model.ErrSteamrelLimit)
	}
	return nil
}
