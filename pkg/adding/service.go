package adding

type AddingService interface {
	AddWeeklyData(*FuelCycle)
}

type StorageService interface {
	FCExistingCheck(string) error
	CreateDBInstance(string)
	AddWeeklyData(*FuelCycle)
	WeekNameExistingCheck(string, int) error
	AddWeekTemplate(string, int)
}

type service struct {
	storage StorageService
}

func NewService(storage StorageService) AddingService {
	return &service{storage: storage}
}

func (s *service) AddWeeklyData(formsData *FuelCycle) {

	//! review required if update feature is enabled
	// //* call fun to check if week name exists
	// err := s.storage.WeekNameExistingCheck(formsData.Name, formsData.WeekName)
	// if err == nil {
	// 	panic("exists!")
	// }
	// err = s.storage.FCExistingCheck(formsData.Name)
	// if err != nil {
	// 	//* call func to create instance
	// 	s.storage.CreateDBInstance(formsData.Name)
	// }
	//* call func to update existing instance by adding new week template and data
	// s.storage.AddWeekTemplate(formsData.Name, formsData.WeekName)
	s.storage.AddWeeklyData(formsData)

}

//
