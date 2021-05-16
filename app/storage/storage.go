package storage

type ServerStorage interface {
	Store(string) error
	Delete(string) error
	GetAll() ([]string, error)
}

type ClientStorage interface {
	GetRandom() (string, error)
}
