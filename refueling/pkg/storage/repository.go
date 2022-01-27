package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage() *Storage {
	dsn := fmt.Sprintf("host=irt-t.ru user=postgres password=postgres dbname=%v port=5433 sslmode=disable", "irt_refueling")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Storage{db: db}
}

func (s *Storage) GetRefuelNamesQuery() []ReactorRefuel {
	var refuels []ReactorRefuel
	// refuels.TableName()
	s.db.Select("id", "refueling_name", "date").Find(&refuels)
	return refuels
}
