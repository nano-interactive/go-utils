# GoUtils

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
    "github.com/nano-interactive/go-utils/config
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