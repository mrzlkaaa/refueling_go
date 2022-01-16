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
	db *mongo.Collection
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
	collection := client.Database("enOutDiary").Collection("Diary")

	return &Storage{db: collection}
}

func (s *Storage) AddWeeklyData(formsData *adding.FuelCycle) {

	//* its actually update of existing FC by adding new weekly data
	var FC FuelCycle
	// var sumTime float64
	// var sumEnOut float64

	// fmt.Println(FC.DetailWeek)

	FC.Name = formsData.Name

	// FC.WeekName = formsData.WeekName
	// for k, v := range formsData.DetailWeek {
	// 	FC.WeeklyData[k].DetailWeek = append(FC.WeeklyData[k].DetailWeek,
	// 		DetailWeek{
	// 			Power:        v.Power,
	// 			FromDate:     v.FromDate,
	// 			ToDate:       v.ToDate,
	// 			Time:         v.Time,
	// 			EnergyOutput: v.EnergyOutput,
	// 		})

	// 	sumTime += v.Time
	// 	sumEnOut += v.EnergyOutput
	// }

	// if err := curr.All(context.TODO(), &results); err != nil {
	// 	panic(err)
	// }

}

func (s *Storage) FCExistingCheck(name string) error {
	filter := bson.M{
		"Name": bson.M{"$eq": name},
	}

	var results bson.D
	err := s.db.FindOne(context.Background(), filter).Decode(&results)
	return err
}

//* will be callced from upper layer if check for existing instance fails
func (s *Storage) CreateDBInstance(name string) {
	var newFCinstance FuelCycle
	var newWDinstance WeeklyData

	newFCinstance.Name = name
	newFCinstance.WeeklyData = append(newFCinstance.WeeklyData,
		newWDinstance)

	newFCinstance.WeeklyData[0].DetailWeek = append(newFCinstance.WeeklyData[0].DetailWeek,
		DetailWeek{
			Power:        0.0,
			FromDate:     "",
			ToDate:       "",
			Time:         0.0,
			EnergyOutput: 0.0,
		})

	fmt.Println("creating new instance", newFCinstance)
	s.db.InsertOne(context.Background(), newFCinstance)

}

// func (s *Storage) AddWeeklyData(formsData *adding.FuelCycle) {

// 	var FC FuelCycle
// 	var sumTimeFC float64
// 	var sumEnOutFC float64

// 	FC.Name = formsData.Name
// 	for k, v := range formsData.WeeklyData {
// 		fmt.Println(v.Name)
// 		FC.WeeklyData[k].Name = v.Name

// 		var sumTimeW float64
// 		var sumEnOutW float64

// 		for _, vv := range v.DetailWeek {
// 			FC.WeeklyData[k].DetailWeek = append(FC.WeeklyData[k].DetailWeek,
// 				DetailWeek{
// 					Power:        vv.Power,
// 					FromDate:     vv.FromDate,
// 					ToDate:       vv.ToDate,
// 					Time:         vv.Time,
// 					EnergyOutput: vv.EnergyOutput,
// 				})

// 			sumTimeW += vv.Time
// 			sumEnOutW += vv.EnergyOutput
// 		}

// 	s.db.InsertOne(context.TODO(), FC)
// }

func (s *Storage) ReadWeeklyData(name int32) {
	curr, err := s.db.Find(context.Background(), bson.D{
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
	curr, err := s.db.Find(context.Background(), bson.D{
		{"name", name},
	})

	if err != nil {
		panic(err)
	}

	var result bson.D
	curr.Decode(&result)
	fmt.Println(result)
	if len(result) > 1 {
		return 10000
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
