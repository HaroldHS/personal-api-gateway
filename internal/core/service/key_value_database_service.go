package service

import (
    "errors"

    "personal-api-gateway/internal/core/port"
)

type KeyValueDatabaseService struct {
    keyValueDatabaseRepository port.KeyValueDatabaseRepository
}

func NewKeyValueDatabase(keyValueDbRepository port.KeyValueDatabaseRepository) *KeyValueDatabaseService {
    return &KeyValueDatabaseService{
        keyValueDatabaseRepository: keyValueDbRepository,
    }
}

// @Summary: Delete a key-value object with the given key.
// @Description: Delete a key-value object with the given key in key-value database repository.
//
// @Param  key   string "e.g. `Test`"
func (kvds *KeyValueDatabaseService) Delete(key string) error {
    if key == "" {
        return errors.New("key should not be empty")
    }

    kvds.keyValueDatabaseRepository.Delete(key)
    return nil
}

// @Summary: Get a value object with the given key.
// @Description: Get a value object with the given key in key-value database repository.
//
// @Param  key   string "e.g. `Test`"
func (kvds *KeyValueDatabaseService) Get(key string) (any, error) {
    if key == "" {
        return nil, errors.New("key should not be empty")
    }

    return kvds.keyValueDatabaseRepository.Get(key)
}

// @Summary: Save key-value object into key value database.
// @Description: Save a new key-value object into key-value database repository.
//
// @Param key   string "e.g. `Test`"
// @Param value any    "e.g. `Hello World` or [`Hello`, `World`]"
// @Param ttl   int64  "e.g. int64(300)"
func (kvds *KeyValueDatabaseService) Save(key string, value any, ttl int64) error {
    if key == "" {
        return errors.New("key should not be empty")
    }

    if err := kvds.keyValueDatabaseRepository.Insert(key, value, ttl); err != nil {
        return errors.New("failed to insert object: " + err.Error())
    }

    return nil
}
