# Increment: Pomodoro Cycle with Breaks and Skip

## User Story

As a focused developer, I want to complete a full pomodoro cycle with work sessions and breaks, so that I can follow the pomodoro technique properly with short breaks after each work session and a longer break after completing 4 work sessions.

## Acceptance Criteria

1. **Work sessions transition to short breaks** — When a 25-minute work session completes, the timer transitions to idle state and displays "05:00" for a short break, waiting for user to start

2. **Short breaks use correct duration** — When user starts a short break, the timer counts down from 5 minutes (05:00 to 00:00)

3. **Cycle tracking persists** — The system tracks which work session the user is on (1st, 2nd, 3rd, or 4th) within the current cycle

4. **Fourth work session unlocks long break** — When the 4th work session completes, the timer transitions to idle state and displays "15:00" for a long break instead of a short break

5. **Long breaks use correct duration** — When user starts a long break, the timer counts down from 15 minutes (15:00 to 00:00)

6. **Cycle resets after long break** — When a long break completes, the cycle counter resets to session 1 for the next work session

7. **Cycle indicator shows progress** — The dropdown UI displays a visual indicator (e.g., "Session 2/4" with 🍅🍅○○) showing which work session the user is currently on or about to start

8. **Skip button is available** — A "Skip" button is visible in the dropdown menu alongside Start/Pause/Reset

9. **Skip advances to next session** — When user clicks Skip during any session, the current session ends early, logs as "skipped", and advances to the appropriate next session (work→break or break→work)

10. **Session type is logged** — Session log entries include the correct session type: "work", "short_break", or "long_break"

11. **Header indicates session type** — The dropdown header displays different text for each session type (e.g., "Work Session", "Short Break", "Long Break")

12. **Tray icon reflects session type** — The tray icon updates to show different states/icons for work sessions, short breaks, and long breaks (as defined in PRD)

## Use Case

### Actors
- **User**: A focused developer using the pomodoro technique
- **System**: GoPomodoro application with timer, cycle tracking, and session logging

### Preconditions
- GoPomodoro application is running with tray icon visible
- User has completed Phase 1 functionality (can start/pause/reset a work session)
- Data directory `~/.gopomodoro/` exists
- Timer is in idle state showing "25:00"
- Cycle counter is at session 1 (start of a fresh cycle)

### Main Flow

1. User clicks the tray icon to open the dropdown panel
2. System displays dropdown showing:
   - Header: "Ready"
   - Timer display: "25:00"
   - Cycle indicator: "Session 1/4  🍅○○○"
   - Enabled "Start" button, disabled "Pause" and "Reset" buttons, enabled "Skip" button
3. User clicks "Start" button
4. System begins first work session (25 minutes)
5. System updates header to "Work Session"
6. System updates tray icon to work state (🍅 red)
7. Timer counts down over 25 minutes
8. Timer reaches 00:00
9. System logs "work,completed,25" to sessions.log
10. System transitions to short break state:
    - Header: "Ready for Break"
    - Timer display: "05:00"
    - Cycle indicator: "Session 1/4  🍅○○○" (still showing session 1 as complete)
    - Tray icon: idle state
11. User clicks "Start" button
12. System begins short break (5 minutes)
13. System updates header to "Short Break"
14. System updates tray icon to short break state (☕ green)
15. System logs "short_break,started,0" to sessions.log
16. Timer counts down 5 minutes to 00:00
17. System logs "short_break,completed,5" to sessions.log
18. System transitions to next work session state:
    - Header: "Ready"
    - Timer display: "25:00"
    - Cycle indicator: "Session 2/4  🍅🍅○○"
    - Tray icon: idle state
19. User starts session 2 (steps 3-18 repeat for sessions 2 and 3)
20. User completes session 3's short break
21. System shows cycle indicator: "Session 4/4  🍅🍅🍅🍅"
22. User starts and completes session 4 (25-minute work session)
23. System logs "work,completed,25" to sessions.log
24. System transitions to long break state:
    - Header: "Ready for Long Break"
    - Timer display: "15:00"
    - Cycle indicator: "Session 4/4  🍅🍅🍅🍅"
    - Tray icon: idle state
25. User clicks "Start" button
26. System begins long break (15 minutes)
27. System updates header to "Long Break"
28. System updates tray icon to long break state (🌟 blue)
29. System logs "long_break,started,0" to sessions.log
30. Timer counts down 15 minutes to 00:00
31. System logs "long_break,completed,15" to sessions.log
32. System resets cycle counter to session 1
33. System transitions back to work session state:
    - Header: "Ready"
    - Timer display: "25:00"
    - Cycle indicator: "Session 1/4  🍅○○○"
34. User can begin a new cycle

### Alternate / Exception Flows

**A1: User skips during work session**
- At any point during steps 4-8 (work session running):
  - User clicks "Skip" button
  - System stops timer immediately
  - System logs "work,skipped,<elapsed_minutes>" to sessions.log
  - System does NOT increment completed work counter
  - System advances to appropriate break state (short break if sessions 1-3, long break if session 4)
  - Flow continues to step 10 or step 24 depending on cycle position

**A2: User skips during break**
- At any point during a short or long break:
  - User clicks "Skip" button
  - System stops timer immediately
  - System logs "short_break,skipped,<elapsed_minutes>" or "long_break,skipped,<elapsed_minutes>"
  - System advances to next work session
  - Cycle indicator increments if it was a short break, resets if it was a long break
  - Flow continues to step 18 (next work session ready state)

**A3: User pauses and resumes during break**
- During step 14 or step 28 (break session running):
  - User clicks "Pause" button
  - System pauses timer
  - System updates header to "Paused"
  - User clicks "Start" or "Resume" button
  - System resumes countdown from paused time
  - Break continues normally

**A4: User resets during any session**
- At any point during a work session or break:
  - User clicks "Reset" button
  - System stops timer immediately
  - System does NOT log the current incomplete session as skipped
  - System resets cycle counter to session 1
  - System returns to initial idle state: "Ready", "25:00", "Session 1/4  🍅○○○"
  - Flow returns to step 2

**E1: Completed work session doesn't advance cycle if skipped**
- If user previously skipped a work session (A1):
  - Cycle indicator does not count that skipped session toward the 4-session requirement
  - Only completed work sessions increment the cycle counter
  - User must complete 4 full work sessions to unlock long break

## Context

Phase 1 (MVP Core) delivered a functional 25-minute work timer. Users can start, pause, and reset a single timer, with session starts and completions logged to `~/.gopomodoro/sessions.log`. The tray icon is visible and the dropdown UI displays timer controls.

However, the current implementation only supports work sessions. There is no way to:
- Take breaks between work sessions
- Track progress through a multi-session cycle
- Follow the complete pomodoro technique as defined in the PRD
- Distinguish between different session types in the UI or logs

The pomodoro technique requires alternating work periods with short breaks, and taking a longer break after completing multiple work sessions. Without this cycle logic, users cannot properly follow the technique and must manually track their own breaks outside the application.

This increment builds on Phase 1 by introducing:
- Two additional session types (short break, long break)
- Cycle tracking to know which session (1-4) the user is on
- Logic to determine the appropriate next session after each completion
- Skip functionality to allow users to advance through sessions flexibly
- Visual feedback showing cycle progress

The increment keeps the manual transition approach (user must click Start for each session) rather than implementing auto-transitions, maintaining user control and keeping complexity manageable.

## Goal

### Outcome
After this increment, users will be able to follow the complete pomodoro technique through the app. They will:
- Complete work sessions followed by appropriate breaks
- See visual feedback about their progress through the 4-session cycle
- Have the option to skip sessions when needed
- Experience different session types (work, short break, long break) with appropriate durations and visual states

### Scope
This increment adds:
- Three distinct session types (work, short break, long break) with correct durations (25m, 5m, 15m)
- Cycle tracking that persists which work session (1-4) the user is on
- State transition logic to determine the next appropriate session after completion
- Skip button functionality to advance through sessions
- Cycle indicator UI element showing progress (e.g., "Session 2/4  🍅🍅○○")
- Session type differentiation in UI (header text) and tray icon states
- Session type logging in sessions.log

### Non-Goals
This increment explicitly does **not** include:
- Auto-transitions between sessions (user must manually start each session)
- Notifications when sessions complete
- Color-coded UI backgrounds or themes
- Progress bars showing time elapsed within a session
- Settings or configuration for custom durations
- Statistics or streak tracking
- Pause/resume during work sessions (already exists from Phase 1, extended to breaks)

### Why This Is a Good Increment
This increment is:
- **Small and focused:** Adds cycle logic and session types without tackling notifications, visual polish, or settings
- **Self-contained:** Builds directly on Phase 1's timer foundation without requiring architectural changes
- **Testable:** Each acceptance criterion can be verified through manual testing (start sessions, observe transitions, check logs)
- **Valuable:** Unlocks the core pomodoro technique workflow, making the app genuinely useful for focus work
- **Low-risk:** Extends existing state machine incrementally; can be tested thoroughly before release

## Tasks

### Task 1: Support Multiple Session Types
- **Task:** The system can run three distinct session types: work sessions (25 minutes), short breaks (5 minutes), and long breaks (15 minutes), each with appropriate duration
- **User/Stakeholder Impact:** Users can take proper breaks as prescribed by the pomodoro technique instead of only running work timers
- **Acceptance Clues:** Starting a short break counts down from 05:00; starting a long break counts down from 15:00; session types are distinguishable in the UI

### Task 2: Track Cycle Position
- **Task:** The system maintains a counter tracking which work session (1, 2, 3, or 4) the user is currently on or about to start within the current cycle
- **User/Stakeholder Impact:** Users can see their progress toward the long break and know where they are in the pomodoro cycle
- **Acceptance Clues:** After completing a work session, the cycle counter increments; after completing the 4th session's long break, the counter resets to 1

### Task 3: Determine Next Session After Completion
- **Task:** When any session completes, the system knows what the appropriate next session should be (work→break, break→work, session 4→long break, long break→session 1)
- **User/Stakeholder Impact:** Users see the correct session type and duration ready to start after each completion, following proper pomodoro structure
- **Acceptance Clues:** After work sessions 1-3 complete, a 5-minute break is ready; after work session 4 completes, a 15-minute break is ready; after any break completes, the next work session is ready

### Task 4: Skip Button Advances Sessions
- **Task:** A "Skip" button is available in the dropdown UI; clicking it ends the current session early (logged as "skipped") and advances to the appropriate next session
- **User/Stakeholder Impact:** Users can move through the cycle flexibly if they need to cut a session short or skip a break
- **Acceptance Clues:** Clicking Skip during a work session logs it as skipped and presents the appropriate break; clicking Skip during a break presents the next work session; skipped work sessions don't count toward the 4-session cycle requirement

### Task 5: Display Cycle Indicator
- **Task:** The dropdown UI shows a visual indicator of cycle progress, displaying the current session number and tomato icons representing completed sessions (e.g., "Session 2/4  🍅🍅○○")
- **User/Stakeholder Impact:** Users can see at a glance how many sessions they've completed and how many remain before their long break
- **Acceptance Clues:** Indicator updates after each completed work session; shows filled tomatoes for completed sessions and empty circles for remaining sessions; resets to "Session 1/4  🍅○○○" after long break

### Task 6: Differentiate Session Types in UI
- **Task:** The dropdown header text clearly indicates the current session type ("Work Session", "Short Break", "Long Break", "Ready", "Ready for Break", etc.)
- **User/Stakeholder Impact:** Users immediately know what type of session is running or available to start without ambiguity
- **Acceptance Clues:** Header displays "Work Session" during work; "Short Break" during 5-minute break; "Long Break" during 15-minute break; "Ready for Break" or "Ready for Long Break" when a break is queued

### Task 7: Update Tray Icon for Session Types
- **Task:** The tray icon displays different visual states for work sessions, short breaks, and long breaks (🍅 red for work, ☕ green for short break, 🌟 blue for long break)
- **User/Stakeholder Impact:** Users can see at a glance from the menu bar what type of session is currently running, even when the dropdown is closed
- **Acceptance Clues:** Tray icon changes when transitioning between session types; shows appropriate icon for each session type as defined in the PRD

### Task 8: Log Session Types Correctly
- **Task:** Session log entries in `sessions.log` record the correct session type for each event: "work", "short_break", or "long_break"
- **User/Stakeholder Impact:** Session logs accurately reflect what type of session was run, enabling future statistics and analysis
- **Acceptance Clues:** Inspecting sessions.log shows entries like "2025-12-25T10:00:00Z,short_break,completed,5" and "2025-12-25T10:30:00Z,long_break,started,0"

## Risks and Assumptions

### Risks
- **State complexity:** Managing cycle position, session type, and next-session logic increases state machine complexity; bugs in transition logic could leave users stuck or confused about where they are in the cycle
- **Skip behavior confusion:** Users might expect Skip to behave differently (e.g., skip all remaining breaks vs. skip current session only); unclear labeling could lead to accidental skips
- **Cycle counter edge cases:** Logic for tracking completed vs. skipped sessions must be correct; off-by-one errors could cause long breaks to appear too early or too late
- **Icon availability:** Emoji icons (🍅☕🌟) may not render consistently across all macOS versions or system tray configurations

### Assumptions
- Users are comfortable with manual transitions (clicking Start for each new session) rather than auto-start
- The existing timer infrastructure from Phase 1 can be extended to support multiple durations without major refactoring
- Logging session types to the existing CSV format is sufficient; no schema changes are needed
- Cycle position can be tracked in memory and reset on app restart (no persistence of cycle counter across restarts initially)
- Users will understand the cycle indicator notation (Session X/4 with tomato icons)

### Mitigations
- Implement state transitions incrementally and test each transition path manually
- Ensure Skip button is clearly labeled and positioned to avoid accidental clicks
- Review cycle logic carefully; consider manual testing scenarios for all possible skip/reset combinations
- If emoji icons don't work well, fall back to text-based indicators or simple colored dots

## Success Criteria and Observability

### Success Criteria
After release, this increment is successful if:
- Users can manually complete a full 4-session pomodoro cycle (work-break-work-break-work-break-work-long break) without errors
- Session log files show correct session types (work, short_break, long_break) with appropriate events (started, completed, skipped)
- Cycle indicator visually updates after each completed work session
- Users can skip sessions and the system correctly advances to the next appropriate session
- No state machine bugs are reported (users getting stuck, incorrect session durations, cycle not resetting)

### Observability
To confirm this increment is working correctly, check:
- **Session logs:** Inspect `~/.gopomodoro/sessions.log` for a completed cycle; should see pattern like: work→short_break→work→short_break→work→short_break→work→long_break
- **Manual testing:** Start the app, complete a full cycle, observe that:
  - Each session type has correct duration
  - Cycle indicator shows Session 1/4 through 4/4
  - After long break completion, cycle resets to Session 1/4
  - Skip button works during both work and break sessions
- **UI inspection:** Open dropdown during different session types; verify header text and tray icon match the session type
- **Edge case testing:** Skip a work session and verify it doesn't count toward the 4-session requirement; reset mid-cycle and verify cycle counter resets to 1

No additional instrumentation or metrics are needed initially. Direct observation of UI behavior and log file contents is sufficient to validate correct operation.

## Process Notes

This increment should be implemented following the project's small-steps approach:
- Build the cycle logic and state transitions first (work→break→work transitions)
- Add the Skip button functionality once basic transitions work
- Implement the cycle indicator UI element after core logic is stable
- Update tray icons for session types as a final polish step
- Manual testing is acceptable; focus testing on state transitions and cycle edge cases (skip, reset, completion)

The increment can be delivered through normal development workflow:
- Small, focused commits as each piece is completed
- Manual testing after each addition to verify no regressions
- Session log inspection to confirm correct event logging
- No special deployment process needed; this is a client-side UI and logic change

If issues arise during implementation:
- Cycle logic can be simplified by deferring Skip functionality to a follow-up increment
- Tray icon updates can be postponed if systray library limitations are encountered
- Cycle indicator can start as text-only ("Session 2/4") and add emoji icons later if rendering is problematic

## Follow-up Increments

After this increment is complete, natural next steps include:

### Notifications for Session Completion
Add native macOS notifications when sessions complete, alerting users to take breaks or return to work without needing to check the tray icon. This builds on the session completion events already logged in this increment.

### Auto-Transition Between Sessions (Optional)
Implement optional automatic transitions so that when a work session completes, the break timer starts automatically (and vice versa), with a configurable setting to enable/disable. This reduces friction but adds complexity around user expectations and interruption.

### Visual Session Type Differentiation
Add color-coded backgrounds or themes to the dropdown UI so work sessions, short breaks, and long breaks are visually distinct at a glance beyond just header text. This enhances the "Visual Clarity" principle from the PRD.

### Progress Bar Within Sessions
Display a visual progress bar in the dropdown showing how much time has elapsed in the current session, updating every 30 seconds. This provides at-a-glance progress feedback without needing to parse the timer countdown.

### Settings for Custom Durations
Allow users to configure work, short break, and long break durations via `settings.json`, overriding the hardcoded 25/5/15 defaults. This supports users with different focus preferences while maintaining the "No Settings UI" principle (manual JSON editing).

## PRD Entry

### Increment ID
`pomodoro-cycle-with-breaks-and-skip`

### Title
Pomodoro Cycle with Breaks and Skip

### Status
Proposed

### Increment Folder
`docs/increments/pomodoro-cycle-with-breaks-and-skip/`

### User Story
As a focused developer, I want to complete a full pomodoro cycle with work sessions and breaks, so that I can follow the pomodoro technique properly with short breaks after each work session and a longer break after completing 4 work sessions.

### Acceptance Criteria
- Work sessions transition to short breaks (user-initiated)
- Short breaks use 5-minute duration
- Cycle tracking persists across sessions
- 4th work session unlocks long break
- Long breaks use 15-minute duration
- Cycle resets after long break
- Cycle indicator shows progress (Session X/4 with 🍅🍅○○)
- Skip button available and functional
- Skip advances to next appropriate session
- Session types logged correctly (work, short_break, long_break)
- Header indicates current session type
- Tray icon reflects session type

### Use Case Summary
Users complete a full 4-session pomodoro cycle by manually starting each session. The system tracks cycle position (Session 1/4 through 4/4), transitions between work sessions (25 min), short breaks (5 min), and long breaks (15 min), and allows skipping sessions. After completing the 4th work session's long break, the cycle resets to Session 1. The UI shows cycle progress, session type in the header, and appropriate tray icons for each session type.
