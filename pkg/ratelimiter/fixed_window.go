package ratelimiter

/*
 * fixed_window.go implements a rate limiter object with fixed window algorithm
 *
 * Benefits / Pros:
 *  - Suitable for stable flow / traffic 
 *
 * Challenges / Cons:
 *  - Burst traffic at the edge / end of the window results in overflown quota
 *
 */

import (
    "errors"
    "sync"
    "time"
)

type FixedWindow struct {
    mu              sync.Mutex
    WindowStartTime time.Time
    MaxRequest      int
    NumOfRequests   int
    WindowSize      int        // in seconds
}

func NewRateLimiterFixedWindow(windowSize int, maxRequest int) *FixedWindow {
    return &FixedWindow{
        MaxRequest: maxRequest,
        WindowSize: windowSize,
        WindowStartTime: time.Now(),
    }
}

func (fw *FixedWindow) AllowRequest() (bool, error) {
    // lock the fixed window to prevent race condition
    fw.mu.Lock()
    defer fw.mu.Unlock()

    currTime := time.Now()
    elapsedTime := currTime.Sub(fw.WindowStartTime).Seconds()

    if elapsedTime >= float64(fw.WindowSize) {
        fw.NumOfRequests = 0
        fw.WindowStartTime = currTime
    } 

    if fw.NumOfRequests < fw.MaxRequest {
        fw.NumOfRequests += 1
        return true, nil
    }

    return false, errors.New("too many requests")
}
