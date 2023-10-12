package storage

type DataStorage interface {
	Init() error
	Add(key string, value string) error
	Find(key string) (string, bool, error)
	Close() error
}
