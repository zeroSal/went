## Project Structure

- `main.go` — Entry point, creates root `cobra.Command`
- `cmd/` — Cobra commands: one file per command (`init.go`, `root.go`)
- `app/kernel.go` — DI kernel definition (`fx.Module`)
- `app/bootstrap/initialization.go` — The initialization function
- `app/bootstrap/module/` — DI providers: one file per domain (`logger.go`, `example.go`)
- `app/config/env.go` — Env struct with `Load()` and `Validate()` methods
- `app/service/` — Internal services, organized by domain (see below)
- `app/model/` — Data models, organized by domain (see below)

---

### Layer conventions

Both `app/service/` and `app/model/` follow the same rule: **one package per domain**. The package name defines the domain context; each filename describes the type of object it contains, never repeating the domain name.

#### `app/service/`

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

#### `app/model/`

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

---

> **Why:** the domain lives in the package, the role lives in the filename.
> `github/github_client.go` ❌ → `github/client.go` ✅
> Imports stay self-documenting: `domain2.Client`, `domain2.Loader`, `release.Entity`.

---

## Tech Stack

- **Cobra** — CLI framework, command routing and flag parsing
- **Uber FX** — Dependency injection container
- **golangci-lint v2** — Linter aggregator

---

## Patterns

#### Interface + Implementation

Each service domain that requires abstraction (e.g. for mocking or DI) splits into two files:

- `service/x/x_interface.go` — defines the interface
- `service/x/x.go` — provides the concrete implementation

Always include a compile-time interface check in the implementation file:

```go
var _ InterfaceType = (*ImplementationType)(nil)
```

#### FX Modules

Providers in `bootstrap/module/` wrap concrete types into FX-injectable named types: `AppXxx` or `XxxClient`. This keeps the DI graph explicit and avoids ambiguity when multiple values share the same underlying type.

#### Client Pattern

External service clients live in `service/<domain>/client.go` and are named `Client` within the package. Consumed as `domain.Client` at the call site.

#### Env Config

`config/env.go` defines a single typed struct. It exposes a `Load()` method (reads from environment) and validation helpers. No raw `os.Getenv` calls outside this file.

#### Cobra Commands

Each command follows the `XxxCmd` struct pattern:

```go
type XxxCmd struct { ... }

func NewXxxCmd(...) *XxxCmd { ... }
func (c *XxxCmd) Command() *cobra.Command { ... }
func (c *XxxCmd) run(...) error { ... }
func (c *XxxCmd) execute(...) error { ... }
```

---

## Lint

- Config file: `.golangci.yml` in v2 format
- `staticcheck` enabled with exclusions: `-ST1005` (error string casing), `-ST1000` (package comments)

---

## Code Conventions

#### One object per file

Each file defines exactly one primary object (struct or interface). The filename must match the role of that object, as described in the layer conventions above.

#### Constructor

Every object must have a constructor named `New<ObjectName>`. Its parameters are exclusively the dependencies to be injected. **No logic, initialization, or side effects of any kind inside the constructor** — it only assigns fields.

```go
func NewClient(logger *zap.Logger, cfg config.Env) *Client {
    return &Client{
        logger: logger,
        cfg:    cfg,
    }
}
```

#### Function naming

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

## Build

- All builds use vendor mode: `GOFLAGS="-mod=vendor"`
- Production builds require `VERSION` and `CHANNEL` env vars to be set
- Build artifacts are output to `build/`