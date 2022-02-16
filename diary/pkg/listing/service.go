package listing

type ListingService interface {
	GetWeeksNum(int) []int
	WeekDetails(int, int) []DetailWeek
	OverallData() []FuelCycle
}

type StorageService interface {
	GetWeeksNum(int) []int
	WeekDetails(int, int) []DetailWeek
	OverallData() []FuelCycle
	// getDataForWeek(string)  NoSQL.FuelCycle
}

type listingService struct {
	storage StorageService
}

func NewListingService(storage StorageService) ListingService {
	return &listingService{storage: storage}
}

func (s listingService) OverallData() []FuelCycle {
	data := s.storage.OverallData()
	return data
}

func (s *listingService) GetWeeksNum(fcName int) []int {
	weekNum := s.storage.GetWeeksNum(fcName)
	return weekNum
}

func (s *listingService) WeekDetails(fcName int, weekName int) []DetailWeek {
	//* check if it's existing instance
	data := s.storage.WeekDetails(fcName, weekName)
	return data
}
