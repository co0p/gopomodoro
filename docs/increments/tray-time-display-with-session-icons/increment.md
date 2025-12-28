# Increment: Tray Time Display with Session Icons

## User Story

As a GoPomodoro user, I want to see the remaining time directly in the tray icon, so that I can track my session progress without opening the dropdown menu.

## Acceptance Criteria

1. **Tomato icon at startup** — When the app launches in idle state, the tray icon displays a tomato (🍅) instead of a gray/neutral icon

2. **Session icon and time visible in tray during sessions** — When a timer is running or paused, the tray displays the session icon and remaining time (e.g., "🍅 24m" for work, "☕ 4m" for short break, "🌟 14m" for long break)

3. **Session icon and time visible in tray when idle** — When the timer is idle/ready, the tray displays the session icon and default duration (e.g., "🍅 25m" for work, "☕ 5m" for short break)

4. **Paused state shows pause icon** — When the timer is paused, the tray displays the pause icon and remaining time (e.g., "⏸️ 12m")

5. **No state header in dropdown** — The dropdown menu does not contain a menu item showing the current state (e.g., "🍅 Work Session", "Ready", "⏸️ Paused")

6. **No timer display in dropdown** — The dropdown menu does not contain a separate menu item showing the time countdown (e.g., "25m", "24m")

7. **Progress bar still visible** — The progress bar menu item remains in the dropdown and continues to function as before

8. **Cycle indicator still visible** — The cycle indicator menu item (e.g., "Session 1/4 🍅○○○") remains in the dropdown and continues to function as before

## Use Case

### Actors
- **User**: A GoPomodoro user who wants to track time without opening the dropdown
- **System**: The GoPomodoro application managing tray display and timer state

### Preconditions
- GoPomodoro is installed and running on macOS
- The app appears in the system tray/menu bar
- Timer logic and session management are functional (from previous increments)

### Main Flow

1. User launches GoPomodoro application
2. System displays tray icon with tomato emoji and default work duration: "🍅 25m"
3. User clicks tray icon to open dropdown menu
4. System displays dropdown menu containing:
   - Progress bar (○○○○○○○○○○)
   - Cycle indicator (Session 1/4 🍅○○○)
   - Start button
   - Pause button (disabled)
   - Reset button (disabled)
   - Skip button
   - Quit button
5. User notices that state header and timer display are not present in dropdown
6. User clicks "Start" button
7. System starts work session timer and updates tray to show: "🍅 24m"
8. System continues counting down, updating tray display: "🍅 23m", "🍅 22m", etc.
9. User glances at menu bar and sees current time without opening dropdown
10. Work session completes (reaches 0m)
11. System transitions to short break and updates tray to show: "☕ 5m"
12. User continues through pomodoro cycle, always seeing current session icon and time in tray

### Alternate / Exception Flows

**A1: User pauses during work session**
- At step 8, user opens dropdown and clicks "Pause"
- System updates tray to show pause icon with frozen time: "⏸️ 18m"
- Time in tray remains static (does not count down)
- User clicks "Start" to resume
- System updates tray back to work session: "🍅 18m" and resumes countdown

**A2: User skips to next session**
- At step 8, user opens dropdown and clicks "Skip"
- System advances to next session type (e.g., short break)
- System updates tray to show new session icon and duration: "☕ 5m"
- Cycle indicator advances in dropdown menu

**A3: User resets during active session**
- At step 8, user opens dropdown and clicks "Reset"
- System resets to idle state with work session ready
- System updates tray to show: "🍅 25m"
- Cycle indicator resets to "Session 1/4 🍅○○○"

**A4: Long break transition**
- After completing 4 work sessions
- System transitions to long break
- System updates tray to show: "🌟 15m"
- User sees different icon indicating special break type

**A5: User quits and relaunches**
- User quits application via "Quit" button
- User relaunches GoPomodoro
- System displays initial state: "🍅 25m" in tray
- All menu items properly initialized without state header or timer display

## Context

GoPomodoro is a minimal macOS native pomodoro timer that lives in the system tray. The application has completed four increments establishing functional timer logic, pomodoro cycle management, tray icon basics, and visual feedback through progress indicators.

### Current Situation

Currently, the tray icon shows different static images based on state (idle, work, break, paused) but does not display numeric time information. The remaining time is shown in the dropdown menu as a separate menu item (e.g., "25m", "24m"). The dropdown also contains a state header menu item (e.g., "🍅 Work Session", "Ready", "⏸️ Paused").

Users must click the tray icon and open the dropdown menu to see how much time remains in the current session. For users checking time frequently, this adds unnecessary interaction steps.

### Why This Matters

- **Reduces interaction friction**: Users can track time with a glance rather than a click
- **Better use of menu bar real estate**: The tray area is always visible, making it ideal for time display
- **Simplifies dropdown menu**: Removing redundant state/time information makes the menu cleaner and easier to scan
- **Improves session awareness**: The session icon (🍅/☕/🌟/⏸️) provides instant context about the current activity type

### Key Constraints

- Must work within macOS menu bar display limitations (length, font rendering)
- Time format should be concise (using "m" notation: "24m" instead of "24:00")
- Changes should not affect existing timer logic or session management
- Progress bar and cycle indicator must remain functional in dropdown

## Goal

### Outcome

After this increment, users will be able to:
- See the remaining session time and session type in the menu bar without any interaction
- Understand at a glance which type of session is active (work, short break, long break, or paused)
- Access a cleaner, more focused dropdown menu with less redundant information

The system will display session-specific emoji icons (🍅/☕/🌟/⏸️) alongside time in minutes format (e.g., "25m") directly in the menu bar, updating every minute as the timer counts down.

### Scope

This increment will:
- Display session icon and remaining time in the menu bar tray title
- Change the idle state icon from neutral/gray to tomato (🍅)
- Remove the state header menu item from the dropdown
- Remove the separate timer display menu item from the dropdown
- Keep all existing timer functionality, buttons, progress bar, and cycle indicator unchanged

### Non-Goals

This increment will not:
- Customize the time display format (stays as "Xm")
- Add user preferences for tray display options
- Show seconds precision in the tray
- Add tooltip customization beyond the standard behavior
- Support displaying multiple timers or complex time information
- Change any timer logic, session duration calculations, or notification behavior

### Why This Is a Good Increment

- **Small and focused**: Only affects tray display and menu item configuration, not core logic
- **Easy to verify**: Visual inspection immediately confirms the changes work correctly
- **Low risk**: Does not modify timer calculations, state management, or data persistence
- **Immediate value**: Users get better time visibility with their next session
- **Reversible**: Changes are UI-only and can be adjusted in future increments if needed

## Tasks

### Task 1: Tray displays session icon and time in all states

**User/Stakeholder Impact:**  
Users see "🍅 25m" / "☕ 5m" / "🌟 15m" / "⏸️ 12m" in the menu bar depending on current session type and state. The display updates as time counts down.

**Acceptance Clues:**
- Visual inspection shows emoji icon followed by time in "Xm" format
- Tray updates every minute during active timer countdown
- Different emoji appears for each session type (work/short break/long break/paused)
- Time format is consistent across all states

### Task 2: Idle state shows tomato instead of neutral icon

**User/Stakeholder Impact:**  
When users launch the app or reset to idle, they see a tomato emoji (🍅) with "25m" rather than a gray/neutral icon, clearly signaling "ready to start work session."

**Acceptance Clues:**
- Fresh app launch displays "🍅 25m" in tray
- After reset button, tray shows "🍅 25m"
- After completing long break (cycle restart), tray shows "🍅 25m"

### Task 3: Dropdown menu no longer shows state header

**User/Stakeholder Impact:**  
Users see a cleaner dropdown menu without the redundant state header (like "🍅 Work Session" or "Ready"). The progress bar becomes the first visible item.

**Acceptance Clues:**
- Opening the dropdown shows progress bar as the first menu item
- No menu item displays session state text or emoji header
- Menu structure flows directly from tray to progress indicators and buttons

### Task 4: Dropdown menu no longer shows timer display

**User/Stakeholder Impact:**  
Users see a streamlined dropdown without a separate countdown timer menu item (like "25m" or "24m"), since time is now visible in the tray itself.

**Acceptance Clues:**
- Opening the dropdown shows no dedicated timer display menu item
- Menu structure shows progress bar, cycle indicator, then buttons
- No numeric time countdown appears in the dropdown

### Task 5: Progress bar and cycle indicator remain functional

**User/Stakeholder Impact:**  
Users continue to see visual session progress (○○○○○●●●●●) and cycle position (Session 2/4 🍅🍅○○) in the dropdown menu exactly as before.

**Acceptance Clues:**
- Progress bar menu item is present and animates during sessions
- Cycle indicator menu item shows correct session count and tomato markers
- Both items update at appropriate times (progress bar every minute, cycle indicator on session transitions)

## Risks and Assumptions

### Risks

- **Text readability**: Time text in the menu bar might be too small or hard to read on certain displays or with certain accessibility settings. Some users may prefer the larger dropdown timer display.

- **Emoji rendering variations**: Emoji appearance (🍅/☕/🌟/⏸️) may vary across macOS versions, potentially affecting visual consistency or clarity.

- **Menu bar space**: On smaller displays or when many menu bar items are present, the tray title might be truncated or hidden.

- **User preference mismatch**: Some users may prefer a minimal tray icon without text, finding the time display cluttered.

### Assumptions

- Users prefer glanceable information in the menu bar over clicking to see time
- The "Xm" format (e.g., "24m") is sufficiently clear and concise for time tracking
- Session emoji icons are visually distinct enough to recognize at menu bar size
- Most users have enough menu bar space to display icon + time without truncation

### Mitigations

- If readability becomes an issue, a future increment could add user preferences for tray display format
- Manual testing across different macOS versions can validate emoji rendering
- The change is reversible—can restore dropdown timer if users report preference issues

## Success Criteria and Observability

### Success Criteria

This increment succeeds when:
- Users can determine remaining time and session type without opening the dropdown menu
- The menu bar tray shows accurate, updating time information throughout a full pomodoro cycle
- The dropdown menu is visually cleaner with fewer redundant items
- All existing functionality (start, pause, reset, skip, progress tracking) continues to work unchanged

### Observability

**What to check after release:**
- Visual inspection during manual testing of a complete pomodoro cycle (work → break → work → break → work → break → work → long break)
- Verify tray displays correct format in all states: idle ("🍅 25m"), running ("🍅 24m"), paused ("⏸️ 18m"), different session types ("☕ 5m", "🌟 15m")
- Confirm dropdown menu shows only: progress bar, cycle indicator, and buttons (no header, no timer display)
- Verify tray time updates every minute during active countdown
- Check that emoji icons render clearly on the test macOS version

**Evidence to collect:**
- Screenshot or screen recording of full cycle showing tray updates
- Visual confirmation that dropdown no longer contains header or timer items
- Confirmation that progress bar and cycle indicator still function correctly

**Where to look:**
- macOS menu bar tray area (primary observation point)
- Dropdown menu structure (verify removal of items)
- Timer behavior logs (confirm timer logic unchanged)

## Process Notes

This increment should:
- Be implemented through small, focused code changes to tray display logic and menu initialization
- Remove the menu item initialization calls for header and timer display
- Update the tray title/text whenever timer state or time changes
- Adjust the idle state icon selection to use tomato instead of neutral icon
- Be tested manually through a complete pomodoro cycle to verify all states

The work should:
- Flow through the normal build and test process defined in the project
- Preserve all existing timer logic, state management, and event handling
- Not require database or file format changes
- Not affect session logging or statistics tracking

Manual testing is sufficient per the project's lite constitution. Focus testing on:
- Launch and idle state appearance
- Transition through all session types (work, short break, long break)
- Pause and resume behavior
- Skip and reset button effects on tray display
- Menu structure after changes

## Follow-up Increments

### Tray Display Customization
Allow users to configure tray display preferences via settings file, such as:
- Showing/hiding emoji icons
- Choosing between "Xm" and "XX:XX" time formats
- Enabling/disabling tray time display entirely (icon only)

### Seconds Precision Option
Add optional seconds display in tray for users who want more precise time tracking (e.g., "🍅 24m 37s" or "🍅 24:37"), configurable via settings.

### Tray Icon Animation
Add subtle visual animation or color changes to the tray icon during the final minute of a session to draw attention as time expires.

### Menu Bar Space Optimization
Detect available menu bar space and automatically shorten or hide tray text when space is limited, falling back to icon-only display.

## PRD Entry (for docs/PRD.md)

**Increment ID:** `tray-time-display-with-session-icons`

**Title:** Tray Time Display with Session Icons

**Status:** Proposed

**Increment Folder:** `docs/increments/tray-time-display-with-session-icons/`

**User Story:**  
As a GoPomodoro user, I want to see the remaining time directly in the tray icon, so that I can track my session progress without opening the dropdown menu.

**Acceptance Criteria:**
- Tomato icon at startup — App launches with "🍅 25m" instead of neutral/gray icon
- Session icon and time visible in tray during sessions — Running timer shows "🍅 24m", "☕ 4m", "🌟 14m", etc.
- Session icon and time visible in tray when idle — Idle state shows "🍅 25m" or appropriate session default
- Paused state shows pause icon — Paused timer shows "⏸️ 12m" with frozen time
- No state header in dropdown — Menu does not contain state/session header item
- No timer display in dropdown — Menu does not contain separate timer countdown item
- Progress bar still visible — Progress bar menu item remains and functions normally
- Cycle indicator still visible — Cycle indicator menu item remains and functions normally

**Use Case Summary:**  
User launches app and sees "🍅 25m" in menu bar. User starts timer and watches tray update to "🍅 24m", "🍅 23m", etc. User glances at menu bar to check time without opening dropdown. On completion, tray transitions to next session type (e.g., "☕ 5m"). Dropdown menu shows progress bar and cycle indicator but no longer includes redundant header or timer items.
