package adding

type Refuel struct {
	RefuelName int    `json:"name"`
	Date       string `json:"date"`
	Acts       []Acts `json:"acts"`
}

type Acts struct {
	Name        string     `json:"fileName"`
	CoreConfig  [][]string `json:"map"`
	PDC         []string   `json:"pdc"`
	Description string     `json:"description"`
}
