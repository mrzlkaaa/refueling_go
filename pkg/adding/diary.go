package adding

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//* structs below more suitable for listing section
// type FuelCycle struct {
// 	ID         primitive.ObjectID `bson:"_id, omitempty"`
// 	Name       string             `json:"fcName"`
// 	WeeklyData []WeeklyData       `json:"weeklyOuts"`
// }

// type WeeklyData struct {
// 	Name       int32        `json:"week"`
// 	DetailWeek []DetailWeek `json:"weeklyDetail"`
// }

// type DetailWeek struct {
// 	Power        float64 `json:"power"`
// 	FromDate     string  `json:"fromDate"`
// 	ToDate       string  `json:"toDate"`
// 	Time         float64 `json:"totalHours"`
// 	EnergyOutput float64 `json:"energyOutput"`
// }

type FuelCycle struct {
	_ID         primitive.ObjectID `bson:"_id, omitempty"`
	Name       string             `json:"fcName"`
	WeekName   int32              `json:"week"`
	DetailWeek []DetailWeek       `json:"weeklyDetail"`
}

type DetailWeek struct {
	Power        float64 `json:"power"`
	FromDate     string  `json:"fromDate"`
	ToDate       string  `json:"toDate"`
	Time         float64 `json:"totalHours"`
	EnergyOutput float64 `json:"energyOutput"`
}
