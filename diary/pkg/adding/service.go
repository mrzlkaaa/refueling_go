package adding

type AddingService interface {
	AddWeeklyData(*FuelCycle)
}

type StorageService interface {
	FCExistingCheck(string) error
	WeekNameExistingCheck(string, int) error
	AddWeeklyData(*FuelCycle)
	UpdateWeeklyData(*FuelCycle)
	AppendWeeklyData(*FuelCycle)
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
	if err := s.storage.FCExistingCheck(formsData.Name); err != nil {
		s.storage.AddWeeklyData(formsData)
		return
	}

	//* take the same db instance as a base
	if err := s.storage.WeekNameExistingCheck(formsData.Name, formsData.WeekName); err == nil {
		s.storage.UpdateWeeklyData(formsData)
	} else {
		s.storage.AppendWeeklyData(formsData)
	}
	// return

	// err = s.storage.FCExistingCheck(formsData.Name)
	// if err != nil {
	// 	//* call func to create instance
	// 	s.storage.CreateDBInstance(formsData.Name)
	// }
	//* call func to update existing instance by adding new week template and data
	// s.storage.AddWeekTemplate(formsData.Name, formsData.WeekName)

}
