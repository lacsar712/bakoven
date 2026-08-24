package api

import (
	"errors"

	"github.com/lacsar712/bakoven/internal/model"
)

func classifyTrayError(err error) (string, bool) {
	if errors.Is(err, model.ErrTrayLevelLow) {
		return "tray_level_low", true
	}
	return "", false
}
