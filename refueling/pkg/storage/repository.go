package storage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	//! core config and pdc are not stored anymore
	s.db.Select("id", "name", "core_config",
		"description", "refuel_id").Where(Act{RefuelID: id}).Find(&binded)

	lenght := len(binded)
	var acts []listing.Act = make([]listing.Act, lenght)
	for i, v := range binded {
		acts[i].ID = v.ID
		acts[i].Name = v.Name
		//* formatting of core config
		// acts[i].CoreConfig = *BackFormatterCoreConfig(&v.CoreConfig)
		acts[i].Description = v.Description
		acts[i].RefuelID = v.RefuelID
	}
	// BackFormatterCoreConfig(&binded[0].CoreConfig)
	fmt.Println(acts)
	return acts
}

//todo (s *Storage)
func ConfigStorageQuery(id int) []string { //* query pdcs and congigs of refuel
	var storedConfigNames []string
	root := "data" //* Storage root
	queryPath := filepath.Join(root, strconv.Itoa(id))
	files, err := ioutil.ReadDir(queryPath)
	if err != nil || len(files) == 0 {
		fmt.Println(err)
		return []string{}
	}
	//* loop over files and create array of
	//todo read as goroutine
	for _, file := range files {
		if contains := strings.Contains(file.Name(), ".txt"); contains {
			storedConfigNames = append(storedConfigNames, file.Name())
			filePath := filepath.Join(queryPath, file.Name())
			file, err := os.Open(filePath)

			if err != nil {
				//* do error handler
			}
			defer file.Close()
			fileScanner := bufio.NewScanner(file)
			for fileScanner.Scan() {
				text := fileScanner.Text()
				fmt.Println(text)
			}
			fmt.Printf("\n")
			//todo add text to 2d array --> return
			// fmt.Println(file.Name(), file.IsDir())
		}
	}

	// var configsArr [][]string

	return storedConfigNames
}

func (s *Storage) RefuelPDC(id int) []string { //* DO refuel
	// 	var binded Act
	// 	uid := uint(id)
	// 	s.db.Select("pdc").Where(Act{ID: uid}).Find(&binded)
	// 	arr := *BackFormatterPDC(&binded.PDC)
	// 	// fmt.Println(arr)
	// 	return arr
	return []string{}
}

func (s *Storage) SavePDC(id int) (string, *[]byte) { //* PDC query and later download
	// 	var binded Act
	// 	uid := uint(id)
	// 	s.db.Select("name", "pdc").Where(Act{ID: uid}).Find(&binded)
	// 	return binded.Name, &binded.PDC
	return "", &[]byte{}
}

func (s *Storage) Adding(refuel adding.Refuel) error { //! only to communicate with db
	var ref Refuel
	var act Act

	//* fill ref
	ref.RefuelName = refuel.RefuelName
	ref.Date = refuel.Date

	//* fill acts
	for _, v := range refuel.Acts {
		act.Name = v.Name
		act.Description = v.Description
		//! Not stored in db anymore
		// act.CoreConfig = *FormatterCoreConfig(&v.CoreConfig)
		// act.PDC = *FormatterPDC(&v.PDC)
		ref.Acts = append(ref.Acts, act)
	}

	//* pdc and core configs settle in folders in parallel and synchronized
	folderName := strconv.Itoa(refuel.RefuelName)
	//* create root folder basePath/RefuelName
	// popErr := make(chan error)
	//todo move fodler creation step from this func to let use for adding refueling act
	err := os.Mkdir("data/"+folderName, 0755)
	if err != nil && os.IsExist(err) {
		return err //* better to re0turn err
	}
	err = PopulatePathTo(&refuel)
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

func (s *Storage) AddingAct(act adding.Act) (error, uint) {
	var ac Act
	ac.Name = act.Name
	//! Not stored in db anymore
	// ac.CoreConfig = *FormatterCoreConfig(&act.CoreConfig)
	// ac.PDC = *FormatterPDC(&act.PDC)
	ac.Description = act.Description
	ac.RefuelID = act.RefuelID
	//todo check if root folder exists --> populate folder

	res := s.db.Create(&ac)
	return res.Error, ac.ID
}

//! Not finished
//* Takes RefuelName as a key path element
func PopulatePathTo(r *adding.Refuel) error { //* *Refuel is instance of adding package
	rand.Seed(time.Now().UnixNano())
	folderName := strconv.Itoa(r.RefuelName)
	//* create ctx
	ctx := context.TODO()
	ctx, cancelCtx := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancelCtx()
	// fmt.Println(<-crErr)
	// //* Populate folder
	popErrPDC := make(chan error)
	popErrCon := make(chan error)
	wg.Add(len(r.Acts) * 2)
	for _, k := range r.Acts {
		//todo pass value (err) to channels
		go PopulateFolderByPDC(ctx, folderName, k, popErrPDC)
		go PopulateFolderByConfig(ctx, folderName, k, popErrCon)
	}
	fmt.Println("check sync")
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
func PopulateFolderByPDC(ctx context.Context, rn string, act adding.Act, popErrPDC chan<- error) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Int63n(1e6)))
	// filePath := filepath.Join(path, act.Name+".PDC")
	filePath := filepath.Join("data", rn, act.Name+".PDC")
	fmt.Println(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		popErrPDC <- err
		return
	}
	defer f.Close()
	for _, v := range act.PDC {
		_, err = f.WriteString(v) //todo add "\n" after this if it is not in string
		if err != nil {
			popErrPDC <- err
			return
		}
	}
	fmt.Println("File", filePath, "created")

}

//* this folder we populate by core configs
func PopulateFolderByConfig(ctx context.Context, rn string, act adding.Act, popErrCon chan<- error) {
	// _, cancel := context.WithCancel(ctx)
	// defer cancel()
	defer wg.Done()
	time.Sleep(time.Duration(rand.Int63n(1e6)))
	// filePath := filepath.Join(path, act.Name+".txt")
	filePath := filepath.Join("data", rn, act.Name+".txt")
	fmt.Println(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		popErrCon <- err
	}
	defer f.Close()
	//* looping in 2d slice
	for _, v1 := range act.CoreConfig {
		for _, v2 := range v1 {
			_, err = f.WriteString(v2) //todo add "\n" after this if it is not in string
			if err != nil {
				popErrCon <- err
			}
		}
	}
	fmt.Println("File", filePath, "created")
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
