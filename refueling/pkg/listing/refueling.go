package listing

type Refuel struct {
	ID         uint `gorm:"primaryKey"`
	RefuelName int  `gorm:"unique"`
	Date       string
	Acts       []Act `gorm:"foreignKey:RefuelName;references:RefuelName"`
}

type Act struct {
	ID         uint `gorm:"foreignKey:RefuelName"`
	Name       string
	CoreConfig [][]string
	// PDC           []string
	Description   string
	RefuelNameRef int
}
