package adding

type Refuel struct {
	RefuelName int    `json:"name"`
	Date       string `json:"date"`
	Acts       []Act  `json:"acts"`
}

type Act struct {
	Name          string     `json:"fileName"`
	CoreConfig    [][]string `json:"map"` //* separate from data that is aimed to store?
	PDC           []string   `json:"pdc"` //* separate from data that is aimed to store?
	Description   string     `json:"description"`
	RefuelNameRef int        `json:"refuelNameRef"`
}
