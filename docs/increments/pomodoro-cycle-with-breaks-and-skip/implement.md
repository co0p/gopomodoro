# Implement: Pomodoro Cycle with Breaks and Skip

## Context

This increment transforms the Phase 1 single-session timer into a full pomodoro cycle system:
- **Goal**: Support work sessions (25m), short breaks (5m), long breaks (15m), cycle tracking (1-4), and skip functionality
- **Design approach**: Timer becomes configurable (accepts session type and duration), UI orchestrates cycle state and transitions, Tray updates icons based on session type
- **Key constraints**: Lite mode (manual testing acceptable for UI flows), volatile cycle state (in-memory, no persistence), user manually starts each session

**Links**:
- [increment.md](increment.md)
- [design.md](design.md)
- [CONSTITUTION.md](../../../CONSTITUTION.md)

**Status**: Not started  
**Next step**: Step 1 – Add sessionType field and GetSessionType method

---

## 1. Workstreams

- **Workstream A** – Timer Package: Session Type Support
- **Workstream B** – Asset Creation: Emoji-Based Icons
- **Workstream C** – Tray Package: Icon State Updates
- **Workstream D** – UI Package: Cycle Orchestration
- **Workstream E** – UI Package: Controls and Display

---

## 2. Steps

### Step 1: Add sessionType field and GetSessionType method

- [ ] **Step 1: Add sessionType field and GetSessionType method**

**Workstream**: A  
**Based on Design**: Design §6 "Contracts and Data – Timer Package API Changes"  
**Files**: `internal/timer/timer.go`, `internal/timer/timer_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestGetSessionType_ReturnsEmptyWhenIdle()` expecting `GetSessionType()` to return empty string `""` from a newly created timer
  - Test will fail because `GetSessionType()` method doesn't exist yet
  
- **Green – Make the test(s) pass**:
  - Add `sessionType string` field to the `Timer` struct
  - Implement `GetSessionType() string` method that returns `t.sessionType`
  - New timer will have empty sessionType by default (zero value)
  
- **Refactor – Clean up with tests green**:
  - None needed for simple getter

**CI / Checks**: `go test ./internal/timer/`

---

### Step 2: Modify Start method to accept sessionType and duration parameters

- [ ] **Step 2: Modify Start method to accept sessionType and duration parameters**

**Workstream**: A  
**Based on Design**: Design §6 "Timer Package API Changes – Start signature"  
**Files**: `internal/timer/timer.go`, `internal/timer/timer_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestStart_SetsSessionTypeAndDuration()` that:
    - Calls `Start("work", 1500)`
    - Expects `GetSessionType()` to return `"work"`
    - Expects `GetRemaining()` to return `1500`
    - Expects state to be `StateRunning`
  - Test will fail because `Start()` doesn't accept parameters yet
  
- **Green – Make the test(s) pass**:
  - Change method signature from `Start()` to `Start(sessionType string, durationSeconds int)`
  - Set `t.sessionType = sessionType`
  - Set `t.remaining = durationSeconds`
  - Remove usage of `defaultDuration` constant in Start method
  
- **Refactor – Clean up with tests green**:
  - Remove or deprecate `defaultDuration` constant (no longer used)
  - Note: Existing tests will break – this is expected and will be addressed in next steps

**CI / Checks**: `go test ./internal/timer/` (some tests will fail – fix in subsequent steps)

---

### Step 3: Update OnStarted callback signature to include session context

- [ ] **Step 3: Update OnStarted callback signature to include session context**

**Workstream**: A  
**Based on Design**: Design §6 "Timer Package API Changes – OnStarted callback"  
**Files**: `internal/timer/timer.go`, `internal/timer/timer_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Update existing timer tests to register callback with new signature: `func(sessionType string, durationSeconds int)`
  - Add test `TestOnStarted_CallbackReceivesSessionContext()` that:
    - Registers callback capturing sessionType and duration parameters
    - Calls `Start("work", 1500)`
    - Expects callback to receive `"work"` and `1500`
  - Tests will fail because callback signature is still `func()`
  
- **Green – Make the test(s) pass**:
  - Change `onStarted func()` field to `onStarted func(string, int)`
  - Update `OnStarted(handler func())` method signature to `OnStarted(handler func(string, int))`
  - In `Start()` method, call `t.onStarted(sessionType, durationSeconds)` instead of `t.onStarted()`
  
- **Refactor – Clean up with tests green**:
  - Update all test callback handlers to use new signature
  - Remove any test scaffolding for old callback style

**CI / Checks**: `go test ./internal/timer/`

---

### Step 4: Update Reset to clear sessionType and set remaining to 0

- [ ] **Step 4: Update Reset to clear sessionType and set remaining to 0**

**Workstream**: A  
**Based on Design**: Design §6 "Timer Package API Changes – Reset behavior"  
**Files**: `internal/timer/timer.go`, `internal/timer/timer_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestReset_ClearsSessionType()` that:
    - Starts timer with `Start("work", 1500)`
    - Calls `Reset()`
    - Expects `GetSessionType()` to return `""`
    - Expects `GetRemaining()` to return `0`
  - Test will fail because Reset currently sets remaining to `defaultDuration` (1500)
  
- **Green – Make the test(s) pass**:
  - In `Reset()` method, change `t.remaining = defaultDuration` to `t.remaining = 0`
  - Add `t.sessionType = ""` to clear session type
  
- **Refactor – Clean up with tests green**:
  - Update any existing reset tests that expect remaining to be 1500 to expect 0 instead
  - Timer package is now fully session-type-aware

**CI / Checks**: `go test ./internal/timer/`

---

### Step 5: Generate emoji-based icon images

- [ ] **Step 5: Generate emoji-based icon images**

**Workstream**: B  
**Based on Design**: Design §6 "Tray Package API Addition – Icon updates", with emoji-based approach  
**Files**: `assets/icon-work.png`, `assets/icon-short-break.png`, `assets/icon-long-break.png`, `assets/icon-paused.png`

**TDD Cycle**:
- **Red – Failing test first**:
  - N/A (asset creation, not code)
  
- **Green – Make the test(s) pass**:
  - Generate or create PNG images from emoji characters:
    - `assets/icon-work.png` from 🍅 (red tomato emoji, 32x32 or 64x64 pixels)
    - `assets/icon-short-break.png` from ☕ (coffee cup emoji)
    - `assets/icon-long-break.png` from 🌟 (star emoji, gold/yellow)
    - `assets/icon-paused.png` from ⏸️ or ⏱️ (pause or timer emoji)
  - Keep existing `assets/icon-idle.png` (or create gray circle/neutral icon if needed)
  - Methods: Use emoji-to-PNG converter tool, screenshot emojis, or simple image editor
  
- **Refactor – Clean up with tests green**:
  - Ensure all icons are same dimensions (recommend 32x32 or 64x64 for tray icons)
  - Verify PNG format and transparency if needed
  - Test that files can be loaded by macOS

**CI / Checks**: Manual visual check – preview icons in Finder, verify they display correctly

---

### Step 6: Add Tray struct and UpdateIcon method

- [ ] **Step 6: Add Tray struct and UpdateIcon method**

**Workstream**: C  
**Based on Design**: Design §6 "Tray Package API Addition"  
**Files**: `internal/tray/tray.go`, `internal/tray/tray_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestUpdateIcon_SelectsCorrectAsset()` that:
    - Calls `UpdateIcon("work", timer.StateRunning)` and expects it to attempt loading `"icon-work.png"`
    - Calls `UpdateIcon("", timer.StateIdle)` and expects it to attempt loading `"icon-idle.png"`
  - Test will fail because `UpdateIcon` method doesn't exist
  - Note: May be difficult to unit test systray calls; manual testing acceptable per lite mode
  
- **Green – Make the test(s) pass**:
  - Define `type Tray struct {}` (if not already present as a struct)
  - Add constructor `func New() *Tray` that returns `&Tray{}`
  - Implement `func (t *Tray) UpdateIcon(sessionType string, state timer.State)` with logic:
    - Determine icon filename based on sessionType and state:
      - If `state == timer.StateRunning` and `sessionType == "work"` → `"icon-work.png"`
      - If `state == timer.StateRunning` and `sessionType == "short_break"` → `"icon-short-break.png"`
      - If `state == timer.StateRunning` and `sessionType == "long_break"` → `"icon-long-break.png"`
      - If `state == timer.StatePaused` → `"icon-paused.png"`
      - Otherwise (idle) → `"icon-idle.png"`
    - Load icon bytes from `assets/` directory (reuse existing `LoadIconFromAssets` pattern or similar)
    - Call `systray.SetIcon(iconData)`
  
- **Refactor – Clean up with tests green**:
  - Extract icon path resolution to helper function `getIconPath(sessionType string, state timer.State) string`
  - Simplify switch/case logic for readability

**CI / Checks**: `go test ./internal/tray/` (or manual test if systray integration is difficult to unit test)

---

### Step 7: Add cycle state fields to Window struct

- [ ] **Step 7: Add cycle state fields to Window struct**

**Workstream**: D  
**Based on Design**: Design §6 "UI Package State Changes"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - N/A (pure data structure change)
  
- **Green – Make the test(s) pass**:
  - Add three new fields to the `Window` struct:
    - `currentSessionType string` – tracks "work", "short_break", or "long_break"
    - `completedWorkSessions int` – tracks 0-4 completed work sessions in current cycle
    - `tray *tray.Tray` – reference to tray for icon updates
  - Initialize in `CreateWindow()`: `currentSessionType: "work"`, `completedWorkSessions: 0`
  
- **Refactor – Clean up with tests green**:
  - None needed

**CI / Checks**: `go build ./...` (verify compilation)

---

### Step 8: Add session type constants and duration mappings

- [ ] **Step 8: Add session type constants and duration mappings**

**Workstream**: D  
**Based on Design**: Design §6 "UI Package Internal Constants"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - N/A (constants definition)
  
- **Green – Make the test(s) pass**:
  - Add constants at package level (top of file, after imports):
    ```go
    const (
        sessionTypeWork       = "work"
        sessionTypeShortBreak = "short_break"
        sessionTypeLongBreak  = "long_break"
        
        durationWork       = 1500  // 25 minutes in seconds
        durationShortBreak = 300   // 5 minutes in seconds
        durationLongBreak  = 900   // 15 minutes in seconds
        
        sessionsPerCycle = 4
    )
    ```
  - Remove old constants: `sessionType = "work"` and `sessionDuration = 25`
  
- **Refactor – Clean up with tests green**:
  - Ensure naming is consistent with design document
  - Group related constants together

**CI / Checks**: `go build ./...`

---

### Step 9: Implement determineNextSessionType transition logic

- [ ] **Step 9: Implement determineNextSessionType transition logic**

**Workstream**: D  
**Based on Design**: Design §6 "UI Package State Changes – Transition logic"  
**Files**: `internal/ui/window.go`, `internal/ui/window_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestDetermineNextSessionType_AfterWork_ReturnsShortBreak()`:
    - Set `w.currentSessionType = sessionTypeWork`, `w.completedWorkSessions = 2`
    - Call `w.determineNextSessionType()`
    - Expect return value `sessionTypeShortBreak`
  - Write test `TestDetermineNextSessionType_AfterWork4_ReturnsLongBreak()`:
    - Set `w.currentSessionType = sessionTypeWork`, `w.completedWorkSessions = 4`
    - Expect return value `sessionTypeLongBreak`
  - Write test `TestDetermineNextSessionType_AfterShortBreak_ReturnsWork()`
  - Write test `TestDetermineNextSessionType_AfterLongBreak_ReturnsWorkAndResetsCounter()` (also verify counter becomes 0)
  - Tests will fail because method doesn't exist
  
- **Green – Make the test(s) pass**:
  - Implement `func (w *Window) determineNextSessionType() string`:
    ```go
    switch w.currentSessionType {
    case sessionTypeWork:
        if w.completedWorkSessions >= sessionsPerCycle {
            return sessionTypeLongBreak
        }
        return sessionTypeShortBreak
    case sessionTypeShortBreak:
        return sessionTypeWork
    case sessionTypeLongBreak:
        w.completedWorkSessions = 0  // Reset cycle
        return sessionTypeWork
    default:
        return sessionTypeWork
    }
    ```
  
- **Refactor – Clean up with tests green**:
  - Simplify conditional logic if possible
  - Ensure reset logic for long break is clear

**CI / Checks**: `go test ./internal/ui/`

---

### Step 10: Update handleTimerStarted to use new callback signature and log session type

- [ ] **Step 10: Update handleTimerStarted to use new callback signature and log session type**

**Workstream**: D  
**Based on Design**: Design §5 "Typical Flows – First work session"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Compilation will fail because timer's OnStarted callback signature changed to `func(string, int)`
  
- **Green – Make the test(s) pass**:
  - Change `handleTimerStarted()` signature to `handleTimerStarted(sessionType string, durationSeconds int)`
  - Update `storage.LogSession()` call to use `sessionType` parameter instead of hardcoded `sessionType` constant
  - Update header text based on session type:
    - If `sessionType == sessionTypeWork` → `w.header.SetTitle("Work Session")`
    - If `sessionType == sessionTypeShortBreak` → `w.header.SetTitle("Short Break")`
    - If `sessionType == sessionTypeLongBreak` → `w.header.SetTitle("Long Break")`
  - Update tray icon: `w.tray.UpdateIcon(sessionType, timer.StateRunning)` (if tray is set)
  
- **Refactor – Clean up with tests green**:
  - Extract header text logic to helper `getRunningHeaderText(sessionType string) string` for clarity
  - Add nil check for tray before calling UpdateIcon

**CI / Checks**: `go build ./...`, manual test starting a work session

---

### Step 11: Update handleTimerCompleted to increment cycle counter and transition

- [ ] **Step 11: Update handleTimerCompleted to increment cycle counter and transition**

**Workstream**: D  
**Based on Design**: Design §5 "Typical Flows – Work completion to break"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Manual test (integration with timer; difficult to unit test in lite mode)
  - Expected behavior: completing work session should transition to short break ready state
  
- **Green – Make the test(s) pass**:
  - In `handleTimerCompleted()`, after logging completion:
    - If `w.currentSessionType == sessionTypeWork`, increment `w.completedWorkSessions++`
    - Call `nextSessionType := w.determineNextSessionType()`
    - Set `w.currentSessionType = nextSessionType`
    - Update display for ready state:
      - Set header based on next session type (e.g., "Ready for Break", "Ready for Long Break", "Ready")
      - Set timer display based on next session duration (e.g., "05:00", "15:00", "25:00")
    - Update button states to idle
    - Update tray icon: `w.tray.UpdateIcon("", timer.StateIdle)`
  
- **Refactor – Clean up with tests green**:
  - Extract display update logic to `updateReadyStateDisplay(sessionType string)` helper that:
    - Sets appropriate header text ("Ready", "Ready for Break", "Ready for Long Break")
    - Sets timer display to formatted duration
    - Updates button states
  - Call this helper from handleTimerCompleted

**CI / Checks**: Manual test – complete a work session, verify transitions to short break ready state

---

### Step 12: Update handleStartClick to call timer.Start with session type and duration

- [ ] **Step 12: Update handleStartClick to call timer.Start with session type and duration**

**Workstream**: D  
**Based on Design**: Design §6 "Timer Package API Changes – Breaking change impact"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Compilation will fail because `timer.Start()` now requires parameters
  
- **Green – Make the test(s) pass**:
  - In `handleStartClick()`, change `w.timer.Start()` to:
    ```go
    w.timer.Start(w.currentSessionType, w.getSessionDuration())
    ```
  - Implement helper `func (w *Window) getSessionDuration() int`:
    ```go
    switch w.currentSessionType {
    case sessionTypeWork:
        return durationWork
    case sessionTypeShortBreak:
        return durationShortBreak
    case sessionTypeLongBreak:
        return durationLongBreak
    default:
        return durationWork
    }
    ```
  
- **Refactor – Clean up with tests green**:
  - None needed

**CI / Checks**: `go build ./...`, manual test starting different session types

---

### Step 13: Update handleResetClick to reset cycle state

- [ ] **Step 13: Update handleResetClick to reset cycle state**

**Workstream**: D  
**Based on Design**: Design §5 "Typical Flows – Reset from any state"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Manual test – reset during session 3 should return to session 1/4
  
- **Green – Make the test(s) pass**:
  - In `handleResetClick()`, after calling `timer.Reset()`:
    - Add `w.completedWorkSessions = 0`
    - Add `w.currentSessionType = sessionTypeWork`
    - Call `w.updateReadyStateDisplay(sessionTypeWork)` to reset display to "Ready", "25:00", session 1/4
  
- **Refactor – Clean up with tests green**:
  - Reuse the `updateReadyStateDisplay()` helper created in Step 11
  - Ensure all state is properly reset

**CI / Checks**: Manual test – reset from various states (work session 3, during break, etc.)

---

### Step 14: Add Skip button menu item

- [ ] **Step 14: Add Skip button menu item**

**Workstream**: E  
**Based on Design**: Design §4 "Component Responsibilities – UI new UI elements"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - N/A (UI element creation)
  
- **Green – Make the test(s) pass**:
  - In `InitializeMenu()`, after the Reset button:
    ```go
    systray.AddSeparator()
    w.btnSkip = systray.AddMenuItem("Skip", "Skip to next session")
    ```
  - Add `btnSkip *systray.MenuItem` field to Window struct
  
- **Refactor – Clean up with tests green**:
  - None needed

**CI / Checks**: `go build ./...`, manual visual check – Skip button appears in dropdown

---

### Step 15: Implement handleSkipClick logic

- [ ] **Step 15: Implement handleSkipClick logic**

**Workstream**: E  
**Based on Design**: Design §5 "Typical Flows – Skip during work session"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Manual test – Skip button won't respond yet
  
- **Green – Make the test(s) pass**:
  - Implement `func (w *Window) handleSkipClick()`:
    ```go
    for range w.btnSkip.ClickedCh {
        if w.timer == nil {
            continue
        }
        
        // Calculate elapsed time
        currentSessionDuration := w.getSessionDuration()
        elapsed := currentSessionDuration - w.timer.GetRemaining()
        elapsedMinutes := elapsed / 60
        
        // Stop timer
        currentSessionType := w.currentSessionType
        w.timer.Reset()
        
        // Log as skipped
        err := storage.LogSession(time.Now(), currentSessionType, "skipped", elapsedMinutes)
        if err != nil {
            log.Printf("[ERROR] Failed to log skipped session: %v", err)
        }
        
        // Determine next session (do NOT increment completedWorkSessions for skipped work)
        nextSessionType := w.determineNextSessionType()
        w.currentSessionType = nextSessionType
        
        // Update display
        w.updateReadyStateDisplay(nextSessionType)
        
        // Update tray icon
        if w.tray != nil {
            w.tray.UpdateIcon("", timer.StateIdle)
        }
    }
    ```
  
- **Refactor – Clean up with tests green**:
  - Extract elapsed calculation to helper `calculateElapsedMinutes() int` if it gets complex
  - Ensure error handling is consistent with other handlers

**CI / Checks**: Manual test – skip during work session, verify logs "skipped" and advances to break

---

### Step 16: Start Skip button click handler goroutine

- [ ] **Step 16: Start Skip button click handler goroutine**

**Workstream**: E  
**Based on Design**: Design pattern from existing click handlers  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Skip button won't respond to clicks (handleSkipClick not started)
  
- **Green – Make the test(s) pass**:
  - In `startClickHandlers()` method, add:
    ```go
    go w.handleSkipClick()
    ```
  
- **Refactor – Clean up with tests green**:
  - None needed

**CI / Checks**: Manual test – click Skip button, verify it responds

---

### Step 17: Add cycle indicator menu item

- [ ] **Step 17: Add cycle indicator menu item**

**Workstream**: E  
**Based on Design**: Design §6 "UI Package State Changes – cycle indicator"  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - N/A (UI element creation)
  
- **Green – Make the test(s) pass**:
  - In `InitializeMenu()`, after the timer display:
    ```go
    w.cycleIndicator = systray.AddMenuItem("Session 1/4  🍅○○○", "Cycle progress")
    w.cycleIndicator.Disable()
    ```
  - Add `cycleIndicator *systray.MenuItem` field to Window struct
  
- **Refactor – Clean up with tests green**:
  - None needed

**CI / Checks**: `go build ./...`, manual visual check – cycle indicator appears in dropdown

---

### Step 18: Implement formatCycleIndicator helper

- [ ] **Step 18: Implement formatCycleIndicator helper**

**Workstream**: E  
**Based on Design**: Acceptance Criteria #7 – Cycle indicator shows progress  
**Files**: `internal/ui/window.go`, `internal/ui/window_test.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Write test `TestFormatCycleIndicator_ShowsCorrectProgress()`:
    - With `completedWorkSessions = 0` and `currentSessionType = "work"` → expect `"Session 1/4  🍅○○○"`
    - With `completedWorkSessions = 1` → expect `"Session 2/4  🍅🍅○○"` (showing next session is 2)
    - With `completedWorkSessions = 3` → expect `"Session 4/4  🍅🍅🍅🍅"`
  - Test will fail because method doesn't exist
  
- **Green – Make the test(s) pass**:
  - Implement `func (w *Window) formatCycleIndicator() string`:
    - Calculate display session number based on `completedWorkSessions` and `currentSessionType`
    - If currently on work session, display session number is `completedWorkSessions + 1`
    - If currently on break, display session number reflects the work session just completed
    - Build tomato string: repeat "🍅" `completedWorkSessions` times, then "○" for remaining sessions
    - Return formatted string like `"Session 2/4  🍅🍅○○"`
  
- **Refactor – Clean up with tests green**:
  - Simplify string building (use `strings.Repeat()` if helpful)
  - Ensure edge cases handled (0 sessions, 4 sessions)

**CI / Checks**: `go test ./internal/ui/`

---

### Step 19: Update cycle indicator on state changes

- [ ] **Step 19: Update cycle indicator on state changes**

**Workstream**: E  
**Based on Design**: Cycle indicator updates throughout flows  
**Files**: `internal/ui/window.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Manual test – cycle indicator doesn't update as sessions progress
  
- **Green – Make the test(s) pass**:
  - Call `w.cycleIndicator.SetTitle(w.formatCycleIndicator())` in:
    - `handleTimerStarted()` – to show current running session
    - `handleTimerCompleted()` – after transition to next session
    - `handleSkipClick()` – after transition to next session
    - `handleResetClick()` – after resetting to session 1/4
    - `updateReadyStateDisplay()` – when preparing for next session
  
- **Refactor – Clean up with tests green**:
  - Consider adding cycle indicator update to `updateReadyStateDisplay()` helper to centralize
  - Ensure indicator is updated consistently across all state transitions

**CI / Checks**: Manual test – complete/skip through full cycle, verify indicator updates correctly

---

### Step 20: Wire tray reference to UI and call UpdateIcon

- [ ] **Step 20: Wire tray reference to UI and call UpdateIcon**

**Workstream**: E  
**Based on Design**: Design §4 "Component Responsibilities – UI extended responsibilities"  
**Files**: `internal/ui/window.go`, `cmd/gopomodoro/main.go`

**TDD Cycle**:
- **Red – Failing test first**:
  - Manual test – tray icon doesn't change based on session type
  
- **Green – Make the test(s) pass**:
  - Add `func (w *Window) SetTray(t *tray.Tray)` method to Window:
    ```go
    func (w *Window) SetTray(t *tray.Tray) {
        w.tray = t
    }
    ```
  - In `cmd/gopomodoro/main.go`, after creating tray and window instances, call:
    ```go
    window.SetTray(trayInstance)
    ```
  - Verify all UI event handlers call `w.tray.UpdateIcon()` at appropriate points:
    - `handleTimerStarted()` – already added in Step 10
    - `handleTimerCompleted()` – already added in Step 11
    - `handleSkipClick()` – already added in Step 15
    - `handlePauseClick()` – add `w.tray.UpdateIcon(w.currentSessionType, timer.StatePaused)`
    - `handleResetClick()` – add `w.tray.UpdateIcon("", timer.StateIdle)`
  
- **Refactor – Clean up with tests green**:
  - Ensure all tray calls have nil checks: `if w.tray != nil { w.tray.UpdateIcon(...) }`
  - Verify icon changes are visible and correct for each session type

**CI / Checks**: `go build ./...`, manual test – start work/break sessions, verify tray icon changes color/emoji

---

## 3. Rollout & Validation Notes

### Suggested Grouping into PRs

- **PR 1: Timer Session Type Support** (Steps 1-4)
  - Timer package becomes session-type-aware
  - Breaking API change isolated to timer package
  - Estimated time: 30-60 minutes
  - Verification: `go test ./internal/timer/` all pass

- **PR 2: Assets and Tray Icon Updates** (Steps 5-6)
  - Emoji-based icons created
  - Tray UpdateIcon method implemented
  - Estimated time: 30 minutes
  - Verification: Manual check icons load and display

- **PR 3: UI Cycle Orchestration Core** (Steps 7-13)
  - Cycle state tracking and transition logic
  - UI handlers updated for new timer API
  - Largest change set
  - Estimated time: 90-120 minutes
  - Verification: Manual test full work session → short break → work cycle

- **PR 4: Skip Button and Cycle Indicator** (Steps 14-19)
  - New UI controls added
  - Skip functionality implemented
  - Cycle progress visualization
  - Estimated time: 60 minutes
  - Verification: Manual test skip scenarios, verify cycle indicator updates

- **PR 5: Final Integration** (Step 20)
  - Wire tray to UI for icon updates
  - End-to-end integration complete
  - Estimated time: 15 minutes
  - Verification: Manual test all session types show correct tray icons

### Validation Checkpoints

**After Step 4 (Timer changes)**:
- Unit tests pass for timer package
- Timer can be started with different session types and durations
- OnStarted callback receives session context

**After Step 13 (UI cycle core)**:
- Can start and complete a work session
- Work completion transitions to short break ready state
- Completing 4 work sessions transitions to long break
- Long break completion resets cycle to session 1
- Reset button returns to session 1/4 from any state

**After Step 19 (Controls added)**:
- Skip button appears in menu
- Skip during work logs "skipped" and advances to break
- Skip during break advances to next work session
- Cycle indicator shows correct session number and tomato progress

**After Step 20 (Final integration)**:
- Tray icon shows red tomato during work sessions
- Tray icon shows coffee/green during short breaks
- Tray icon shows star/blue during long breaks
- Tray icon shows paused state when paused
- Tray icon shows idle/gray when between sessions

### Manual Testing Scenarios

1. **Full cycle walkthrough**:
   - Start work session 1 → verify timer counts down, tray is red, indicator shows 1/4
   - Complete → verify transitions to short break, timer shows 05:00
   - Start short break → complete → verify transitions to work session 2, indicator shows 2/4
   - Repeat for sessions 3 and 4
   - Complete session 4 → verify long break appears (15:00)
   - Complete long break → verify resets to session 1/4

2. **Skip scenarios**:
   - Start work session → wait 10 minutes → Skip → verify logs "work,skipped,10"
   - Verify does NOT increment cycle counter
   - Start short break → Skip after 2 minutes → verify advances to next work session

3. **Reset scenarios**:
   - During work session 3 → Reset → verify returns to session 1/4, 25:00, "Ready"
   - During long break → Reset → verify returns to session 1/4

4. **Session log verification**:
   - After testing, inspect `~/.gopomodoro/sessions.log`
   - Verify contains entries for work, short_break, long_break
   - Verify skipped entries show correct elapsed time
   - Verify completed entries show full duration

### Success Criteria

Implementation is complete when:
- ✅ All unit tests pass (`go test ./...`)
- ✅ Application builds without errors
- ✅ User can complete a full 4-session pomodoro cycle with breaks
- ✅ Skip button advances to next session without counting toward cycle
- ✅ Cycle indicator shows correct session progress (1/4 through 4/4)
- ✅ Tray icon visually differentiates work/short break/long break
- ✅ Session log contains all three session types with correct durations
- ✅ Reset button returns to cycle start from any state
