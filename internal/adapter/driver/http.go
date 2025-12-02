package driver

import (
    "context"
    "net/http"
    "reflect"
    "time"

    "personal-api-gateway/internal/core/domain"
    "personal-api-gateway/internal/core/port"
    "personal-api-gateway/internal/core/util"
    "personal-api-gateway/pkg/log"
    "personal-api-gateway/pkg/ratelimiter"
)

type HttpDriver struct {
    JsonConfig             *domain.JsonConfig
    HttpProxies            *util.BasicHttpProxies
    HttpModifier           port.HttpModifierPortInterface
    KeyValueDatabase       port.KeyValueDatabasePortInterface
    TokenBucketRateLimiter *ratelimiter.TokenBucket
}

func NewHttpDriver(
    jsonConfig *domain.JsonConfig,
    httpProxies *util.BasicHttpProxies,
    httpModifier port.HttpModifierPortInterface,
    keyValueDb port.KeyValueDatabasePortInterface,
    tokenBucketRateLimiter *ratelimiter.TokenBucket) *HttpDriver {

    httpDriver := &HttpDriver{
        HttpModifier: httpModifier,
        KeyValueDatabase: keyValueDb,
    }

    if !reflect.DeepEqual(jsonConfig, domain.JsonConfig{}) {
        httpDriver.JsonConfig = jsonConfig
    }

    if !reflect.DeepEqual(httpProxies, util.BasicHttpProxies{}) {
        httpDriver.HttpProxies = httpProxies
    }

    if !reflect.DeepEqual(tokenBucketRateLimiter, ratelimiter.TokenBucket{}) {
        httpDriver.TokenBucketRateLimiter = tokenBucketRateLimiter
    }

    return httpDriver
}

func (hd *HttpDriver) HttpBasicEntryPoint(res http.ResponseWriter, req *http.Request) {
    logger := log.GetLoggerInstance()

    if hd.JsonConfig.LogEndpointUrl{
        logger.Info("[*] Incoming request: " + req.Method + " " + req.URL.Path)
    }

    // Define and set context timeout to prevent timing related error
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    req = req.WithContext(ctx)

    hd.TokenBucketRateLimiter.SetRuleConfig(req.URL.Path, 10, 10)
    if allowed, err := hd.TokenBucketRateLimiter.AllowRequest(req.URL.Path); !allowed || (err != nil) {
        http.Error(res, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    destProxy, err := hd.HttpProxies.GetBasicHttpProxy(req.URL.Path)
    if err != nil {
        http.Error(res, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    cfgEndpointObj := hd.JsonConfig.Endpoints[req.URL.Path]
    if cfgEndpointObj.DestinationScheme != "http" && cfgEndpointObj.DestinationScheme != "https" {
        http.Error(res, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    requestedHost := req.Host

    // Modify HTTP request
    hd.HttpModifier.ModifyRequestHeader(req, cfgEndpointObj)

    // Define custom proxy error handler
    destProxy.ErrorHandler = func (res http.ResponseWriter, req *http.Request, err error) {
        logger.Error("[HttpBasicEntryPoint] Proxy error occurred: %v", err)
        http.Error(res, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Modify HTTP response
    destProxy.ModifyResponse = func(res *http.Response) error {
        res.Header.Set("Host", requestedHost)
        return nil
    }

    // Serve the request through the reverse proxy
    destProxy.ServeHTTP(res, req)
}

