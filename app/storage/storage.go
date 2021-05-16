package storage

type ServerStorage interface {
	Store(string) error
	Delete(string) error
	GetAll() []string
}

type ClientStorage interface {
	GetRandom() (string, error)
}
