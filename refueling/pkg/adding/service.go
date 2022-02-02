package adding

import (
	"errors"
	"refueling/refueling/pkg/listing"
)

var ErrDuplicate = errors.New("Refuel already exists")

type AddingService interface {
	Adding(refuel Refuel) error
}

type Storage interface {
	Adding(refuel Refuel) error
	GetRefuelNames() []listing.Refuel
}

type addingSerivce struct {
	storage Storage
}

func NewService(storage Storage) AddingService {
	return &addingSerivce{storage: storage}
}

func (s addingSerivce) Adding(refuel Refuel) error {
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
