# CONSTITUTION.md

## Architectural Decisions

### Layering (Hexagonal Architecture)

- Core domain in `internal/domain/` — pure timer logic, no infrastructure dependencies
- Ports (interfaces) in `internal/ports/` — define boundaries (TimerDriver, Clock)
- Adapters in `internal/adapters/` — tray UI, real clock, test driver
- Entry point in `cmd/gopomodoro/`
- Domain MUST NOT import adapters; adapters depend on ports

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

- Clock abstracted via port interface (enables testing)
- UI abstracted via driver interface (enables test driver)
- No direct third-party imports in domain layer

## Testing Expectations

- Test location: Colocated (`*_test.go` next to implementation)
- Coverage: State transitions must be tested; startup must be tested
- Test driver: Implements same interface as UI, exercises domain through ports
- Mocking: Clock mocked for deterministic tests; real clock only in adapters

## Artifact Layout

- **CONSTITUTION.md**: Project root
- **ADRs**: `docs/adr/0001-title.md` (sequential numbering)
- **Other docs**: `docs/`
- **Working context**: `.4dc/current/` (gitignored)

## Delivery Practices

- Build: `make build` — compile binary
- Test: `make test` — run all tests
- Install: `make install` — copy binary to `/usr/local/bin`
