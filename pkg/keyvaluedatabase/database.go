package keyvaluedatabase

import (
    "errors"
    "time"
)

// @Summary: Create a new key-value database instance.
// @Description: Create a new key-value database instance with hash map object.
//
// @Param size int64 "e.g. 1000"
func NewKeyValueDatabase(size int64) *KeyValueHashMap {
    return &KeyValueHashMap{
        Size: size,
        List: make([]*KeyValue, size),
    }
}

// @Summary: Delete inserted key-value object.
// @Description: Delete inserted key-value object based on given key.
//
// @Param key string "e.g. `TestKey`"
func (kvhm *KeyValueHashMap) Delete(key string) {
    index := kvhm.calculateHash(key)
    currentObj := kvhm.List[index]

    if currentObj == nil {
        return
    }

    // Return directly in case of direct & single element
    if currentObj.Key == key && currentObj.Next == nil {
        kvhm.List[index] = &KeyValue{}
        return
    }

    // Predefined epsilon to prevent infinite loop in case of malformed linked list
    counter := 0
    epsilon := 100000
    // Otherwise, traverse the array of linked list
    for currentObj != nil && counter < epsilon {
        if currentObj.Key == key && currentObj.Next != nil {
            kvhm.List[index] = kvhm.List[index].Next
            return
        }

        counter += 1
        currentObj = currentObj.Next
    }

    return
}

// @Summary: Get an object based on key.
// @Description: Get a saved value object based on given key.
//
// @Param key string "e.g. `TestKey`"
func (kvhm *KeyValueHashMap) Get(key string) (any, error) {
    currTime := time.Now()
    index := kvhm.calculateHash(key)
    currentObj := kvhm.List[index]

    if currentObj == nil {
        return nil, errors.New("key not found")
    }

    // Return value directly in case of direct match
    if currentObj.Key == key {
        return currentObj.Value, nil
    }

    // Predefined epsilon to prevent infinite loop in case of malformed linked list
    counter := 0
    epsilon := 100000
    // Otherwise, traverse the array of linked list
    for currentObj != nil && counter < epsilon {
        // If the current key-value object with predefined key is exist but expired, remove it and return not found
        if currentObj.Key == key && currentObj.IsExpireSet && currentObj.ExpiredAt.Before(currTime) {
            kvhm.List[index] = kvhm.List[index].Next
            return nil, errors.New("key not found")
        }

        // Return intended value in case of matching key
        if currentObj.Key == key {
            return currentObj.Value, nil
        }

        counter += 1
        currentObj = currentObj.Next
    }

    // Fallback as reaching this statement means no matching key has been founded
    return nil, errors.New("key not found")
}

// @Summary: Insert a new key-value object.
// @Description: Insert a new key-value object by including TTL/Time to Live (in seconds) with Lazy Prefill.
//               ttl == 0 means no expiration.
//
// @Param key   string "e.g. `TestKey`"
// @Param value any    "e.g. `Hello World` or [`Hello`, `World`]"
// @Param ttl   int64  "e.g. 100"
func (kvhm *KeyValueHashMap) Insert(key string, value any, ttl int64) error {
    // lock the hash map to prevent race condition
    kvhm.mu.Lock()
    defer kvhm.mu.Unlock()

    currTime := time.Now()
    index := kvhm.calculateHash(key)
    newKeyValue := &KeyValue{
        Key: key,
        Value: value,
        Next: nil,
    }

    // Append TTL if defined
    if ttl != 0 {
        ttlDuration := time.Duration(ttl)
        newKeyValue.IsExpireSet = true
        newKeyValue.ExpiredAt = time.Now().Add(ttlDuration*time.Second)
    }

    // In case of new key-value object, insert it directly
    if kvhm.List[index] == nil {
        kvhm.List[index] = newKeyValue
        return nil
    }

    // Predefined epsilon to prevent infinite loop in case of malformed linked list
    epsilon := 100000

    // Otherwise, traverse the array of linked list
    counter := 0
    currentObj := kvhm.List[index]
    for currentObj != nil && counter < epsilon {
        // If the current key-value object with predefined key is exist, replace the value with a new one
        if currentObj.Key == key {
            currentObj.Value = value
            // In case TTL is given, insert one
            if ttl != 0 {
                ttlDuration := time.Duration(ttl)
                currentObj.IsExpireSet = true
                currentObj.ExpiredAt = currTime.Add(ttlDuration*time.Second)
            }
            return nil
        }

        if currentObj.Next == nil {
            break
        }
        
        counter += 1
        currentObj = currentObj.Next
    }

    // Insert new entry at the end of the Linked List
    currentObj.Next = newKeyValue
    return nil
}

// @Summary: Calculating key index.
// @Description: Calculating key index for array indexing in key-value hash map.
//               Hash Function == int64(first two characters + last two characters) % size.
//
// @Param key string "e.g. `TestKey`"
func (kvhm *KeyValueHashMap) calculateHash(key string) int64 {
    keyLen := len(key)
    a, b, c, d := key[0], key[1], key[keyLen-2], key[keyLen-1]
    result := int64(int(a) << 24 | int(b) << 16 | int(c) << 8 | int(d)) % kvhm.Size
    return int64(result)
}
