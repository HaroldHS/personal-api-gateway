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
    "personal-api-gateway/pkg/log"
    "personal-api-gateway/pkg/ratelimiter"
)

func main() {
    logger := log.GetLoggerInstance()

    logger.Info("[main] Load JSON config.")
    configFile, err := os.Open("./example_config.json")
    if err != nil {
        logger.ErrorFormat("Failed to open JSON config file: %s.", err)
        return
    }
    defer configFile.Close()

    var config domain.JsonConfig
    err = json.NewDecoder(configFile).Decode(&config)
    if err != nil {
        logger.ErrorFormat("[main] Failed to decode JSON config file: %s.", err)
        return
    }

    logger.Info("[main] Starting built in key value database.")
    builtInKeyValueDbRepo := driven.NewBuiltInKeyValueDatabase()
    keyValueDb := service.New(builtInKeyValueDbRepo)

    logger.Info("[main] Starting rate limiter.")
    rateLimiter := ratelimiter.NewRateLimiterTokenBucket()

    logger.Info("[main] Initializing handlers.")
    httpDriver := driver.NewHttpDriver(config, keyValueDb, rateLimiter)

    http.HandleFunc("/", httpDriver.HttpBasicEntryPoint)

    logger.Info("[main] Starting API gateway.")
    httpServer := &http.Server{
        Addr: ":8080",
    }

    go func() {
        if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            logger.ErrorFormat("[main] Failed to serve connection: %v.", err)
        }
        logger.Info("[main] Stop serving new connections.")
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 20*time.Second)
    defer shutdownRelease()

    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        logger.ErrorFormat("[main] Failed to shutdown HTTP server: %v.", err)
    }

    logger.Info("[main] Graceful shutdown completed.")
}
