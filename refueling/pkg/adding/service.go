package adding

import (
	"errors"
	"refueling/refueling/pkg/listing"
)

var ErrDuplicate = errors.New("Refuel already exists")

type AddingService interface {
	Adding(refuel Refuel) error
	AddingAct(act Act) (error, uint)
	Deleting(id int) error
	DeletingAct(id int) error
}

type Storage interface {
	Adding(refuel Refuel) error
	AddingAct(act Act) (error, uint)
	GetRefuelNames() []listing.Refuel
	Deleting(id int) error
	DeletingAct(id int) error
}

type addingSerivce struct {
	storage Storage
}

func NewService(storage Storage) AddingService {
	return &addingSerivce{storage: storage}
}

func (s *addingSerivce) Adding(refuel Refuel) error {
	// check duplicates
	listNames := s.storage.GetRefuelNames()
	for _, v := range listNames {
		if v.RefuelName == refuel.RefuelName {
			return ErrDuplicate
		}
	}
	//
	err := s.storage.Adding(refuel)
	return err
}

func (s *addingSerivce) AddingAct(act Act) (error, uint) {
	err, id := s.storage.AddingAct(act)
	if err != nil {
		panic(err)
	}
	return err, id
}

func (s *addingSerivce) Deleting(id int) error {
	err := s.storage.Deleting(id)
	return err
}

func (s *addingSerivce) DeletingAct(id int) error {
	err := s.storage.DeletingAct(id)
	return err
}
