package NoSQL

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelCycle struct {
	_ID          primitive.ObjectID
	Name        string
	TotalTime   float64
	TotalEnOuts float64
	WeeklyData  []WeeklyData
}

type WeeklyData struct {
	WeekName    int32	
	TotalTime   float64
	TotalEnOuts float64
	DetailWeek  []DetailWeek
}

type DetailWeek struct {
	Power        float64
	FromDate     string
	ToDate       string
	Time         float64
	EnergyOutput float64
}
