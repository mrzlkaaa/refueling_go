package listing

type ListingService interface {
	GetWeeksNum(int) map[string]int
	WeekDetails(int, int) WeeklyData
	OverallData() []FuelCycle
}

type StorageService interface {
	GetWeeksNum(int) map[string]int
	WeekDetails(int, int) WeeklyData
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

func (s *listingService) GetWeeksNum(fcName int) map[string]int {
	weekNum := s.storage.GetWeeksNum(fcName)
	return weekNum
}

func (s *listingService) WeekDetails(fcName int, weekName int) WeeklyData {
	//* check if it's existing instance
	data := s.storage.WeekDetails(fcName, weekName)
	return data
}
