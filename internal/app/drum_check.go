package app

import (
	"fmt"

	"github.com/lacsar712/bakoven/internal/model"
)

func (a *App) CheckTrayLevel(snap model.PlantSnapshot) error {
	if snap.Tray.LevelPercent < model.MinTrayLevelPercent {
		return fmt.Errorf("%w", model.ErrTrayLevelLow)
	}
	return nil
}
