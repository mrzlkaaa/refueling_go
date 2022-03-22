package listing

type service struct {
	storage Storage
}

type Service interface {
}

type Storage interface {
	//*funcs of storage package
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}
