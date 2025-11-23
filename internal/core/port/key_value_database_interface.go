package port

//import (
//    "personal-api-gateway/internal/core/domain"
//)

// key-value database in "pkg" module
type KeyValueDatabaseRepository interface {
    Insert(key string, value any, ttl int64) error
}

type KeyValueDatabaseService interface {
    Save(key string, value any, ttl int64) error
}
