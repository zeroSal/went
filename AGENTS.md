## Project Structure

- `main.go` — Entry point, creates root `cobra.Command`, uses embed.FS for templates/static
- `cmd/` — Command implementations using custom `went-command` pattern (`serve.go`)
- `app/kernel.go` — DI kernel that creates and runs the fx.App
- `app/container.go` — FX module container with DI providers
- `app/build_specs.go` — Build metadata (version, channel, build date)
- `app/bootstrap/` — Bootstrap providers:
  - `init.go` — Initialization function, validates env, creates working dirs, initializes loggers
  - `module/loggers.go` — Logger providers (AuditLogger, ErrorLogger)
  - `module/iris.go` — Iris web framework provider (uses Django template engine, enables hot-reload in dev)
- `app/config/env.go` — Env struct with `Load()` and `Validate()` methods
- `templates/` — Embedded template files (banner.template)
- `static/` — Embedded static files
- `.env.dist` — Environment template file (must be copied to .env and customized)

---

## Template Engine

Iris uses the Django template engine (`.html.django` extension). Static files are served from the embedded filesystem at `/static`.

---

## Environment Configuration

Copy `.env.dist` to `.env` and configure:

| Variable  | Default    | Description                      |
|-----------|------------|----------------------------------|
| `ENV`     | `dev`      | Environment (`dev` or `prod`)    |
| `VAR_DIR` | `var`      | Base directory for runtime files |
| `HOST`    | `127.0.0.1`| Server bind address              |
| `PORT`    | `3096`     | Server port                      |

The `Validate()` method enforces `ENV` as `dev` or `prod`, and port range `1-65535`.

---

## Tech Stack

- **Iris v12** — Web framework
- **Cobra** — Root command framework (integrated via `went-command`)
- **Uber FX** — Dependency injection container
- **went-clio** — CLI output formatting (console input/output)
- **went-logger** — File logging (AuditLogger, ErrorLogger)
- **went-command** — Custom command pattern
- **air** — Auto-reload during development
- **golangci-lint v2** — Linter aggregator

---

## Patterns

### Command Pattern

Commands implement `command.Interface` and use `GetHeader()` to define cobra metadata:

```go
var _ command.Interface = (*ServeCmd)(nil)
type ServeCmd struct {
    command.Base
}

func NewServeCmd() command.Interface { ... }

func (c *ServeCmd) GetHeader() command.Header {
    return command.Header{
        Use:   "serve",
        Short: "Run the web server",
        Long:  "...",
    }
}

func (c *ServeCmd) Invoke() any {
    return c.run
}

func (c *ServeCmd) run(...) error { ... }
```

### Kernel Pattern

The `Kernel` wraps fx.App creation and execution:

```go
kernel := app.NewKernel(EmbedFS, specs, clio)
kernel.Run(instance.Invoke())
```

### FX Modules

Providers in `container.go` wrap types into FX-injectable named types:

```go
var Container = fx.Module(
    "container",
    fx.Provide(module.IrisProvider),
    fx.Provide(module.AuditLoggerProvider),
    fx.Provide(module.ErrorLoggerProvider),
    fx.Provide(config.LoadEnv),
)
```

### Env Config

`config/env.go` defines a single typed struct with `Load()` (reads from environment) and `Validate()` methods. It uses `godotenv` to load `.env` file. No raw `os.Getenv` calls outside this file.

### Runtime Directories

The bootstrap initialization creates the following directories at startup:
- `var/logs/` — Log files (audit.log, error.log)
- `var/uploads/` — Uploaded files

---

## Layer Conventions

Both `app/service/` and `app/model/` follow the same rule: **one package per domain**. The package name defines the domain context; each filename describes the type of object it contains, never repeating the domain name.

### `app/service/`

```
app/service/
├── domain_1/
│   └── client.go            # package domain_1
└── domain_2/
    ├── client.go            # package domain_2
    ├── client_interface.go  # package domain_2
    └── loader.go            # package domain_2
```

Typical filenames: `client.go`, `client_interface.go`, `loader.go`, `resolver.go`, `cache.go`, `parser.go`

### `app/model/`

```
app/model/
├── domain_1/
│   └── entity.go            # package domain_1
└── domain_2/
    ├── entity.go            # package domain_2
    ├── dto.go               # package domain_2
    ├── request.go           # package domain_2
    └── response.go          # package domain_2
```

Typical filenames: `entity.go`, `dto.go`, `request.go`, `response.go`, `enum.go`, `event.go`

> **Why:** the domain lives in the package, the role lives in the filename.
> `github/github_client.go` ❌ → `github/client.go` ✅
> Imports stay self-documenting: `domain2.Client`, `domain2.Loader`, `release.Entity`.

---

## Code Conventions

### One object per file

Each file defines exactly one primary object (struct or interface). The filename must match the role of that object, as described in the layer conventions above.

### Constructor

Every object must have a constructor named `New<ObjectName>`. Its parameters are exclusively the dependencies to be injected. **No logic, initialization, or side effects of any kind inside the constructor** — it only assigns fields.

```go
func NewKernel(embedFS embed.FS, specs *Specs, clio *clio.Clio) *Kernel {
    return &Kernel{
        embedFS,
        specs,
        clio,
    }
}
```

### Function naming

Function and method names must always be a **verb** or a **verbNoun** (`load`, `fetchReleases`, `parseResponse`). Names must be self-explanatory in isolation.

Context already provided by the receiver type must not be repeated in the method name:

```go
// receiver is ElementManager
func (m *ElementManager) getAll() { ... }      // ✅
func (m *ElementManager) getElements() { ... } // ❌ redundant

// receiver is ReleaseClient
func (c *ReleaseClient) fetch() { ... }        // ✅
func (c *ReleaseClient) fetchRelease() { ... } // ❌ redundant
```

---

## Lint

- Config file: `.golangci.yml` in v2 format
- `staticcheck` enabled with exclusions: `-ST1000` (package comments)

---

## Build

- All builds use vendor mode: `GOFLAGS="-mod=vendor"`
- Production builds require `VERSION` and `CHANNEL` env vars to be set
- Build artifacts are output to `build/`
- Use `make build-dev`, `make build-staging`, or `make build` (production)