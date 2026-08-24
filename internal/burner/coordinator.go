package burner

import (
	"context"
	"fmt"
	"math"

	"github.com/lacsar712/bakoven/internal/clock"
	"github.com/lacsar712/bakoven/internal/model"
)

type Coordinator struct {
	clk     clock.ProcessClock
	burner  *BurnerController
	airflow *AirflowBalancer
	dough    *DoughRegulator
	proof   *clock.ProofWindow
	ignition *clock.IgnitionDelayWindow
	warmup  *clock.BurnerWarmupWindow
}

func NewCoordinator(clk clock.ProcessClock) *Coordinator {
	return &Coordinator{
		clk:      clk,
		burner:   NewBurnerController(clk),
		airflow:  NewAirflowBalancer(clk),
		dough:     NewDoughRegulator(clk),
		proof:    clock.NewProofWindow(clk),
		ignition: clock.NewIgnitionDelayWindow(clk),
		warmup:   clock.NewBurnerWarmupWindow(clk),
	}
}

func (c *Coordinator) Burner() *BurnerController  { return c.burner }
func (c *Coordinator) Airflow() *AirflowBalancer { return c.airflow }
func (c *Coordinator) Dough() *DoughRegulator     { return c.dough }

func (c *Coordinator) StartProof(ctx context.Context, snap model.PlantSnapshot) (model.BurnerReading, error) {
	select {
	case <-ctx.Done():
		return snap.Burner, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerProof
	out.ProofStartedAt = c.clk.Now()
	out.DoughFlowTPH = 0
	out.AirflowTPH = c.airflow.ProofRate()
	return out, nil
}

func (c *Coordinator) CompleteProof(snap model.BurnerReading) error {
	return c.proof.Require(snap.ProofStartedAt)
}

func (c *Coordinator) Ignite(ctx context.Context, snap model.PlantSnapshot) (model.BurnerReading, error) {
	select {
	case <-ctx.Done():
		return snap.Burner, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if err := c.proof.Require(snap.Burner.ProofStartedAt); err != nil {
		return snap.Burner, err
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerIgnition
	out.IgnitionAt = c.clk.Now()
	out.DoughFlowTPH = c.dough.IgnitionRate(snap.Settings)
	out.AirflowTPH = c.airflow.IgnitionRate(snap.Settings)
	out.RackframeTempF = 400
	return out, nil
}

func (c *Coordinator) Stabilize(snap model.PlantSnapshot) (model.BurnerReading, error) {
	if err := c.ignition.Require(snap.Burner.IgnitionAt); err != nil {
		return snap.Burner, err
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerStable
	out.DoughFlowTPH = snap.Settings.DoughFlowTPH * 0.5
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.RackframeTempF = c.burner.EstimateRackframeTemp(out)
	return out, nil
}

func (c *Coordinator) RampToLoad(snap model.PlantSnapshot, loadPct float64) model.BurnerReading {
	out := snap.Burner
	out.DoughFlowTPH = snap.Settings.DoughFlowTPH * loadPct
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.RackframeTempF = c.burner.EstimateRackframeTemp(out)
	return out
}

func (c *Coordinator) Trip(snap model.BurnerReading) model.BurnerReading {
	out := snap
	out.BurnerPhase = model.BurnerTrip
	out.DoughFlowTPH = 0
	out.RackframeTempF = math.Max(200, out.RackframeTempF*0.5)
	return out
}

func (c *Coordinator) WarmupReady(snap model.BurnerReading) bool {
	return c.warmup.Ready(snap.IgnitionAt)
}
