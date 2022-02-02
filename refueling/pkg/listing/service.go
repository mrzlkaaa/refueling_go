package listing

type ListingService interface {
	GetRefuelNames() map[string][]int
}

type StorageService interface {
	GetRefuelNames() []Refuel
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
