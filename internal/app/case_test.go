package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/lacsar712/bakoven/internal/app"
	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/config"
	"github.com/lacsar712/bakoven/internal/model"
)

func TestCase(t *testing.T) {
	clk := clock.NewManual(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	cfg := config.Default("BURN-1")
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	snap := a.Snapshot()
	comb := snap.Burner
	comb.IgnitionAt = clk.Now().Add(-model.BurnerWarmupWindow)
	if err := a.Store().UpdateBurner(cfg.UnitID, comb); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		_ = a.RunWarmupBurnScheduler(ctx, comb.IgnitionAt)
		close(done)
	}()
	clk.Advance(50 * time.Millisecond)
	cancel()
	<-done
	if a.SchedulerItemCount() != 0 {
		t.Fatalf("cancelled warmup should not queue burn plan items, got %d", a.SchedulerItemCount())
	}
}
