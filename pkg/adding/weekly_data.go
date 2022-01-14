package adding

type WeeklyData struct {
	ID          int32
	Name        int32 `json:"week"`
	TotalTime   float64
	TotalEnOuts float64
	Detail      []DetailWeek
}

type FormsData struct {
	FCName     string       `json:"fcName"`
	WeekName   int32        `json:"week"`
	DetailWeek []DetailWeek `json:"weeklyDetail"`
}

type DetailWeek struct {
	Power        float64 `json:"power"`
	FromDate     string  `json:"fromDate"`
	ToDate       string  `json:"toDate"`
	Time         float64 `json:"totalHours"`
	EnergyOutput float64 `json:"energyOutput"`
}
