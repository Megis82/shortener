package storage

import (
	"context"
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

func (m *MemoryStorage) Add(ctx context.Context, key string, value string) error {
	m.data[key] = value
	return nil
}

func (m *MemoryStorage) AddBatch(ctx context.Context, values map[string]string) error {
	for key, val := range values {
		m.data[key] = val
	}
	return nil
}

func (m *MemoryStorage) Find(ctx context.Context, key string) (string, bool, error) {
	value, ok := m.data[key]
	return value, ok, nil
}

func (m *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func NewMemoryStorage(fileStorage string) (*MemoryStorage, error) {

	mem := &MemoryStorage{data: make(map[string]string), fileStorage: fileStorage}

	if fileStorage == "" {
		return mem, nil
	}

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
		mem.Add(context.Background(), str.ShortURL, str.OriginalURL)
	}

	return mem, nil
}

func (m *MemoryStorage) Close() error {
	if m.fileStorage == "" {
		return nil
	}

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
