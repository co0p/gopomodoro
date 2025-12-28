# Implement: Tray Time Display with Session Icons

## Context

This increment enables users to see remaining session time and session type directly in the macOS menu bar without opening the dropdown menu. The tray will display formatted text like "🍅 24m" alongside the icon, and redundant state/timer menu items will be removed from the dropdown.

**Links:**
- [increment.md](increment.md) – Product requirements and acceptance criteria
- [design.md](design.md) – Technical design and architecture
- [../../CONSTITUTION.md](../../CONSTITUTION.md) – Project principles (mode: lite)

**Key Design Decisions:**
- Add `UpdateDisplay(sessionType, state, remainingSeconds)` method to tray package
- Use `systray.SetTitle()` to display emoji + time string in menu bar
- Change idle state from neutral icon to tomato (🍅)
- Remove `header` and `timerDisplay` menu items from dropdown
- Time format: "Xm" (e.g., "25m", "5m")

**Constitution Constraints (lite mode):**
- Keep steps small and testable
- Prefer simple solutions over complex abstractions
- Manual testing is acceptable for UI/integration

**Status:** Not started  
**Next step:** Step 1 – Add emoji selection helper

## 1. Workstreams

- **Workstream A** – Tray Package Extensions (display capabilities)
- **Workstream B** – UI Package Simplification (remove redundant items, update tray calls)
- **Workstream C** – Main Package Initialization (startup tray title)
- **Workstream D** – Timer Package Verification (timing accuracy test)

## 2. Steps

### Step 1: Add emoji selection helper to tray package

**Workstream:** A  
**Based on Design:** Design §4 "Contracts and Data – Emoji selection logic"  
**Files:** `internal/tray/tray.go`, `internal/tray/tray_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Add `TestGetEmojiForState` in `tray_test.go`
- Test cases for each state/session combination:
  - Idle state → expect "🍅"
  - Running work → expect "🍅"
  - Running short break → expect "☕"
  - Running long break → expect "🌟"
  - Paused (any session) → expect "⏸️"
- Run `go test ./internal/tray` – tests fail (function doesn't exist)

**Green – Make the test(s) pass:**
- Add `getEmojiForState(sessionType string, state timer.State) string` function to `tray.go`
- Implement switch logic based on state and sessionType
- Return appropriate emoji string for each case
- Run `go test ./internal/tray` – tests pass

**Refactor – Clean up with tests green:**
- Extract emoji constants if mapping table feels repetitive
- Ensure function is clear and maintainable
- Verify tests still pass

**CI / Checks:**
- `go test ./internal/tray`
- `go build ./...`

---

### Step 2: Add time formatting helper to tray package

**Workstream:** A  
**Based on Design:** Design §4 "Contracts and Data – Time formatting"  
**Files:** `internal/tray/tray.go`, `internal/tray/tray_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Add `TestFormatMinutes` in `tray_test.go`
- Test cases:
  - 1500 seconds → expect "25m"
  - 300 seconds → expect "5m"
  - 900 seconds → expect "15m"
  - 60 seconds → expect "1m"
  - 30 seconds → expect "0m"
  - 0 seconds → expect "0m"
- Run `go test ./internal/tray` – tests fail (function doesn't exist)

**Green – Make the test(s) pass:**
- Add `formatMinutes(seconds int) string` function to `tray.go`
- Convert seconds to minutes: `minutes := seconds / 60`
- Return `fmt.Sprintf("%dm", minutes)`
- Run `go test ./internal/tray` – tests pass

**Refactor – Clean up with tests green:**
- No significant refactoring expected for this simple helper
- Verify edge cases (negative values if applicable)
- Ensure tests still pass

**CI / Checks:**
- `go test ./internal/tray`
- `go build ./...`

---

### Step 3: Add UpdateDisplay method to tray package

**Workstream:** A  
**Based on Design:** Design §4 "Contracts and Data – Tray Package New Method"  
**Files:** `internal/tray/tray.go`, `internal/tray/tray_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Add `TestUpdateDisplay` in `tray_test.go`
- Test that method combines emoji + formatted time correctly:
  - Call `UpdateDisplay("work", StateRunning, 1500)`
  - Verify it would call `systray.SetTitle("🍅 25m")` (may need to mock or test indirectly)
  - Verify it selects correct icon data
- Since systray is external, consider testing the string formatting logic separately
- Run `go test ./internal/tray` – tests fail (method doesn't exist)

**Green – Make the test(s) pass:**
- Add `UpdateDisplay(sessionType string, state timer.State, remainingSeconds int)` method to `Tray` struct
- Call `getEmojiForState(sessionType, state)` from Step 1
- Call `formatMinutes(remainingSeconds)` from Step 2
- Combine: `title := emoji + " " + timeStr`
- Call `systray.SetTitle(title)`
- Call existing `getIconData()` and `systray.SetIcon()`
- Run `go test ./internal/tray` – tests pass

**Refactor – Clean up with tests green:**
- Extract title formatting to `formatDisplayText(emoji, remainingSeconds)` helper if clearer
- Ensure method is idempotent and safe to call repeatedly
- Verify tests still pass

**CI / Checks:**
- `go test ./internal/tray`
- `go build ./...`

---

### Step 4: Change idle state to use tomato icon

**Workstream:** A  
**Based on Design:** Design §4 "Contracts and Data – Emoji selection logic table"  
**Files:** `internal/tray/tray.go`, `internal/tray/tray_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Update `TestGetIconData` (or add if missing) in `tray_test.go`
- Change expectation: `getIconData("", StateIdle)` should return `iconWork` (not `iconIdle`)
- Run `go test ./internal/tray` – test fails (currently returns iconIdle)

**Green – Make the test(s) pass:**
- Modify `getIconData()` in `tray.go`
- Change the default/idle case to return `iconWork` instead of `iconIdle`
- Run `go test ./internal/tray` – test passes

**Refactor – Clean up with tests green:**
- Consider if `iconIdle` is still needed (may remove in follow-up, acceptable to leave for now)
- Verify emoji selection in Step 1 also returns "🍅" for idle state
- Ensure tests still pass

**CI / Checks:**
- `go test ./internal/tray`
- `go build ./...`

---

### Step 5: Update formatTime to use "m" suffix in UI package

**Workstream:** B  
**Based on Design:** Design §4 "Contracts and Data – UI Package Modified Helper"  
**Files:** `internal/ui/window.go`, `internal/ui/window_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Add or update `TestFormatTime` in `window_test.go`
- Test cases expecting "m" suffix:
  - 1500 seconds → expect "25m" (not "25min")
  - 300 seconds → expect "5m"
  - 0 seconds → expect "0m"
- Run `go test ./internal/ui` – tests fail (currently returns "Xmin")

**Green – Make the test(s) pass:**
- Modify `formatTime()` in `window.go`
- Change `fmt.Sprintf("%dmin", minutes)` to `fmt.Sprintf("%dm", minutes)`
- Run `go test ./internal/ui` – tests pass

**Refactor – Clean up with tests green:**
- No significant refactoring expected
- This function is now only used for progress bar (timerDisplay will be removed in Step 8)
- Ensure tests still pass

**CI / Checks:**
- `go test ./internal/ui`
- `go build ./...`

---

### Step 6: Update updateTrayIcon signature and all call sites

**Workstream:** B  
**Based on Design:** Design §4 "Contracts and Data – UI Package Modified Method"  
**Files:** `internal/ui/window.go`

**TDD Cycle:**

**Red – Failing test first:**
- Identify all call sites of `updateTrayIcon()` in `window.go`
- Update method signature to accept `remainingSeconds int` parameter
- Compilation will fail at all call sites (they don't pass the new parameter)
- Run `go build ./...` – build fails

**Green – Make the test(s) pass:**
- Update `updateTrayIcon(state timer.State, remainingSeconds int)` signature
- Change implementation to call `tray.UpdateDisplay(sessionType, state, remainingSeconds)` instead of `tray.UpdateIcon(sessionType, state)`
- Update all call sites to pass `w.timer.GetRemaining()`:
  - In `UpdateButtonStates()` – pass `w.timer.GetRemaining()`
  - Handle StateIdle case: use session default duration (1500 for work, 300 for short break, 900 for long break)
- Run `go build ./...` – build succeeds
- Run `go test ./internal/ui` – tests pass

**Refactor – Clean up with tests green:**
- Ensure all call sites are consistent
- Consider extracting "get remaining time for display" logic if repetitive
- Verify compilation and tests still pass

**CI / Checks:**
- `go test ./internal/ui`
- `go build ./...`

---

### Step 7: Remove header menu item field and usage

**Workstream:** B  
**Based on Design:** Design §2 "Scope and Non-Scope – UI simplification"  
**Files:** `internal/ui/window.go`

**TDD Cycle:**

**Red – Failing test first:**
- Remove `header *systray.MenuItem` field from `Window` struct
- Compilation fails at all locations that reference `w.header`
- Run `go build ./...` – build fails

**Green – Make the test(s) pass:**
- Remove all `w.header.SetTitle()` calls throughout the file:
  - In `InitializeMenu()` – remove header initialization
  - In `handleStartClick()` – remove header updates
  - In `handlePauseClick()` – remove header update
  - In `handleResetClick()` – remove header update
  - In `handleSkipClick()` – remove header updates
  - In `handleTimerStarted()` – remove header update
  - In `handleTimerCompleted()` – remove header updates
- Run `go build ./...` – build succeeds

**Refactor – Clean up with tests green:**
- Verify no other references to `header` remain
- Ensure button states and other menu items still work correctly
- Manual test: Open dropdown, verify progress bar is now first item (no header above it)

**CI / Checks:**
- `go test ./internal/ui`
- `go build ./...`
- Manual verification: Launch app, open dropdown, verify no state header

---

### Step 8: Remove timerDisplay menu item field and usage

**Workstream:** B  
**Based on Design:** Design §2 "Scope and Non-Scope – UI simplification"  
**Files:** `internal/ui/window.go`

**TDD Cycle:**

**Red – Failing test first:**
- Remove `timerDisplay *systray.MenuItem` field from `Window` struct
- Compilation fails at all locations that reference `w.timerDisplay`
- Run `go build ./...` – build fails

**Green – Make the test(s) pass:**
- Remove all `w.timerDisplay.SetTitle()` calls throughout the file:
  - In `InitializeMenu()` – remove timerDisplay initialization
  - In `handleResetClick()` – remove timerDisplay update
  - In `handleSkipClick()` – remove timerDisplay update
  - In `handleTimerTick()` – remove timerDisplay update
  - In `handleTimerCompleted()` – remove timerDisplay update
- Run `go build ./...` – build succeeds

**Refactor – Clean up with tests green:**
- Verify no other references to `timerDisplay` remain
- Note: `formatTime()` is still used for progress bar, so keep it
- Manual test: Open dropdown, verify no timer countdown item

**CI / Checks:**
- `go test ./internal/ui`
- `go build ./...`
- Manual verification: Launch app, open dropdown, verify no timer display

---

### Step 9: Set initial tray title on startup

**Workstream:** C  
**Based on Design:** Design §2 "Proposed Solution – Main Package"  
**Files:** `cmd/gopomodoro/main.go`

**TDD Cycle:**

**Red – Failing test first:**
- Manual observation: Currently tray shows only icon, no text on startup
- Run app and observe menu bar

**Green – Make the test(s) pass:**
- In `onReady()` function in `main.go`
- After `systray.SetTooltip("GoPomodoro")` line
- Add `systray.SetTitle("🍅 25m")`
- Build and run app
- Observe tray displays "🍅 25m" on startup

**Refactor – Clean up with tests green:**
- Verify initial title matches idle work session state
- Consider if this should use the new tray package method instead (acceptable either way for lite mode)
- Ensure app launches correctly

**CI / Checks:**
- `go build ./...`
- Manual verification: Launch app, verify tray shows "🍅 25m" immediately

---

### Step 10: Add timing accuracy test

**Workstream:** D  
**Based on Design:** Design §5 "Testing Strategy – Timer Package"  
**Files:** `internal/timer/timer_test.go`

**TDD Cycle:**

**Red – Failing test first:**
- Add `TestTimingAccuracy` in `timer_test.go`
- Create 60-second timer and measure real elapsed time:
  ```go
  func TestTimingAccuracy(t *testing.T) {
      tmr := New()
      completed := make(chan time.Time, 1)
      tmr.OnCompleted(func() {
          completed <- time.Now()
      })
      
      startTime := time.Now()
      tmr.Start("work", 60)
      
      completionTime := <-completed
      elapsed := completionTime.Sub(startTime).Seconds()
      
      // Should complete within 60 ± 2 seconds
      if elapsed < 58 || elapsed > 62 {
          t.Errorf("Expected ~60s, got %.2fs", elapsed)
      }
  }
  ```
- Run `go test ./internal/timer` – test may pass or fail depending on actual timing accuracy
- This test helps diagnose the user-reported perception that timers feel slow

**Green – Make the test(s) pass:**
- If test fails (timer takes significantly longer than 60 seconds):
  - Investigate timer tick logic and countdown algorithm
  - Adjust if needed (out of scope for this increment, but test provides data)
- If test passes:
  - Timing is accurate; perceived slowness may be UX issue (10-second update interval)
  - Document findings in test comments

**Refactor – Clean up with tests green:**
- Ensure test is reliable (not flaky due to system scheduling)
- Consider tolerance range (±2 seconds is reasonable)
- May want to test longer durations (5 minutes) for more realistic validation

**CI / Checks:**
- `go test ./internal/timer -timeout 90s` (allow time for 60-second test)
- Document test results to inform follow-up work on timing improvements

---

## 3. Rollout & Validation Notes

### Suggested Grouping into PRs

**PR 1: Tray display capabilities (Steps 1-4)**
- Self-contained tray package changes
- All tests passing
- Can be merged without affecting UI

**PR 2: UI simplification and integration (Steps 5-8)**
- Depends on PR 1
- Removes redundant menu items and wires new tray display
- All tests passing, manual UI verification

**PR 3: Final integration (Step 9)**
- Depends on PR 2
- Completes user-facing feature
- Manual verification of startup behavior

**PR 4: Timing verification (Step 10)**
- Independent of other PRs
- Can be merged anytime
- Provides diagnostic data for follow-up work

### Suggested Validation Checkpoints

**After Step 4:**
- `go test ./internal/tray` passes
- `go build ./...` succeeds
- Tray package ready for integration

**After Step 8:**
- `go test ./...` passes
- `go build ./...` succeeds
- Manual test: Launch app, observe tray updates during timer cycle
- Verify dropdown menu is cleaner (no header, no timer display)
- Verify progress bar and cycle indicator still work

**After Step 9:**
- Manual test: Fresh launch shows "🍅 25m" in tray
- Click Start: Tray updates to "🍅 24m" after first tick
- Pause: Tray shows "⏸️ Xm"
- Resume: Tray shows "🍅 Xm"
- Complete work session: Tray shows "☕ 5m"
- Complete 4 work sessions: Tray shows "🌟 15m"
- Reset: Tray returns to "🍅 25m"

**After Step 10:**
- Use stopwatch to time a 5-minute break
- Verify it takes approximately 5 real minutes (±6 seconds)
- Document any timing drift observed
- Consider follow-up increment for 1-second ticks if needed

### Manual Integration Test Flow

1. Build: `make build`
2. Run: `./bin/gopomodoro`
3. Verify tray shows "🍅 25m" on startup
4. Open dropdown: Verify progress bar is first item (no header)
5. Verify no timer display in dropdown
6. Click Start: Observe tray updates to "🍅 24m" after ~10 seconds
7. Let timer run for 1 minute: Verify tray shows decreasing minutes
8. Click Pause: Verify tray shows "⏸️ Xm" and stops updating
9. Click Start (Resume): Verify tray shows "🍅 Xm" and resumes
10. Let work session complete: Verify tray shows "☕ 5m"
11. Time break with stopwatch: Verify 5-minute break takes ~5 real minutes
12. Complete 4 work sessions: Verify tray shows "🌟 15m" for long break
13. Click Reset: Verify tray returns to "🍅 25m"
14. Verify all buttons, progress bar, and cycle indicator still functional

### Acceptance Criteria Verification

- ✓ Tomato icon at startup (AC 1)
- ✓ Session icon and time visible in tray during sessions (AC 2)
- ✓ Session icon and time visible in tray when idle (AC 3)
- ✓ Paused state shows pause icon (AC 4)
- ✓ No state header in dropdown (AC 5)
- ✓ No timer display in dropdown (AC 6)
- ✓ Progress bar still visible (AC 7)
- ✓ Cycle indicator still visible (AC 8)
