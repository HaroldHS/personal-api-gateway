package port

// key-value database in "pkg" module
type KeyValueDatabaseRepository interface {
    Delete(key string)
    Get(key string) (any, error)
    Insert(key string, value any, ttl int64) error
}

type KeyValueDatabasePortInterface interface {
    Delete(key string) error
    Get(key string) (any, error)
    Save(key string, value any, ttl int64) error
}
