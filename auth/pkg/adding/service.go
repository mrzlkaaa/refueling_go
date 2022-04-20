package adding

import (
	"refueling/auth/pkg/listing"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	storage Storage
}

type Service interface {
	AddUser(User) error
	UpdateUsers(*[]listing.User) error
	DeleteUser(uint) error
	HashPassword(string) ([]byte, error)
}

type Storage interface {
	AddUser(User, []byte) error
	UpdateUser(*listing.User) error
	DeleteUser(uint) error
	IfUserExists(string) error
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) AddUser(userForm User) error {
	err := s.storage.IfUserExists(userForm.Username)
	if err != nil {
		return err
	}
	hashedPassword, err := s.HashPassword(userForm.Password)
	if err != nil {
		return err
	}
	err = s.storage.AddUser(userForm, hashedPassword)
	return err
}

func (s *service) UpdateUsers(users *[]listing.User) error {
	for _, user := range *users {
		err := s.storage.UpdateUser(&user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteUser(id uint) error {
	err := s.storage.DeleteUser(id)
	return err
}

func (s *service) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return hashedPassword, err
}
