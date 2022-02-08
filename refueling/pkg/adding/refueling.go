package adding

type Refuel struct {
	RefuelName int    `json:"name"`
	Date       string `json:"date"`
	Acts       []Act  `json:"acts"`
}

type Act struct {
	Name        string     `json:"fileName"`
	CoreConfig  [][]string `json:"map"`
	PDC         []string   `json:"pdc"`
	Description string     `json:"description"`
	RefuelID    int        `json:"refuelId"`
}
