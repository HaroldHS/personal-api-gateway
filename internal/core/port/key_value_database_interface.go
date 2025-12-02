package port

// key-value database in "pkg" module
type KeyValueDatabaseRepository interface {
    Insert(key string, value any, ttl int64) error
}

type KeyValueDatabasePortInterface interface {
    Save(key string, value any, ttl int64) error
}
