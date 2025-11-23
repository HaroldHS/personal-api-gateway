package main

import (
    "fmt"

    "personal-api-gateway/internal/adapter/driven"
    "personal-api-gateway/internal/core/service"
)

func main() {
    fmt.Println("[*] Starting API gateway....")

    builtInKeyValueDbRepo := driven.NewBuiltInKeyValueDatabase()
    builtInKeyValueDbService := service.New(builtInKeyValueDbRepo)
    builtInKeyValueDbService.Save("test", "Hello World", 10)

    fmt.Println("[*] Graceful shutdown....")
}
