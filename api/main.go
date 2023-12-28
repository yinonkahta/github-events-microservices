package main

import (
	"fmt"
	"github-events-microservices/api/config"
	"github-events-microservices/api/net"
	"github-events-microservices/logging"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(logging.Create())
	slog.Info("Github Events API started")

	handler := net.NewRequestsHandler(config.ApiConfiguration.MongoDbUrl, config.ApiConfiguration.MongoDbPort)

	http.HandleFunc("/list", handler.List)
	http.HandleFunc("/count", handler.Count)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error(fmt.Sprintf("github events api failed. reason: %s", err.Error()))
		os.Exit(1)
	}
}
