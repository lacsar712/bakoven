package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/bakoven/internal/model"
)

type OvenlineFSM struct {
	mu            sync.RWMutex
	state         model.PlantState
	doughPermissive bool
	proofComplete  bool
	hooks          *HookChain
}

func NewOvenlineFSM(unitID string) *OvenlineFSM {
	_ = unitID
	return &OvenlineFSM{state: model.StateColdStandby, hooks: NewHookChain()}
}

func (f *OvenlineFSM) Hooks() *HookChain { return f.hooks }

func (f *OvenlineFSM) State() model.PlantState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state
}

func (f *OvenlineFSM) SetDoughPermissive(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.doughPermissive = ok
}

func (f *OvenlineFSM) SetProofComplete(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.proofComplete = ok
}

func (f *OvenlineFSM) DoughPermissive() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.doughPermissive
}

func (f *OvenlineFSM) Dispatch(ctx context.Context, event PlantEvent) (model.PlantState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	select {
	case <-ctx.Done():
		return f.state, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if event == EvTrip {
		from := f.state
		if f.hooks != nil {
			if err := f.hooks.RunBefore(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		f.state = model.StateTrip
		if f.hooks != nil {
			if err := f.hooks.RunAfter(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		return f.state, nil
	}
	next, ok := NextState(f.state, event)
	if !ok {
		return f.state, fmt.Errorf("%s from %s: %w", event, f.state, ErrIllegalTransition)
	}
	if event == EvIgnite && !f.doughPermissive {
		return f.state, fmt.Errorf("%w", model.ErrDoughPermissive)
	}
	if event == EvProofComplete && !f.proofComplete {
		return f.state, fmt.Errorf("%w", model.ErrProofIncomplete)
	}
	from := f.state
	if f.hooks != nil {
		if err := f.hooks.RunBefore(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	f.state = next
	if f.hooks != nil {
		if err := f.hooks.RunAfter(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	return f.state, nil
}

func (f *OvenlineFSM) ForceState(state model.PlantState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = state
}
