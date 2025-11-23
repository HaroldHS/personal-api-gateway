package service

import (
    "errors"

    "personal-api-gateway/internal/core/port"
)

type KeyValueDatabaseService struct {
    keyValueDatabaseRepository port.KeyValueDatabaseRepository
}

func New(keyValueDbRepository port.KeyValueDatabaseRepository) *KeyValueDatabaseService {
    return &KeyValueDatabaseService{
        keyValueDatabaseRepository: keyValueDbRepository,
    }
}

func (kvds *KeyValueDatabaseService) Save(key string, value any, ttl int64) error {
    if key == "" {
        return errors.New("key should not be empty")
    }

    if err := kvds.keyValueDatabaseRepository.Insert(key, value, ttl); err != nil {
        return errors.New("failed to insert object: " + err.Error())
    }

    return nil
}
