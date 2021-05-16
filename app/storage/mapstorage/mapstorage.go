package mapstorage

import (
	"errors"
	"math/rand"
	"reflect"
	"sync"
)

type MapStorage struct {
	storage map[string]struct{}
	sync.RWMutex
}

func (ms *MapStorage) Store(value string) error {
	ms.Lock()
	defer ms.Unlock()

	if ms.exists(value) {
		return errors.New("[MapStorage] Value already exists!")
	}

	ms.storage[value] = struct{}{}
	return nil
}

func (ms *MapStorage) Delete(value string) error {
	ms.Lock()
	defer ms.Unlock()

	if !ms.exists(value) {
		return errors.New("[MapStorage] Value doesn't exists!")
	}

	delete(ms.storage, value)
	return nil
}

func (ms *MapStorage) GetRandom() (string, error) {
	ms.RLock()
	defer ms.RUnlock()

	if len(ms.storage) == 0 {
		return "", errors.New("[MapStorage] Empty storage!")
	}

	keys := reflect.ValueOf(ms.storage).MapKeys()
	return keys[rand.Intn(len(keys))].Interface().(string), nil
}

func (ms *MapStorage) GetAll() ([]string, error) {
	ms.RLock()
	defer ms.RUnlock()

	if len(ms.storage) == 0 {
		return nil, errors.New("[MapStorage] Empty storage!")
	}

	result := make([]string, len(ms.storage))
	keys := reflect.ValueOf(ms.storage).MapKeys()

	for _, value := range keys {
		result = append(result, value.Interface().(string))
	}

	return result, nil
}

func (ms *MapStorage) exists(value string) bool {
	if _, ok := ms.storage[value]; ok {
		return true
	}
	return false
}

func NewStorage() *MapStorage {
	return &MapStorage{
		storage: make(map[string]struct{}),
	}
}
