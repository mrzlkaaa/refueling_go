package listing

import (
	"refueling/pkg/storage/SQL"
)

type ListingService interface {
	GetRefuelNames() map[string][]string
	GetNewWeekNum(string)
}

type StorageService interface {
	GetRefuelNamesQuery() []SQL.ReactorRefuel
	GetNewWeekNum(string)
}

type listingService struct {
	storage StorageService
}

func NewListingService(storage StorageService) ListingService {
	return &listingService{storage: storage}
}

func (s *listingService) GetRefuelNames() map[string][]string {
	mapQuery := make(map[string][]string)
	data := s.storage.GetRefuelNamesQuery()
	var names []string
	for _, v := range data {
		names = append(names, v.Refueling_name)
	}
	mapQuery["names"] = names
	return mapQuery
}

func (s *listingService) GetNewWeekNum(name string) {
	s.storage.GetNewWeekNum(name)
}