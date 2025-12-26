# Increment: Visual Feedback with Progress Indicators

## User Story

As a GoPomodoro user, I want to see clear visual indicators of my current session state and progress, so that I can quickly understand where I am in my pomodoro cycle without reading text carefully.

## Acceptance Criteria

1. **Tray icon changes per session type** — Different icon files are displayed for idle, work, short break, long break, and paused states (using existing icon infrastructure in tray.go)

2. **Progress bar shows session advancement** — The dropdown displays a visual progress bar using Unicode block characters (e.g., `●●●●○○○○○○`) that updates as the session progresses

3. **Header text includes state emoji/symbols** — The header menu item uses emoji and formatting to clearly indicate state (e.g., "🍅 Work Session" vs "☕ Short Break")

4. **Timer displays minutes only** — The countdown shows minutes in a clean format (e.g., "12min", "5min") instead of MM:SS format

5. **Timer updates every 10 seconds** — The countdown display refreshes every 10 seconds, balancing CPU efficiency with responsive feedback

6. **Cycle indicator uses visual symbols** — The session progress (1/4, 2/4, etc.) is shown with filled/empty tomato emoji (🍅○○○) making it scannable at a glance

7. **Paused state is visually distinct** — When paused, the UI clearly indicates this through both the tray icon and header text (e.g., "⏸️ Paused")

## Use Case

### Actors
- **User**: A GoPomodoro user who has the app running and wants to see clear visual feedback about their pomodoro session

### Preconditions
- GoPomodoro is installed and running on macOS
- Phase 1 (MVP Core) and Phase 2 (Pomodoro Cycle) are complete and functional
- Timer and session logic work correctly
- Basic tray icon and dropdown menu are operational
- Icon asset files exist in the `assets/` directory

### Main Flow

1. User launches GoPomodoro or has it already running
2. System displays the idle tray icon (gray/neutral) in the menu bar
3. User clicks the tray icon to open the dropdown
4. Dropdown displays:
   - Header: "Ready" or "○ Idle"
   - Timer display: "25min"
   - Cycle indicator: "Session 1/4  🍅○○○"
   - Progress bar: Empty (○○○○○○○○○○)
   - Control buttons (Start, Pause, Reset, Skip, Quit)
5. User clicks "Start" button
6. System transitions to work session and updates display:
   - Tray icon changes to work session icon (red/tomato themed)
   - Header updates to: "🍅 Work Session"
   - Timer shows: "25min"
   - Progress bar shows: "○○○○○○○○○○" (empty at start)
7. After 10 seconds elapse
8. System updates timer display to "24min"
9. After 2.5 minutes total (10% of 25min session)
10. System updates progress bar to: "●○○○○○○○○○" (1 filled segment)
11. User observes the visual progress every 10 seconds
12. After 12.5 minutes (50% complete)
13. Progress bar displays: "●●●●●○○○○○" (halfway filled)
14. Timer shows: "12min"
15. After 25 minutes, session completes
16. System shows: "0min" and progress bar: "●●●●●●●●●●" (fully filled)
17. System auto-transitions to short break
18. Tray icon changes to break icon (green/coffee themed)
19. Header updates to: "☕ Short Break"
20. Timer resets to: "5min"
21. Progress bar resets to empty: "○○○○○○○○○○"
22. Cycle indicator updates to: "Session 1/4  🍅○○○"
23. User continues through cycles, seeing visual feedback for each session type

### Alternate / Exception Flows

**A1: User pauses during work session**
- At step 11, user clicks "Pause"
- Tray icon changes to paused icon (gray pause symbol)
- Header updates to: "⏸️ Paused"
- Timer freezes at current value (e.g., "18min")
- Progress bar maintains current state
- User can click "Resume" to continue from same point

**A2: Long break reached after 4 work sessions**
- After completing 4th work session
- Tray icon changes to long break icon (blue/star themed)
- Header updates to: "🌟 Long Break"
- Timer displays: "15min"
- Cycle indicator shows: "Session 4/4  🍅🍅🍅🍅"
- Progress bar starts empty and fills over 15 minutes

**A3: User skips current session**
- At any point during active session
- User clicks "Skip"
- System logs the skip event
- System advances to next session in cycle
- All visual indicators update immediately to new session type

**A4: Timer shows less than 1 minute remaining**
- When timer reaches 0 minutes remaining (< 1 min left)
- System displays: "0min" instead of negative or fractional values
- Progress bar shows fully filled: "●●●●●●●●●●"

## Context

GoPomodoro has completed its foundational phases: Phase 1 (MVP Core) delivered a functional timer with tray integration, and Phase 2 (Pomodoro Cycle) implemented the complete 4-session workflow with automatic transitions between work and break periods.

**Current State:**
The app works reliably. Users can start, pause, reset, and skip through pomodoro sessions. The cycle logic correctly handles short breaks after work sessions and long breaks after 4 completed cycles. Session events are logged to CSV files for future statistics.

However, the UI is entirely text-based, relying on systray's menu items. Users see:
- Generic text like "Ready" or "Work Session"
- Time in MM:SS format (e.g., "24:37")
- Cycle indicator as text: "Session 2/4  🍅🍅○○"
- Timer updates every second

**The Problem:**
Users must carefully read and interpret text to understand their current state. There's no immediate visual distinction between a work session and a break. Progress through a session isn't visible—users must mentally calculate how much of "24:37" represents completion. The MM:SS format provides more precision than needed for a 25-minute session. Every-second updates consume CPU unnecessarily.

**Why Now:**
With core functionality stable and tested, this is the right moment to improve user experience. Visual feedback transforms the app from "functional but plain" to "functional and pleasant to use." This increment delivers clear value without touching the stable timer and session logic.

**Technical Constraints:**
The app uses the `systray` library for macOS tray integration, which provides text-based menu items only—no native support for colored backgrounds, rich formatting, or animations. All visual enhancements must work within these constraints using Unicode characters, emoji, and icon file swapping.

**Related Work:**
- Phase 1 established the tray icon infrastructure and basic menu structure
- Phase 2 implemented the session state machine and cycle tracking
- Existing icon assets in `assets/` directory provide foundation for state-based icon changes

## Goal

**Outcome:**
Users can glance at either the tray icon or the dropdown menu and immediately understand three key pieces of information without reading text carefully:
1. What session type is currently active (work, short break, long break, or paused)
2. Approximately how much time remains in the current session
3. Where they are in the 4-session pomodoro cycle

This transforms the UI from "requires interpretation" to "instantly scannable."

**Scope:**
This increment focuses exclusively on visual presentation improvements:
- Enhanced tray icon states that clearly distinguish session types
- Unicode-based progress bar showing session advancement
- Emoji-enriched header text for quick state recognition
- Simplified timer display (minutes only)
- Optimized update frequency (10 seconds instead of 1 second)
- Visual distinction for paused state

All changes are presentational—no modifications to timer logic, session state machine, or data persistence.

**Non-Goals:**
- Not migrating to a different UI framework (staying with systray)
- Not implementing true color-coded backgrounds (systray limitation)
- Not adding dark mode support (deferred to future work)
- Not implementing smooth animations or transitions (not possible with static menu items)
- Not changing any timer duration logic or session behavior
- Not adding new functionality beyond visual feedback

**Why This Is a Good Increment:**
This increment is small, focused, and low-risk. It builds directly on stable foundation without touching core logic. Changes are purely presentational, making them easy to test and validate through visual inspection. The scope is clear and bounded—there's no risk of scope creep into behavior changes. Users will immediately notice and appreciate the improved experience, making this a high-value, low-risk change that can be delivered quickly.

## Tasks

### Task 1: Progress Bar Visualization
**Task:** A visual progress bar appears in the dropdown menu that fills proportionally as the current session advances.

**User/Stakeholder Impact:** Users can see at a glance how much of their current session has elapsed without doing mental math on the remaining time. This provides a satisfying sense of progress and helps users pace their work or break time.

**Acceptance Clues:**
- A menu item displays a 10-segment bar using Unicode characters (●○○○○○○○○○)
- The bar updates as the session progresses (e.g., empty at start, half-filled at 50%, fully filled at 100%)
- The bar resets to empty when a new session begins
- The bar maintains its current state when paused

### Task 2: Minutes-Only Timer Display
**Task:** The timer display shows only minutes in a clean, simple format instead of minutes and seconds.

**User/Stakeholder Impact:** Users see a simplified time display that's easier to scan and read at a glance. For 25-minute sessions, knowing "12min remaining" is more useful than "12:34 remaining."

**Acceptance Clues:**
- Timer menu item shows formats like "25min", "12min", "5min", "0min"
- Format is consistent across all session types (work, short break, long break)
- When less than 1 minute remains, display shows "0min"

### Task 3: 10-Second Update Frequency
**Task:** The timer display and progress bar update every 10 seconds instead of every second.

**User/Stakeholder Impact:** Reduced CPU usage and battery consumption while maintaining useful feedback. Users still get timely updates without the unnecessary overhead of second-by-second refreshes.

**Acceptance Clues:**
- Timer display changes approximately every 10 seconds during active sessions
- Progress bar updates at the same 10-second intervals
- CPU usage is measurably lower compared to 1-second update frequency
- Users perceive the timer as responsive without noticing the reduced update rate

### Task 4: Session-Specific Header Emoji
**Task:** The header menu item includes emoji symbols that clearly indicate the current session type.

**User/Stakeholder Impact:** Users instantly recognize whether they're in a work session, break, or paused state through familiar visual symbols. This makes the app scannable even in peripheral vision.

**Acceptance Clues:**
- Work session shows: "🍅 Work Session"
- Short break shows: "☕ Short Break"
- Long break shows: "🌟 Long Break"
- Paused state shows: "⏸️ Paused"
- Idle state shows: "○ Idle" or similar neutral indicator

### Task 5: State-Based Tray Icon Changes
**Task:** The tray icon image changes based on the current session type and timer state.

**User/Stakeholder Impact:** Users can see their current state in the macOS menu bar without opening the dropdown. This provides at-a-glance awareness even when focused on other applications.

**Acceptance Clues:**
- Different icon files are displayed for idle, work, short break, long break, and paused states
- Icon changes are visible in the menu bar
- Icon updates happen in sync with session transitions
- Icons are visually distinct enough to recognize quickly

### Task 6: Distinct Paused State Appearance
**Task:** When the timer is paused, the UI clearly indicates this through both the tray icon and menu text.

**User/Stakeholder Impact:** Users immediately know when they've paused the timer versus when it's actively running. This prevents confusion about whether time is counting down.

**Acceptance Clues:**
- Paused state shows "⏸️ Paused" in header
- Tray icon displays a paused-specific image
- Timer value freezes at current minutes
- Progress bar maintains its current fill state
- Paused appearance is visually distinct from running states

### Task 7: Enhanced Cycle Indicator
**Task:** The cycle progress indicator uses filled/empty tomato emoji to show position in the 4-session cycle.

**User/Stakeholder Impact:** Users can quickly see how many work sessions they've completed in the current cycle and how many remain before their long break.

**Acceptance Clues:**
- Format displays as: "Session 1/4  🍅○○○", "Session 2/4  🍅🍅○○", etc.
- Indicator updates when advancing to the next session
- After 4 sessions, shows: "Session 4/4  🍅🍅🍅🍅"
- Resets to "Session 1/4  🍅○○○" after long break completes

## Risks and Assumptions

### Risks

**Unicode Rendering Inconsistency:**
Different macOS versions or system font settings might render Unicode progress bar characters (● ○) inconsistently. Some systems might show different-sized circles or spacing issues.

**10-Second Updates Feel Too Slow:**
Some users might perceive 10-second update intervals as unresponsive, especially if they're used to seeing second-by-second countdown. This could make the timer feel "stuck" or unresponsive.

**Progress Bar Segment Alignment:**
With 10 segments, some session durations won't divide evenly (e.g., 25 minutes = 2.5 minutes per segment). This might cause the progress bar to appear to "jump" or update irregularly.

**Emoji Font Availability:**
Session-type emoji (🍅☕🌟⏸️) rely on macOS emoji support. Older macOS versions or specific system configurations might not render these correctly.

**Icon File Requirements:**
Implementation assumes icon asset files for different states already exist or are trivial to create. If icons need design work, this could expand scope.

### Assumptions

**systray Unicode Support:**
Assumes the systray library correctly displays Unicode characters and emoji in menu item text on macOS. (Prior phases have used emoji successfully in cycle indicators, suggesting this is safe.)

**10-Second Granularity is Sufficient:**
Assumes users don't need second-level precision for pomodoro sessions. For 25-minute sessions, updates every 10 seconds provide enough feedback without feeling sluggish.

**Icon Assets Exist:**
Assumes the `assets/` directory already contains or can easily accommodate icon files for: idle, work, short break, long break, and paused states.

**Existing Tray Icon Infrastructure:**
Assumes the `tray.go` package's `UpdateIcon()` method works correctly and just needs to be called at appropriate times with correct session type parameters.

**No Accessibility Concerns:**
Assumes emoji-based visual feedback is acceptable for the target user base. (Future work might add accessibility considerations for users who rely on screen readers.)

## Success Criteria and Observability

### Success Criteria

This increment is successful when users can demonstrate or report that they:
- **Instantly recognize session type** by glancing at the tray icon or dropdown header without reading text
- **Understand progress** through their current session by looking at the progress bar
- **Feel satisfied** with the visual polish and find the app more pleasant to use than the text-only version

### What Will Be Observed After Release

**Visual Inspection:**
- Open the dropdown during idle, work, short break, long break, and paused states
- Verify emoji appear correctly in headers ("🍅 Work Session", "☕ Short Break", etc.)
- Confirm progress bar displays using Unicode circles (●○○○○○○○○○)
- Check that progress bar fills proportionally during a session
- Observe that timer displays in minutes-only format ("25min", "12min")

**Tray Icon Verification:**
- Visually confirm different icons appear for each state in the macOS menu bar
- Verify icon changes happen in sync with session transitions
- Check that paused state shows distinct icon

**Performance Check:**
- Monitor CPU usage with Activity Monitor during active session
- Compare CPU usage to previous 1-second update implementation (if measurable baseline exists)
- Verify updates happen approximately every 10 seconds

**Functional Verification:**
- Start a work session, observe visual updates over several minutes
- Pause mid-session, confirm "⏸️ Paused" appears and timer freezes
- Resume, verify timer continues and progress bar advances from paused position
- Complete a full session, verify progress bar fills completely before transitioning
- Skip a session, observe immediate visual updates to next state
- Complete 4 work sessions, verify long break indicator appears correctly

**Edge Cases:**
- Test behavior when timer shows "0min" (final minute)
- Verify progress bar doesn't overflow or show more than 10 filled segments
- Check that rapid session transitions (skip repeatedly) update visuals correctly

**No Regression:**
- All existing timer functionality works: Start, Pause, Resume, Reset, Skip
- Session cycle logic continues working (short breaks, long breaks, auto-transitions)
- Session logging still records events correctly

## Process Notes

This increment should move through the standard build-test-release workflow established for GoPomodoro:

**Development Approach:**
- Implement changes incrementally: start with timer format, then progress bar, then header emoji, then icon updates
- Test each piece independently before integrating
- Manual testing is appropriate for visual changes (automated testing of UI strings would be brittle)
- Commit working changes frequently to avoid large, risky commits

**Testing Strategy:**
- Manual visual verification is the primary testing method
- Run through complete pomodoro cycles to verify all states appear correctly
- Test pause/resume and skip functionality to ensure visual updates work in all scenarios
- No need for extensive automated UI tests given the constitution's "lite" approach

**Deployment:**
- Build and test locally first
- Run through at least one complete 4-session cycle before considering it ready
- Since changes are visual-only and don't affect data or core logic, rollback risk is minimal
- If rendering issues appear on specific macOS versions, may need to adjust Unicode characters or emoji choices

**Rollback Plan:**
- If visual changes cause issues, can easily revert to previous text-only display
- No data migration or backwards compatibility concerns since no persistence changes

## Follow-up Increments

### Future Visual Enhancements
If this increment succeeds and users want additional polish, consider:
- Dark mode support (detect system theme, adjust icon set)
- More sophisticated progress visualization (different bar styles, animations if framework allows)
- Customizable emoji or symbols via settings file
- Tray icon badge showing time remaining or session number

### Settings-Based Customization
Allow users to configure visual preferences:
- Toggle between minutes-only and MM:SS display formats
- Adjust update frequency (5s, 10s, 30s options)
- Choose alternative emoji sets for session types
- Enable/disable progress bar if they prefer minimal UI

### Icon Design Improvements
Enhance the icon set for better visual distinction:
- Create higher-resolution icon assets
- Add subtle color variations (within grayscale or monochrome constraints for menu bar)
- Design animated icon sequences (if macOS supports, likely requires custom framework)

### Accessibility Features
Make visual feedback accessible to more users:
- Add text-based alternatives to emoji for screen reader compatibility
- Provide high-contrast icon variants
- Add audio feedback for state transitions (subtle clicks or chimes)

### Extended Progress Information
Show additional context in the dropdown:
- Estimated completion time (e.g., "Work ends at 3:45 PM")
- Time elapsed in current session (in addition to time remaining)
- Daily progress summary at bottom of dropdown

## PRD Entry (for docs/PRD.md)

### Increment: Visual Feedback with Progress Indicators

**Increment ID:** `visual-feedback-with-progress-indicators`

**Title:** Visual Feedback with Progress Indicators

**Status:** Proposed

**Increment Folder:** `docs/increments/visual-feedback-with-progress-indicators/`

**User Story:**
As a GoPomodoro user, I want to see clear visual indicators of my current session state and progress, so that I can quickly understand where I am in my pomodoro cycle without reading text carefully.

**Acceptance Criteria:**
- Tray icon changes per session type (idle, work, short break, long break, paused)
- Progress bar shows session advancement using Unicode characters (●●●●○○○○○○)
- Header text includes state-specific emoji (🍅 Work Session, ☕ Short Break, 🌟 Long Break, ⏸️ Paused)
- Timer displays minutes only in clean format (e.g., "12min")
- Timer updates every 10 seconds for CPU efficiency
- Cycle indicator uses visual symbols (🍅○○○)
- Paused state is visually distinct

**Use Case Summary:**
User launches app → sees idle icon → starts work session → tray icon changes to work state → header shows "🍅 Work Session" → timer displays "25min" → every 10 seconds timer updates and progress bar fills gradually (○○○○○ → ●○○○○ → ●●○○○ → etc.) → at completion, auto-transitions to break → icon changes to break state → cycle indicator updates → process continues through 4-session cycle. Paused state shows distinct "⏸️ Paused" appearance. Long break after 4 sessions shows "🌟 Long Break".
