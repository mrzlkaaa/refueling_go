package listingDiary

type ListingService interface {
	GetNewWeekNum(string) int32
}

type StorageService interface {
	GetNewWeekNum(string) int32
}

type listingService struct {
	storage StorageService
}

func NewListingService(storage StorageService) ListingService {
	return &listingService{storage: storage}
}

func (s *listingService) GetNewWeekNum(name string) int32 {
	weekNum := s.storage.GetNewWeekNum(name)
	return weekNum
}
