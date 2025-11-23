package driven

import (
    "personal-api-gateway/pkg/keyvaluedatabase"
)

const (
    DefaultStorageSize int64 = 1000
)

func NewBuiltInKeyValueDatabase() *keyvaluedatabase.KeyValueHashMap {
    return keyvaluedatabase.NewKeyValueDatabase(DefaultStorageSize)
}
