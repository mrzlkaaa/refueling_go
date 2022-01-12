package adding

type DetailWeek struct {
	Power        float64 `json:"power"`
	FromDate     string  `json:"fromDate"`
	ToDate       string  `json:"toDate"`
	Time         float64 `json:"totalHours"`
	EnergyOutput float64 `json:"energyOutput"`
}
