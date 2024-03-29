# GoUtils

[Documentation](https://pkg.go.dev/github.com/nano-interactive/go-utils/v2)

[![Testing](https://github.com/nano-interactive/go-utils/actions/workflows/test.yml/badge.svg)](https://github.com/nano-interactive/go-utils/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/nano-interactive/go-utils/branch/master/graph/badge.svg?token=JQTAGQ11DS)](https://codecov.io/gh/nano-interactive/go-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/nano-interactive/go-utils)](https://goreportcard.com/report/github.com/nano-interactive/go-utils)

<img src="https://github.com/egonelbre/gophers/blob/master/vector/superhero/standing.svg" width="300" />
## This is a collection of useful packages including:

- Config
- Environment
- Logging
- Signals
- Testing

## Config

```go
package main

import (
    "github.com/nano-interactive/go-utils/v2/config
)

// Defaults

var DefaultConfig = Config {
    Env: "development",
    Name: "config",
    Type: "yaml",
}

func main() {
    config, err := config.New(config.Config)

    if err != nil {
        // Failed to load configuration
    }
}

```

## Environment

## Logging

## Signals

## Testing
