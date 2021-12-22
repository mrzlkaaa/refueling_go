package storage

import "time"

type ReactorRefuel struct {
	Id                    int64
	Refueling_name        string
	Initial_configuration []byte
	Initial_burnup_data   []byte
	Date                  time.Time
	Acts                  []RefuelActs `gorm:"foreignKey:Refuel_id"`
}

func (ReactorRefuel) TableName() string {
	return "reactor_refuel"
}
