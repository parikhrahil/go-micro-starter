# go-micro-starter

A lightweight, high-performance Go microservice starter kit configured with structured logging, configuration management, native HTTP middleware, and pre-configured PostgreSQL and Redis data-layer abstractions.

## Core Features

* **HTTP Router & Server:** Thin wrapper over high-performance routing with structured logging context injection.
* **Structured Logging:** JSON-based contextual logger supporting trace IDs across asynchronous contexts.
* **Database Management (`sqldb`):** PostgreSQL database connection pooling with health check primitives.
* **Caching Engine (`cache`):** Built-in Redis client wrapping connection configurations and utility operations (K/V strings, hashes, queues).
* **Environment Configuration:** Type-safe extraction of parameters from dynamic system states or `.env` files.

---

## Architecture Blueprint

```text
  [ Client Request ]
          │
          ▼
   ┌──────────────┐
   │  Middleware  │ ──► Request ID Generation & Logger Context Mapping
   └──────┬───────┘
          │
          ▼
   ┌──────────────┐
   │  Web Router  │ ──► Native Multiplexer Wrapper
   └──────┬───────┘
          │
          ▼
   ┌──────────────┐
   │ Domain Logic │
   └────┬────┬────┘
        │    │
        │    └──────────────────────┐
        ▼                           ▼
 ┌──────────────┐            ┌──────────────┐
 │  SQL Engine  │            │ Cache Engine │
 │ (PostgreSQL) │            │   (Redis)    │
 └──────────────┘            └──────────────┘
```

---

## Configuration Settings

The starter library consumes system variables natively. Copy the `.env.example` file to target local defaults:

```bash
cp .env.example .env
```

### Required Fields Breakdown

| Variable | Description | Default Vector |
| :--- | :--- | :--- |
| `APP_ENV` | Application environment state (`development`, `production`) | `development` |
| `HTTP_PORT` | Port the internal HTTP API engine binds to | `8080` |
| `DB_DSN` | PostgreSQL Data Source Name string | `postgres://...` |
| `REDIS_ADDR` | Network location of the Redis instance | `localhost:6379` |

---

## Getting Started

### 1. Bootstrapping a Core Server

```go
package main

import (
	"context"
	"net/http"

	"github.com/parikhrahil/go-micro-starter/config"
	"github.com/parikhrahil/go-micro-starter/logger"
	"github.com/parikhrahil/go-micro-starter/server"
)

func main() {
	// Initialize context and logger
	ctx := context.Background()
	log := logger.NewJSONLogger()

	// Parse settings
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(ctx, "failed to parse environment variables", "error", err)
	}

	// Instantiate HTTP engine
	srv := server.New(server.Config{
		Port: cfg.HTTPPort,
	})

	log.Info(ctx, "starting HTTP execution runner", "port", cfg.HTTPPort)
	if err := srv.Start(ctx); err != nil {
		log.Fatal(ctx, "http server loop broke down", "error", err)
	}
}
```

### 2. Using Data Access Backends

```go
// Initializing Postgres Connectivity via your package
dbClient, err := sqldb.Connect(cfg.DB_DSN)
if err != nil {
    log.Fatal(ctx, "database handshake aborted", "error", err)
}

// Initializing Redis Core Backends
redisClient, err := cache.Connect(cache.Config{
    Address:  cfg.RedisAddr,
    Password: cfg.RedisPassword,
    DB:       0,
})
```
