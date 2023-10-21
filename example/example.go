package main

import (
	"fmt"
	"log"
	"time"

	"log/slog"

	"github.com/fluent/fluent-logger-golang/fluent"
	slogfluentd "github.com/samber/slog-fluentd/v2"
)

func main() {
	// docker-compose up -d
	client, err := fluent.New(fluent.Config{
		FluentHost:    "localhost",
		FluentPort:    24224,
		FluentNetwork: "tcp",
		Async:         true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := slog.New(slogfluentd.Option{Level: slog.LevelDebug, Client: client}.NewFluentdHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}
