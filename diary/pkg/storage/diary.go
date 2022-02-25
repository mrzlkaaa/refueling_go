package storage

type FuelCycle struct {
	Name        int
	TotalTime   float64
	TotalEnOuts float64
	WeeklyData  []WeeklyData
}

type WeeklyData struct {
	WeekName     int
	TotalTime    float64
	RodsPosition RodsPosition
	TotalEnOuts  float64
	DetailWeek   []DetailWeek
}

type DetailWeek struct {
	Power        float64
	FromDate     string
	ToDate       string
	Time         float64
	EnergyOutput float64
}

type RodsPosition struct {
	AR   int
	KS1  int
	KS2  int
	KS3  int
	Temp float64
}
