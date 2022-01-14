package adding

type AddingService interface {
	AddWeeklyData(*FormsData)
}

type StorageService interface {
	AddWeeklyData(*FormsData)
}

type service struct {
	storage StorageService
}

func NewService(storage StorageService) AddingService {
	return &service{storage: storage}
}

func (s *service) AddWeeklyData(formsData *FormsData) {
	//* do some dublicates checks
	s.storage.AddWeeklyData(formsData)
}

//
