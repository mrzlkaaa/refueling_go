package NoSQL

import (
	"context"
	"fmt"
	"log"

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
	collection := client.Database("EnOutDiary").Collection("FuelCycle")

	fmt.Printf("Type of collection interface: %T\n", collection)

	return &Storage{db: collection}
}

func (s *Storage) AddFC() {
	// replace on incoming struct -->
	testInsertion := FuelCycle{
		Name:        "FC001",
		TotalTime:   144,
		TotalEnOuts: 432,
	}
	insertRes, err := s.db.InsertOne(context.TODO(), testInsertion)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertRes.InsertedID)
}

func (s *Storage) ReadFCOne() {
	var testRead FuelCycle
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	curr, err := s.db.Find(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	for curr.Next(context.Background()) {
		err := curr.Decode(&testRead)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(testRead)
}
