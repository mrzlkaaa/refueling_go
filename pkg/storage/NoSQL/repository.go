package NoSQL

import (
	"context"
	"fmt"
	"log"
	"errors"
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

func (s *Storage) FCExistingCheck(name string) error {
	filter := bson.M{
		"name": bson.M{"$eq": name},
	}

	var results bson.D
	err := s.db.FindOne(context.Background(), filter).Decode(&results)
	fmt.Println(err)
	return err
}

func (s *Storage) WeekNameExistingCheck(name string, weekName int32) error {
	if weekName == 1 {
		return errors.New("empty template")
	}
	
	filter := bson.M{
		"name": bson.M{"$eq": name},
		"weeklydata": bson.M{"$elemMatch": 
			bson.M{"weekname":weekName}},
	}

	var results bson.D
	err := s.db.FindOne(context.Background(), filter).Decode(&results)
	fmt.Println(err)
	return err
}

//* will be callced from upper layer if check for existing instance fails
func (s *Storage) CreateDBInstance(name string) {
	var newFCinstance FuelCycle
	var newWDinstance WeeklyData = WeeklyData{WeekName:1}
	var newDWinstance DetailWeek

	newFCinstance.Name = name
	newFCinstance.WeeklyData = append(newFCinstance.WeeklyData,
		newWDinstance)

	newFCinstance.WeeklyData[0].DetailWeek = append(newFCinstance.WeeklyData[0].DetailWeek,
		newDWinstance)

	fmt.Println("creating new instance", newFCinstance)
	s.db.InsertOne(context.Background(), newFCinstance)

}

func (s *Storage) AddWeekTemplate(name string, weekName int32) {

	filter := bson.M{
		"name": bson.M{"$eq": name},
	}

	update := bson.D{
		{"$push", bson.D{{"weeklydata", WeeklyData{WeekName: weekName}}}},
	}

	s.db.UpdateOne(context.Background(),filter, update)
}

func (s *Storage) AddWeeklyData(formsData *adding.FuelCycle) {

	//* its actually update of existing FC by adding new weekly data
	var FC FuelCycle
	var totalTime float64
	var totalEnOut float64

	FC.Name = formsData.Name
	
	for _, v := range formsData.DetailWeek {
		totalTime += v.Time
		totalEnOut += v.EnergyOutput
	}

	filter := bson.M{
		"name": bson.M{"$eq": FC.Name},
		"weeklydata": bson.M{"$elemMatch": 
			bson.M{"weekname":formsData.WeekName}},
	}

	update := bson.D{
		{"$inc",  bson.D{
			{"weeklydata.$.totaltime", totalTime},
			{"weeklydata.$.totalenouts", totalEnOut},
			{"totaltime", totalTime},
			{"totalenouts", totalEnOut},
		}},
		{"$set", bson.D{{"weeklydata.$.detailweek", formsData.DetailWeek}}},
	}

	_, err := s.db.UpdateOne(context.TODO(), filter, update); if err != nil {
		panic(err)
	}
}

func (s *Storage) GetNewWeekNum(name string) int32 {
	var result FuelCycle
	err := s.db.FindOne(context.Background(), bson.D{
		{"name", name}}).Decode(&result); if err != nil {
			fmt.Println(err)
		}

	fmt.Println(result)
	if result.TotalTime > 0.0 {
		//* getting the highest val of week
		var maxValue int32
		for _, v := range result.WeeklyData {
			if v.WeekName > maxValue {
				maxValue = v.WeekName
				fmt.Println(maxValue)
			}
		}
		return maxValue+1
	}
	return 1
}
