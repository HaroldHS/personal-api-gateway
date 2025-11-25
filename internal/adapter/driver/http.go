package driver

import (
    "fmt"
    "net/http"
    "reflect"

    "personal-api-gateway/internal/core/domain"
    "personal-api-gateway/internal/core/service"
    "personal-api-gateway/pkg/log"
    "personal-api-gateway/pkg/ratelimiter"
)

type HttpDriver struct {
    JsonConfig              domain.JsonConfig
    KeyValueDatabase        *service.KeyValueDatabaseService
    BasicLogger             *log.BasicLogger
    TokenBucketRateLimiter  *ratelimiter.TokenBucket
}

func NewHttpDriver(
    jsonConfig domain.JsonConfig,
    keyValueDb *service.KeyValueDatabaseService,
    tokenBucketRateLimiter *ratelimiter.TokenBucket) *HttpDriver {

    httpDriver := &HttpDriver{}

    if !reflect.DeepEqual(jsonConfig, domain.JsonConfig{}) {
        httpDriver.JsonConfig = jsonConfig
    }

    if !reflect.DeepEqual(keyValueDb, service.KeyValueDatabaseService{}) {
        httpDriver.KeyValueDatabase = keyValueDb
    }

    if !reflect.DeepEqual(tokenBucketRateLimiter, ratelimiter.TokenBucket{}) {
        httpDriver.TokenBucketRateLimiter = tokenBucketRateLimiter
    }

    return httpDriver
}

func (hd *HttpDriver) HttpBasicEntryPoint(res http.ResponseWriter, req *http.Request) {
    /*
    if hd.JsonConfig.LogEndpointUrl{
        hd.BasicLogger.Info(req.URL.Path)
    }
    */

    fmt.Fprintf(res, "Hello World %s", req.URL.Path)
}

