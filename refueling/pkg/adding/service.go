package adding

import (
	"errors"
	"fmt"
	"refueling/refueling/pkg/listing"
)

var ErrDuplicate = errors.New("Refuel already exists")

type AddingService interface {
	Adding(Refuel) error
	AddingAct(Act) (uint, error)
	Delete(int) error
	DeleteAct(int, string) error
}

type Storage interface {
	Adding(Refuel) error
	AddingAct(Act) (uint, error)
	GetRefuelNames() []listing.Refuel
	Delete(int) error
	DeleteAct(int, string) error
}

type addingSerivce struct {
	storage Storage
}

func NewService(storage Storage) AddingService {
	return &addingSerivce{storage: storage}
}

func (s *addingSerivce) Adding(refuel Refuel) error {
	//* check duplicates in db and storage
	listNames := s.storage.GetRefuelNames()
	for _, v := range listNames {
		if v.RefuelName == refuel.RefuelName {
			return ErrDuplicate
		}
	}
	//
	// db_err := make(chan bool)
	fmt.Println(refuel.Acts[0].CoreConfig)
	err := s.storage.Adding(refuel)

	// for i:=0; i<2; i++ {
	// 	select {
	// 		case db_err
	// 	}
	// }

	return err
}

func (s *addingSerivce) AddingAct(act Act) (uint, error) {
	id, err := s.storage.AddingAct(act)
	if err != nil {
		panic(err)
	}
	return id, err
}

func (s *addingSerivce) Delete(refuelName int) error {
	err := s.storage.Delete(refuelName)
	return err
}

func (s *addingSerivce) DeleteAct(refuelName int, name string) error {
	err := s.storage.DeleteAct(refuelName, name)
	return err
}
