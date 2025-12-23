# Increment: Functional Timer with Session Logging

## User Story

As a macOS user who wants to track my focus time, I want to start a 25-minute work session and see it count down, so that I can complete a pomodoro and have it recorded for later review.

## Acceptance Criteria

1. **Timer starts on button click** — When the user clicks "Start" in idle mode, a 25-minute work session begins and the timer display shows "25:00"

2. **Timer counts down visibly** — The timer display updates every second, showing the remaining time in MM:SS format (e.g., "24:59", "24:58", etc.)

3. **Pause stops the countdown** — When the user clicks "Pause" during an active session, the timer stops counting but retains the current time

4. **Resume continues from paused time** — When the user clicks "Resume" (or "Start") after pausing, the timer continues counting down from where it stopped

5. **Reset returns to initial state** — When the user clicks "Reset", the timer stops and returns to idle mode with "25:00" displayed

6. **Session completes automatically** — When the timer reaches "00:00", the session is marked as complete and the timer returns to idle mode

7. **Completed sessions are recorded** — When a session completes naturally (reaches 00:00), a record is appended to `~/.gopomodoro/sessions.log` in CSV format with timestamp, session type (work), event (completed), and duration

8. **Skipped sessions are not recorded** — If the user clicks Reset before the timer reaches 00:00, no session record is written

9. **Data directory is created** — The `~/.gopomodoro/` directory is created automatically if it doesn't exist when the first session needs to be logged

10. **Button states reflect timer state** — Start/Pause/Reset buttons are enabled or disabled appropriately based on whether the timer is idle, running, or paused

11. **App fails fast on storage errors** — If the `~/.gopomodoro/` directory cannot be created or is not writeable, the app displays an error and exits rather than continuing in a degraded state

## Use Case

### Actors
- **User**: A macOS user who wants to track focus time using pomodoro sessions
- **System**: The GoPomodoro application and macOS file system

### Preconditions
- GoPomodoro application is running with tray icon visible (from previous increment)
- User has clicked the tray icon and the dropdown panel is open
- Timer is in idle state showing "25:00"
- User has file system permissions to create directories in their home folder

### Main Flow

1. User clicks the "Start" button in the dropdown panel
2. System begins a 25-minute work session
3. System updates the timer display to show "25:00"
4. System enables the "Pause" button and disables the "Start" button
5. System starts counting down, updating the display every second ("24:59", "24:58", etc.)
6. Timer continues counting down over the next 25 minutes
7. Timer reaches "00:00"
8. System checks if `~/.gopomodoro/` directory exists
9. System creates the directory if it doesn't exist
10. System appends a session record to `~/.gopomodoro/sessions.log` with:
    - Current timestamp in ISO 8601 format
    - Session type: "work"
    - Event: "completed"
    - Duration: 25 minutes
11. System transitions timer back to idle state
12. System updates display to show "25:00"
13. System enables the "Start" button and disables the "Pause" button
14. User can start another session or close the dropdown

### Alternate / Exception Flows

**A1: User pauses the timer**
- At step 5, while timer is counting down:
  - User clicks the "Pause" button
  - System stops the countdown at current time (e.g., "18:42")
  - System changes "Pause" button to show "Resume" (or re-enables "Start")
  - System keeps "Reset" button enabled
  - Timer display remains frozen at paused time
  - User can resume (continue to A2) or reset (continue to A3)

**A2: User resumes from pause**
- After A1, user clicks "Resume" or "Start":
  - System continues counting down from the paused time
  - System returns to step 5 of main flow
  - Countdown continues normally

**A3: User resets the timer**
- At any point during steps 5-6 (timer running or paused):
  - User clicks the "Reset" button
  - System stops the countdown immediately
  - System does NOT log any session record (session was not completed)
  - System transitions back to idle state showing "25:00"
  - Flow returns to step 13 of main flow

**A4: User starts a session and lets it complete (started event)**
- At step 2, after user clicks "Start":
  - System may optionally log a "started" event to sessions.log with:
    - Timestamp, session_type: "work", event: "started", duration: 0
  - This provides a complete audit trail
  - Flow continues with step 3 of main flow

**E1: Storage directory cannot be created**
- At step 9, if system cannot create `~/.gopomodoro/`:
  - System displays an error message to the user (via log or simple dialog)
  - System exits the application immediately
  - No degraded mode operation is allowed

**E2: Sessions log file is not writeable**
- At step 10, if `sessions.log` cannot be written to:
  - System displays an error message to the user
  - System exits the application immediately
  - Session completion is not silently ignored

**E3: User closes dropdown while timer is running**
- At step 6, if user closes the dropdown panel:
  - Timer continues running in the background
  - Tray icon may update to show timer is active (future enhancement)
  - User can reopen dropdown to see current countdown
  - Timer completes normally and logs the session

## Context

GoPomodoro is building toward a complete pomodoro timer with cycles, breaks, and statistics. This increment adds the foundational timer mechanics.

### Current Situation
- The tray icon and dropdown UI structure exist from the previous increment
- Buttons are present but non-functional placeholders
- No timer logic, state management, or persistence exists yet
- Users can launch the app and see the UI, but cannot actually run a pomodoro session

### Why This Matters
- This is the minimum viable timer—users can finally complete an actual pomodoro session
- Session logging creates the foundation for statistics and streak tracking in future increments
- Hardcoded 25-minute duration keeps this increment simple while validating the timer engine
- Failing fast on storage errors prevents silent data loss and builds user trust

### Key Constraints
- Timer duration is hardcoded to 25 minutes (no configuration yet)
- Only work sessions are supported (no break timers)
- No notifications when sessions complete
- No statistics display (just logging for now)
- CSV format for sessions.log as specified in PRD
- Must follow the constitution's "small, safe steps" principle

## Goal

### Outcome
After this increment, users will experience:
- A functional countdown timer that runs for 25 minutes
- Visual feedback every second as the timer counts down
- The ability to pause, resume, and reset the timer
- Automatic return to idle state when a session completes
- Confidence that completed sessions are being recorded for future reference

### Scope
This increment focuses narrowly on:
- A single timer that counts down from 25:00 to 00:00
- Start, pause, resume, and reset controls
- Updating the display every second
- Detecting session completion
- Creating the data directory and logging completed sessions in CSV format
- Failing gracefully if storage is unavailable

### Non-Goals
This increment explicitly does NOT include:
- Notifications or alerts when sessions complete
- Short break or long break timers
- The full 4-session pomodoro cycle
- Configurable timer durations
- Reading or displaying statistics
- Settings persistence
- Tray icon state changes during active sessions
- Sound effects
- Progress bars or visual indicators beyond the countdown

### Why This Is a Good Increment
- **Small and focused**: Adds one core capability (working timer) without sprawling into cycles, breaks, or notifications
- **Testable**: Easy to verify the countdown works and sessions are logged correctly
- **Builds on previous work**: Uses the UI structure from the tray-icon increment
- **Enables future work**: Session logging is required for statistics, and the timer engine is needed for breaks and cycles
- **Follows constitution**: Hardcoded values ("make it work"), simple CSV persistence, small safe step

## Tasks

### Task 1: Timer Countdown Engine
**Task:** A timer can count down from 25 minutes to zero, updating every second, and detect when it reaches completion.

**User/Stakeholder Impact:** Users see the time remaining in their work session decrease second by second, creating a sense of progress and urgency.

**Acceptance Clues:**
- Timer display changes from "25:00" to "24:59" to "24:58" and so on
- After 25 minutes (1500 seconds), the timer reaches "00:00"
- The timer automatically stops at "00:00" without going negative

### Task 2: Start, Pause, and Resume Controls
**Task:** Users can start the countdown, pause it mid-session, and resume from the paused time.

**User/Stakeholder Impact:** Users have control over their pomodoro session, allowing them to handle interruptions without losing progress.

**Acceptance Clues:**
- Clicking "Start" when idle begins the countdown
- Clicking "Pause" stops the countdown at the current time
- Clicking "Start" or "Resume" after pausing continues the countdown
- The timer remembers the paused time accurately

### Task 3: Reset to Idle State
**Task:** Users can reset the timer at any point, returning it to idle state without logging a session.

**User/Stakeholder Impact:** Users can abandon a session that was interrupted or started accidentally, with a clean slate for the next attempt.

**Acceptance Clues:**
- Clicking "Reset" during a running or paused timer stops the countdown
- Display returns to "25:00"
- Timer is ready to start a fresh session
- No session record is written for the abandoned session

### Task 4: Automatic Idle Transition on Completion
**Task:** When the timer reaches "00:00", the session is considered complete and the timer automatically returns to idle state.

**User/Stakeholder Impact:** Users don't need to manually reset after completing a session—the timer is ready for the next pomodoro.

**Acceptance Clues:**
- At "00:00", the timer stops counting
- Display shows "25:00" again (idle state)
- Start button is available for the next session

### Task 5: Session Logging to CSV File
**Task:** When a session completes naturally (reaches "00:00"), a record is appended to `~/.gopomodoro/sessions.log` in CSV format.

**User/Stakeholder Impact:** Users build a history of completed pomodoros that can be reviewed or analyzed later, creating a foundation for progress tracking.

**Acceptance Clues:**
- After a session completes, a new line appears in `sessions.log`
- The line contains: timestamp (ISO 8601), "work", "completed", and "25"
- Multiple completed sessions result in multiple log entries
- The file is human-readable CSV format

### Task 6: Data Directory Creation
**Task:** The `~/.gopomodoro/` directory is created if it doesn't exist when the first session needs to be logged.

**User/Stakeholder Impact:** Users don't need to manually create folders—the app handles its own storage setup.

**Acceptance Clues:**
- Before the first completed session, `~/.gopomodoro/` may not exist
- After the first completed session, `~/.gopomodoro/` exists
- The directory is created in the user's home folder with appropriate permissions

### Task 7: Fail Fast on Storage Errors
**Task:** If the data directory cannot be created or the sessions.log file cannot be written, the app displays an error and exits immediately.

**User/Stakeholder Impact:** Users are alerted to storage problems immediately rather than losing session data silently, building trust in the app's reliability.

**Acceptance Clues:**
- If `~/.gopomodoro/` cannot be created (e.g., permissions issue), the app shows an error before or when trying to log
- If `sessions.log` is not writeable, the app shows an error and exits
- The app does not continue running in a degraded state
- Error messages are clear enough for the user to understand the problem

### Task 8: Button State Management
**Task:** Start, Pause, and Reset buttons are enabled or disabled based on the current timer state (idle, running, paused).

**User/Stakeholder Impact:** Users receive clear visual feedback about which actions are available, reducing confusion and accidental clicks.

**Acceptance Clues:**
- When idle: Start is enabled, Pause is disabled (or hidden)
- When running: Pause is enabled, Start is disabled
- When paused: Start/Resume is enabled, Pause is disabled
- Reset is enabled whenever the timer is not idle
- Disabled buttons are visually distinct (grayed out or hidden)

### Task 9: Background Timer Operation
**Task:** The timer continues counting down even when the dropdown panel is closed, and displays the correct time when reopened.

**User/Stakeholder Impact:** Users can minimize the dropdown and focus on work without stopping the timer, checking back as needed.

**Acceptance Clues:**
- User starts timer, closes dropdown, waits, reopens dropdown
- Timer display shows the expected elapsed time
- Timer completes and logs session even if dropdown is closed at completion time

## Risks and Assumptions

### Risks
- **Timer accuracy**: System timers may drift over 25 minutes if not implemented carefully, leading to sessions that are slightly longer or shorter than expected
- **File I/O failures**: Disk full, permissions issues, or file corruption could prevent session logging
- **State management complexity**: Coordinating timer state, UI state, and button states could introduce bugs if not carefully designed
- **User confusion**: If buttons don't clearly indicate current state, users may not understand whether the timer is running or paused

### Assumptions
- Users have standard file system permissions in their home directory
- The `~/.gopomodoro/` directory name does not conflict with existing files or directories
- CSV format is sufficient for session logging (no need for JSON or database yet)
- Users understand MM:SS time format
- One-second update granularity is sufficient (no need for sub-second precision)
- 25 minutes is an acceptable hardcoded default for initial testing

### Mitigations
- If timer accuracy becomes an issue, later increments can improve the timing mechanism
- Fail-fast behavior on storage errors ensures users are aware of problems immediately
- Manual testing of state transitions can catch most button state bugs before release
- Clear button labels and visual states reduce user confusion

## Success Criteria and Observability

### Success Criteria
After this increment is complete, we will know it succeeded if:
- Users can start a timer and watch it count down for 25 minutes
- The timer display updates smoothly every second
- Pause, resume, and reset controls work as expected
- Completed sessions appear in `~/.gopomodoro/sessions.log` with correct timestamps and format
- Abandoned (reset) sessions do NOT appear in the log
- The app exits with a clear error if storage is unavailable

### Observability
After release, check:
- **Manual verification**: Run a full 25-minute session and inspect `~/.gopomodoro/sessions.log` to confirm the entry is correct
- **Pause/resume testing**: Pause at various points, wait, resume, and verify the timer continues correctly
- **Reset testing**: Start a session, reset it, verify no log entry was created
- **Storage error simulation**: Remove write permissions on `~/.gopomodoro/` and verify the app fails fast with an error
- **Background operation**: Start timer, close dropdown, wait, reopen, verify time is correct

Since this is a lite constitution project, automated tests are optional, but manual testing of these scenarios is essential before considering the increment complete.

## Process Notes

### Delivery Approach
This increment should be implemented via small, safe changes:
- Build the timer countdown logic first, get it ticking reliably
- Add start/pause/reset state management next
- Wire up the button click handlers to control the timer
- Implement session completion detection
- Add CSV logging last, with fail-fast error handling

Each piece can be tested independently before moving to the next.

### Testing Strategy
- Manual testing is acceptable per the lite constitution
- Test the full 25-minute session at least once to verify accuracy
- Test pause/resume at various time points (early, middle, late in session)
- Test reset from running and paused states
- Test session logging by inspecting the CSV file
- Test storage error handling by simulating permission issues

### Rollout
- This increment flows through normal build and run processes
- No special deployment steps required
- If the increment introduces bugs, they will be caught during manual testing
- The fail-fast behavior ensures storage issues are caught immediately, not silently

## Follow-up Increments

### Short Break Timer
After completing a work session, transition to a 5-minute short break timer instead of returning to idle. This begins implementing the full pomodoro cycle.

### Long Break and Cycle Logic
Implement the 4-session cycle with long breaks, tracking session count and offering appropriate break durations.

### Notifications on Session Complete
Add native macOS notifications when work sessions and breaks complete, with action buttons to start the next phase.

### Configurable Timer Durations
Read timer durations from a settings file instead of hardcoding them, allowing users to customize work/break lengths.

### Statistics Display
Calculate and display today's completed pomodoros, current streak, and other stats in the dropdown footer.

### Tray Icon State Indication
Update the tray icon color or display to show when a timer is actively running, paused, or on break.

## PRD Entry (for docs/PRD.md)

**Increment ID:** `functional-timer-with-session-logging`

**Title:** Functional Timer with Session Logging

**Status:** Proposed

**Increment Folder:** `docs/increments/functional-timer-with-session-logging/`

**User Story:** As a macOS user who wants to track my focus time, I want to start a 25-minute work session and see it count down, so that I can complete a pomodoro and have it recorded for later review.

**Acceptance Criteria:**
- Timer starts on button click, displaying "25:00"
- Timer counts down visibly in MM:SS format, updating every second
- Pause stops countdown; resume continues from paused time
- Reset returns timer to idle state without logging
- Session completes automatically at "00:00" and transitions to idle
- Completed sessions are recorded in `~/.gopomodoro/sessions.log` (CSV format)
- Skipped sessions are not recorded
- Data directory (`~/.gopomodoro/`) is created automatically if missing
- Button states reflect timer state (idle, running, paused)
- App fails fast with error if storage directory cannot be created or is not writeable

**Use Case Summary:**
User clicks "Start" → timer counts down from 25:00 for 25 minutes → timer reaches 00:00 → session is logged to CSV file → timer returns to idle state. User can pause/resume or reset at any time. Timer runs in background even when dropdown is closed.
