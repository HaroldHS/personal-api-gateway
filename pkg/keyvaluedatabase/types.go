package keyvaluedatabase

import (
    "sync"
    "time"
)

type KeyValue struct {
    ExpiredAt   time.Time // For lazy deletion
    Key         string
    Value       any
    Next        *KeyValue // Chaining mechanism like "Linked List" in case of collision
    IsExpireSet bool
}

type KeyValueHashMap struct {
    mu   sync.Mutex
    Size int64
    List []*KeyValue // Array of Linked List
}

