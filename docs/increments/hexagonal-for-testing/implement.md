# Implement: Hexagonal Architecture for Testing

## Context

This increment refactors GoPomodoro from callback-based architecture to hexagonal (ports and adapters), enabling:
- Automated testing of complete pomodoro cycles without UI
- Elimination of deadlock issues from circular timer ↔ UI dependencies
- Clear separation between business logic (service) and infrastructure (adapters)

**Key constraints:**
- **Mode**: `lite` - pragmatic steps, manual UI testing acceptable
- **Migration**: Keep app working at each phase boundary
- **Non-goals**: No UI redesign, no new features beyond testability

**Links:**
- [design.md](design.md) - Technical approach and component contracts
- [../../CONSTITUTION.md](../../CONSTITUTION.md) - Project principles (mode: lite)

**Status:** Not started  
**Next step:** Step 1 - Create port interfaces

---

## 1. Workstreams

- **Workstream A** - Core Service & Ports (`internal/pomodoro/`)
- **Workstream B** - Test Infrastructure (`internal/adapters/test/`)
- **Workstream C** - UI Adapter (`internal/adapters/ui/`)
- **Workstream D** - Main Wiring & Cleanup (`cmd/gopomodoro/`, cleanup)

---

## 2. Steps

### Step 1: Create port interfaces

- **Workstream:** A
- **Based on Design:** §"Contracts and Data - Outbound Port Interfaces"
- **Files:** `internal/pomodoro/ports.go`
- **Description:**
  
  Define the three outbound port interfaces that the service will depend on: `Notifier` (for publishing events), `Storage` (for persisting sessions), and `Clock` (for time source). Also define the `State` type and constants.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    N/A - Pure interface definitions, no behavior to test yet.
  
  - **Green – Make the test(s) pass:**  
    Create `internal/pomodoro/ports.go` with:
    ```go
    type Notifier interface {
        SessionStarted(sessionType string, duration int)
        SessionTick(remainingSeconds int)
        SessionCompleted(sessionType string)
        StateChanged(state State)
    }
    
    type Storage interface {
        LogSession(timestamp time.Time, sessionType, status string, duration int) error
    }
    
    type Clock interface {
        Ticker(duration time.Duration) *clock.Ticker
    }
    
    type State int
    const (
        StateIdle State = iota
        StateRunning
        StatePaused
    )
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go build ./internal/pomodoro`

---

### Step 2: Create Service struct with state management

- **Workstream:** A
- **Based on Design:** §"Component Changes - service.go Key Behaviors", §"Service Internal State"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Create the `Service` struct with dependencies (timer, session, clock, notifier, storage) and internal state. Implement `NewService()` constructor that initializes the service in idle state.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestNewService` in `service_test.go`:
    ```go
    func TestNewService(t *testing.T) {
        service := NewService(mockClock, nil, nil)
        if service.GetState() != StateIdle {
            t.Errorf("expected StateIdle, got %v", service.GetState())
        }
    }
    ```
    This will fail because `NewService()` and `GetState()` don't exist yet.
  
  - **Green – Make the test(s) pass:**  
    In `service.go`, implement:
    ```go
    type Service struct {
        session  *session.Session
        timer    *timer.Timer
        clock    Clock
        notifier Notifier
        storage  Storage
        state    State
        mu       sync.RWMutex
    }
    
    func NewService(clk Clock, storage Storage, notifier Notifier) *Service {
        sess := session.New()
        tmr := timer.NewWithClock(clk)
        return &Service{
            session:  sess,
            timer:    tmr,
            clock:    clk,
            storage:  storage,
            notifier: notifier,
            state:    StateIdle,
        }
    }
    
    func (s *Service) GetState() State {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.state
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract mock helper for test setup if creating multiple test cases.

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 3: Implement StartSession command

- **Workstream:** A
- **Based on Design:** §"Data Flow Examples - User Starts First Work Session", §"Inbound Port Interface - StartSession()"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Implement `StartSession()` method that determines session type and duration from session package, starts timer, updates state to running, and publishes `SessionStarted` event.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestStartSession_FirstWorkSession`:
    ```go
    func TestStartSession_FirstWorkSession(t *testing.T) {
        mockNotifier := &MockNotifier{}
        service := NewService(clock.NewMock(), nil, mockNotifier)
        
        err := service.StartSession()
        assert.NoError(t, err)
        assert.Equal(t, StateRunning, service.GetState())
        assert.Equal(t, 1, mockNotifier.SessionStartedCalls)
        assert.Equal(t, "work", mockNotifier.LastSessionType)
        assert.Equal(t, 1500, mockNotifier.LastDuration)
    }
    ```
    Test fails because `StartSession()` doesn't exist.
  
  - **Green – Make the test(s) pass:**  
    Implement in `service.go`:
    ```go
    func (s *Service) StartSession() error {
        sessionType := s.session.CurrentType
        duration := s.session.GetDuration()
        
        s.timer.Start(sessionType, duration)
        
        s.mu.Lock()
        s.state = StateRunning
        s.mu.Unlock()
        
        if s.notifier != nil {
            s.notifier.SessionStarted(sessionType, duration)
        }
        return nil
    }
    ```
    Create `MockNotifier` in `service_test.go` to record calls.
  
  - **Refactor – Clean up with tests green:**  
    Consider extracting session type/duration query into helper method if reused.

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 4: Implement state validation (already running error)

- **Workstream:** A
- **Based on Design:** §"Error Handling", §"Implementation Notes - Avoiding Deadlocks"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Add state machine validation to `StartSession()` to prevent starting when already running. Return `ErrAlreadyRunning` error.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestStartSession_WhenAlreadyRunning_ReturnsError`:
    ```go
    func TestStartSession_WhenAlreadyRunning_ReturnsError(t *testing.T) {
        service := NewService(clock.NewMock(), nil, nil)
        service.StartSession()
        
        err := service.StartSession()
        assert.Equal(t, ErrAlreadyRunning, err)
    }
    ```
    Test fails because `StartSession()` doesn't check state yet.
  
  - **Green – Make the test(s) pass:**  
    Update `StartSession()`:
    ```go
    var ErrAlreadyRunning = errors.New("session already running")
    
    func (s *Service) StartSession() error {
        s.mu.Lock()
        if s.state != StateIdle {
            s.mu.Unlock()
            return ErrAlreadyRunning
        }
        s.state = StateRunning
        s.mu.Unlock()
        
        sessionType := s.session.CurrentType
        duration := s.session.GetDuration()
        s.timer.Start(sessionType, duration)
        
        if s.notifier != nil {
            s.notifier.SessionStarted(sessionType, duration)
        }
        return nil
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 5: Implement PauseSession and ResumeSession commands

- **Workstream:** A
- **Based on Design:** §"Inbound Port Interface - PauseSession/ResumeSession", §"Component Changes - Key Behaviors"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Implement pause and resume with state validation. Pause only works when running, resume only when paused. Emit `StateChanged` events.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add tests:
    ```go
    func TestPauseSession_WhenRunning_Succeeds(t *testing.T) { ... }
    func TestPauseSession_WhenIdle_ReturnsError(t *testing.T) { ... }
    func TestResumeSession_WhenPaused_Succeeds(t *testing.T) { ... }
    func TestResumeSession_WhenNotPaused_ReturnsError(t *testing.T) { ... }
    ```
    Tests fail because methods don't exist.
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    var (
        ErrNotRunning = errors.New("no session running")
        ErrNotPaused = errors.New("session not paused")
    )
    
    func (s *Service) PauseSession() error {
        s.mu.Lock()
        if s.state != StateRunning {
            s.mu.Unlock()
            return ErrNotRunning
        }
        s.state = StatePaused
        s.mu.Unlock()
        
        s.timer.Pause()
        if s.notifier != nil {
            s.notifier.StateChanged(StatePaused)
        }
        return nil
    }
    
    func (s *Service) ResumeSession() error {
        s.mu.Lock()
        if s.state != StatePaused {
            s.mu.Unlock()
            return ErrNotPaused
        }
        s.state = StateRunning
        s.mu.Unlock()
        
        s.timer.Resume()
        if s.notifier != nil {
            s.notifier.StateChanged(StateRunning)
        }
        return nil
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract state transition pattern into helper if repeated.

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 6: Implement SkipSession command with cycle advancement

- **Workstream:** A
- **Based on Design:** §"Component Changes - Key Behaviors - SkipSession", §"Session Integration"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Implement skip that stops current timer, increments cycle counter, determines next session, logs skip event, and returns to idle state.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestSkipSession_AdvancesCycle`:
    ```go
    func TestSkipSession_AdvancesCycle(t *testing.T) {
        service := NewService(clock.NewMock(), nil, nil)
        service.StartSession() // Start work session 1
        
        err := service.SkipSession()
        assert.NoError(t, err)
        assert.Equal(t, StateIdle, service.GetState())
        assert.Equal(t, 1, service.session.CompletedWorkSessions)
        assert.Equal(t, "short_break", service.session.CurrentType)
    }
    ```
    Test fails because `SkipSession()` doesn't exist.
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (s *Service) SkipSession() error {
        s.mu.Lock()
        if s.state == StateIdle {
            s.mu.Unlock()
            return ErrIdle
        }
        currentType := s.session.CurrentType
        s.mu.Unlock()
        
        s.timer.Stop()
        s.session.IncrementCycle()
        nextType, nextDuration := s.session.DetermineNext()
        s.session.CurrentType = nextType
        
        s.mu.Lock()
        s.state = StateIdle
        s.mu.Unlock()
        
        if s.notifier != nil {
            s.notifier.SessionCompleted(currentType)
        }
        return nil
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract cycle advancement into `advanceCycle()` helper method since it's reused in timer completion handler.

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 7: Implement ResetCycle command

- **Workstream:** A
- **Based on Design:** §"Inbound Port Interface - ResetCycle", §"Component Changes - Key Behaviors"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Implement reset that stops timer, resets session counter to 0, returns to work session 1, and sets state to idle.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestResetCycle_ResetsToInitialState`:
    ```go
    func TestResetCycle_ResetsToInitialState(t *testing.T) {
        service := NewService(clock.NewMock(), nil, nil)
        service.session.CompletedWorkSessions = 3
        service.session.CurrentType = "short_break"
        
        err := service.ResetCycle()
        assert.NoError(t, err)
        assert.Equal(t, StateIdle, service.GetState())
        assert.Equal(t, 0, service.session.CompletedWorkSessions)
        assert.Equal(t, "work", service.session.CurrentType)
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (s *Service) ResetCycle() error {
        s.timer.Stop()
        s.session.Reset()
        
        s.mu.Lock()
        s.state = StateIdle
        s.mu.Unlock()
        
        return nil
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 8: Implement timer callbacks (internal handlers)

- **Workstream:** A
- **Based on Design:** §"Implementation Notes - Timer Integration", §"Component Changes - Service methods called by timer"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Register service as timer event receiver. Implement internal handlers for `OnStarted`, `OnTick`, and `OnCompleted`. The completion handler advances the cycle and auto-starts the next session.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestTimerCompletion_AdvancesToNextSession`:
    ```go
    func TestTimerCompletion_AdvancesToNextSession(t *testing.T) {
        mockClock := clock.NewMock()
        mockNotifier := &MockNotifier{}
        service := NewService(mockClock, nil, mockNotifier)
        
        service.StartSession() // Start work session
        mockClock.Add(25 * time.Minute) // Complete it
        time.Sleep(50 * time.Millisecond) // Wait for goroutine
        
        assert.Equal(t, 1, service.session.CompletedWorkSessions)
        assert.Equal(t, "short_break", service.session.CurrentType)
        assert.Contains(t, mockNotifier.Events, EventSessionCompleted)
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    In `NewService()`, register handlers:
    ```go
    func NewService(...) *Service {
        // ... existing code ...
        tmr.OnStarted(s.handleTimerStarted)
        tmr.OnTick(s.handleTimerTick)
        tmr.OnCompleted(s.handleTimerCompleted)
        return s
    }
    
    func (s *Service) handleTimerStarted(sessionType string, duration int) {
        // Already published in StartSession, no-op here
    }
    
    func (s *Service) handleTimerTick(remaining int) {
        if s.notifier != nil {
            s.notifier.SessionTick(remaining)
        }
    }
    
    func (s *Service) handleTimerCompleted() {
        currentType := s.session.CurrentType
        s.session.IncrementCycle()
        nextType, nextDuration := s.session.DetermineNext()
        s.session.CurrentType = nextType
        
        s.mu.Lock()
        s.state = StateIdle
        s.mu.Unlock()
        
        if s.notifier != nil {
            s.notifier.SessionCompleted(currentType)
        }
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract common cycle advancement logic from `SkipSession()` and `handleTimerCompleted()` into `advanceCycle()` helper.

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 9: Implement query methods (GetState, GetRemaining, etc.)

- **Workstream:** A
- **Based on Design:** §"Inbound Port Interface - Queries", §"Concurrency Safety"
- **Files:** `internal/pomodoro/service.go`, `internal/pomodoro/service_test.go`
- **Description:**
  
  Implement safe query methods: `GetRemainingSeconds()`, `GetCurrentSessionType()`, `GetCompletedSessions()`, `GetCycleProgress()`. Use read locks for concurrency safety.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestQueryMethods_ReturnCorrectValues`:
    ```go
    func TestQueryMethods_ReturnCorrectValues(t *testing.T) {
        service := NewService(clock.NewMock(), nil, nil)
        service.StartSession()
        
        assert.Equal(t, 1500, service.GetRemainingSeconds())
        assert.Equal(t, "work", service.GetCurrentSessionType())
        assert.Equal(t, 0, service.GetCompletedSessions())
        assert.Contains(t, service.GetCycleProgress(), "Session 1/4")
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (s *Service) GetRemainingSeconds() int {
        return s.timer.GetRemaining()
    }
    
    func (s *Service) GetCurrentSessionType() string {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.session.CurrentType
    }
    
    func (s *Service) GetCompletedSessions() int {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.session.CompletedWorkSessions
    }
    
    func (s *Service) GetCycleProgress() string {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.session.FormatCycleIndicator()
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go test ./internal/pomodoro`

---

### Step 10: Create TestDriver with event recording

- **Workstream:** B
- **Based on Design:** §"Component Changes - Test Driver", §"Testing Strategy - Integration Tests with TestDriver"
- **Files:** `internal/adapters/test/test_driver.go`, `internal/adapters/test/test_driver_test.go`
- **Description:**
  
  Create `TestDriver` struct that implements `Notifier` interface and records all events in a slice for later assertion.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestDriverRecordsEvents`:
    ```go
    func TestDriverRecordsEvents(t *testing.T) {
        driver := NewTestDriver()
        driver.SessionStarted("work", 1500)
        
        events := driver.GetEvents()
        assert.Equal(t, 1, len(events))
        assert.Equal(t, EventSessionStarted, events[0].Type)
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    Implement in `test_driver.go`:
    ```go
    type EventType int
    const (
        EventSessionStarted EventType = iota
        EventSessionTick
        EventSessionCompleted
        EventStateChanged
    )
    
    type Event struct {
        Type EventType
        Timestamp time.Time
        Data interface{}
    }
    
    type TestDriver struct {
        events []Event
        mu sync.Mutex
    }
    
    func NewTestDriver() *TestDriver {
        return &TestDriver{events: []Event{}}
    }
    
    func (d *TestDriver) SessionStarted(sessionType string, duration int) {
        d.mu.Lock()
        defer d.mu.Unlock()
        d.events = append(d.events, Event{
            Type: EventSessionStarted,
            Timestamp: time.Now(),
            Data: map[string]interface{}{"type": sessionType, "duration": duration},
        })
    }
    
    func (d *TestDriver) GetEvents() []Event {
        d.mu.Lock()
        defer d.mu.Unlock()
        return append([]Event{}, d.events...)
    }
    ```
    Implement `SessionTick`, `SessionCompleted`, `StateChanged` similarly.
  
  - **Refactor – Clean up with tests green:**  
    Extract event creation into helper method.

- **CI / Checks:**
  - `go test ./internal/adapters/test`

---

### Step 11: Add test helper methods (WaitForEvent, AssertState)

- **Workstream:** B
- **Based on Design:** §"Component Changes - Test Driver - Test helpers"
- **Files:** `internal/adapters/test/test_driver.go`, `internal/adapters/test/test_driver_test.go`
- **Description:**
  
  Add convenience methods for tests: `WaitForEvent()` with timeout, `ClearEvents()`, and `AssertSessionType()`.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestWaitForEvent_TimesOut`:
    ```go
    func TestWaitForEvent_TimesOut(t *testing.T) {
        driver := NewTestDriver()
        err := driver.WaitForEvent(EventSessionCompleted, 100*time.Millisecond)
        assert.Error(t, err)
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (d *TestDriver) WaitForEvent(eventType EventType, timeout time.Duration) error {
        deadline := time.Now().Add(timeout)
        for time.Now().Before(deadline) {
            d.mu.Lock()
            for _, e := range d.events {
                if e.Type == eventType {
                    d.mu.Unlock()
                    return nil
                }
            }
            d.mu.Unlock()
            time.Sleep(10 * time.Millisecond)
        }
        return fmt.Errorf("event %v not received within timeout", eventType)
    }
    
    func (d *TestDriver) ClearEvents() {
        d.mu.Lock()
        defer d.mu.Unlock()
        d.events = []Event{}
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Use channels instead of polling if goroutine synchronization becomes an issue.

- **CI / Checks:**
  - `go test ./internal/adapters/test`

---

### Step 12: Write full cycle integration test

- **Workstream:** B
- **Based on Design:** §"Testing Strategy - Integration Tests", §"Scenario: Test Verifies Complete Cycle"
- **Files:** `internal/adapters/test/integration_test.go`
- **Description:**
  
  Write end-to-end test that drives service through 4 work sessions + breaks, verifies cycle progression, uses mock clock to fast-forward time, confirms long break after 4th session.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestFullPomodoroCycle_FourSessions`:
    ```go
    func TestFullPomodoroCycle_FourSessions(t *testing.T) {
        mockClock := clock.NewMock()
        driver := NewTestDriver()
        service := pomodoro.NewService(mockClock, nil, driver)
        
        // Work session 1
        service.StartSession()
        driver.WaitForEvent(EventSessionStarted, 1*time.Second)
        mockClock.Add(25 * time.Minute)
        driver.WaitForEvent(EventSessionCompleted, 1*time.Second)
        assert.Equal(t, 1, service.GetCompletedSessions())
        
        // Short break 1
        mockClock.Add(5 * time.Minute)
        driver.WaitForEvent(EventSessionCompleted, 1*time.Second)
        
        // ... repeat for sessions 2-4 ...
        
        // Verify final state
        assert.Equal(t, 4, service.GetCompletedSessions())
        assert.Equal(t, "long_break", service.GetCurrentSessionType())
    }
    ```
    Test initially fails due to timing or event sequence issues.
  
  - **Green – Make the test(s) pass:**  
    Ensure timer completion handler auto-starts next session as per design. Adjust test timing to account for goroutine delays. Add explicit `ClearEvents()` between assertions if needed.
  
  - **Refactor – Clean up with tests green:**  
    Extract helper function `completeSession(mockClock, duration, driver)` to reduce duplication.

- **CI / Checks:**
  - `go test ./internal/adapters/test`

---

### Step 13: Create SystrayAdapter struct

- **Workstream:** C
- **Based on Design:** §"Component Changes - Systray UI - systray_adapter.go"
- **Files:** `internal/adapters/ui/systray_adapter.go`
- **Description:**
  
  Create `SystrayAdapter` struct that holds reference to service and systray menu items. Set up basic constructor.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    N/A - Struct definition and constructor setup.
  
  - **Green – Make the test(s) pass:**  
    Create file:
    ```go
    package ui
    
    import (
        "github.com/co0p/gopomodoro/internal/pomodoro"
        "github.com/getlantern/systray"
    )
    
    type SystrayAdapter struct {
        service        *pomodoro.Service
        progressBar    *systray.MenuItem
        cycleIndicator *systray.MenuItem
        btnStart       *systray.MenuItem
        btnPause       *systray.MenuItem
        btnReset       *systray.MenuItem
        btnSkip        *systray.MenuItem
        btnQuit        *systray.MenuItem
    }
    
    func NewSystrayAdapter(service *pomodoro.Service) *SystrayAdapter {
        return &SystrayAdapter{service: service}
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go build ./internal/adapters/ui`

---

### Step 14: Implement Notifier interface methods

- **Workstream:** C
- **Based on Design:** §"Component Changes - Systray UI - Implements Notifier interface"
- **Files:** `internal/adapters/ui/systray_adapter.go`, `internal/adapters/ui/systray_adapter_test.go`
- **Description:**
  
  Implement `SessionStarted()`, `SessionTick()`, `SessionCompleted()`, and `StateChanged()` to update menu items. Use formatter helpers for progress bar and cycle indicator.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Add test `TestSessionStarted_UpdatesCycleIndicator`:
    ```go
    func TestSessionStarted_UpdatesCycleIndicator(t *testing.T) {
        adapter := &SystrayAdapter{
            cycleIndicator: &MockMenuItem{},
        }
        adapter.SessionStarted("work", 1500)
        assert.Contains(t, adapter.cycleIndicator.Title, "Session 1/4")
    }
    ```
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (a *SystrayAdapter) SessionStarted(sessionType string, duration int) {
        if a.cycleIndicator != nil {
            // Use session.FormatCycleIndicator() via service query
            a.cycleIndicator.SetTitle(a.service.GetCycleProgress())
        }
        if a.progressBar != nil {
            a.progressBar.SetTitle("○○○○○○○○○○")
        }
    }
    
    func (a *SystrayAdapter) SessionTick(remaining int) {
        if a.progressBar != nil {
            a.progressBar.SetTitle(formatProgressBar(remaining, a.service.GetCurrentSessionType()))
        }
    }
    
    func (a *SystrayAdapter) SessionCompleted(sessionType string) {
        // Could show notification here in future
    }
    
    func (a *SystrayAdapter) StateChanged(state pomodoro.State) {
        // Update button states
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract `formatProgressBar()` helper from old window.go.

- **CI / Checks:**
  - `go test ./internal/adapters/ui`

---

### Step 15: Implement button click handlers calling service

- **Workstream:** C
- **Based on Design:** §"Component Changes - Systray UI - Button handlers call service commands"
- **Files:** `internal/adapters/ui/systray_adapter.go`
- **Description:**
  
  Wire systray menu item click channels to service command methods. Handle errors by logging.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Manual testing required (systray click channels are difficult to unit test).
  
  - **Green – Make the test(s) pass:**  
    Implement:
    ```go
    func (a *SystrayAdapter) startClickHandlers() {
        go func() {
            for range a.btnStart.ClickedCh {
                a.handleStartClick()
            }
        }()
        
        go func() {
            for range a.btnPause.ClickedCh {
                a.handlePauseClick()
            }
        }()
        
        // ... similar for reset, skip, quit
    }
    
    func (a *SystrayAdapter) handleStartClick() {
        if err := a.service.StartSession(); err != nil {
            log.Printf("[WARN] Start failed: %v", err)
        }
    }
    
    func (a *SystrayAdapter) handlePauseClick() {
        state := a.service.GetState()
        if state == pomodoro.StateRunning {
            a.service.PauseSession()
        } else if state == pomodoro.StatePaused {
            a.service.ResumeSession()
        }
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    Extract error logging helper.

- **CI / Checks:**
  - `go build ./internal/adapters/ui`

---

### Step 16: Migrate InitializeMenu from window.go

- **Workstream:** C
- **Based on Design:** §"Migration Strategy Phase 3", existing `internal/ui/window.go InitializeMenu()`
- **Files:** `internal/adapters/ui/systray_adapter.go`
- **Description:**
  
  Copy menu initialization code from `window.go` into `SystrayAdapter.InitializeMenu()`. Update to use adapter's menu item fields.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Manual verification - menu appears with correct items.
  
  - **Green – Make the test(s) pass:**  
    Copy and adapt:
    ```go
    func (a *SystrayAdapter) InitializeMenu() {
        a.progressBar = systray.AddMenuItem("○○○○○○○○○○", "Session progress")
        a.progressBar.Disable()
        
        a.cycleIndicator = systray.AddMenuItem("Session 1/4  🍅○○○", "Cycle progress")
        a.cycleIndicator.Disable()
        
        systray.AddSeparator()
        
        a.btnStart = systray.AddMenuItem("Start", "Start timer")
        a.btnPause = systray.AddMenuItem("Pause", "Pause timer")
        a.btnReset = systray.AddMenuItem("Reset", "Reset timer")
        
        systray.AddSeparator()
        
        a.btnSkip = systray.AddMenuItem("Skip", "Skip to next session")
        
        systray.AddSeparator()
        
        a.btnQuit = systray.AddMenuItem("Quit", "Quit the application")
        
        // Start click handlers
        a.startClickHandlers()
    }
    ```
  
  - **Refactor – Clean up with tests green:**  
    N/A

- **CI / Checks:**
  - `go build ./internal/adapters/ui`

---

### Step 17: Update main.go to use Service and adapters

- **Workstream:** D
- **Based on Design:** §"Component Changes - Main Package", §"Migration Strategy Phase 4"
- **Files:** `cmd/gopomodoro/main.go`
- **Description:**
  
  Replace old window-based wiring with new hexagonal architecture: create service, create UI adapter, wire as notifier, initialize menu.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Manual smoke test: `go run ./cmd/gopomodoro --smoke`
  
  - **Green – Make the test(s) pass:**  
    Update `onReady()`:
    ```go
    func onReady() {
        // ... existing tray icon setup ...
        
        // Create infrastructure adapters
        mockClock := clock.New()
        fileStorage := storage.NewFileStorage()
        
        // Create core service
        service := pomodoro.NewService(mockClock, fileStorage, nil)
        
        // Create UI adapter
        uiAdapter := ui.NewSystrayAdapter(service)
        
        // Wire UI as notifier
        service.SetNotifier(uiAdapter)
        
        // Initialize menu
        uiAdapter.InitializeMenu()
        
        // ... smoke test handling ...
    }
    ```
    Add `SetNotifier()` method to service if not already present.
  
  - **Refactor – Clean up with tests green:**  
    Remove unused imports from old window approach.

- **CI / Checks:**
  - `go build -o bin/gopomodoro ./cmd/gopomodoro`
  - `./bin/gopomodoro --smoke`

---

### Step 18: Delete old window.go

- **Workstream:** D
- **Based on Design:** §"Migration Strategy Phase 4 - Remove old code"
- **Files:** Delete `internal/ui/window.go`, `internal/ui/window_test.go`
- **Description:**
  
  Remove legacy callback-based UI code now that adapter pattern is in place.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    N/A - deletion operation.
  
  - **Green – Make the test(s) pass:**  
    Run:
    ```bash
    rm internal/ui/window.go
    rm internal/ui/window_test.go
    ```
    Verify build still succeeds:
    ```bash
    go build ./...
    ```
  
  - **Refactor – Clean up with tests green:**  
    Remove `internal/ui/` directory entirely if now empty, or keep for future UI helpers.

- **CI / Checks:**
  - `go build ./...`
  - `go test ./...`

---

### Step 19: Run full manual UI test

- **Workstream:** D
- **Based on Design:** §"Testing Strategy - UI Manual Tests"
- **Files:** N/A
- **Description:**
  
  Manually verify complete application behavior through systray UI.

- **TDD Cycle:**
  
  - **Red – Failing test first:**  
    Manual checklist (failures indicate issues to fix).
  
  - **Green – Make the test(s) pass:**  
    Test scenarios:
    1. Launch app, verify tray icon appears
    2. Click Start, verify work session starts (progress bar updates)
    3. Click Pause, verify timer pauses
    4. Click Resume, verify timer resumes
    5. Let work session complete, verify auto-transition to short break
    6. Skip during break, verify cycle advances
    7. Complete full 4-session cycle, verify long break
    8. Click Reset mid-cycle, verify return to work session 1
    9. Verify cycle indicator shows correct tomatoes
    10. Click Quit, verify clean shutdown
  
  - **Refactor – Clean up with tests green:**  
    Fix any UI bugs discovered during testing.

- **CI / Checks:**
  - Manual verification checklist complete

---

## 3. Rollout & Validation Notes

### Suggested PR Grouping

- **PR 1: Core Service** (Steps 1-9)  
  Self-contained business logic with full unit test coverage. Can merge without affecting UI.

- **PR 2: Test Infrastructure** (Steps 10-12)  
  Demonstrates end-to-end testability. Integration tests prove architecture works.

- **PR 3: UI Adapter** (Steps 13-16)  
  New adapter code alongside old window.go. Not yet wired, so safe to merge.

- **PR 4: Migration Complete** (Steps 17-19)  
  Switch main.go, delete old code, final validation.

### Validation Checkpoints

**After Step 9:** Core service unit tests pass
```bash
go test ./internal/pomodoro -v
```
Verify: All state transitions, error handling, event publishing works in isolation.

**After Step 12:** Integration test passes
```bash
go test ./internal/adapters/test -v -run TestFullPomodoroCycle
```
Verify: Complete cycle can be driven programmatically without UI.

**After Step 16:** UI adapter builds
```bash
go build ./internal/adapters/ui
```
Verify: New adapter compiles, old window.go still present and working.

**After Step 19:** Manual UI acceptance
```bash
./bin/gopomodoro
```
Verify: All scenarios from manual test checklist pass.

### Rollback Strategy

If issues found after deployment:
- **Steps 1-12**: Isolated, no impact on running app
- **Steps 13-16**: UI adapter not yet wired, no risk
- **Step 17+**: If main.go changes cause issues, revert PR 4 to restore old window.go wiring

---

## Acceptance

Implementation complete when:
- ✅ All unit tests pass (`go test ./...`)
- ✅ Integration test drives full 4-session cycle without UI
- ✅ Manual UI test checklist completed successfully
- ✅ No deadlocks observed during extended usage
- ✅ Old callback-based code removed
- ✅ App builds and runs on macOS
