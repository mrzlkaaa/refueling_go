package listing

type ListingService interface {
	GetRefuelNames() map[string][]int
	Refuels() []Refuel
	RefuelDetails(int) []Act
	PDCStorageQuery(int, string) []string
}

type StorageService interface {
	GetRefuelNames() []Refuel
	RefuelDetails(int) []Act
	PDCStorageQuery(int, string) []string
}

type listingService struct {
	storage StorageService
}

func NewListingService(storage StorageService) ListingService {
	return &listingService{storage: storage}
}

func (s *listingService) GetRefuelNames() map[string][]int {
	mapQuery := make(map[string][]int)
	data := s.storage.GetRefuelNames()
	var names []int
	for _, v := range data {
		names = append(names, v.RefuelName)
	}
	mapQuery["names"] = names
	return mapQuery
}

func (s *listingService) Refuels() []Refuel {
	data := s.storage.GetRefuelNames()
	return data
}

func (s *listingService) RefuelDetails(refuelName int) []Act {
	data := s.storage.RefuelDetails(refuelName)
	return data
}

func (s *listingService) PDCStorageQuery(refuelName int, name string) []string {
	data := s.storage.PDCStorageQuery(refuelName, name)
	return data
}
