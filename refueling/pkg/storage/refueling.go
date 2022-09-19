package storage

type Refuel struct {
	ID         uint `gorm:"primaryKey"`
	RefuelName int  `gorm:"unique"`
	Date       string
	Acts       []Act `gorm:"foreignKey:RefuelNameRef;references:RefuelName"`
}

type Act struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string //*if merge PathTo and Name -->  will recieve FullPath to both Config and PDC
	Description   string
	RefuelNameRef int //* keeps reference to parent scheme
}
