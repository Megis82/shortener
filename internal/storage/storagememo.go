package storage

type MemoryStorage struct {
	data map[string]string
}

func (m *MemoryStorage) Init() error {
	m.data = make(map[string]string)
	return nil
}

func (m *MemoryStorage) Add(key string, value string) error {
	m.data[key] = value
	return nil
}

func (m *MemoryStorage) Find(key string) (string, bool, error) {
	value, ok := m.data[key]
	return value, ok, nil
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{data: make(map[string]string)}, nil
}
