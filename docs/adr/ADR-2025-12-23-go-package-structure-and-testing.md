# ADR: Go Package Structure and Testing Approach

**Date:** 2025-12-23  
**Status:** Accepted  
**Context:** First increment (tray-icon-and-dropdown-ui)

---

## Context

GoPomodoro is a minimal macOS native pomodoro timer application written in Go. During the first increment, we needed to establish:

1. How to structure Go packages for a menu bar application
2. What testing strategy fits the `lite` constitution mode
3. Which external libraries to use and how to depend on them

This is a single-developer tool focused on simplicity, small increments, and effectiveness over comprehensive testing infrastructure.

---

## Decision

### 1. Package Structure

We adopt a **standard Go project layout** with clear separation between application entry point and internal components:

```
gopomodoro/
├── cmd/gopomodoro/          # Application entry point
│   ├── main.go              # Wires components together
│   └── main_test.go         # Basic integration smoke tests
├── internal/                # Private application code
│   ├── tray/                # System tray icon management
│   │   ├── tray.go
│   │   └── tray_test.go
│   └── ui/                  # UI/menu management
│       ├── window.go
│       └── window_test.go
└── assets/                  # Static resources (icons, etc.)
```

**Key principles:**

- **`cmd/` for executables**: Contains only the main package and application wiring
- **`internal/` for components**: Private packages that cannot be imported by external projects
- **Component isolation**: Each component (`tray`, `ui`) is self-contained with its own tests
- **Flat hierarchy**: Avoid deep nesting; prefer simple, obvious structure

### 2. Testing Approach

We adopt a **pragmatic testing strategy**:

- **Unit tests for core logic**: Each `internal/` package has basic unit tests verifying key functions work
- **Manual testing for UI**: Visual and interaction testing done manually
- **Smoke test flag**: `--smoke` flag allows automated verification that the app initializes and shuts down cleanly
- **Test coverage focus**: Prioritize testing where bugs would hurt (initialization, state management when added), not trivial getters/setters
- **No test fixtures complexity**: Keep test data simple and inline where possible

**Testing principles:**

- Manual testing for UI and integration behavior
- Automated tests added when logic stabilizes (e.g., timer state machine in future increments)
- Fast, simple tests preferred over comprehensive coverage
- Tests run via `make test` without requiring CI infrastructure

**Test package naming:**

- **Tests live in `_test` packages**: Test files use the `package <name>_test` convention (e.g., `storage_test` for the `storage` package)
- **Black-box testing**: This enforces testing through the public API only, ensuring tests don't depend on internal implementation details
- **Exceptions allowed**: When testing internal helpers or private functions is necessary, tests may use the same package name
- **Example structure**:
  ```
  internal/storage/
    ├── storage.go          (package storage)
    └── storage_test.go     (package storage_test)
  internal/timer/
    ├── timer.go            (package timer)
    └── timer_test.go       (package timer_test)
  ```

### 3. Library Usage and Dependencies

We adopt a **minimal, pragmatic approach** to external dependencies:

**Core choice: `getlantern/systray`**

- **Why**: Mature, cross-platform Go library for system tray integration
- **How**: Used directly in `internal/tray` and `internal/ui` packages
- **Not wrapped**: Direct usage keeps the codebase simple; wrapping would add unnecessary abstraction
- **Alternative considered**: Native macOS bindings via cgo - rejected as more complex and less portable

**Dependency principles:**

- **Pin versions explicitly**: Use `go.mod` to lock dependency versions
- **Minimal dependencies**: Only add libraries that solve real problems
- **No framework lock-in**: Avoid heavy frameworks; prefer libraries that solve specific needs
- **Direct usage**: Wrapping third-party libraries in adapters adds complexity without clear benefit for this project

**Current dependencies:**

- `getlantern/systray` - System tray integration (required)
- Standard library - File I/O, logging, flags
- No web frameworks, no ORMs, no heavy abstractions

---

## Consequences

### Benefits

- **Clear structure**: New developers (or future-you) can quickly understand where code lives
- **Standard Go idioms**: Follows Go community conventions (`cmd/`, `internal/`)
- **Simple testing**: Tests are straightforward, no complex mocking or fixtures
- **Low ceremony**: Testing strategy matches project scale
- **Minimal dependencies**: Less to maintain, update, or debug
- **Fast builds**: Simple dependency graph keeps compile times low

### Drawbacks

- **Manual testing burden**: UI changes require manual verification
- **Direct library coupling**: Changes to `systray` API require updates throughout codebase
- **Limited cross-platform abstraction**: If we add Windows/Linux support, may need to introduce platform-specific packages

### Future Considerations

- **When to add abstraction**: If we need multiple UI backends or switch systray libraries, introduce adapter pattern then (not speculatively)
- **When to add CI**: If the project grows or gains users, GitHub Actions for automated testing would be next step
- **When to wrap libraries**: If a third-party library becomes problematic (frequent breaking changes, bugs), wrap it behind an interface

---

## Alternatives Considered

### Alternative 1: Wrap systray library behind interface

```go
type TrayManager interface {
    Initialize() error
    SetIcon(data []byte) error
    // ...
}
```

**Rejected because:**
- Adds complexity without clear benefit
- Systray API is stable and simple
- Can refactor later if needed (YAGNI principle)

### Alternative 2: Single `app/` package with all code

**Rejected because:**
- Violates separation of concerns
- Makes testing harder (can't test tray logic independently)
- Less maintainable as codebase grows

### Alternative 3: Comprehensive test coverage with mocks

**Rejected because:**
- Over-engineering for this project's needs
- Mocking `systray` calls adds complexity without proportional value
- Manual testing provides sufficient confidence for UI behavior

### Alternative 4: Native macOS bindings via cgo

**Rejected because:**
- More complex to maintain
- Ties us to macOS only (PRD mentions potential Linux support)
- Harder to build and test
- `systray` provides same functionality with simpler API

---

## References

- [CONSTITUTION.md](../../CONSTITUTION.md) - Project principles and values
- [PRD.md](../../PRD.md) - macOS primary, potential cross-platform future
- [increment.md](../../docs/increments/tray-icon-and-dropdown-ui/increment.md) - First increment scope
- Go project layout: https://github.com/golang-standards/project-layout
- getlantern/systray: https://github.com/getlantern/systray
