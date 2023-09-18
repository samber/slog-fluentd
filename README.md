# slog: Fluentd handler

[![tag](https://img.shields.io/github/tag/samber/slog-fluentd.svg)](https://github.com/samber/slog-fluentd/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20.3-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/slog-fluentd?status.svg)](https://pkg.go.dev/github.com/samber/slog-fluentd)
![Build Status](https://github.com/samber/slog-fluentd/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/slog-fluentd)](https://goreportcard.com/report/github.com/samber/slog-fluentd)
[![Coverage](https://img.shields.io/codecov/c/github/samber/slog-fluentd)](https://codecov.io/gh/samber/slog-fluentd)
[![Contributors](https://img.shields.io/github/contributors/samber/slog-fluentd)](https://github.com/samber/slog-fluentd/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/slog-fluentd)](./LICENSE)

A [Fluentd](https://www.fluentd.org/) Handler for [slog](https://pkg.go.dev/log/slog) Go library.

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): `slog.Handler` chaining, fanout, routing, failover, load balancing...
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-sampling](https://github.com/samber/slog-sampling): `slog` sampling policy
- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-echo](https://github.com/samber/slog-echo): Echo middleware for `slog` logger
- [slog-fiber](https://github.com/samber/slog-fiber): Fiber middleware for `slog` logger
- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-rollbar](https://github.com/samber/slog-rollbar): A `slog` handler for `Rollbar`
- [slog-sentry](https://github.com/samber/slog-sentry): A `slog` handler for `Sentry`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-telegram](https://github.com/samber/slog-telegram): A `slog` handler for `Telegram`
- [slog-mattermost](https://github.com/samber/slog-mattermost): A `slog` handler for `Mattermost`
- [slog-microsoft-teams](https://github.com/samber/slog-microsoft-teams): A `slog` handler for `Microsoft Teams`
- [slog-webhook](https://github.com/samber/slog-webhook): A `slog` handler for `Webhook`
- [slog-kafka](https://github.com/samber/slog-kafka): A `slog` handler for `Kafka`
- [slog-parquet](https://github.com/samber/slog-parquet): A `slog` handler for `Parquet` + `Object Storage`
- [slog-zap](https://github.com/samber/slog-zap): A `slog` handler for `Zap`
- [slog-zap](https://github.com/samber/slog-zerolog): A `slog` handler for `Zerolog`
- [slog-zap](https://github.com/samber/slog-logrus): A `slog` handler for `Logrus`

## 🚀 Install

```sh
go get github.com/samber/slog-fluentd
```

**Compatibility**: go >= 1.21

No breaking changes will be made to exported APIs before v2.0.0.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/slog-fluentd](https://pkg.go.dev/github.com/samber/slog-fluentd)

### Fluentd settings

```
<source>
    @type forward
    bind 0.0.0.0
    port 24224
</source>
```

### Handler options

```go
type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// connection to Fluentd
	Client *fluentd.Fluentd
    Tag    string

	// optional: customize json payload builder
	Converter Converter
}
```

Attributes will be injected in log payload.

Fluentd `tag` can be inserted in logger options or in a record attribute of type string.

### Example

```go
import (
	"github.com/fluent/fluent-logger-golang/fluent"
	slogfluentd "github.com/samber/slog-fluentd"
	"log/slog"
)

func main() {
	// docker-compose up -d
	client, err := fluent.New(fluent.Config{
		FluentHost:    "localhost",
		FluentPort:    24224,
		FluentNetwork: "tcp",
		MarshalAsJSON: true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}


    logger := slog.New(
        slogfluentd.Option{
            Level: slog.LevelDebug,
            Client: client,
            Tag: "api",
        }.NewFluentdHandler(),
    )
    logger = logger.
        With("environment", "dev").
        With("release", "v1.0.0")

    // log error
    logger.
        With("tag", "api.sql").
        With("query.statement", "SELECT COUNT(*) FROM users;").
        With("query.duration", 1*time.Second).
        With("error", fmt.Errorf("could not count users")).
        Error("caramba!")

    // log user signup
    logger.
        With(
            slog.Group("user",
                slog.String("id", "user-123"),
                slog.Time("created_at", time.Now()),
            ),
        ).
        Info("user registration")
}
```

Output:

```json
// tag: api.sql
{
    "timestamp":"2023-04-10T14:00:0.000000+00:00",
    "level":"ERROR",
    "message":"caramba!",
    "tag":"api.sql",
    "error":{
        "error":"could not count users",
        "kind":"*errors.errorString",
        "stack":null
    },
    "extra":{
        "environment":"dev",
        "release":"v1.0.0",
        "query.statement":"SELECT COUNT(*) FROM users;",
        "query.duration": "1s"
    }
}


// tag: api
{
    "timestamp":"2023-04-10T14:00:0.000000+00:00",
    "level":"INFO",
    "message":"user registration",
    "tag":"api",
    "error":null,
    "extra":{
        "environment":"dev",
        "release":"v1.0.0",
        "user":{
            "id":"user-123",
            "created_at":"2023-04-10T14:00:0.000000+00:00"
        }
    }
}
```

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/slog-fluentd)
- Fix [open issues](https://github.com/samber/slog-fluentd/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/slog-fluentd)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
