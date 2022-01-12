package NoSQL

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FuelCycle struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Name        string
	TotalTime   float64
	TotalEnOuts float64
	WeeklyOuts  []WeeklyData
}
