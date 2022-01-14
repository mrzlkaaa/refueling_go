package NoSQL

type WeeklyData struct {
	Name        int32
	TotalTime   float64
	TotalEnOuts float64
	Detail      []DetailWeek
	FCbackref   string
}

type DetailWeek struct {
	Power        float64
	FromDate     string
	ToDate       string
	Time         float64
	EnergyOutput float64
}
