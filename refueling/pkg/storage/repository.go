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
	var binded []Refuel
	s.db.Select("id", "refuel_name", "date").Find(&binded)
	lenght := len(binded)
	var refuels []listing.Refuel = make([]listing.Refuel, lenght)
	for i, v := range binded {
		refuels[i].ID = v.ID
		refuels[i].RefuelName = v.RefuelName
		refuels[i].Date = v.Date
	}
	return refuels
}

func (s *Storage) RefuelDetails(id int) []listing.Act {
	var binded []Act
	s.db.Select("id", "name", "core_config",
		"description", "refuel_id").Where(Act{RefuelID: id}).Find(&binded)

	lenght := len(binded)
	var acts []listing.Act = make([]listing.Act, lenght)
	for i, v := range binded {
		acts[i].ID = v.ID
		acts[i].Name = v.Name
		acts[i].CoreConfig = *BackFormatterCoreConfig(&v.CoreConfig)
		acts[i].Description = v.Description
		acts[i].RefuelID = v.RefuelID
	}
	// BackFormatterCoreConfig(&binded[0].CoreConfig)
	fmt.Println(acts)
	return acts
}

func (s *Storage) RefuelPDC(id int) []string {
	var binded Act
	uid := uint(id)
	s.db.Select("pdc").Where(Act{ID: uid}).Find(&binded)
	arr := *BackFormatterPDC(&binded.PDC)
	// fmt.Println(arr)
	return arr
}

func (s *Storage) Adding(refuel adding.Refuel) error {
	var ref Refuel
	var act Act

	//* fill ref
	ref.RefuelName = refuel.RefuelName
	ref.Date = refuel.Date

	//* fill acts
	for _, v := range refuel.Acts {
		act.Name = v.Name
		act.Description = v.Description
		act.CoreConfig = FormatterCoreConfig(v.CoreConfig)
		act.PDC = *FormatterPDC(&v.PDC)
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

func BackFormatterCoreConfig(coreConfig *[]byte) *[][]string {
	var str string
	var firstArr []string
	var arr2D [][]string

	for _, v := range *coreConfig { //*conv to str
		str += string(v)
	}
	array := strings.Split(str, ",")
	for i, v := range array {
		i += 1
		firstArr = append(firstArr, v)
		if i%4 == 0 {
			arr2D = append(arr2D, firstArr)
			firstArr = []string{}
		}
	}
	return &arr2D
}

func FormatterPDC(pdc *[]string) *[]byte {
	joined := strings.Join(*pdc, "")
	formattedPDC := []byte(joined)
	return &formattedPDC
}

func BackFormatterPDC(pdc *[]byte) *[]string {
	var str string
	var arr []string
	// for _, v := range *pdc { //*conv to str
	// 	str += string(v)
	// }
	str = string(*pdc)
	arr = strings.Split(str, "\n")
	for i := 0; i < len(arr); i++ {
		arr[i] += "\n"
	}
	// arr = append(arr, str)
	return &arr
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
