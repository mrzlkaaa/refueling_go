package listing

import "go.mongodb.org/mongo-driver/bson/primitive"

// type FuelCycle struct {
// 	ID          primitive.ObjectID
// 	Name        string
// 	TotalTime   float64
// 	TotalEnOuts float64
// 	WeeklyOuts  []WeeklyData
// }

// type WeeklyData struct {
// 	Name        int32
// 	TotalTime   float64
// 	TotalEnOuts float64
// 	Detail      []DetailWeek
// 	FCbackref   string
// }

// type DetailWeek struct {
// 	Power        float64
// 	FromDate     string
// 	ToDate       string
// 	Time         float64
// 	EnergyOutput float64
// }

//* structs below more suitable for listing section
type FuelCycle struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Name        string             `bson:"name"`
	TotalTime   float64            `bson:"totaltime"`
	TotalEnOuts float64            `bson:"totalenouts"`
	WeeklyData  []WeeklyData       `bson:"weeklydata"`
}

type WeeklyData struct {
	Name        int          `bson:"weekname"`
	TotalTime   float64      `bson:"totaltime"`
	TotalEnOuts float64      `bson:"totalenouts"`
	DetailWeek  []DetailWeek `bson:"detailweek"`
}

type DetailWeek struct {
	Power        float64 `bson:"power"`
	FromDate     string  `bson:"fromdate"`
	ToDate       string  `bson:"todate"`
	Time         float64 `bson:"time"`
	EnergyOutput float64 `bson:"energyoutput"`
}
