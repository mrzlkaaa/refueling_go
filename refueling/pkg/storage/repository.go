package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var NotFoundErr error = errors.New("Requested data not found")

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
	s.db.Select("id", "refuel_name", "date").Order("refuel_name desc").Find(&binded)
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

func (s *Storage) SavePDC(id int) (string, *[]byte) {
	var binded Act
	uid := uint(id)
	s.db.Select("name", "pdc").Where(Act{ID: uid}).Find(&binded)
	return binded.Name, &binded.PDC
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
		act.CoreConfig = *FormatterCoreConfig(&v.CoreConfig)
		act.PDC = *FormatterPDC(&v.PDC)
		ref.Acts = append(ref.Acts, act)
	}
	res := s.db.Create(&ref)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Storage) AddingAct(act adding.Act) (error, uint) {
	var ac Act
	ac.Name = act.Name
	ac.CoreConfig = *FormatterCoreConfig(&act.CoreConfig)
	ac.PDC = *FormatterPDC(&act.PDC)
	ac.Description = act.Description
	ac.RefuelID = act.RefuelID
	res := s.db.Create(&ac)
	return res.Error, ac.ID
}

func (s *Storage) Deleting(id int) error {
	var act Act
	// act.RefuelID = id
	res := s.db.Where("refuel_id = ?", id).Delete(&act)
	fmt.Println(res.Error)
	if res.RowsAffected == 0 {
		return NotFoundErr
	}
	res = s.db.Delete(&Refuel{}, uint(id))
	if res.RowsAffected == 0 {
		return NotFoundErr
	}
	return res.Error
}

func (s *Storage) DeletingAct(id int) error {
	res := s.db.Delete(&Act{}, uint(id))
	if res.RowsAffected == 0 {
		return NotFoundErr
	}
	return res.Error
}
