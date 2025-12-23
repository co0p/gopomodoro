# Implement: Functional Timer with Session Logging

## Context

This increment transforms the placeholder UI into a working 25-minute pomodoro timer.

**Goal:**
- Countdown timer from 25:00 to 00:00 with per-second display updates
- Start, Pause, Resume, and Reset controls
- Automatic session completion detection
- CSV session logging to `~/.gopomodoro/sessions.log`
- Fail-fast error handling for storage failures

**Key Non-Goals:**
- Break timers, cycles, notifications
- Configurable durations or settings
- Tray icon state changes
- Statistics display

**Design Approach:**
- New `internal/timer/` package: Pure countdown logic with event callbacks (OnStarted, OnTick, OnCompleted)
- New `internal/storage/` package: CSV persistence with fail-fast errors
- Extended `internal/ui/`: Orchestrates timer events, updates display, triggers logging

**Constitution Mode:** `lite` - Pragmatic steps, manual UI testing acceptable, focus on core functionality

**Links:**
- [increment.md](increment.md)
- [design.md](design.md)
- [CONSTITUTION.md](../../../CONSTITUTION.md)

**Status:** Not started
**Next step:** Step 1 – Storage directory creation

---

## 1. Workstreams

- **Workstream A** – Storage package (session persistence)
- **Workstream B** – Timer package (countdown engine)
- **Workstream C** – UI integration (event handling and display)
- **Workstream D** – Main wiring (initialization and orchestration)

---

## 2. Steps

### Step 1: Storage directory creation

- [ ] **Workstream:** A
- **Based on Design:** §6 Contracts and Data – Storage Package API (EnsureDataDir)
- **Files:** `internal/storage/storage.go`, `internal/storage/storage_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Create `internal/storage/storage_test.go`
    - Write `TestEnsureDataDir()` that calls `EnsureDataDir()` and verifies directory exists at `$HOME/.gopomodoro/` using `os.Stat()`
    - Test should fail with "undefined: EnsureDataDir"
  - **Green – Make the test(s) pass:**
    - Create `internal/storage/storage.go`
    - Implement `EnsureDataDir() error` using `os.MkdirAll()` with `0755` permissions
    - Use `os.UserHomeDir()` to get home directory path
    - Return error from `MkdirAll()` directly (no retry logic)
  - **Refactor – Clean up with tests green:**
    - Extract directory path construction to private `getDataDir() string` helper
    - Add package-level constant for directory name `".gopomodoro"`
- **CI / Checks:**
  - Run: `go test ./internal/storage/`
  - Verify test passes and directory is created

---

### Step 2: Session logging to CSV

- [ ] **Workstream:** A
- **Based on Design:** §6 Contracts and Data – LogSession() API, §6 CSV Format
- **Files:** `internal/storage/storage.go`, `internal/storage/storage_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestLogSession()` in `storage_test.go`
    - Use `t.TempDir()` to create isolated test directory
    - Call `LogSession(timestamp, "work", "completed", 25)` with a known timestamp
    - Read file contents and verify CSV line matches format: `timestamp,session_type,event,duration_minutes`
    - Test should fail with "undefined: LogSession"
  - **Green – Make the test(s) pass:**
    - Implement `LogSession(timestamp time.Time, sessionType, event string, durationMinutes int) error`
    - Open `sessions.log` in append mode using `os.OpenFile()` with `O_APPEND|O_CREATE|O_WRONLY`
    - Format CSV line: `timestamp.Format(time.RFC3339), sessionType, event, durationMinutes`
    - Write line with `fmt.Fprintf()` and newline
    - Close file and return any error
  - **Refactor – Clean up with tests green:**
    - Extract CSV line formatting to `formatCSVLine()` helper function
    - Extract full file path construction to `getSessionsLogPath() string` helper
- **CI / Checks:**
  - Run: `go test ./internal/storage/`
  - Verify CSV format exactly matches design specification

---

### Step 3: CSV header creation for new files

- [ ] **Workstream:** A
- **Based on Design:** §6 CSV Format – header line for new files
- **Files:** `internal/storage/storage.go`, `internal/storage/storage_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestLogSessionCreatesFileWithHeader()` in `storage_test.go`
    - Use fresh `t.TempDir()` where `sessions.log` doesn't exist
    - Call `LogSession()` once
    - Read file and verify first line is header: `timestamp,session_type,event,duration_minutes`
    - Verify second line is the actual session record
    - Test should fail (no header written currently)
  - **Green – Make the test(s) pass:**
    - In `LogSession()`, before opening file, check if it exists using `os.Stat()`
    - If error is `os.IsNotExist()`, set flag `needsHeader = true`
    - After opening file, if `needsHeader`, write header line first
    - Then write session record
  - **Refactor – Clean up with tests green:**
    - Extract header string as package constant `csvHeader`
    - Ensure header logic doesn't affect append-only behavior for existing files
- **CI / Checks:**
  - Run: `go test ./internal/storage/`
  - Verify both new-file and existing-file scenarios work

---

### Step 4: Storage error handling

- [ ] **Workstream:** A
- **Based on Design:** §8 Testing and Safety Net – fail-fast scenarios
- **Files:** `internal/storage/storage.go`, `internal/storage/storage_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestEnsureDataDirFailure()` that creates read-only parent directory in `t.TempDir()`
    - Attempt to call `EnsureDataDir()` (modified to accept custom path for testing)
    - Verify error is returned (not nil)
    - Write `TestLogSessionFailure()` that makes `sessions.log` directory (conflict)
    - Verify `LogSession()` returns error when file can't be opened
    - Tests should fail if errors aren't properly returned
  - **Green – Make the test(s) pass:**
    - Ensure all `os` package errors are returned directly from `EnsureDataDir()`
    - Ensure file open/write errors are returned from `LogSession()`
    - No silent error swallowing
  - **Refactor – Clean up with tests green:**
    - Wrap errors with `fmt.Errorf()` to add context: `"failed to create data directory: %w"`
    - Add similar wrapping for file operations: `"failed to log session: %w"`
- **CI / Checks:**
  - Run: `go test ./internal/storage/`
  - Verify error messages are descriptive and include original error

---

### Step 5: Timer initialization and state

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – Timer Package API (New, GetState, GetRemaining)
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Create `internal/timer/timer_test.go`
    - Write `TestNew()` that calls `New()` and verifies:
      - `GetState()` returns `StateIdle`
      - `GetRemaining()` returns `1500` (25 minutes in seconds)
    - Test should fail with "undefined: New"
  - **Green – Make the test(s) pass:**
    - Create `internal/timer/timer.go`
    - Define `State` type as `int` with constants: `StateIdle`, `StateRunning`, `StatePaused`
    - Define `Timer` struct with fields: `state State`, `remaining int`
    - Implement `New() *Timer` that returns `&Timer{state: StateIdle, remaining: 1500}`
    - Implement `GetState() State` and `GetRemaining() int` getter methods
  - **Refactor – Clean up with tests green:**
    - Add package constant `defaultDuration = 1500`
    - Use constant in `New()` for initial remaining value
- **CI / Checks:**
  - Run: `go test ./internal/timer/`
  - Verify initial state is correct

---

### Step 6: Start timer and state transition

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – Start() method, OnStarted event
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestStart()` in `timer_test.go`
    - Create timer, register `OnStarted()` callback that sets flag `startedCalled = true`
    - Call `Start()`
    - Verify state changed to `StateRunning`
    - Verify `startedCalled == true`
    - Test should fail with "undefined: Start"
  - **Green – Make the test(s) pass:**
    - Add `onStarted func()` field to `Timer` struct
    - Implement `OnStarted(handler func())` that stores handler
    - Implement `Start()` that:
      - Checks if state is `StateIdle` (no-op otherwise per design)
      - Sets `state = StateRunning`
      - Calls `onStarted()` if not nil
  - **Refactor – Clean up with tests green:**
    - Extract state transition to `setState(newState State)` helper that could call state-change hooks
    - Ensure `Start()` is idempotent when already running (no-op)
- **CI / Checks:**
  - Run: `go test ./internal/timer/`
  - Verify state transition and callback work

---

### Step 7: Timer tick mechanism

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – OnTick event, §9 Risks – timer drift mitigation
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestTick()` in `timer_test.go`
    - Create timer with modified remaining (e.g., 10 seconds for faster test)
    - Register `OnTick()` callback that captures `remainingSeconds` value
    - Call `Start()`
    - Wait ~1 second (or use fast-forwarded test with shorter tick interval)
    - Verify `OnTick` was called with decremented value (9, 8, 7...)
    - Test should fail (no tick mechanism exists)
  - **Green – Make the test(s) pass:**
    - Add `onTick func(int)` and `ticker *time.Ticker` fields to `Timer` struct
    - Implement `OnTick(handler func(int))` that stores handler
    - In `Start()`, create `time.NewTicker(1 * time.Second)` and start goroutine:
      - Loop on `<-ticker.C`
      - Decrement `remaining`
      - Call `onTick(remaining)` if not nil
    - Store ticker reference for later cleanup
  - **Refactor – Clean up with tests green:**
    - Extract tick interval as constant `tickInterval = 1 * time.Second`
    - Add mutex (`sync.Mutex`) to protect `remaining` and `state` fields from concurrent access
    - Ensure ticker is stopped when leaving `StateRunning`
- **CI / Checks:**
  - Run: `go test ./internal/timer/`
  - Verify countdown works and no race conditions (run with `go test -race`)

---

### Step 8: Pause and resume functionality

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – Pause(), Resume() methods
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestPauseAndResume()` in `timer_test.go`
    - Start timer, wait for one tick
    - Call `Pause()`, verify state is `StatePaused`
    - Capture current `remaining` value
    - Wait another second, verify `remaining` hasn't changed (ticks stopped)
    - Call `Resume()`, verify state is `StateRunning`
    - Wait for tick, verify `remaining` is decrementing again
    - Test should fail with "undefined: Pause"
  - **Green – Make the test(s) pass:**
    - Implement `Pause()`:
      - Check state is `StateRunning` (no-op otherwise)
      - Stop ticker with `ticker.Stop()`
      - Set `state = StatePaused`
    - Implement `Resume()`:
      - Check state is `StatePaused` (no-op otherwise)
      - Set `state = StateRunning`
      - Create new ticker and restart goroutine (reuse logic from `Start()`)
  - **Refactor – Clean up with tests green:**
    - Extract ticker start logic to `startTicker()` helper used by both `Start()` and `Resume()`
    - Extract ticker stop logic to `stopTicker()` helper used by `Pause()` and `Reset()`
    - Ensure goroutine cleanup prevents leaks
- **CI / Checks:**
  - Run: `go test ./internal/timer/ -race`
  - Verify pause/resume works correctly without races

---

### Step 9: Reset to idle

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – Reset() method
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestReset()` in `timer_test.go`
    - Start timer, let it tick a few times
    - Call `Reset()`
    - Verify `GetState()` returns `StateIdle`
    - Verify `GetRemaining()` returns `1500`
    - Verify ticker is stopped (no more tick callbacks)
    - Test from paused state as well
    - Test should fail with "undefined: Reset"
  - **Green – Make the test(s) pass:**
    - Implement `Reset()`:
      - Stop ticker if running (call `stopTicker()` helper)
      - Set `state = StateIdle`
      - Set `remaining = defaultDuration` (1500)
  - **Refactor – Clean up with tests green:**
    - Ensure `Reset()` can be called from any state safely
    - Verify no goroutine leaks with ticker cleanup
- **CI / Checks:**
  - Run: `go test ./internal/timer/ -race`
  - Verify reset works from all states

---

### Step 10: Session completion detection

- [ ] **Workstream:** B
- **Based on Design:** §6 Contracts and Data – OnCompleted event
- **Files:** `internal/timer/timer.go`, `internal/timer/timer_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestCompletion()` in `timer_test.go`
    - Create timer, set `remaining = 1` manually for fast test (or use short tick for test)
    - Register `OnCompleted()` callback that sets `completedCalled = true`
    - Call `Start()`
    - Wait for tick to process
    - Verify `completedCalled == true`
    - Verify `GetState()` returns `StateIdle`
    - Verify ticker stopped
    - Test should fail (no completion logic)
  - **Green – Make the test(s) pass:**
    - Add `onCompleted func()` field to `Timer` struct
    - Implement `OnCompleted(handler func())` that stores handler
    - In tick goroutine, after decrementing `remaining`:
      - Check if `remaining == 0`
      - If true: call `onCompleted()` if not nil, stop ticker, set `state = StateIdle`
  - **Refactor – Clean up with tests green:**
    - Extract completion logic to `handleCompletion()` method
    - Ensure ticker cleanup happens properly in completion path
- **CI / Checks:**
  - Run: `go test ./internal/timer/ -race`
  - Verify completion detection works and state transitions correctly

---

### Step 11: Button state helper logic

- [ ] **Workstream:** C
- **Based on Design:** §6 Contracts and Data – UI Updates (button state table)
- **Files:** `internal/ui/window.go`, `internal/ui/window_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestButtonStateLogic()` in `window_test.go`
    - Import `timer` package for `State` constants
    - Test cases for each state:
      - `StateIdle` → Start enabled, Pause disabled, Reset disabled
      - `StateRunning` → Start disabled, Pause enabled, Reset enabled
      - `StatePaused` → Start enabled, Pause disabled, Reset enabled
    - Test should fail with "undefined: shouldStartBeEnabled"
  - **Green – Make the test(s) pass:**
    - In `window.go`, implement helper functions:
      - `shouldStartBeEnabled(state timer.State) bool` - returns true for Idle or Paused
      - `shouldPauseBeEnabled(state timer.State) bool` - returns true for Running
      - `shouldResetBeEnabled(state timer.State) bool` - returns true for Running or Paused
  - **Refactor – Clean up with tests green:**
    - Consider consolidating to single `getButtonStates(state timer.State) (start, pause, reset bool)` if cleaner
- **CI / Checks:**
  - Run: `go test ./internal/ui/`
  - Verify button state logic matches design table exactly

---

### Step 12: Display formatting

- [ ] **Workstream:** C
- **Based on Design:** §6 Contracts and Data – UI Updates (MM:SS format)
- **Files:** `internal/ui/window.go`, `internal/ui/window_test.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestFormatTime()` in `window_test.go`
    - Test cases:
      - `1500` seconds → `"25:00"`
      - `1499` seconds → `"24:59"`
      - `61` seconds → `"01:01"`
      - `3` seconds → `"00:03"`
      - `0` seconds → `"00:00"`
    - Test should fail with "undefined: formatTime"
  - **Green – Make the test(s) pass:**
    - Implement `formatTime(seconds int) string`:
      - Calculate minutes: `minutes := seconds / 60`
      - Calculate remaining seconds: `secs := seconds % 60`
      - Return `fmt.Sprintf("%02d:%02d", minutes, secs)`
  - **Refactor – Clean up with tests green:**
    - None needed (straightforward implementation)
- **CI / Checks:**
  - Run: `go test ./internal/ui/`
  - Verify all time formatting edge cases work

---

### Step 13: Timer event handlers in UI

- [ ] **Workstream:** C
- **Based on Design:** §5 Architecture – UI subscribes to timer events
- **Files:** `internal/ui/window.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestUIHandlesTimerEvents()` in `window_test.go`
    - Create mock or test double for timer that can trigger callbacks
    - Register UI event handlers
    - Trigger events and verify UI state changes (via test accessors)
    - Test should fail with "undefined: handleTimerStarted"
  - **Green – Make the test(s) pass:**
    - Add `timer *timer.Timer` field to `Window` struct
    - Implement event handler methods:
      - `handleTimerStarted()` - logs start event
      - `handleTimerTick(remaining int)` - updates timer display using `formatTime(remaining)`
      - `handleTimerCompleted()` - logs completion, resets display to "25:00"
    - Add `SetTimer(t *timer.Timer)` method that:
      - Stores timer reference
      - Registers handlers: `t.OnStarted(w.handleTimerStarted)`, etc.
  - **Refactor – Clean up with tests green:**
    - Extract display update logic to `updateTimerDisplay(text string)` helper
    - Ensure each handler has clear single responsibility
- **CI / Checks:**
  - Run: `go test ./internal/ui/`
  - Verify handlers update UI state correctly

---

### Step 14: Button click handlers

- [ ] **Workstream:** C
- **Based on Design:** §5 Architecture – UI calls timer methods on button clicks
- **Files:** `internal/ui/window.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestButtonClickHandlers()` in `window_test.go`
    - Create window with mock timer
    - Simulate button clicks (may need test helpers)
    - Verify timer methods called (Start, Pause, Resume, Reset)
    - Test should fail (no click handling implemented)
  - **Green – Make the test(s) pass:**
    - Add `startClickHandlers()` method to `Window`
    - For each button, launch goroutine:
      ```go
      go func() {
          for range w.btnStart.ClickedCh {
              if w.timer.GetState() == timer.StateIdle {
                  w.timer.Start()
              } else if w.timer.GetState() == timer.StatePaused {
                  w.timer.Resume()
              }
          }
      }()
      ```
    - Similar goroutines for Pause and Reset buttons
    - Call `startClickHandlers()` after timer is set
  - **Refactor – Clean up with tests green:**
    - Extract click logic to individual handler methods: `handleStartClick()`, `handlePauseClick()`, `handleResetClick()`
    - Ensure goroutines are tracked for potential cleanup
    - Add safety checks (nil timer)
- **CI / Checks:**
  - Run: `go test ./internal/ui/`
  - Verify click handlers call correct timer methods

---

### Step 15: Storage calls from UI

- [ ] **Workstream:** C
- **Based on Design:** §5 Flow – UI calls storage on Started/Completed events
- **Files:** `internal/ui/window.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Write `TestUILogsSessionsOnEvents()` in `window_test.go`
    - Create window with mock storage
    - Trigger timer Started and Completed events
    - Verify `LogSession()` called with correct parameters:
      - Started: `("work", "started", 0)`
      - Completed: `("work", "completed", 25)`
    - Test should fail (no storage calls)
  - **Green – Make the test(s) pass:**
    - Add `storage *storage.Storage` or storage interface field to `Window` struct
    - In `handleTimerStarted()`:
      - Call `w.storage.LogSession(time.Now(), "work", "started", 0)`
      - Log error if returned (per fail-fast requirement): `log.Fatalf("[ERROR] Failed to log session: %v", err)`
    - In `handleTimerCompleted()`:
      - Call `w.storage.LogSession(time.Now(), "work", "completed", 25)`
      - Log error and exit on failure
    - Add `SetStorage(s *storage.Storage)` method
  - **Refactor – Clean up with tests green:**
    - Extract session type `"work"` as constant
    - Extract duration `25` as constant (matches timer default)
    - Consider extracting error handling to helper method
- **CI / Checks:**
  - Run: `go test ./internal/ui/`
  - Verify storage calls happen on correct events

---

### Step 16: Component initialization and wiring

- [ ] **Workstream:** D
- **Based on Design:** §5 Architecture – Main creates and connects components
- **Files:** `cmd/gopomodoro/main.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Manual test: Run app with incomplete wiring
    - Verify error or non-functional buttons
  - **Green – Make the test(s) pass:**
    - In `onReady()` function, after creating window:
      - Create timer: `tmr := timer.New()`
      - Create storage: no constructor needed (uses package functions)
      - Wire timer to window: `window.SetTimer(tmr)`
      - Wire storage to window: `window.SetStorage()` or pass storage reference
      - Start click handlers: `window.StartClickHandlers()`
  - **Refactor – Clean up with tests green:**
    - Extract wiring logic to `setupComponents(window *ui.Window)` function
    - Keep `onReady()` high-level and readable
    - Ensure initialization order is correct (timer/storage before wiring)
- **CI / Checks:**
  - Run: `make build && make run`
  - Manual test: Verify app launches without errors

---

### Step 17: Startup storage check

- [ ] **Workstream:** D
- **Based on Design:** §7 CI/CD – fail-fast on storage errors, increment AC 11
- **Files:** `cmd/gopomodoro/main.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Manual test: Make `~/.gopomodoro` unwritable (`chmod 000` on parent)
    - Run app
    - Expect error message and exit
    - Currently would start without checking
  - **Green – Make the test(s) pass:**
    - In `onReady()`, before creating window:
      ```go
      if err := storage.EnsureDataDir(); err != nil {
          log.Fatalf("[ERROR] Failed to initialize storage: %v", err)
      }
      log.Println("[INFO] Data directory ensured: /Users/username/.gopomodoro")
      ```
  - **Refactor – Clean up with tests green:**
    - Consider logging actual path from `storage.GetDataDir()` helper
    - Ensure error message is clear and actionable for user
- **CI / Checks:**
  - Run: `make run`
  - Manual tests:
    - Normal case: directory created successfully
    - Error case: read-only parent directory → clean error and exit

---

### Step 18: Enable buttons after wiring

- [ ] **Workstream:** D
- **Based on Design:** §6 UI Updates – Start enabled in idle state
- **Files:** `cmd/gopomodoro/main.go`, `internal/ui/window.go`
- **TDD Cycle:**
  - **Red – Failing test first:**
    - Manual test: Launch app, observe all buttons disabled
    - Expected: Start button should be enabled in idle state
  - **Green – Make the test(s) pass:**
    - Add `UpdateButtonStates(state timer.State)` method to `Window`:
      - Use helper functions from Step 11
      - Call `w.btnStart.Enable()` or `Disable()` based on state
      - Same for Pause and Reset buttons
    - In `onReady()` after wiring:
      ```go
      window.UpdateButtonStates(tmr.GetState())
      ```
    - In timer event handlers (`handleTimerStarted`, etc.), call `UpdateButtonStates()` after state changes
  - **Refactor – Clean up with tests green:**
    - Ensure button state updates happen consistently on all timer state transitions
    - Consider adding state change callback to timer to automatically trigger UI updates
- **CI / Checks:**
  - Run: `make run`
  - Manual test checklist:
    - Launch app → Start enabled, others disabled
    - Click Start → Pause and Reset enabled, Start disabled
    - Click Pause → Start and Reset enabled, Pause disabled
    - Click Reset → back to initial state

---

## 3. Rollout & Validation Notes

### Suggested Grouping into PRs

**PR 1: Storage foundation (Steps 1-4)**
- Complete storage package with all tests
- Can be reviewed and merged independently
- Validates CSV format and fail-fast behavior

**PR 2: Timer engine (Steps 5-10)**
- Complete timer package with all tests
- Independent of storage and UI
- Validates countdown logic and state machine

**PR 3: UI integration (Steps 11-15)**
- Extends UI to handle timer events and call storage
- Depends on timer and storage packages
- Includes unit tests for formatting and state logic

**PR 4: Final wiring and manual validation (Steps 16-18)**
- Connects all components in main.go
- Enables full end-to-end manual testing
- Should be small and focused on wiring only

### Validation Checkpoints

**After Step 4 (Storage complete):**
- Run: `go test ./internal/storage/`
- Manually inspect test coverage
- Verify CSV format matches PRD exactly
- Test fail-fast with permission errors

**After Step 10 (Timer complete):**
- Run: `go test ./internal/timer/ -race`
- Manually verify countdown accuracy (no significant drift)
- Test all state transitions with unit tests
- Consider writing simple main.go test harness to watch timer tick in terminal

**After Step 15 (UI integration complete):**
- Run: `go test ./internal/ui/`
- Verify button state logic comprehensive
- Review event handler test coverage

**After Step 18 (Full system wired):**
- Manual testing checklist:
  - [ ] Launch app → tray icon appears
  - [ ] Click tray → dropdown shows
  - [ ] Timer shows "25:00", Start enabled
  - [ ] Click Start → countdown begins, display updates every second
  - [ ] Let countdown run for ~30 seconds, verify display accurate
  - [ ] Click Pause → countdown stops at current time
  - [ ] Click Start/Resume → countdown continues
  - [ ] Click Reset → returns to "25:00" and idle state
  - [ ] Let full 25-minute session complete → timer returns to idle automatically
  - [ ] Check `~/.gopomodoro/sessions.log`:
    - [ ] Directory created
    - [ ] File has header: `timestamp,session_type,event,duration_minutes`
    - [ ] Contains "started" record with duration 0
    - [ ] Contains "completed" record with duration 25
    - [ ] Timestamps are ISO 8601 format
  - [ ] Delete `~/.gopomodoro/`, restart app → directory recreated
  - [ ] Make `~/.gopomodoro` read-only, restart app → error displayed and app exits

### Success Criteria

Implementation complete when:
- All unit tests pass (`go test ./...`)
- No race conditions detected (`go test -race ./...`)
- Manual validation checklist fully passed
- `sessions.log` format matches PRD specification
- Fail-fast behavior verified (storage errors cause immediate exit with clear message)
- Timer accuracy verified (no drift >1 second over 25 minutes)
- All acceptance criteria from [increment.md](increment.md) satisfied
