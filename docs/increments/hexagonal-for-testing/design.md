# Design: Hexagonal Architecture for Testing

## Context and Problem

### Increment Goal

Enable automated testing of complete pomodoro cycles without requiring UI interaction or manual button clicks. This refactoring applies hexagonal (ports and adapters) architecture to separate business logic from infrastructure, making the application fully testable and eliminating the deadlock issues caused by circular dependencies between timer callbacks and UI updates.

### Why This Change Now

After implementing the complete pomodoro cycle with UI, we've encountered critical issues:
1. **Deadlocks**: Timer callbacks re-entering timer methods while holding locks causes UI freeze
2. **Untestable business logic**: Cannot test full pomodoro cycles without starting the systray GUI
3. **Tight coupling**: Business logic intertwined with systray implementation details
4. **Circular dependencies**: UI → Timer → UI callback cycles create fragile code

The current architecture makes it impossible to verify that a 4-session pomodoro cycle works correctly through automated tests.

### Current System Behavior

**Current Architecture (Callback-driven):**
```
┌─────────────────────────────────────────────┐
│  cmd/gopomodoro/main.go                    │
│  ├─ Creates: Timer, Session, Window, Tray  │
│  └─ Wires: window.SetTimer(timer)          │
└─────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────┐
│  internal/ui/window.go (Orchestrator)       │
│  ├─ Owns: Timer, Session, Tray             │
│  ├─ Registers callbacks on Timer           │
│  └─ Updates: systray menu items            │
└─────────────────────────────────────────────┘
                    │
      ┌─────────────┴─────────────┐
      ▼                           ▼
┌─────────────────┐    ┌─────────────────────┐
│ internal/timer  │    │ internal/session    │
│ - OnStarted()   │    │ - DetermineNext()   │
│ - OnTick()      │◄───│ - IncrementCycle()  │
│ - OnCompleted() │    │ - CurrentType       │
└─────────────────┘    └─────────────────────┘
      │
      └─► Callbacks to UI (DEADLOCK RISK)
```

**Problems with Current Design:**

1. **Inverted Control Flow**: Timer calls UI, UI calls Timer → circular dependency
2. **Lock Contention**: 
   - `timer.Start()` holds lock, calls `onStarted()`
   - `onStarted()` → `UI.handleTimerStarted()` → `UI.UpdateButtonStates()`
   - `UpdateButtonStates()` → `GetRemaining()` tries to acquire same lock
   - **Result**: DEADLOCK
3. **Untestable**: 
   - Business logic buried in UI event handlers
   - Cannot drive app without systray running
   - Tests would need to simulate button clicks on actual menu items
4. **Mixed Responsibilities**:
   - Window handles UI rendering AND business logic orchestration
   - Timer manages time AND notifies about state changes
   - No clear separation between "what" (business rules) and "how" (UI/infrastructure)

### Links

- [increment.md](increment.md) - Product requirements and acceptance criteria
- [../../CONSTITUTION.md](../../CONSTITUTION.md) - Project principles (mode: lite)
- [../../ARCHITECTURE.md](../../ARCHITECTURE.md) - Current architecture
- [../../PRD.md](../../PRD.md) - Original product requirements

## Proposed Solution (Technical Overview)

### High-Level Approach

Apply **Hexagonal Architecture** (Ports and Adapters) to invert dependencies: create a central `PomodoroService` that owns the business logic and exposes command methods (inbound ports) that any driver (UI, tests, CLI) can call. The service publishes events through interfaces (outbound ports) that adapters implement.

**Key Principles:**
- **Dependency Inversion**: Business logic depends on abstractions (interfaces), not concrete implementations
- **No Callbacks to Drivers**: Service never calls back into UI/tests - only emits events to interfaces
- **Testable Core**: Service can be tested with mock adapters, no GUI required
- **Single Responsibility**: Each layer has one clear job

### Target Architecture (Hexagonal)

```
┌────────────────────────────────────────────────────────┐
│              PRIMARY ADAPTERS (Drivers)                │
│                                                        │
│  ┌──────────────────────┐    ┌────────────────────┐  │
│  │   SystrayAdapter     │    │   TestDriver       │  │
│  │  (UI Implementation) │    │  (Test Harness)    │  │
│  │                      │    │                    │  │
│  │  - Button handlers   │    │  - Programmatic    │  │
│  │    call commands     │    │    commands        │  │
│  │  - Implements        │    │  - Records events  │  │
│  │    Notifier          │    │  - Asserts state   │  │
│  └──────────────────────┘    └────────────────────┘  │
└────────────────────────────────────────────────────────┘
              │                           │
              └───────────┬───────────────┘
                          │
                          ▼
         ┌────────────────────────────────────┐
         │      INBOUND PORTS (Commands)      │
         │                                    │
         │  • StartSession()                  │
         │  • PauseSession()                  │
         │  • ResumeSession()                 │
         │  • SkipSession()                   │
         │  • ResetCycle()                    │
         │  • GetState() State                │
         └────────────────────────────────────┘
                          │
                          ▼
         ┌────────────────────────────────────┐
         │    CORE DOMAIN (Business Logic)    │
         │                                    │
         │         PomodoroService            │
         │  ┌──────────────────────────────┐ │
         │  │ - session: Session           │ │
         │  │ - timer: Timer               │ │
         │  │ - state machine logic        │ │
         │  │ - cycle progression          │ │
         │  │ - duration selection         │ │
         │  └──────────────────────────────┘ │
         │                                    │
         │  Business Rules:                   │
         │  • 4 work sessions → long break    │
         │  • Work → short break → work       │
         │  • Skip increments cycle counter   │
         │  • Pause only when running         │
         └────────────────────────────────────┘
                          │
                          ▼
         ┌────────────────────────────────────┐
         │     OUTBOUND PORTS (Events)        │
         │                                    │
         │  Notifier:                         │
         │    • SessionStarted(type, duration)│
         │    • SessionTick(remaining)        │
         │    • SessionCompleted(type)        │
         │    • StateChanged(state)           │
         │                                    │
         │  Clock:                            │
         │    • Ticker(duration) Ticker       │
         │                                    │
         │  Storage:                          │
         │    • LogSession(...)               │
         └────────────────────────────────────┘
                          │
                          ▼
┌────────────────────────────────────────────────────────┐
│           SECONDARY ADAPTERS (Driven)                  │
│                                                        │
│  ┌──────────────┐  ┌────────────┐  ┌──────────────┐  │
│  │   Systray    │  │   Clock    │  │    File      │  │
│  │   Notifier   │  │  (real/    │  │   Storage    │  │
│  │              │  │   mock)    │  │              │  │
│  │ - Updates    │  │ - time.    │  │ - CSV logs   │  │
│  │   menu bar   │  │   Ticker   │  │              │  │
│  │ - Updates    │  │ - benbjohn │  │              │  │
│  │   progress   │  │   son/clock│  │              │  │
│  └──────────────┘  └────────────┘  └──────────────┘  │
└────────────────────────────────────────────────────────┘
```

### Component Changes

#### New: Core Domain Package (`internal/pomodoro/`)

**`service.go`** - The heart of the application
```go
type Service struct {
    session  *session.Session
    timer    *timer.Timer
    notifier Notifier
    storage  Storage
    state    State
}

type State int
const (
    StateIdle State = iota
    StateRunning
    StatePaused
)

// Inbound port: Commands that drivers call
func (s *Service) StartSession() error
func (s *Service) PauseSession() error
func (s *Service) ResumeSession() error
func (s *Service) SkipSession() error
func (s *Service) ResetCycle() error
func (s *Service) GetState() State
func (s *Service) GetRemainingSeconds() int
func (s *Service) GetCurrentSessionType() string
func (s *Service) GetCompletedSessions() int
```

**Key Behaviors:**
- `StartSession()`: Determines duration from session type, starts timer, logs "started"
- `PauseSession()`: Only works if running, logs pause (future)
- `SkipSession()`: Stops timer, logs "skipped" with elapsed time, advances cycle
- `ResetCycle()`: Stops timer, resets session counter to 0, returns to work session 1
- **State Machine**: Prevents invalid transitions (e.g., pause when idle)
- **Cycle Logic**: Determines next session based on completed counter
- **Event Publishing**: Calls `notifier.SessionStarted()` etc. after state changes

**`ports.go`** - Interface definitions
```go
// Outbound port: Service publishes events to this
type Notifier interface {
    SessionStarted(sessionType string, duration int)
    SessionTick(remainingSeconds int)
    SessionCompleted(sessionType string)
    StateChanged(state State)
}

// Outbound port: Service logs to this
type Storage interface {
    LogSession(timestamp time.Time, sessionType, status string, duration int) error
}

// Outbound port: Service uses this for time
type Clock interface {
    Ticker(duration time.Duration) *clock.Ticker
}
```

#### Modified: Timer Package (`internal/timer/`)

**Changes:**
- Already uses `Clock` interface ✓ (from mock clock refactoring)
- Keep callbacks for now, but service will register them (not UI)
- No structural changes needed - timer is already dependency-injected

**Integration with Service:**
```go
// Service registers itself as timer callbacks
timer.OnStarted(s.handleTimerStarted)
timer.OnTick(s.handleTimerTick)
timer.OnCompleted(s.handleTimerCompleted)

// Service methods called by timer (internal, not exposed to drivers)
func (s *Service) handleTimerStarted(sessionType string, duration int) {
    s.state = StateRunning
    s.notifier.SessionStarted(sessionType, duration)
    s.storage.LogSession(time.Now(), sessionType, "started", 0)
}
```

#### Modified: Session Package (`internal/session/`)

**No changes needed** - Already pure domain logic with no dependencies

#### New: Primary Adapter - Systray UI (`internal/adapters/ui/`)

**`systray_adapter.go`** - Implements Notifier, drives Service
```go
type SystrayAdapter struct {
    service      *pomodoro.Service
    // systray menu items
    progressBar    *systray.MenuItem
    cycleIndicator *systray.MenuItem
    btnStart       *systray.MenuItem
    btnPause       *systray.MenuItem
    // etc.
}

// Implements Notifier interface (receives events from service)
func (a *SystrayAdapter) SessionStarted(sessionType string, duration int) {
    // Update UI elements
    a.cycleIndicator.SetTitle(...)
    a.progressBar.SetTitle("○○○○○○○○○○")
}

func (a *SystrayAdapter) SessionTick(remaining int) {
    a.progressBar.SetTitle(formatProgressBar(...))
}

// Button handlers call service commands
func (a *SystrayAdapter) handleStartClick() {
    if err := a.service.StartSession(); err != nil {
        log.Printf("Start failed: %v", err)
    }
}

func (a *SystrayAdapter) handlePauseClick() {
    if err := a.service.PauseSession(); err != nil {
        log.Printf("Pause failed: %v", err)
    }
}
```

**Key Points:**
- **No direct timer access** - Only calls service methods
- **No business logic** - Pure UI presentation
- **No deadlocks** - Service never calls back into UI synchronously

#### New: Primary Adapter - Test Driver (`internal/adapters/test/`)

**`test_driver.go`** - For automated testing
```go
type TestDriver struct {
    service *pomodoro.Service
    events  []Event
    mu      sync.Mutex
}

type Event struct {
    Type      EventType
    Timestamp time.Time
    Data      interface{}
}

// Implements Notifier
func (d *TestDriver) SessionStarted(sessionType string, duration int) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.events = append(d.events, Event{
        Type: EventSessionStarted,
        Data: SessionStartedData{sessionType, duration},
    })
}

// Test helpers
func (d *TestDriver) GetEvents() []Event
func (d *TestDriver) ClearEvents()
func (d *TestDriver) WaitForEvent(eventType EventType, timeout time.Duration) error
func (d *TestDriver) AssertState(expected pomodoro.State) error
```

#### Modified: Main Package (`cmd/gopomodoro/`)

**`main.go`** - Dependency injection and wiring
```go
func main() {
    // Create infrastructure adapters (secondary/driven)
    mockClock := clock.New()
    fileStorage := storage.NewFileStorage()
    
    // Create core service with outbound ports
    service := pomodoro.NewService(
        mockClock,
        fileStorage,
        nil, // notifier set later
    )
    
    // Run systray with UI adapter (primary/driver)
    systray.Run(func() {
        // Create UI adapter
        uiAdapter := ui.NewSystrayAdapter(service)
        
        // Wire UI as notifier
        service.SetNotifier(uiAdapter)
        
        // Initialize menu
        uiAdapter.InitializeMenu()
    }, onExit)
}
```

### Responsibilities After This Change

**PomodoroService (domain logic - orchestration)**
- Enforce pomodoro cycle rules (4 work → long break)
- Manage state transitions (idle → running → paused)
- Coordinate between timer and session packages
- Determine next session type and duration
- Publish events when state changes occur
- **Does NOT** know about UI, systray, or how events are handled

**SystrayAdapter (primary adapter - UI presentation)**
- Render menu bar and dropdown UI elements
- Convert button clicks into service commands
- Listen to service events and update visual elements
- Format times/cycles for user display
- **Does NOT** contain business logic or make decisions

**TestDriver (primary adapter - test harness)**
- Drive service programmatically from tests
- Record events for assertion
- Provide test-friendly synchronization primitives
- **Does NOT** require actual UI to run

**Timer (infrastructure - time management)**
- Count down seconds and emit tick events
- Use injected Clock for time source (real or mock)
- Remain unaware of pomodoro-specific logic
- **Does NOT** know about sessions or cycles

**Session (domain logic - cycle tracking)**
- Track completed work sessions counter
- Determine next session type based on current type and counter
- Provide cycle display information
- **Does NOT** know about timers or UI

**Storage (infrastructure - persistence)**
- Write CSV logs to filesystem
- **Does NOT** know about timers or UI

### Data Flow Examples

**Scenario: User Starts First Work Session**

```
1. User clicks Start button in systray
   └─► SystrayAdapter.handleStartClick()

2. UI adapter calls service command
   └─► service.StartSession()

3. Service determines session details
   ├─► session.GetDuration() → 1500 seconds
   ├─► session.CurrentType → "work"
   └─► timer.Start("work", 1500)

4. Timer starts, calls back to service
   └─► service.handleTimerStarted("work", 1500)

5. Service publishes events outbound
   ├─► notifier.SessionStarted("work", 1500)
   │   └─► SystrayAdapter.SessionStarted(...)
   │       └─► Updates cycle indicator: "Session 1/4 🍅○○○"
   │
   └─► storage.LogSession(..., "work", "started", 0)
       └─► Writes: "2025-12-28 13:00:00,work,started,0"
```

**Scenario: Test Verifies Complete Cycle**

```go
func TestFullPomodoroCycle(t *testing.T) {
    // Arrange: Create service with mock clock and test driver
    mockClock := clock.NewMock()
    testDriver := NewTestDriver()
    service := pomodoro.NewService(mockClock, nil, testDriver)
    
    // Act: Start first work session
    service.StartSession()
    
    // Assert: Session started event
    assert(testDriver.HasEvent(EventSessionStarted, "work", 1500))
    
    // Act: Fast-forward 25 minutes
    mockClock.Add(25 * time.Minute)
    time.Sleep(50 * time.Millisecond) // Let goroutine process
    
    // Assert: Session completed
    assert(testDriver.HasEvent(EventSessionCompleted, "work"))
    assert(service.GetCurrentSessionType() == "short_break")
    
    // Continue cycle...
    service.StartSession() // Start break
    mockClock.Add(5 * time.Minute)
    
    // ... repeat for full 4-session cycle
    
    // Assert final state
    assert(service.GetCompletedSessions() == 4)
    assert(service.GetCurrentSessionType() == "long_break")
}
```

### Migration Strategy

Since this is a significant refactoring, we'll migrate incrementally:

**Phase 1: Create Core Service**
- Create `internal/pomodoro/` package with Service and ports
- Service wraps existing timer and session packages
- Write unit tests for service in isolation
- **No UI changes yet**

**Phase 2: Create Test Adapter**
- Create `internal/adapters/test/` with TestDriver
- Write integration tests using TestDriver
- Verify full cycle works programmatically
- **Still no UI changes**

**Phase 3: Create Systray Adapter**
- Create `internal/adapters/ui/` with SystrayAdapter
- Move UI code from `internal/ui/window.go` to adapter
- Implement Notifier interface
- **Keep old UI working during transition**

**Phase 4: Switch Main Wiring**
- Update `cmd/gopomodoro/main.go` to use new architecture
- Remove old `internal/ui/window.go`
- **Complete migration**

Each phase is deployable and testable independently.

---

## Contracts and Data

### Inbound Port Interface (Commands)

```go
// Service commands - called by drivers (UI, tests)
type Service interface {
    // Session lifecycle
    StartSession() error              // Returns error if invalid state
    PauseSession() error              // Returns error if not running
    ResumeSession() error             // Returns error if not paused
    SkipSession() error               // Returns error if idle
    ResetCycle() error                // Always succeeds
    
    // Queries (safe to call anytime)
    GetState() State                  // Current state (idle/running/paused)
    GetRemainingSeconds() int         // 0 if idle
    GetCurrentSessionType() string    // "work", "short_break", "long_break"
    GetCompletedSessions() int        // 0-4
    GetCycleProgress() string         // "Session 2/4 🍅🍅○○"
}

// State represents the timer state
type State int
const (
    StateIdle State = iota    // Not running
    StateRunning              // Counting down
    StatePaused               // Paused mid-session
)
```

### Outbound Port Interfaces (Events)

```go
// Notifier publishes session events
type Notifier interface {
    SessionStarted(sessionType string, duration int)
    SessionTick(remainingSeconds int)
    SessionCompleted(sessionType string)
    StateChanged(state State)
}

// Storage persists session data
type Storage interface {
    LogSession(
        timestamp time.Time,
        sessionType string,  // "work", "short_break", "long_break"
        status string,       // "started", "completed", "skipped"
        duration int,        // Actual duration in minutes
    ) error
}

// Clock provides time source (real or mock)
type Clock interface {
    Ticker(duration time.Duration) *clock.Ticker
}
```

### Service Internal State

```go
type Service struct {
    // Dependencies (injected)
    session  *session.Session
    timer    *timer.Timer
    clock    Clock
    notifier Notifier
    storage  Storage
    
    // Internal state
    state    State
    mu       sync.RWMutex  // Protects concurrent access
}
```

### Error Handling

```go
var (
    ErrAlreadyRunning = errors.New("session already running")
    ErrNotRunning     = errors.New("no session running")
    ErrNotPaused      = errors.New("session not paused")
    ErrIdle           = errors.New("no active session")
)
```

---

## Testing Strategy

### Unit Tests (Service in Isolation)

Test `PomodoroService` with mock dependencies:

```go
func TestStartSession_FirstWorkSession(t *testing.T) {
    mockClock := clock.NewMock()
    mockNotifier := &MockNotifier{}
    mockStorage := &MockStorage{}
    
    service := NewService(mockClock, mockStorage, mockNotifier)
    
    err := service.StartSession()
    assert.NoError(err)
    assert.Equal(StateRunning, service.GetState())
    assert.Equal("work", service.GetCurrentSessionType())
    
    // Verify notifier was called
    assert.Equal(1, mockNotifier.SessionStartedCalls)
    assert.Equal("work", mockNotifier.LastSessionType)
    assert.Equal(1500, mockNotifier.LastDuration)
}

func TestStartSession_WhenAlreadyRunning_ReturnsError(t *testing.T) {
    service := NewService(...)
    service.StartSession()
    
    err := service.StartSession() // Try to start again
    assert.Equal(ErrAlreadyRunning, err)
}
```

### Integration Tests (with TestDriver)

Test complete cycles programmatically:

```go
func TestFullPomodoroCycle_FourSessions(t *testing.T) {
    mockClock := clock.NewMock()
    driver := NewTestDriver()
    service := NewServiceWithDriver(mockClock, driver)
    
    // Execute full cycle
    for i := 0; i < 4; i++ {
        // Work session
        service.StartSession()
        driver.AssertSessionType("work")
        mockClock.Add(25 * time.Minute)
        driver.WaitForCompletion(1 * time.Second)
        
        // Break
        service.StartSession()
        if i < 3 {
            driver.AssertSessionType("short_break")
            mockClock.Add(5 * time.Minute)
        } else {
            driver.AssertSessionType("long_break")
            mockClock.Add(15 * time.Minute)
        }
        driver.WaitForCompletion(1 * time.Second)
    }
    
    // Verify cycle completed
    assert.Equal(4, service.GetCompletedSessions())
    assert.Equal("work", service.GetCurrentSessionType())
    
    // Verify all events recorded
    events := driver.GetEvents()
    assert.Equal(16, len(events)) // 8 starts + 8 completions
}
```

### UI Manual Tests

After migration, manually test systray adapter:
1. Start/pause/resume work session
2. Complete full 4-session cycle
3. Skip during work session
4. Reset mid-cycle
5. Verify menu bar displays correctly
6. Verify dropdown menu updates

---

## Implementation Notes

### Avoiding Deadlocks

**Problem solved**: Service never calls back into adapters synchronously while holding locks.

**Pattern**:
```go
func (s *Service) StartSession() error {
    s.mu.Lock()
    // ... update internal state ...
    s.mu.Unlock()
    
    // Publish events AFTER releasing lock
    s.notifier.SessionStarted(sessionType, duration)
    return nil
}
```

### Concurrency Safety

- Service uses `sync.RWMutex` for state access
- Read operations (GetState, GetRemaining) use read lock
- Write operations (Start, Pause) use write lock
- Events published after lock release

### Timer Integration

Service owns timer lifecycle:
```go
// Service creates timer with injected clock
s.timer = timer.NewWithClock(s.clock)

// Service registers as timer callback receiver
s.timer.OnStarted(s.handleTimerStarted)
s.timer.OnTick(s.handleTimerTick)
s.timer.OnCompleted(s.handleTimerCompleted)

// Internal handlers translate to outbound events
func (s *Service) handleTimerTick(remaining int) {
    s.notifier.SessionTick(remaining)
}
```

### Backward Compatibility

During migration, both old and new code paths can coexist:
- Old: `cmd/gopomodoro/main.go` uses `internal/ui/window.go`
- New: `cmd/gopomodoro/main_hex.go` uses `internal/adapters/ui/`
- Build flag switches between them

---

## Deployment Considerations

### Build Variants

Create two build targets during migration:

```makefile
# Current version (callback-based)
build:
    go build -o bin/gopomodoro ./cmd/gopomodoro

# New version (hexagonal)
build-hex:
    go build -tags hex -o bin/gopomodoro-hex ./cmd/gopomodoro
```

### Feature Flag

Use environment variable for gradual rollout:

```go
if os.Getenv("GOPOMODORO_HEX") == "1" {
    runHexagonalVersion()
} else {
    runLegacyVersion()
}
```

### Rollback Plan

If issues found after deployment:
1. Switch build back to legacy version
2. Keep new packages in codebase for future retry
3. Fix issues in isolation
4. Re-deploy when stable

---

## Future Enhancements Enabled

Once hexagonal architecture is in place, these become trivial:

1. **CLI Adapter**: `gopomodoro start --duration 25m`
2. **HTTP API Adapter**: RESTful endpoint for remote control
3. **Desktop Notification Adapter**: OS-level notifications on completion
4. **Multiple UI Adapters**: macOS + Linux simultaneously
5. **Analytics Adapter**: Stream events to analytics service
6. **AI Adapter**: Auto-adjust durations based on focus patterns

All without touching core business logic.
