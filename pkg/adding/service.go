package adding

type AddingService interface {
	AddWeeklyData(*FuelCycle)
}

type StorageService interface {
	AddWeeklyData(*FuelCycle)
	FCExistingCheck(string) error
	CreateDBInstance(string)
}

type service struct {
	storage StorageService
}

func NewService(storage StorageService) AddingService {
	return &service{storage: storage}
}

func (s *service) AddWeeklyData(formsData *FuelCycle) {
	//* do some dublicates checks or existing checks

	err := s.storage.FCExistingCheck(formsData.Name)
	if err != nil {
		//TODO call func to create instance
		s.storage.CreateDBInstance(formsData.Name)
	}

	s.storage.AddWeeklyData(formsData)
}

//
