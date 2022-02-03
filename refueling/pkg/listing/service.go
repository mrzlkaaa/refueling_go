package listing

type ListingService interface {
	GetRefuelNames() map[string][]int
	Refuels() []Refuel
	RefuelDetails(int) []Act
	RefuelPDC(int) []string
}

type StorageService interface {
	GetRefuelNames() []Refuel
	RefuelDetails(int) []Act
	RefuelPDC(int) []string
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

func (s *listingService) RefuelDetails(id int) []Act {
	data := s.storage.RefuelDetails(id)
	return data
}

func (s *listingService) RefuelPDC(id int) []string {
	data := s.storage.RefuelPDC(id)
	return data
}
