package ratelimiter

/*
 * token_bucket.go implements a rate limiter object with token bucket algorithm
 *
 * Benefits / Pros:
 *  - Allow burst traffic in short period of time (as the tokens are available)
 *
 * Challenges / Cons:
 *  - Need to tune `bucket size` and `refill rate` properly
 *
 */

import (
    "errors"
    "sync"
    "time"
)

type TokenBucketLazyRefill struct {
    LastRefill time.Time // timestamp for lazy refill calculation
    BucketSize int64     // maximum throttle size
    RefillRate int64     // number of tokens to be prefilled every second
    Token      int64     // number of tokens container
}

type TokenBucket struct {
    mu    sync.Mutex
    Rules map[string]*TokenBucketLazyRefill
}

func NewRateLimiterTokenBucket() *TokenBucket {
    return &TokenBucket{}
}

func (tb *TokenBucket) SetRuleConfig(rule string, refillRate int64, bucketSize int64) {
    if _, ok := tb.Rules[rule]; !ok {
        tb.Rules[rule] = &TokenBucketLazyRefill{
            RefillRate: refillRate,
            BucketSize: bucketSize,
            Token:      bucketSize,
            LastRefill: time.Now(),
        }
    }
}

func (tb *TokenBucket) AllowRequest(rule string) (bool, error) {
    // lock the token bucket to prevent race condition
    tb.mu.Lock()
    defer tb.mu.Unlock()

    ruleTokenBucket, ok := tb.Rules[rule]
    if !ok {
        return false, errors.New("rule not found")
    }

    currTime := time.Now()
    elapsedTime := currTime.Sub(ruleTokenBucket.LastRefill).Seconds()

    // Accept if token is still available
    if ruleTokenBucket.Token >= 1 {
        // Consume token of the given rule
        tb.Rules[rule].Token -= 1
        return true, nil
    }

    // Refill the bucket when the lazy refill time is bigger than 1 second
    if elapsedTime > 1.00 {
        tb.Rules[rule].Token = min(ruleTokenBucket.Token + ruleTokenBucket.RefillRate, ruleTokenBucket.BucketSize)
        tb.Rules[rule].LastRefill = currTime
    }

    return false, errors.New("too many requests")
}
