package main

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "personal-api-gateway/internal/adapter/driven"
    "personal-api-gateway/internal/adapter/driver"
    "personal-api-gateway/internal/core/domain"
    "personal-api-gateway/internal/core/service"
    "personal-api-gateway/internal/core/util"
    "personal-api-gateway/pkg/log"
    "personal-api-gateway/pkg/ratelimiter"
)

func main() {
    logger := log.GetLoggerInstance()

    logger.Info("[main] Load JSON config.")
    configBytes, err := os.ReadFile("./example_config.json")
    if err != nil {
        logger.Error("Failed to read JSON config file: %v.", err)
        return
    }

    var config domain.JsonConfig
    err = json.Unmarshal(configBytes, &config)
    if err != nil {
        logger.Error("[main] Failed to unmarshal JSON config: %v.", err)
        return
    }

    logger.Info("[main] Starting built in key value database.")
    builtInKeyValueDbRepo := driven.NewBuiltInKeyValueDatabase()
    keyValueDb := service.New(builtInKeyValueDbRepo)

    logger.Info("[main] Starting rate limiter.")
    rateLimiter := ratelimiter.NewRateLimiterTokenBucket(1000)

    logger.Info("[main] Starting HTTP proxies.")
    httpProxies, err := util.NewBasicHttpProxies(config)
    if err != nil {
        logger.Error("[main] Failed to start proxies instances: %v.", err)
        return
    }

    logger.Info("[main] Initializing handlers.")
    httpDriver := driver.NewHttpDriver(&config, httpProxies, keyValueDb, rateLimiter)

    // Route all HTTP requests to a single entry point
    http.HandleFunc("/", httpDriver.HttpBasicEntryPoint)

    logger.Info("[main] Starting API gateway.")
    httpServer := &http.Server{
        Addr: ":8080",
    }

    go func() {
        if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            logger.Error("[main] Failed to serve connection: %v.", err)
        }
        logger.Info("[main] Stop serving new connections.")
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 20*time.Second)
    defer shutdownRelease()

    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        logger.Error("[main] Failed to shutdown HTTP server: %v.", err)
    }

    logger.Info("[main] Graceful shutdown completed.")
}
