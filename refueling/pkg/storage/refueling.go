package storage

// type Refuel struct {
// 	ID         uint `gorm:"primaryKey"`
// 	RefuelName int
// 	Date       string
// 	Acts       []Act `gorm:"foreignKey:RefuelID"`
// }

// type Act struct {
// 	ID          uint `gorm:"primaryKey"`
// 	Name        string
// 	CoreConfig  []byte
// 	PDC         []byte
// 	Description string
// 	RefuelID    int
// }

type Refuel struct {
	ID         uint `gorm:"primaryKey"`
	RefuelName int
	Date       string
	PathTo     string
	Acts       []Act `gorm:"foreignKey:RefuelID"`
}

type Act struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string //*if merge PathTo and Name -->  will recieve FullPath to both Config and PDC
	CoreConfig  string //! remove
	PDC         string //! remove
	Description string
	RefuelID    int //* keep reference to parent scheme
}
