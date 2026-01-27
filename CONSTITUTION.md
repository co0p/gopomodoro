# CONSTITUTION.md

## Architectural Decisions

### Package Layout

- `pkg/` — Core domain: Cycle, Ticker interface, state machine
- `pkg/tray/` — System tray implementation (getlantern/systray)
- `pkg/ticker/` — Real ticker implementation (time.Ticker)
- `cmd/gopomodoro/` — Entry point, wires dependencies

### Error Handling

- Standard Go error handling: return `error` as last value
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- No panics for expected failures

### State Management

- Timer state machine lives in domain layer
- States: Idle, Running (Pomodoro/ShortBreak/LongBreak), Paused
- State transitions are the core domain logic
- Cycle counting (4 pomodoros → long break) is domain responsibility

### Dependencies

**Direction:**
- Domain layer (`pkg/`) has zero imports
- Domain defines interfaces (Ticker, CycleObserver)
- Adapters (`pkg/tray/`, `pkg/ticker/`) import and implement domain interfaces
- Main (`cmd/`) wires concrete implementations together
- This prevents cyclic dependencies through dependency inversion

**Abstractions:**
- Clock abstracted via port interface (enables testing)
- UI abstracted via driver interface (enables test driver)
- No direct third-party imports in domain layer

### Construction Patterns

- Prefer public fields with direct struct construction: `Cycle{Ticker: t, Observer: o}`
- Avoid constructors and setter methods unless validation is needed
- Wiring happens in `main()`, not in constructors
- This is idiomatic Go: simple, explicit, testable

## Testing Expectations

- Test location: Colocated (`*_test.go` next to implementation)
- Test package: Black-box testing via `package <name>_test` (tests only access exported API)
- Coverage: State transitions must be tested; startup must be tested
- Test driver: Implements same interface as UI, exercises domain through ports
- Mocking: Clock mocked for deterministic tests; real clock only in adapters

### Test Helpers

- Mock implementations: `pkg/testing/` subdirectory
- Keeps test files focused on behavior, not mock boilerplate
- Mocks implement domain interfaces (e.g., MockTicker, MockObserver)

## Artifact Layout

- **CONSTITUTION.md**: Project root
- **ADRs**: `docs/adr/0001-title.md` (sequential numbering)
- **Other docs**: `docs/`
- **Working context**: `.4dc/current/` (gitignored)

## Delivery Practices

- Build: `make build` — compile binary
- Test: `make test` — run all tests
- Install: `make install` — copy binary to `/usr/local/bin`
