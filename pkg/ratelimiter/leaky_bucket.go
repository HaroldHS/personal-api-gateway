package ratelimiter

/*
 * leaky_bucket.go implements a rate limiter object with leaky bucket algorithm
 *
 * Benefits / Pros:
 *  - Memory efficient due to limited capicity
 *  - Suitable for stable outflow rate
 *
 * Challenges / Cons:
 *  - Burst traffic fills up the queue with old requests and overthrown recent requests
 *  - Need to tune `capacity` and `leak rate` properly
 *
 */

import (
    "errors"
    "sync"
    "time"
)

type LeakyBucket struct {
    mu          sync.Mutex
    LastUpdated time.Time
    Capacity    int
    CurrentSize int
    LeakRate    int        // number of leaked units per second
}

func NewRateLimiterLeakyBucket(capacity int, leakRate int) *LeakyBucket {
    return &LeakyBucket{
        Capacity: capacity,
        LastUpdated: time.Now(),
        LeakRate: leakRate,
    }
}

func (lb *LeakyBucket) AddData(dataSize int) (bool, error) {
    // lock the leaky bucket to prevent race condition
    lb.mu.Lock()
    defer lb.mu.Unlock()

    currTime := time.Now()
    elapsedTime := currTime.Sub(lb.LastUpdated).Seconds()

    lb.LasUpdated = currTime

    lb.CurrentSize -= lb.LeakRate * elapsedTime                 // remove data based on leak rate and previous data
    lb.CurretSize = min(lb.CurrentSize + dataSize, lb.Capacity) // add new data for calculation

    if lb.CurrentSize >= dataSize {
        lb.CurrentSize -= dataSize
        return true, nil
    }

    return false, errors.New("too many requests")
}

