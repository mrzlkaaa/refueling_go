package adding

type AddingService interface {
	AddWeeklyData(*FuelCycle)
}

type StorageService interface {
	AddWeeklyData(*FuelCycle)
	FCExistingCheck(string) error
	WeekNameExistingCheck(string, int32) error
	CreateDBInstance(string)
	AddWeekTemplate(string, int32)
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
		//* call func to create instance
		s.storage.CreateDBInstance(formsData.Name)
	}

	//* call fun to check if week name exists
	err = s.storage.WeekNameExistingCheck(formsData.Name, formsData.WeekName)
	if err == nil {
		panic("exists!")
	}
	//* call func to update existing instance by adding new week template and data
	s.storage.AddWeekTemplate(formsData.Name, formsData.WeekName)
	s.storage.AddWeeklyData(formsData)
	
	
	
	
}

//
