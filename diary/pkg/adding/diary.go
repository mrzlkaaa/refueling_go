package adding

type FuelCycle struct {
	Name         int          `json:"fcName"`
	WeekName     int          `json:"week"`
	RodsPosition RodsPosition `json:"rodsPosition"`
	DetailWeek   []DetailWeek `json:"weeklyDetail"`
}

type DetailWeek struct {
	Power        float64 `json:"power"`
	FromDate     string  `json:"fromDate"`
	ToDate       string  `json:"toDate"`
	Time         float64 `json:"totalHours"`
	EnergyOutput float64 `json:"energyOutput"`
}

type RodsPosition struct {
	AR   int     `json:"AR"`
	KS1  int     `json:"KS1"`
	KS2  int     `json:"KS2"`
	KS3  int     `json:"KS3"`
	Temp float64 `json:"Temp"`
}
