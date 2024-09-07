package storage

import "errors"

type RaiStorage struct {
	data map[string]string
}

func NewRaiStorage() *RaiStorage {
	return &RaiStorage{
		data: make(map[string]string),
	}
}

func (rs *RaiStorage) Get(key string) (*string, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return &value, nil
}

func (rs *RaiStorage) Put(key string, value string) error {
	rs.data[key] = value
	return nil
}

func (rs *RaiStorage) Post(key string, value string) error {
	if _, exists := rs.data[key]; exists {
		return errors.New("key already exists")
	}
	rs.data[key] = value
	return nil
}

func (rs *RaiStorage) Delete(key string) error {
	if _, exists := rs.data[key]; !exists {
		return errors.New("key not found")
	}
	delete(rs.data, key)
	return nil
}
