package service

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
	ErrNotFound      = errors.New("record not found")
	ErrParamNotEmpty = errors.New("key or value params not be empty")
	ErrModelEmpty    = errors.New("model is empty")
)

var (
	memoryModel map[string]string
)



type IMemoryStorage interface {
	SaveKey(key, val string) error
	GetKey(key string) (string, error)
	GetAll() (map[string]string, error)
	FlushMemory()
}

type MemoryStorage struct {
	mutex sync.RWMutex
}

func NewMemoryStorage(jsonData *map[string]string) *MemoryStorage {
	if jsonData != nil {
		memoryModel = *jsonData
	} else {
		memoryModel = make(map[string]string)
	}
	return &MemoryStorage{}
}

func (m *MemoryStorage) SaveKey(key, val string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if key == "" || val == "" {
		return ErrParamNotEmpty
	}

	if memoryModel[key] != "" {
		return ErrAlreadyExists
	}

	memoryModel[key] = val
	return nil
}

func (m *MemoryStorage) GetKey(key string) (string, error) {
	if key == "" {
		return "", ErrParamNotEmpty
	}

	if memoryModel[key] == "" {
		return "", ErrNotFound
	}
	return memoryModel[key], nil
}

func (m *MemoryStorage) FlushMemory() {
	for k := range memoryModel {
		delete(memoryModel, k)
	}
}
func (m *MemoryStorage) GetAll() (map[string]string, error) {
	if len(memoryModel) > 0 {
		return memoryModel, nil
	}

	return nil, ErrModelEmpty
}
