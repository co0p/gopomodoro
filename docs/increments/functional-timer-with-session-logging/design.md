# Design: Functional Timer with Session Logging

## Context and Problem

This design implements the core timer functionality for GoPomodoro, transforming the placeholder UI from the previous increment into a working 25-minute pomodoro timer with session tracking.

### Current State

The codebase has:
- A functional tray icon and systray-based dropdown menu ([tray-icon-and-dropdown-ui increment](../tray-icon-and-dropdown-ui/increment.md))
- Placeholder UI elements (header, timer display, Start/Pause/Reset buttons) that are disabled
- No timer logic, state management, or persistence

### Problem

Users need to actually run a pomodoro session. This requires:
1. A countdown timer that runs for 25 minutes and updates every second
2. Controls to start, pause, resume, and reset the timer
3. Automatic session completion detection when timer reaches 00:00
4. Persistent session logging to `~/.gopomodoro/sessions.log` in CSV format
5. Fail-fast error handling if session data cannot be written

### References

- [increment.md](increment.md) - Full acceptance criteria and use cases
- [CONSTITUTION.md](../../../CONSTITUTION.md) - Lite mode principles
- [PRD.md](../../../PRD.md) - CSV format specification and full product vision
- [ADR-2025-12-23-go-package-structure-and-testing.md](../../adr/ADR-2025-12-23-go-package-structure-and-testing.md) - Package structure conventions

---

## Proposed Solution (Technical Overview)

### Architecture

We introduce two new internal packages that remain decoupled:

**`internal/timer/`** - Pure countdown timer logic
- Manages countdown state (idle → running → paused → completed)
- Counts down from 1500 seconds (25 minutes) to 0
- Emits events when state changes (Started, Tick, Completed)
- No knowledge of UI or storage

**`internal/storage/`** - Session persistence
- Creates and validates `~/.gopomodoro/` directory
- Appends session records to `sessions.log` in CSV format
- Fails fast with clear errors if filesystem operations fail
- No knowledge of timer or UI

**`internal/ui/`** - Orchestration (extended from previous increment)
- Listens to timer events and updates display
- Handles button clicks by calling timer methods
- Calls storage when timer events indicate logging is needed
- Manages button enable/disable state based on timer state

**`cmd/gopomodoro/main.go`** - Wiring
- Creates timer, storage, and UI components
- Connects event handlers between components
- Handles fail-fast errors during initialization

### Flow

**Starting a session:**
1. User clicks Start → UI calls `timer.Start()`
2. Timer transitions from `idle` to `running`, emits `Started` event
3. UI receives `Started` event → calls `storage.LogSession()` with "started" event
4. Timer begins ticking every second, emitting `Tick` events with remaining time
5. UI receives `Tick` events → updates display to show countdown (e.g., "24:59", "24:58")

**Completing a session:**
1. Timer counts down to 0 seconds
2. Timer emits `Completed` event, transitions to `idle`
3. UI receives `Completed` event → calls `storage.LogSession()` with "completed" event and 25-minute duration
4. UI updates display back to "25:00" and re-enables Start button

**Pausing and resuming:**
1. User clicks Pause → UI calls `timer.Pause()`
2. Timer transitions to `paused`, stops emitting `Tick` events
3. UI disables Pause button, enables Resume button
4. User clicks Resume → UI calls `timer.Resume()`
5. Timer transitions back to `running`, resumes ticking from paused time

**Resetting:**
1. User clicks Reset → UI calls `timer.Reset()`
2. Timer stops, transitions to `idle`, remaining time set back to 1500 seconds
3. No logging occurs (session was not completed)
4. UI updates display to "25:00", re-enables Start button

---

## Scope and Non-Scope (Technical)

### In Scope

This design covers:
- Timer countdown engine (25 minutes hardcoded)
- State machine with three states: idle, running, paused
- Start, Pause, Resume, Reset commands
- Event emission for Started, Tick, Completed
- CSV session logging with two event types: "started" and "completed"
- `~/.gopomodoro/` directory creation with fail-fast error handling
- Button state management (enable/disable based on timer state)
- Display updates every second during countdown

### Explicitly Out of Scope

Deferred to future increments:
- Break timers (short 5-minute, long 15-minute)
- Pomodoro cycle tracking (4 sessions before long break)
- Notifications when sessions complete
- Statistics calculation or display
- Configurable timer durations
- Tray icon visual state changes during active sessions
- Sound effects
- Reading from sessions.log

### How This Fits the Roadmap

This increment establishes:
- The timer engine that will be extended for break timers
- The session logging foundation needed for statistics
- The event-driven architecture that notifications will hook into

---

## Architecture and Boundaries

### Component Diagram

```
┌─────────────────────────────────────────────────┐
│ cmd/gopomodoro/main.go                          │
│  - Initializes all components                   │
│  - Wires event handlers                         │
│  - Calls storage.EnsureDataDir() at startup     │
└──────────┬──────────────────────────────────────┘
           │
           ├──────────────────┬──────────────────────┐
           │                  │                      │
           ▼                  ▼                      ▼
┌────────────────────┐ ┌─────────────────┐ ┌──────────────────┐
│ internal/timer/    │ │ internal/ui/    │ │ internal/storage/│
│                    │ │                 │ │                  │
│ - State machine    │ │ - Menu items    │ │ - EnsureDataDir()│
│ - Countdown logic  │ │ - Event handlers│ │ - LogSession()   │
│ - Event callbacks  │ │ - Display update│ │ - CSV formatting │
└────────────────────┘ └─────────────────┘ └──────────────────┘
         │                     │
         │  Events             │  Calls storage
         │  (Started, Tick,    │  when logging needed
         │   Completed)        │
         └─────────────────────┘
```

### Package Responsibilities

**`internal/timer/`**
- Maintain timer state (idle, running, paused)
- Track remaining seconds (starts at 1500)
- Tick every second when running
- Detect completion (remaining == 0)
- Provide callback registration for events
- Expose query methods: `GetState()`, `GetRemaining()`

**`internal/storage/`**
- Ensure `~/.gopomodoro/` directory exists
- Append CSV records to `sessions.log`
- Format: `timestamp,session_type,event,duration_minutes`
- Return errors immediately (no retry logic)

**`internal/ui/`** (extends existing)
- Subscribe to timer events via callbacks
- Update timer display menu item on Tick events
- Call `storage.LogSession()` on Started and Completed events
- Handle button click events from systray menu items
- Enable/disable buttons based on timer state

**`cmd/gopomodoro/main.go`**
- Create timer, storage, and UI instances
- Call `storage.EnsureDataDir()` during startup
- Wire UI button handlers to timer methods
- Wire timer event handlers to UI updates and storage calls
- Exit with clear error if storage initialization fails

### Dependency Rules

Following the existing ADR and constitution principles:

- **Timer** has zero dependencies on UI or storage
- **Storage** has zero dependencies on timer or UI
- **UI** depends on timer (calls methods, subscribes to events) and storage (calls logging)
- **Main** depends on all three (orchestration layer)
- No circular dependencies
- Direct use of `systray` library continues (no wrapping)

### Guardrails Respected

From [CONSTITUTION.md](../../../CONSTITUTION.md):
- Small, safe steps: timer logic separate from persistence
- Simple is better than complex: straightforward state machine, no event bus
- Make it work: hardcoded 25-minute duration first

From [ADR-2025-12-23-go-package-structure-and-testing.md](../../adr/ADR-2025-12-23-go-package-structure-and-testing.md):
- `internal/` packages for new components
- Flat hierarchy: `internal/timer/`, `internal/storage/`
- Minimal dependencies: only standard library beyond existing systray

---

## Contracts and Data

### Timer Package API

```go
package timer

import "time"

// State represents the timer's current state
type State int

const (
    StateIdle State = iota
    StateRunning
    StatePaused
)

// Timer manages a countdown timer
type Timer struct {
    // unexported fields
}

// New creates a timer initialized to 25 minutes (1500 seconds)
func New() *Timer

// Start begins countdown from current remaining time
// Transitions: idle → running
// Emits: Started event
func (t *Timer) Start()

// Pause stops the countdown, preserving remaining time
// Transitions: running → paused
func (t *Timer) Pause()

// Resume continues countdown from paused time
// Transitions: paused → running
func (t *Timer) Resume()

// Reset stops timer and sets remaining time back to 1500 seconds
// Transitions: any → idle
func (t *Timer) Reset()

// GetState returns current state
func (t *Timer) GetState() State

// GetRemaining returns remaining seconds
func (t *Timer) GetRemaining() int

// Event handler registration
func (t *Timer) OnStarted(handler func())
func (t *Timer) OnTick(handler func(remainingSeconds int))
func (t *Timer) OnCompleted(handler func())
```

**Behavioral contracts:**
- Timer ticks every 1 second when in `StateRunning`
- `OnTick` handler called with decremented remaining time
- When remaining reaches 0, `OnCompleted` is called and state → `StateIdle`
- `Start()` called when already running is a no-op
- `Pause()` called when not running is a no-op
- `Reset()` can be called from any state

### Storage Package API

```go
package storage

import "time"

// EnsureDataDir creates ~/.gopomodoro/ if it doesn't exist
// Returns error if directory cannot be created or is not writable
// Should be called once during app initialization
func EnsureDataDir() error

// LogSession appends a session record to sessions.log
// CSV format: timestamp,session_type,event,duration_minutes
// timestamp: ISO 8601 format (e.g., "2025-12-23T14:30:00Z")
// sessionType: "work" (breaks in future increments)
// event: "started" or "completed"
// durationMinutes: 0 for started, 25 for completed work sessions
func LogSession(timestamp time.Time, sessionType, event string, durationMinutes int) error
```

**Behavioral contracts:**
- `EnsureDataDir()` is idempotent (safe to call multiple times)
- Directory path: `$HOME/.gopomodoro/`
- File path: `$HOME/.gopomodoro/sessions.log`
- Appends to existing file or creates new file with header if missing
- No retry logic - errors returned immediately
- Timestamps formatted in UTC with ISO 8601 format

### CSV Format

**File**: `~/.gopomodoro/sessions.log`

**Format**:
```csv
timestamp,session_type,event,duration_minutes
2025-12-23T14:30:00Z,work,started,0
2025-12-23T14:55:00Z,work,completed,25
2025-12-23T15:00:00Z,work,started,0
2025-12-23T15:25:00Z,work,completed,25
```

**Field specifications:**
- `timestamp`: RFC3339 format (ISO 8601 compatible), UTC timezone
- `session_type`: `"work"` for this increment (future: `"short_break"`, `"long_break"`)
- `event`: `"started"` when timer begins, `"completed"` when timer reaches 00:00
- `duration_minutes`: `0` for started events, `25` for completed work sessions

**Notes:**
- CSV header included if file is newly created
- Reset or skipped sessions are NOT logged
- File is append-only (no edits or deletions by app)
- Human-readable for manual inspection and debugging

### UI Updates

**Menu item state changes:**

| Timer State | Display   | Start Button | Pause Button | Reset Button |
|-------------|-----------|--------------|--------------|--------------|
| Idle        | "25:00"   | Enabled      | Disabled     | Disabled     |
| Running     | "24:59"   | Disabled     | Enabled      | Enabled      |
| Paused      | "18:42"   | Enabled*     | Disabled     | Enabled      |

*Resume functionality (clicking Start when paused calls `timer.Resume()`)

**Display format**: `MM:SS` (e.g., "25:00", "24:59", "00:03", "00:00")

---

## Testing and Safety Net

### Unit Tests

**`internal/timer/timer_test.go`:**
- Test state transitions:
  - `TestStart()`: idle → running, Started event fires
  - `TestPause()`: running → paused, Tick events stop
  - `TestResume()`: paused → running, Tick events resume
  - `TestReset()`: any state → idle, remaining = 1500
  - `TestCompletion()`: running with remaining=1 → Tick → idle, Completed event fires
- Test timer countdown:
  - Verify remaining decrements each tick
  - Verify tick interval is approximately 1 second (allow small drift)
- Test event callbacks:
  - Verify OnStarted called when Start() invoked
  - Verify OnTick called each second with correct remaining value
  - Verify OnCompleted called when timer reaches 0

**`internal/storage/storage_test.go`:**
- `TestEnsureDataDir()`: Directory created if missing, no error if exists
- `TestLogSession()`: CSV record appended with correct format
- `TestLogSessionCreatesFile()`: File created with header if missing
- `TestEnsureDataDirFailure()`: Error returned if directory cannot be created (mock filesystem constraint)
- `TestLogSessionFailure()`: Error returned if file not writable (mock filesystem constraint)

**`internal/ui/window_test.go`** (extended from existing):
- Test button state logic (not systray interaction):
  - Helper function `shouldButtonBeEnabled(timerState State) bool` tests

### Manual Testing

Per lite constitution, UI behavior verified manually:

- Launch app, click Start, observe timer counts down from 25:00
- Verify display updates every second (24:59, 24:58, ...)
- Click Pause mid-session, verify countdown stops
- Click Start again (Resume), verify countdown continues from paused time
- Click Reset, verify timer returns to 25:00 and idle state
- Let timer run to completion (00:00), verify:
  - Timer returns to idle state automatically
  - `~/.gopomodoro/sessions.log` contains two records (started, completed)
- Delete `~/.gopomodoro/` directory, launch app, verify it's recreated
- Make `~/.gopomodoro/` read-only, launch app, verify error message and exit

### Test Data and Fixtures

- Timer tests use time mocking or fast-forward techniques (tick every 10ms for speed)
- Storage tests use `t.TempDir()` for isolated filesystem operations
- No external test fixtures required

---

## CI/CD and Rollout

### CI

No CI required per lite constitution. Tests run locally via `make test`.

### Build and Run

Existing Makefile targets continue to work:
```bash
make build  # Compiles to bin/gopomodoro
make run    # Builds and runs
make test   # Runs all unit tests
```

No changes to build process needed.

### Rollout Plan

**First-time users:**
1. User builds and runs app
2. App creates `~/.gopomodoro/` directory automatically
3. User starts first session, `sessions.log` is created

**Existing users (from previous increment):**
1. User rebuilds app
2. App creates `~/.gopomodoro/` directory on first session start
3. Previously non-functional buttons now work

**No migration needed** - this is the first increment with persistence.

### Rollback Plan

If bugs are found post-deployment:
- Revert commits to previous increment
- User can manually delete `~/.gopomodoro/sessions.log` to clear session history
- No schema migrations or data compatibility concerns

### Feature Flags

None needed - this is new functionality, not a replacement of existing behavior.

---

## Observability and Operations

### Logging

Using Go's standard `log` package (existing pattern from previous increment):

**Startup:**
```
[INFO] GoPomodoro starting...
[INFO] Initializing tray icon...
[INFO] Tray icon initialized successfully
[INFO] Dropdown window created
[INFO] Data directory ensured: /Users/username/.gopomodoro
```

**Timer lifecycle:**
```
[INFO] Timer started (25:00)
[INFO] Timer paused (18:42 remaining)
[INFO] Timer resumed (18:42)
[INFO] Timer reset to idle
[INFO] Session completed (25:00 → 00:00)
```

**Session logging:**
```
[INFO] Logged session: started at 2025-12-23T14:30:00Z
[INFO] Logged session: completed at 2025-12-23T14:55:00Z, duration 25 minutes
```

**Error conditions:**
```
[ERROR] Failed to create data directory: permission denied
[ERROR] Failed to log session: disk full
```

**Error handling behavior:**
- If `EnsureDataDir()` fails during startup: log error, exit immediately
- If `LogSession()` fails during session: log error, exit immediately (fail-fast per increment AC)

### Metrics

None required per lite constitution. Session data accumulates in CSV for future statistics increment.

### Operational Considerations

**Storage requirements:**
- Minimal: ~100 bytes per session pair (started + completed)
- At 10 sessions/day: ~1KB/day, ~365KB/year

**Performance:**
- Timer ticks every second (negligible CPU)
- CSV append is O(1) file operation (minimal I/O)
- No performance concerns for single-user desktop app

**Known operational risks:**
- User's disk full → fail-fast with clear error
- User deletes `~/.gopomodoro/` while app running → next log attempt will fail, app exits
  - Acceptable for lite mode (user can restart app)

### Monitoring and Alerts

None required for lite mode. User observes:
- Timer countdown in UI (visual confirmation)
- Session records in `~/.gopomodoro/sessions.log` (manual inspection)

---

## Risks, Trade-offs, and Alternatives

### Known Risks

**Timer drift:**
- Over 25 minutes, small drift may accumulate due to tick interval precision
- Mitigation: Use `time.Ticker` (Go standard library handles drift correction)
- Acceptable: Drift would be <1 second over 25 minutes

**CSV file corruption:**
- If app crashes mid-write, partial line may be written
- Mitigation: CSV append is atomic at OS level for small writes
- Recovery: User can manually edit/delete corrupted line
- Acceptable for lite mode

**Storage failure during session:**
- App exits immediately if logging fails (fail-fast per AC)
- User loses current session progress
- Mitigation: Future increments could add "retry" or "degraded mode"
- Acceptable: Fail-fast prevents silent data loss, builds trust

### Trade-offs

**Event-driven architecture vs. polling:**
- Chosen: Event callbacks (OnStarted, OnTick, OnCompleted)
- Alternative: UI polls timer every second for state/remaining
- Trade-off: Events add slight complexity but cleaner separation, easier testing
- Justification: Events make timer completely UI-agnostic, better for future features (notifications)

**Fail-fast vs. degraded mode:**
- Chosen: Exit on storage errors (per increment AC 11)
- Alternative: Continue running without logging
- Trade-off: Less graceful UX, but prevents silent data loss
- Justification: User trust is critical for focus tool, better to fail loudly

**Started event logging:**
- Chosen: Log both "started" and "completed" events
- Alternative: Only log completed sessions
- Trade-off: More CSV records, but complete audit trail
- Justification: Enables future features (interrupted session tracking, accurate session duration calculation)

### Alternatives Considered

**Alternative 1: JSON logging instead of CSV**
- Pro: Easier to extend with nested data
- Con: Less human-readable, PRD specifies CSV
- Rejected: CSV meets current needs, can migrate later if needed

**Alternative 2: SQLite database**
- Pro: Easier querying for statistics
- Con: Adds dependency, more complex setup
- Rejected: Violates "simple is better than complex" principle, CSV sufficient for now

**Alternative 3: Timer as goroutine with channel communication**
- Pro: "More idiomatic Go"
- Con: Harder to test, more complex lifecycle management
- Rejected: Callback pattern simpler, testable, sufficient for single-threaded UI app

**Alternative 4: Combined timer+storage package**
- Pro: Fewer packages
- Con: Tight coupling, harder to test independently
- Rejected: Violates separation of concerns, makes storage logic untestable without timer

---

## Follow-up Work

### Future Increments

**Short-term (next 1-2 increments):**
- Break timers (5-minute short break, 15-minute long break)
  - Extend `timer.New()` to accept duration and session type
  - Add break-related buttons to UI
- Notifications when sessions complete
  - Hook into `OnCompleted` event
  - Use macOS notification APIs

**Medium-term (next 3-5 increments):**
- Pomodoro cycle tracking (4 work sessions → long break)
  - Add cycle state to UI
  - Automatic transition from work → break → work
- Statistics display
  - Read and parse `sessions.log`
  - Compute daily/weekly totals, streaks
- Configurable durations
  - Read from `settings.json` (per PRD)
  - Extend timer to accept custom durations

**Long-term (future):**
- Tray icon state changes during active sessions (color-coded per PRD)
- Sound effects on completion
- Pause/resume persistence across app restarts

### Tech Debt

None anticipated - this is a fresh implementation following established patterns.

### Post-Rollout Tasks

After implementation:
1. Manual testing checklist (see Testing and Safety Net section)
2. Verify `sessions.log` format matches PRD exactly
3. Test fail-fast behavior with read-only directory
4. Run for a full 25-minute session to verify no drift

---

## References

- [increment.md](increment.md) - Acceptance criteria, use cases, tasks
- [CONSTITUTION.md](../../../CONSTITUTION.md) - Lite mode, small safe steps
- [PRD.md](../../../PRD.md) - CSV format, full product vision
- [ADR-2025-12-23-go-package-structure-and-testing.md](../../adr/ADR-2025-12-23-go-package-structure-and-testing.md) - Package structure
- [tray-icon-and-dropdown-ui increment](../tray-icon-and-dropdown-ui/increment.md) - Previous work
- Go time package: https://pkg.go.dev/time
- CSV RFC: https://datatracker.ietf.org/doc/html/rfc4180
