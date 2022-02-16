package adding

type FuelCycle struct {
	Name       int          `json:"fcName"`
	WeekName   int          `json:"week"`
	DetailWeek []DetailWeek `json:"weeklyDetail"`
}

type DetailWeek struct {
	Power        float64 `json:"power"`
	FromDate     string  `json:"fromDate"`
	ToDate       string  `json:"toDate"`
	Time         float64 `json:"totalHours"`
	EnergyOutput float64 `json:"energyOutput"`
}
