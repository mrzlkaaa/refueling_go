package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"refueling/auth/pkg/adding"
	"refueling/auth/pkg/listing"
	"refueling/auth/pkg/loggining"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbUserExists = "User already exists"
	UserNotFound = "User not found"
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
	var user User

	user.Name = UserForm.Name
	user.Surname = UserForm.Surname
	user.Email = UserForm.Email
	user.Username = UserForm.Username
	user.Password = password
	user.Admin = UserForm.Admin

	res := s.db.Create(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Storage) UpdateUser(user *listing.User) error {
	var userStored User

	res := s.db.Find(&userStored, "id = ?", user.ID)
	if res.RowsAffected == 0 {
		return errors.New(UserNotFound)
	}

	//* do replacement
	userStored.Moderator = user.Moderator
	userStored.Admin = user.Admin

	res = s.db.Save(&userStored)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Storage) DeleteUser(id uint) error {
	res := s.db.Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Storage) FindUser(username string) (loggining.UserData, error) {
	var userStored User
	var userData loggining.UserData
	res := s.db.Find(&userStored, &User{Username: username})
	if res.RowsAffected == 0 {
		return userData, errors.New(UserNotFound)
	}

	userData.ID = userStored.ID
	userData.Name = userStored.Name
	userData.Surname = userStored.Surname
	userData.Email = userStored.Email
	userData.Username = userStored.Username
	userData.PswdHash = userStored.Password
	userData.Moderator = userStored.Moderator
	userData.Admin = userStored.Admin

	return userData, nil
}

func (s *Storage) FindUserID(ID uint) (loggining.UserData, error) {
	var userStored User
	var userData loggining.UserData
	res := s.db.Where("id =?", ID).Find(&userStored)
	if res.RowsAffected == 0 {
		return userData, errors.New(UserNotFound)
	}

	fmt.Println(userStored)
	userData.ID = userStored.ID
	userData.Name = userStored.Name
	userData.Surname = userStored.Surname
	userData.Email = userStored.Email
	userData.Moderator = userStored.Moderator
	userData.Admin = userStored.Admin

	return userData, nil
}

func (s *Storage) GetAllUsers() ([]listing.User, error) {
	var dbUsers []User

	err := s.db.Select("id", "name", "surname", "username", "email", "moderator", "admin").Find(&dbUsers)
	if err.Error != nil {
		return []listing.User{}, nil
	}

	users := make([]listing.User, len(dbUsers))
	for k, v := range dbUsers {
		users[k].ID = v.ID
		users[k].Name = v.Name
		users[k].Surname = v.Surname
		users[k].Username = v.Username
		users[k].Email = v.Email
		users[k].Moderator = v.Moderator
		users[k].Admin = v.Admin
	}
	fmt.Println(users)

	return users, nil
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
