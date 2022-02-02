package storage

import (
	"fmt"
	"log"
	"os"
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}

func NewStorage() *Storage {
	LoadEnv()
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("PSQL_USER"),
		os.Getenv("PSWD"),
		os.Getenv("DB"),
		os.Getenv("PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Refuel{}, &Act{})
	return &Storage{db: db}
}

func (s *Storage) GetRefuelNames() []listing.Refuel {
	var refuels []listing.Refuel
	// refuels.TableName()
	s.db.Select("id", "refuel_name", "date").Find(&refuels)
	fmt.Println(refuels)
	return refuels
}

func (s *Storage) Adding(refuel adding.Refuel) error {
	var ref Refuel
	var act Act

	// fill ref
	ref.RefuelName = refuel.RefuelName
	ref.Date = refuel.Date

	// fill acts
	for _, v := range refuel.Acts {
		act.Name = v.Name
		act.Description = v.Description
		act.CoreConfig = FormatterCoreConfig(v.CoreConfig)
		act.PDC = FormatterPDC(v.PDC)
		ref.Acts = append(ref.Acts, act)
	}
	err := s.db.Create(&ref)
	if err != nil {
		return err.Error
	}
	return nil
}

func FormatterCoreConfig(coreConfig [][]string) []byte {
	var str string
	sliceLen := len(coreConfig)
	for i, vv := range coreConfig {
		for _, elem := range vv {
			str += elem + ","
		}
		if sliceLen == i+1 {
			str = str[:len(str)-1]
		}
	}
	formattedConfig := []byte(str)

	return formattedConfig
}

func FormatterPDC(pdc []string) []byte {
	joined := strings.Join(pdc, "")
	formattedPDC := []byte(joined)
	return formattedPDC
}

// var configs []string //* change to [][]byte for latter ease assign
// 	for _, v := range refuel.Acts {
// 		var str string
// 		sliceLen := len(v.CoreConfig)
// 		for i, vv := range v.CoreConfig {
// 			for _, elem := range vv {
// 				str += elem + ","
// 			}
// 			if sliceLen == i+1 {
// 				str = str[:len(str)-1]
// 			}
// 			// fmt.Println(str)
// 		}
// 		configs = append(configs, str)
// 	}
// 	fmt.Println(configs)
