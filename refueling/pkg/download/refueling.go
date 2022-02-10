package download

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
