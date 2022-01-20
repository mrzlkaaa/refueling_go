package listingDiary

type ListingService interface {
	GetWeeksNum(string) []int
	WeekDetails(string, int) []DetailWeek
}

type StorageService interface {
	GetWeeksNum(string) []int
	WeekDetails(string, int) []DetailWeek
	// getDataForWeek(string)  NoSQL.FuelCycle
}

type listingService struct {
	storage StorageService
}

func NewListingService(storage StorageService) ListingService {
	return &listingService{storage: storage}
}

func (s *listingService) GetWeeksNum(fcName string) []int {

	weekNum := s.storage.GetWeeksNum(fcName)
	return weekNum
}

func (s *listingService) WeekDetails(fcName string, weekName int) []DetailWeek {
	//* check if it's existing instance
	data := s.storage.WeekDetails(fcName, weekName)
	return data
}
