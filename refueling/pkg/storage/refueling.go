package storage

type Refuel struct {
	ID         uint `gorm:"primaryKey"`
	RefuelName int
	Date       string
	Acts       []Act `gorm:"foreignKey:RefuelID"`
}

type Act struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	CoreConfig  []byte
	PDC         []byte
	Description string
	RefuelID    int
}

// type ReactorRefuel struct {
// 	Id                    int64
// 	Refueling_name        string
// 	Initial_configuration []byte
// 	Initial_burnup_data   []byte
// 	Date                  time.Time
// 	Acts                  []RefuelActs `gorm:"foreignKey:Refuel_id"`
// }

// func (ReactorRefuel) TableName() string {
// 	return "reactor_refuel"
// }

// type RefuelActs struct {
// 	Id                    int64
// 	Description           string
// 	Current_configuration byte
// 	Burnup_data           byte
// 	Refuel_id             int64
// }

// func (ReactorRefuel) TableName() string {
// 	return "reactor_refuel"
// }
