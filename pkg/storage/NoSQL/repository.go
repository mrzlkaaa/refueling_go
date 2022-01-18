package NoSQL

import (
	"context"
	"errors"
	"fmt"
	"log"
	"refueling/pkg/adding"
	listingDiary "refueling/pkg/listing/diary"

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

func (s *Storage) WeekNameExistingCheck(name string, weekName int) error {
	if weekName == 1 {
		return errors.New("empty template")
	}

	filter := bson.M{
		"name":       bson.M{"$eq": name},
		"weeklydata": bson.M{"$elemMatch": bson.M{"weekname": weekName}},
	}

	var results bson.D
	err := s.db.FindOne(context.Background(), filter).Decode(&results)
	fmt.Println(err)
	return err
}

//* will be callced from upper layer if check for existing instance fails
func (s *Storage) CreateDBInstance(name string) {
	var newFCinstance FuelCycle
	var newWDinstance WeeklyData = WeeklyData{WeekName: 1}
	var newDWinstance DetailWeek

	newFCinstance.Name = name
	newFCinstance.WeeklyData = append(newFCinstance.WeeklyData,
		newWDinstance)

	newFCinstance.WeeklyData[0].DetailWeek = append(newFCinstance.WeeklyData[0].DetailWeek,
		newDWinstance)

	fmt.Println("creating new instance", newFCinstance)
	s.db.InsertOne(context.Background(), newFCinstance)

}

//! remove func
func (s *Storage) AddWeekTemplate(name string, weekName int) {

	filter := bson.M{
		"name": bson.M{"$eq": name},
	}

	update := bson.D{
		{"$push", bson.D{{"weeklydata", WeeklyData{WeekName: weekName}}}},
	}

	s.db.UpdateOne(context.Background(), filter, update)
}

func (s *Storage) AddWeeklyData(formsData *adding.FuelCycle) {

	var FC FuelCycle
	var WD WeeklyData
	// var DW []DetailWeek
	DW := make([]DetailWeek, len(formsData.DetailWeek))

	FC.Name = formsData.Name
	WD.WeekName = formsData.WeekName

	fmt.Println(len(formsData.DetailWeek))

	for k, v := range formsData.DetailWeek {

		WD.TotalTime += v.Time
		WD.TotalEnOuts += v.EnergyOutput

		DW[k].Power = v.Power
		DW[k].FromDate = v.FromDate
		DW[k].ToDate = v.ToDate
		DW[k].Time = v.Time
		DW[k].EnergyOutput = v.EnergyOutput

		WD.DetailWeek = append(WD.DetailWeek, DW[k])
	}

	FC.TotalTime = WD.TotalTime
	FC.TotalEnOuts = WD.TotalEnOuts

	FC.WeeklyData = append(FC.WeeklyData, WD)

	fmt.Println(FC)

	s.db.InsertOne(context.Background(), FC)

}

func (s *Storage) GetWeeksNum(fcName string) []int {
	var result FuelCycle
	err := s.db.FindOne(context.Background(), bson.D{
		{"name", fcName}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return []int{1}
	}

	fmt.Println(result)

	var values []int
	var lastValue int
	for _, v := range result.WeeklyData {
		values = append(values, v.WeekName)
		lastValue = v.WeekName + 1
	}
	values = append(values, lastValue)
	return values
}

func (s *Storage) WeekDetails(fcName string, weekName int) []listingDiary.DetailWeek {

	filter := bson.M{
		"name":       bson.M{"$eq": fcName},
		"weeklydata": bson.M{"$elemMatch": bson.M{"weekname": weekName}},
	}

	var result listingDiary.FuelCycle
	err := s.db.FindOne(context.Background(), filter).Decode(&result)
	fmt.Println(err, result)

	//* retrieving only WeekDetails for specific week
	for _, v := range result.WeeklyData {
		if v.Name == weekName {
			return v.DetailWeek
		}
	}

	return []listingDiary.DetailWeek{}
}
