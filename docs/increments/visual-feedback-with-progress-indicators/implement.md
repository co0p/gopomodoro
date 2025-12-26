# Implement: Visual Feedback with Progress Indicators

## Context

This increment transforms GoPomodoro's UI from text-only to visually scannable by adding progress bars, emoji indicators, and optimized update frequency.

**Increment Goal:**
- Add Unicode-based progress bar showing session advancement
- Simplify timer display to minutes-only format
- Reduce update frequency from 1s to 10s for CPU efficiency
- Add session-specific emoji to headers (🍅 ☕ 🌟 ⏸️)
- Verify tray icon state changes work correctly

**Design Approach:**
All changes are presentational, confined to UI and Timer packages. No modifications to timer state machine, session cycle logic, or data persistence.

**Key Constraints from Constitution:**
- `lite` mode: Keep implementation pragmatic, manual testing for visual changes
- Small, safe steps that can be tested and committed independently
- Simple solutions over complex abstractions

**References:**
- [increment.md](./increment.md)
- [design.md](./design.md)
- [CONSTITUTION.md](../../../CONSTITUTION.md)

**Status:** Not started  
**Next step:** Step 1 – Change timer display to minutes-only format

---

## 1. Workstreams

- **Workstream A** – Display Formatting (UI Package)
- **Workstream B** – Timer Frequency Optimization (Timer Package)
- **Workstream C** – Integration and Verification

---

## 2. Steps

### Step 1: Change timer display to minutes-only format

**Workstream:** A  
**Based on Design:** Design §4.1 Minutes-Only Timer Display

**Files:** `internal/ui/window.go`, `internal/ui/window_test.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Create `window_test.go` if it doesn't exist
  - Add test cases for `formatTime()`:
    - `formatTime(1500)` expects `"25min"` (currently returns `"25:00"`)
    - `formatTime(720)` expects `"12min"` (currently returns `"12:00"`)
    - `formatTime(60)` expects `"1min"` (currently returns `"01:00"`)
    - `formatTime(59)` expects `"0min"` (currently returns `"00:59"`)
    - `formatTime(0)` expects `"0min"` (currently returns `"00:00"`)
  - Run tests, observe failures showing current MM:SS format

- **Green – Make the test(s) pass:**
  - Modify `formatTime(seconds int) string` in `window.go`:
    - Calculate `minutes := seconds / 60`
    - Return `fmt.Sprintf("%dmin", minutes)`
  - Run tests, verify all pass

- **Refactor – Clean up with tests green:**
  - Remove unused `secs` variable calculation
  - Verify function signature remains unchanged (backward compatible)
  - Update any comments to reflect new format

**CI / Checks:**
- `go test ./internal/ui/...`
- `go build ./...`

---

### Step 2: Add progress bar menu item and tracking fields

**Workstream:** A  
**Based on Design:** Design §4.1 Progress Bar Visualization

**Files:** `internal/ui/window.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Manual verification: Open dropdown, confirm no progress bar menu item exists between timer display and cycle indicator

- **Green – Make the test(s) pass:**
  - Add fields to `Window` struct:
    - `progressBar *systray.MenuItem`
    - `sessionStartTime time.Time`
    - `sessionDuration int`
  - In `InitializeMenu()`, after `timerDisplay` and before `cycleIndicator`:
    - `w.progressBar = systray.AddMenuItem("○○○○○○○○○○", "Session progress")`
    - `w.progressBar.Disable()`
  - Initialize tracking fields to zero values in `CreateWindow()`

- **Refactor – Clean up with tests green:**
  - Verify menu item order is correct: header, timer, **progress bar**, cycle indicator, separator, buttons
  - Ensure field initialization is consistent with other menu items

**CI / Checks:**
- `go build ./...`
- Manual visual verification: Open dropdown, see progress bar menu item

---

### Step 3: Implement progress bar calculation and rendering

**Workstream:** A  
**Based on Design:** Design §5 Contracts and Data - Progress Bar Logic

**Files:** `internal/ui/window.go`, `internal/ui/window_test.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Add test cases for new `formatProgressBar(elapsed, duration int) string` function:
    - `formatProgressBar(0, 1500)` expects `"○○○○○○○○○○"` (0% progress)
    - `formatProgressBar(150, 1500)` expects `"●○○○○○○○○○"` (10% progress)
    - `formatProgressBar(750, 1500)` expects `"●●●●●○○○○○"` (50% progress)
    - `formatProgressBar(1500, 1500)` expects `"●●●●●●●●●●"` (100% progress)
    - `formatProgressBar(1600, 1500)` expects `"●●●●●●●●●●"` (over 100% should cap at 10)
  - Run tests, observe failures (function doesn't exist)

- **Green – Make the test(s) pass:**
  - Implement `formatProgressBar(elapsed, duration int) string`:
    ```go
    func formatProgressBar(elapsed, duration int) string {
        if duration == 0 {
            return "○○○○○○○○○○"
        }
        fillPercentage := float64(elapsed) / float64(duration)
        filledSegments := int(fillPercentage * 10)
        if filledSegments > 10 {
            filledSegments = 10
        }
        
        result := ""
        for i := 0; i < filledSegments; i++ {
            result += "●"
        }
        for i := filledSegments; i < 10; i++ {
            result += "○"
        }
        return result
    }
    ```
  - Run tests, verify all pass

- **Refactor – Clean up with tests green:**
  - Extract constants for segment count (10) and circle characters if helpful
  - Consider rounding logic (currently truncates) - test edge cases
  - Add comments explaining calculation

**CI / Checks:**
- `go test ./internal/ui/...`

---

### Step 4: Update progress bar on timer ticks

**Workstream:** A  
**Based on Design:** Design §4.1 Progress Bar Visualization

**Files:** `internal/ui/window.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Manual verification: Start a session, observe that progress bar remains empty despite timer ticking

- **Green – Make the test(s) pass:**
  - Modify `handleTimerTick(remaining int)`:
    - Calculate `elapsed := w.sessionDuration - remaining`
    - Generate progress bar: `progressStr := formatProgressBar(elapsed, w.sessionDuration)`
    - Update menu item: `w.progressBar.SetTitle(progressStr)`
  - Verify progress bar now updates alongside timer display

- **Refactor – Clean up with tests green:**
  - Ensure elapsed calculation handles edge cases (remaining = 0, negative values)
  - Verify progress bar updates smoothly during session
  - Check that progress bar doesn't flicker or show inconsistent state

**CI / Checks:**
- `go test ./internal/ui/...`
- Manual run: Start session, verify progress bar fills over time

---

### Step 5: Track session start time and duration

**Workstream:** A  
**Based on Design:** Design §4.1 Progress Bar Visualization - needs elapsed time calculation

**Files:** `internal/ui/window.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Manual verification: Progress bar calculation relies on `w.sessionDuration` which is currently uninitialized (0)
  - Start a session, observe progress bar might not work correctly or show divide-by-zero behavior

- **Green – Make the test(s) pass:**
  - In `handleTimerStarted(sessionType string, durationSeconds int)`:
    - Store session metadata: `w.sessionStartTime = time.Now()`
    - Store duration: `w.sessionDuration = durationSeconds`
  - In `handleTimerCompleted()`:
    - Reset tracking: `w.sessionDuration = 0`
  - In `handleResetClick()`:
    - Reset tracking: `w.sessionDuration = 0`
    - Reset progress bar: `w.progressBar.SetTitle("○○○○○○○○○○")`
  - In `handleSkipClick()`:
    - Reset progress bar: `w.progressBar.SetTitle("○○○○○○○○○○")`

- **Refactor – Clean up with tests green:**
  - Verify all state transitions properly maintain timing fields
  - Ensure progress bar resets to empty on new session start
  - Check pause/resume behavior maintains progress correctly

**CI / Checks:**
- `go test ./internal/ui/...`
- Manual verification: Start session, pause, resume, reset, skip - verify progress bar behaves correctly

---

### Step 6: Add session-specific emoji to headers

**Workstream:** A  
**Based on Design:** Design §4.4 Session-Specific Header Emoji

**Files:** `internal/ui/window.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Manual verification: Start work session, observe header shows "Work Session" without emoji
  - Start break session, observe header shows "Short Break" without emoji

- **Green – Make the test(s) pass:**
  - Update `handleTimerStarted()` header text:
    - Work session: `w.header.SetTitle("🍅 Work Session")`
    - Short break: `w.header.SetTitle("☕ Short Break")`
    - Long break: `w.header.SetTitle("🌟 Long Break")`
  - Update `handleTimerCompleted()` and `handleSkipClick()` for "Ready for..." states:
    - Work: `w.header.SetTitle("🍅 Ready for Work")`
    - Short break: `w.header.SetTitle("☕ Ready for Break")`
    - Long break: `w.header.SetTitle("🌟 Ready for Long Break")`

- **Refactor – Clean up with tests green:**
  - Consider extracting emoji constants if headers are updated in multiple places
  - Ensure consistency between running and ready states
  - Verify emoji render correctly on macOS

**CI / Checks:**
- Manual visual verification: Start sessions of each type, verify emoji appear correctly

---

### Step 7: Add paused state emoji and visual distinction

**Workstream:** A  
**Based on Design:** Design §4.6 Paused State Visual Distinction

**Files:** `internal/ui/window.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Manual verification: Start a session, click pause, observe header shows "Paused" without emoji

- **Green – Make the test(s) pass:**
  - In `handlePauseClick()`:
    - Change header to: `w.header.SetTitle("⏸️ Paused")`
  - In `handleStartClick()` resume branch (when `state == timer.StatePaused`):
    - Restore session-specific emoji based on `w.session.CurrentType`:
      ```go
      switch w.session.CurrentType {
      case session.TypeWork:
          w.header.SetTitle("🍅 Work Session")
      case session.TypeShortBreak:
          w.header.SetTitle("☕ Short Break")
      case session.TypeLongBreak:
          w.header.SetTitle("🌟 Long Break")
      default:
          w.header.SetTitle("Running")
      }
      ```

- **Refactor – Clean up with tests green:**
  - Ensure paused → resumed header restoration is correct
  - Test multiple pause/resume cycles
  - Verify progress bar maintains state during pause

**CI / Checks:**
- Manual pause/resume testing across all session types
- Verify header emoji consistency

---

### Step 8: Change tick interval to 10 seconds

**Workstream:** B  
**Based on Design:** Design §4.3 10-Second Update Frequency

**Files:** `internal/timer/timer.go`, potentially `internal/timer/timer_test.go`

**TDD Cycle:**

- **Red – Failing test first:**
  - Review existing timer tests in `timer_test.go`
  - Check if any tests rely on 1-second precision (might fail with 10s interval)
  - If timing-sensitive tests exist, note they will need adjustment

- **Green – Make the test(s) pass:**
  - In `timer.go`, change constant:
    - From: `tickInterval = 1 * time.Second`
    - To: `tickInterval = 10 * time.Second`
  - Run tests: `go test ./internal/timer/...`
  - If any tests fail due to timing expectations, adjust test delays/assertions to accommodate 10-second intervals

- **Refactor – Clean up with tests green:**
  - Review timer behavior: verify callbacks still fire correctly
  - Check that 10-second granularity is acceptable for 25-minute sessions
  - Document the tick interval choice in comments if helpful

**CI / Checks:**
- `go test ./internal/timer/...`
- `go test ./...`
- Manual CPU observation with Activity Monitor (verify lower CPU usage during active sessions)

---

### Step 9: Verify tray icon updates work correctly

**Workstream:** C  
**Based on Design:** Design §4.5 Tray Icon Updates

**Files:** Verification only, no code changes expected

**TDD Cycle:**

- **Red – Failing test first:**
  - N/A (verification step)

- **Green – Make the test(s) pass:**
  - Manual verification checklist:
    - Confirm `updateTrayIcon()` is called in `UpdateButtonStates()`
    - Verify icon assets exist in `assets/` directory:
      - `icon-idle.png`
      - `icon-work.png`
      - `icon-short-break.png`
      - `icon-long-break.png`
      - `icon-paused.png`
    - Start work session: verify work icon appears in menu bar
    - Start short break: verify break icon appears
    - Pause session: verify paused icon appears
    - Resume session: verify icon returns to appropriate session icon

- **Refactor – Clean up with tests green:**
  - N/A

**CI / Checks:**
- Manual visual verification in macOS menu bar
- Test all state transitions: idle → work → pause → resume → break → long break

---

### Step 10: End-to-end visual and functional verification

**Workstream:** C  
**Based on Design:** Design §6 Testing and Safety Net - Manual Visual Verification

**Files:** N/A (testing only)

**TDD Cycle:**

- **Red – Failing test first:**
  - N/A (validation step)

- **Green – Make the test(s) pass:**
  - **Complete pomodoro cycle test:**
    - Start work session (25 min)
    - Verify: header shows "🍅 Work Session"
    - Verify: timer shows "25min" and decrements to "24min", "23min", etc. every 10 seconds
    - Verify: progress bar starts empty "○○○○○○○○○○" and fills: "●○○○○○○○○○", "●●○○○○○○○○", etc.
    - Verify: work icon appears in menu bar
    - Pause mid-session
    - Verify: header shows "⏸️ Paused"
    - Verify: timer and progress bar freeze
    - Verify: paused icon appears
    - Resume
    - Verify: header restores to "🍅 Work Session"
    - Verify: timer and progress bar continue from paused position
    - Complete work session or skip
    - Verify: auto-transition to short break with "☕ Short Break" header
    - Complete 4 work sessions
    - Verify: long break shows "🌟 Long Break"
  
  - **Edge cases:**
    - Timer showing "0min" during final minute (< 60 seconds remaining)
    - Progress bar at 100% shows exactly 10 filled segments: "●●●●●●●●●●"
    - Rapid skip clicks update all visuals correctly
    - Reset button clears progress bar and resets timer
  
  - **No regression:**
    - All buttons (Start, Pause, Reset, Skip, Quit) work correctly
    - Session cycle logic continues: 4 work → long break
    - Session logging still records events (check CSV files)

- **Refactor – Clean up with tests green:**
  - Fix any visual glitches discovered during testing
  - Adjust emoji or Unicode characters if rendering issues appear
  - Fine-tune progress bar segment calculation if jumps are too jarring

**CI / Checks:**
- Complete manual test scenarios from [increment.md](./increment.md)
- CPU usage check with Activity Monitor
- Visual consistency across all states

---

## 3. Rollout & Validation Notes

**Suggested grouping into PRs:**

- **PR 1:** Steps 1-5 (Timer format + Complete progress bar feature)
  - After Step 5: Verify timer shows minutes-only and progress bar fills during sessions
  
- **PR 2:** Steps 6-7 (Header emoji enhancements)
  - After Step 7: Verify all session types show correct emoji, paused state is distinct
  
- **PR 3:** Step 8 (Timer frequency optimization)
  - After Step 8: Verify 10-second updates feel responsive, CPU usage decreases

**Validation checkpoints:**

- After Step 1: Timer display shows clean "25min" format instead of "25:00"
- After Step 5: Progress bar appears and fills proportionally during active session
- After Step 7: All states (idle, work, breaks, paused) have appropriate emoji
- After Step 8: Updates happen every 10 seconds, CPU usage is noticeably lower
- After Step 10: Complete pomodoro cycle works perfectly with all visual enhancements

**Final acceptance criteria:**
- User can glance at dropdown and instantly recognize session type without reading text
- Progress bar provides clear visual feedback of session advancement
- Timer updates feel responsive despite 10-second intervals
- All existing timer functionality continues working (no regressions)
- Visual elements render correctly on macOS (Unicode circles, emoji)
