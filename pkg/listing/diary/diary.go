package listingDiary

import "go.mongodb.org/mongo-driver/bson/primitive"

//* structs below more suitable for listing section
type FuelCycle struct {
	ID          primitive.ObjectID
	Name        string
	TotalTime   float64
	TotalEnOuts float64
	WeeklyOuts  []WeeklyData
}

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
