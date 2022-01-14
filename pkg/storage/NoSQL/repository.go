package NoSQL

import (
	"context"
	"fmt"
	"log"
	"refueling/pkg/adding"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	dbFC *mongo.Collection
	dbWD *mongo.Collection
}

func NewStorage() *Storage {
	clientOptions := options.Client().ApplyURI("mongodb://irt-t.ru:3000")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collectionFC := client.Database("EnOutDiary").Collection("FuelCycle")
	collectionWD := client.Database("EnOutDiary").Collection("WeeklyData")

	return &Storage{dbFC: collectionFC, dbWD: collectionWD}
}

func (s *Storage) AddWeeklyData(formsData *adding.FormsData) {

	var sumTime float64
	var sumEnOut float64
	var Details []DetailWeek

	for _, v := range formsData.DetailWeek {
		sumTime += v.Time
		sumEnOut += v.EnergyOutput
		Details = append(Details, DetailWeek{
			Power:        v.Power,
			FromDate:     v.FromDate,
			ToDate:       v.ToDate,
			Time:         v.Time,
			EnergyOutput: v.EnergyOutput,
		})
	}

	insert := WeeklyData{
		Name:        formsData.WeekName,
		TotalTime:   sumTime,
		TotalEnOuts: sumEnOut,
		Detail:      Details,
		FCbackref:   formsData.FCName,
	}

	s.dbWD.InsertOne(context.TODO(), insert)
}

func (s *Storage) ReadWeeklyData(name int32) {
	curr, err := s.dbWD.Find(context.Background(), bson.D{
		{"name", 1},
	})

	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err := curr.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println(results)

}

func (s *Storage) GetNewWeekNum(name string) int32 {
	curr, err := s.dbFC.Find(context.Background(), bson.D{
		{"name", name},
	})

	if err != nil {
		panic(err)
	}

	var result bson.D
	curr.Decode(&result)
	fmt.Println(result)
	if len(result) > 0 {
		return 1
	}
	return 1
}

// func (s *Storage) ReadFCOne() {
// 	var testRead FuelCycle
// 	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	// defer cancel()
// 	curr, err := s.db.Find(context.Background(), bson.D{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for curr.Next(context.Background()) {
// 		err := curr.Decode(&testRead)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	fmt.Println(testRead)
// }
