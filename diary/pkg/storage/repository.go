package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"refueling/diary/pkg/adding"
	"refueling/diary/pkg/listing"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Collection
}

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}

func NewStorage() *Storage {
	LoadEnv()
	URI := fmt.Sprintf("mongodb://%v:%v", os.Getenv("HOST"), os.Getenv("PORT"))
	clientOptions := options.Client().ApplyURI(URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("COLLECTION"))

	return &Storage{db: collection}
}

func (s *Storage) FCExistingCheck(fcName int) error {
	filter := bson.M{
		"name": bson.M{"$eq": fcName}}

	var result bson.D
	err := s.db.FindOne(context.Background(), filter).Decode(&result)
	return err
}

func (s *Storage) WeekNameExistingCheck(fcName int, weekName int) error {
	filter := bson.M{
		"name":       bson.M{"$eq": fcName},
		"weeklydata": bson.M{"$elemMatch": bson.M{"weekname": weekName}},
	}

	var result FuelCycle
	err := s.db.FindOne(context.Background(), filter).Decode(&result)
	fmt.Println(err)
	return err
}

func (s *Storage) AddWeeklyData(formsData *adding.FuelCycle) {

	var FC FuelCycle
	var WD WeeklyData
	var DW DetailWeek
	var RP RodsPosition

	FC.Name = formsData.Name
	WD.WeekName = formsData.WeekName

	fmt.Println(formsData)

	fmt.Println(len(formsData.DetailWeek))

	PopulatingWD_DW(&WD, &DW, &RP, formsData)

	FC.TotalTime += WD.TotalTime
	FC.TotalEnOuts += WD.TotalEnOuts

	FC.WeeklyData = append(FC.WeeklyData, WD)

	fmt.Println("added new: ", FC)

	// s.db.InsertOne(context.Background(), FC)

}

func (s *Storage) AppendWeeklyData(formsData *adding.FuelCycle) {

	var WD WeeklyData
	var DW DetailWeek
	var RP RodsPosition
	WD.WeekName = formsData.WeekName
	PopulatingWD_DW(&WD, &DW, &RP, formsData)

	filter := bson.M{
		"name": bson.M{"$eq": formsData.Name},
	}
	update := bson.M{
		"$inc": bson.D{
			{"totaltime", WD.TotalTime},
			{"totalenouts", WD.TotalEnOuts},
		},
		"$push": bson.D{{"weeklydata", WD}},
	}
	fmt.Println(formsData)
	fmt.Println("appended new: ", WD)
	s.db.UpdateOne(context.Background(), filter, update)
}

func (s *Storage) UpdateWeeklyData(formsData *adding.FuelCycle) {

	var FCcurrent FuelCycle
	var WD WeeklyData
	var DW DetailWeek
	var RP RodsPosition
	PopulatingWD_DW(&WD, &DW, &RP, formsData)

	arrayIndex := formsData.WeekName - 1
	fmt.Println("weekanme from form is :", formsData.WeekName)

	filter := bson.M{
		"name":       bson.M{"$eq": formsData.Name},
		"weeklydata": bson.M{"$elemMatch": bson.M{"weekname": formsData.WeekName}},
	}

	s.db.FindOne(context.Background(), filter).Decode(&FCcurrent)

	//* calc diff in time and enOuts for FC
	//* replace existing DetailWeek on new one
	WD.WeekName = formsData.WeekName
	FCcurrent.TotalTime += WD.TotalTime - FCcurrent.WeeklyData[arrayIndex].TotalTime
	FCcurrent.TotalEnOuts += WD.TotalEnOuts - FCcurrent.WeeklyData[arrayIndex].TotalEnOuts
	FCcurrent.WeeklyData[arrayIndex] = WD

	fmt.Println(FCcurrent.WeeklyData[arrayIndex])

	update := bson.M{
		"$set": bson.D{
			{"totaltime", FCcurrent.TotalTime},
			{"totalenouts", FCcurrent.TotalEnOuts},
			{"weeklydata", FCcurrent.WeeklyData},
		},
	}
	fmt.Println("updated to new: ", FCcurrent)
	if res, err := s.db.UpdateOne(context.Background(), filter, update); err != nil {
		fmt.Println(res, err)
	}

}

func (s Storage) OverallData() []listing.FuelCycle {
	var FC []listing.FuelCycle
	// var WD []WeeklyData
	// var DW []DetailWeek

	// FC.WeeklyData = append(FC.WeeklyData, WD)
	opts := options.Find().SetSort(bson.D{{"name", -1}})
	crs, err := s.db.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	if err = crs.All(context.Background(), &FC); err != nil {
		panic(err)
	}
	fmt.Println(FC)
	return FC
}

func (s *Storage) GetWeeksNum(fcName int) map[string]int {
	var result FuelCycle
	err := s.db.FindOne(context.Background(), bson.D{
		{"name", fcName}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		new := make(map[string]int)
		new["New"] = 1
		return new
	}

	fmt.Println(result)

	values := make(map[string]int)
	var started string
	var finished string
	var lastValue int
	for _, v := range result.WeeklyData {
		started = v.DetailWeek[0].FromDate
		finished = v.DetailWeek[len(v.DetailWeek)-1].ToDate
		values[fmt.Sprintf("%v - %v", started, finished)] = v.WeekName
		// values = append(values, v.WeekName)
		lastValue = v.WeekName + 1
	}
	values["New"] = lastValue
	// values = append(values, lastValue)
	return values
}

func (s *Storage) WeekDetails(fcName int, weekName int) listing.WeeklyData {

	filter := bson.M{
		"name":       bson.M{"$eq": fcName},
		"weeklydata": bson.M{"$elemMatch": bson.M{"weekname": weekName}},
	}

	var result listing.FuelCycle
	err := s.db.FindOne(context.Background(), filter).Decode(&result)
	fmt.Println(err, result)

	//* retrieving only WeekDetails for specific week
	for i, v := range result.WeeklyData {
		if v.Name == weekName {
			return result.WeeklyData[i]
		}
	}
	return listing.WeeklyData{} //* change return to listing.WeeklyData ---> that already include []listing.DetailWeek{}
}

func PopulatingWD_DW(WD *WeeklyData, DW *DetailWeek, RP *RodsPosition, formsData *adding.FuelCycle) {

	RP.AR = formsData.RodsPosition.AR
	RP.KS1 = formsData.RodsPosition.KS1
	RP.KS2 = formsData.RodsPosition.KS2
	RP.KS3 = formsData.RodsPosition.KS3
	RP.Temp = formsData.RodsPosition.Temp
	WD.RodsPosition = *RP
	for _, v := range formsData.DetailWeek {

		WD.TotalTime += v.Time
		WD.TotalEnOuts += v.EnergyOutput

		DW.Power = v.Power
		DW.FromDate = v.FromDate
		DW.ToDate = v.ToDate
		DW.Time = v.Time
		DW.EnergyOutput = v.EnergyOutput

		WD.DetailWeek = append(WD.DetailWeek, *DW)
	}
}
