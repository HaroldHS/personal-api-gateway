package ratelimiter

/*
 * sliding_window.go implements a rate limiter object with sliding window algorithm
 *
 * Benefits / Pros:
 *  - 
 *
 * Challenges / Cons:
 *  - 
 *
 */

import (
    "errors"
    "sync"
    "time"
)

type SlidingWindow struct {
    mu              sync.Mutex
    MaxRequest      int
    WindowSize      int        // in seconds
}

func NewRateLimiterSlidingWindow(windowSize int, maxRequest int) *SlidingWindow {
    return &SlidingWindow{
        MaxRequest: maxRequest,
        WindowSize: windowSize,
    }
}

func (sw *SlidingWindow) AllowRequest() (bool, error) {
    // lock the fixed window to prevent race condition
    sw.mu.Lock()
    defer sw.mu.Unlock()

    currTime := time.Now()

    return false, errors.New("too many requests")
}
