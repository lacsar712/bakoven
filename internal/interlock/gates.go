package interlock

import (
	"fmt"

	"github.com/lacsar712/bakoven/internal/model"
)

type PermissiveSet struct {
	doughOK       bool
	ignitionOK   bool
	trayOK       bool
	pressureOK   bool
	burnerOK bool
}

func NewPermissiveSet() *PermissiveSet { return &PermissiveSet{} }

func (p *PermissiveSet) SetDough(ok bool)       { p.doughOK = ok }
func (p *PermissiveSet) SetIgnition(ok bool)   { p.ignitionOK = ok }
func (p *PermissiveSet) SetTray(ok bool)       { p.trayOK = ok }
func (p *PermissiveSet) SetPressure(ok bool)   { p.pressureOK = ok }
func (p *PermissiveSet) SetBurner(ok bool) { p.burnerOK = ok }

func (p *PermissiveSet) DoughOK() bool       { return p.doughOK }
func (p *PermissiveSet) IgnitionOK() bool   { return p.ignitionOK }
func (p *PermissiveSet) TrayOK() bool       { return p.trayOK }
func (p *PermissiveSet) PressureOK() bool   { return p.pressureOK }
func (p *PermissiveSet) BurnerOK() bool { return p.burnerOK }

func (p *PermissiveSet) AllFiring() bool {
	return p.doughOK && p.ignitionOK && p.trayOK && p.pressureOK && p.burnerOK
}

func (p *PermissiveSet) CheckIgnition() error {
	if !p.doughOK {
		return fmt.Errorf("%w", model.ErrDoughPermissive)
	}
	if !p.ignitionOK {
		return fmt.Errorf("%w", model.ErrIgnitionBlocked)
	}
	return nil
}

func CheckBakeLoss(reading model.BurnerReading) error {
	if reading.BurnerPhase == model.BurnerStable && reading.RackframeTempF < 600 {
		return fmt.Errorf("%w", model.ErrBakeLoss)
	}
	return nil
}

func (p *PermissiveSet) CheckFiring() error {
	if err := p.CheckIgnition(); err != nil {
		return err
	}
	if !p.trayOK {
		return fmt.Errorf("%w", model.ErrTrayLevelTrip)
	}
	if !p.pressureOK {
		return fmt.Errorf("%w", model.ErrPressureTrip)
	}
	if !p.burnerOK {
		return fmt.Errorf("%w", model.ErrBurnerTrip)
	}
	return nil
}

type CoordinationLock struct {
	holder string
	held   bool
}

func NewCoordinationLock() *CoordinationLock { return &CoordinationLock{} }

func (c *CoordinationLock) Acquire(holder string) error {
	if c.held {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	c.holder = holder
	c.held = true
	return nil
}

func (c *CoordinationLock) Release(holder string) {
	if c.held && c.holder == holder {
		c.held = false
		c.holder = ""
	}
}

func (c *CoordinationLock) Require(holder string) error {
	if !c.held || c.holder != holder {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	return nil
}

func (c *CoordinationLock) Held() bool { return c.held }
