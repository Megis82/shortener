package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemoryStorage struct {
	data        map[string]string
	fileStorage string
}

type MemoryStorageSave struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func (m *MemoryStorage) Ping() error {
	return nil
}

func NewMemoryStorage(fileStorage string) (*MemoryStorage, error) {

	mem := &MemoryStorage{data: make(map[string]string), fileStorage: fileStorage}

	file, err := os.OpenFile(fileStorage, os.O_RDONLY, 0644)

	if os.IsNotExist(err) {
		return mem, nil
	}

	if err != nil {
		return nil, err
	}

	buff := make([]byte, 0)
	byteCount, err := file.Read(buff)

	if err != nil {
		return nil, err
	}

	if byteCount == 0 {
		return mem, nil
	}

	data := make([]MemoryStorageSave, 0)

	err = json.Unmarshal(buff, &data)

	if err != nil {
		return nil, err
	}

	for _, str := range data {
		mem.Add(str.ShortURL, str.OriginalURL)
	}

	return mem, nil
}

func (m *MemoryStorage) Close() error {

	_ = os.Remove(m.fileStorage)

	data := make([]MemoryStorageSave, 0)

	idx := 1
	for key, val := range m.data {
		data = append(data, MemoryStorageSave{
			UUID:        fmt.Sprint(idx),
			ShortURL:    key,
			OriginalURL: val,
		})
		idx++
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return err
	}

	os.WriteFile(m.fileStorage, dataJSON, 0644)

	return nil
}
