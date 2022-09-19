package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"refueling/refueling/pkg/adding"
	"refueling/refueling/pkg/listing"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//todo
//* create refuel folder
//* create and populate PDC and config files
//* query PDC and config files
//* delete PDC and config files
//* delete refuel folder

var NotFoundErr error = errors.New("Requested data not found")
var wg sync.WaitGroup

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
	db.AutoMigrate(&Refuel{}, &Act{}) //* keep db up to date
	return &Storage{db: db}
}

func (s *Storage) Adding(refuel adding.Refuel) error { //! only to communicate with db
	var ref Refuel
	var act Act

	//* fill ref
	ref.RefuelName = refuel.RefuelName
	ref.Date = refuel.Date

	folderName := strconv.Itoa(refuel.RefuelName)
	//* fill acts
	for _, v := range refuel.Acts {
		act.Name = v.Name
		act.Description = v.Description
		ref.Acts = append(ref.Acts, act)
	}

	//* pdc and core configs settle in folders in parallel and synchronized

	//* create root folder basePath/RefuelName
	// popErr := make(chan error)
	//todo move fodler creation step from this func to let use for adding refueling act
	err := os.Mkdir("data/"+folderName, 0755)
	dir, _ := os.Getwd()
	fmt.Println(err, dir)
	if err != nil && os.IsExist(err) {
		return err //* better to re0turn err
	}
	// err = PopulatePathTo(&refuel)
	err = PopulatePathTo(refuel.RefuelName, refuel.Acts)
	if err != nil {
		return err
	}
	//* new instance successfully created and added -->
	res := s.db.Create(&ref)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Storage) AddingAct(act adding.Act) (uint, error) {
	var ac Act
	fmt.Println(act.RefuelNameRef)
	ac.Name = act.Name
	ac.Description = act.Description
	ac.RefuelNameRef = act.RefuelNameRef
	//todo check if root folder exists --> populate folder
	_, _, err := CheckPathToStoredFiles(act.RefuelNameRef)
	if err != nil {
		return 0, err
	}
	var acts []adding.Act
	acts = append(acts, act)
	err = PopulatePathTo(act.RefuelNameRef, acts)
	if err != nil {
		return 0, err
	}
	res := s.db.Create(&ac)
	return ac.ID, res.Error
}

//! Not finished
//* Takes RefuelName as a key path element
func PopulatePathTo(id int, acts []adding.Act) error { //* *Refuel is instance of adding package
	rand.Seed(time.Now().UnixNano())

	//* create ctx
	ctx := context.TODO()
	ctx, cancelCtx := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancelCtx()
	// fmt.Println(<-crErr)
	// //* Populate folder
	popErrPDC := make(chan error)
	popErrCon := make(chan error)
	wg.Add(len(acts) * 2)
	for _, v := range acts {
		//todo pass value (err) to channels
		go PopulateFolderByPDC(ctx, id, v, popErrPDC)
		go PopulateFolderByConfig(ctx, id, v, popErrCon)
	}
	for {
		select {
		case err := <-popErrPDC:
			return err
			// return errors.New("error has happended while populating root folder by PDCs")
		case err := <-popErrCon:
			return err
			// return errors.New("error has happended while populating root folder by coreConfigs")
			// fmt.Println("Successfully populated")
		case <-ctx.Done():
			return nil
		}

	}
}

//* this folder we populate by pdcs
func PopulateFolderByPDC(ctx context.Context, id int, act adding.Act, popErrPDC chan<- error) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Int63n(1e6)))
	folderName := strconv.Itoa(id)
	filePath := filepath.Join("data", folderName, GetFileName(id, act.Name, "PDC"))
	// fileName :=
	fmt.Println(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		popErrPDC <- err
		return
	}
	defer f.Close()
	for _, v := range act.PDC {
		_, err = f.WriteString(v)
		if err != nil {
			popErrPDC <- err
			return
		}
	}
	fmt.Println("File", filePath, "created")

}

//* this folder we populate by core configs
func PopulateFolderByConfig(ctx context.Context, id int, act adding.Act, popErrCon chan<- error) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Int63n(1e6)))
	folderName := strconv.Itoa(id)
	filePath := filepath.Join("data", folderName, GetFileName(id, act.Name, "txt"))
	// filePath := filepath.Join("data", rn, act.Name+".txt")
	fmt.Println(filePath)
	f, err := os.Create(filePath)

	if err != nil {
		popErrCon <- err
	}

	defer f.Close()
	//* looping in 2d slice
	for _, v1 := range act.CoreConfig {
		for _, v2 := range v1 {
			_, err = f.WriteString(v2 + "    ")
			if err != nil {
				popErrCon <- err
			}
		}
		_, err := f.WriteString("\n")
		if err != nil {
			popErrCon <- err
		}
		//todo add "\n" after this if it is not in string
	}
	fmt.Println("File", filePath, "created")
}

//* returns map of refuel names
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

//* query refuel details
func (s *Storage) RefuelDetails(refuelName int) []listing.Act {
	var bindedActs []Act
	var refuel Refuel

	// res := s.db.First(&refuel, id)
	// if res.RowsAffected == 0 {
	// 	return []listing.Act{}
	// }
	err := s.IsRefuelExists(&refuel, refuelName)
	if err != nil {
		return []listing.Act{}
	}

	res := s.db.Where(Act{RefuelNameRef: refuelName}).Find(&bindedActs) //!todo use refuelName

	if res.RowsAffected == 0 {
		return []listing.Act{}
	}

	lenght := len(bindedActs)
	var acts []listing.Act = make([]listing.Act, lenght)
	for i, v := range bindedActs {
		fmt.Println(v.RefuelNameRef)
		acts[i].ID = v.ID
		acts[i].Name = v.Name
		//* formatting of core config
		acts[i].CoreConfig = ConfigStorageQuery(v.RefuelNameRef, v.Name)
		acts[i].Description = v.Description
		acts[i].RefuelNameRef = v.RefuelNameRef
	}
	//* query both PDC and config

	// BackFormatterCoreConfig(&binded[0].CoreConfig)
	fmt.Println(acts)
	return acts
}

//todo (s *Storage)
func ConfigStorageQuery(refuelName int, name string) [][]string { //* query configs of refuel

	files, path, err := CheckPathToStoredFiles(refuelName)
	if err != nil {
		return [][]string{}
	}

	//* loop over files and create array of configs
	//todo open and read as goroutine
	fileData := make(chan [][]string)
	defer close(fileData)
	var configArr [][]string
	fileName := GetFileName(refuelName, name, "txt")
	fmt.Println(fileName)
	for _, file := range files {
		if contains := strings.Contains(file.Name(), fileName); contains {
			go OpenReadStoredConfig(filepath.Join(path, file.Name()), fileData)
			// storedConfigNames = append(storedConfigNames, file.Name())
			fmt.Println(file.Name())
			// path <- filepath.Join(queryPath, file.Name())
			configArr = <-fileData

		}
	}
	fmt.Println(configArr)
	return configArr
}

//* find and read pdc in refuel folder
func (s *Storage) PDCStorageQuery(refuelName int, name string) []string { //* query pdc of refuel
	//todo fetch db by PK and retrieve refuelName for further use!

	var refuel Refuel

	err := s.IsRefuelExists(&refuel, refuelName)
	if err != nil {
		return []string{}
	}

	files, path, err := CheckPathToStoredFiles(refuel.RefuelName)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	//* loop over files and create array of configs
	//todo open and read as goroutine
	pdcChan := make(chan []string)
	defer close(pdcChan)
	var pdc []string
	fileName := GetFileName(refuel.RefuelName, name, "PDC")
	fmt.Println(fileName, path)
	for _, file := range files {
		if contains := strings.Contains(file.Name(), fileName); contains {
			go OpenReadStoredPDC(filepath.Join(path, file.Name()), pdcChan)
			// storedConfigNames = append(storedConfigNames, file.Name())
			fmt.Println(file.Name())
			// path <- filepath.Join(queryPath, file.Name())
			pdc = <-pdcChan
			// PDCArr = append(PDCArr, <-lineData)

		}
	}
	fmt.Println(len(pdc))
	return pdc
}

//* Delete config and PDc of requested refuel act
func ConfigPDCStoredDelete(refuelName int, name string) error {
	files, path, err := CheckPathToStoredFiles(refuelName)
	if err != nil {
		return err
	}

	fileName := GetFileName(refuelName, name, "")
	for _, file := range files {
		if contains := strings.Contains(file.Name(), fileName); contains {
			fmt.Printf("File %v removed\n", file.Name())
			err = os.Remove(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//* Delete refuel folder with all children contains
func RefuelStoredDelete(refuelName int) error {
	_, path, err := CheckPathToStoredFiles(refuelName)
	if err != nil {
		return err
	}
	if err = os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

//* delete refueling and all binded acts
func (s *Storage) Delete(refuelName int) error {
	var act Act
	var refuel Refuel
	// act.RefuelID = id
	// res := s.db.Where("refuel_id = ?", id).Delete(&act)
	// res := s.db.First(&refuel, uint(id)) //* fetch refuel by PK
	// if res.RowsAffected == 0 {
	// 	return NotFoundErr
	// }
	err := s.IsRefuelExists(&refuel, refuelName)
	if err != nil {
		return NotFoundErr
	}

	//* delete all pdc and core configs in root folder
	err = RefuelStoredDelete(refuelName)
	if err != nil {
		return err
	}

	res := s.db.Where("refuel_name_ref = ?", refuel.RefuelName).Delete(&act) //*delete all Acts with matched FK
	if res.RowsAffected == 0 {
		return NotFoundErr
	}

	res = s.db.Delete(&Refuel{}, refuel.ID) //* delete by PK
	if res.RowsAffected == 0 {
		return NotFoundErr
	}

	return nil
}

//* delete act of refueling
func (s *Storage) DeleteAct(refuelName int, actName string) error {
	var act Act
	var refuel Refuel

	// res := s.db.First(&refuel, uint(id)) //* fetch refuel by PK
	err := s.IsRefuelExists(&refuel, refuelName)
	if err != nil {
		return NotFoundErr
	}

	err = ConfigPDCStoredDelete(refuelName, actName)
	if err != nil {
		return err
	}

	res := s.db.Where(&Act{Name: actName, RefuelNameRef: refuel.RefuelName}).Delete(&act) //*delete all Acts with matched FK
	if res.RowsAffected == 0 {
		return NotFoundErr
	}
	return res.Error
}

func (s *Storage) IsRefuelExists(refuel *Refuel, refuelName int) error {

	// res := s.db.First(refuel, uint(id))
	res := s.db.Where(&Refuel{RefuelName: refuelName}).Find(refuel)
	fmt.Println(*refuel)
	if res.RowsAffected == 0 {
		return NotFoundErr
	}

	return nil
}
