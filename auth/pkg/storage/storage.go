package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"refueling/auth/pkg/adding"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbUserExists = "User already exists"
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
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})

	return &Storage{db: db}
}

func (s *Storage) AddUser(UserForm adding.User, password []byte) error {
	var user *User

	user.Username = UserForm.Username
	user.Password = password
	user.Admin = UserForm.Admin

	res := s.db.Create(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Storage) UpdRole() {
	//* do ipdate via member id?
}

func (s *Storage) IfUserExists(username string) error {
	var user User

	res := s.db.Where(&User{Username: username}).Find(&user)
	if res.RowsAffected != 0 {
		return errors.New(dbUserExists)
	}

	return nil
}