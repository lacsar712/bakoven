package store

import "github.com/lacsar712/bakoven/internal/model"

type TraySnapshotView struct {
	UnitID   string
	Tray     model.TrayReading
	Alarms   []model.AlarmEvent
	Revision uint64
}

func CloneTraySnapshot(s model.PlantSnapshot) TraySnapshotView {
	out := TraySnapshotView{
		UnitID:   s.UnitID,
		Tray:     s.Tray,
		Revision: s.Revision,
	}
	out.Alarms = make([]model.AlarmEvent, len(s.Alarms))
	copy(out.Alarms, s.Alarms)
	return out
}
