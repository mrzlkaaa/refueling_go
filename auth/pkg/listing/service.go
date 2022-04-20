package listing

type service struct {
	storage Storage
}

type Service interface {
	GetAllUsers() ([]User, error)
}

type Storage interface {
	GetAllUsers() ([]User, error)
	//*funcs of storage package
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetAllUsers() ([]User, error) {
	users, err := s.storage.GetAllUsers()
	return users, err
}
