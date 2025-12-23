# Increment: Tray Icon and Dropdown UI

## User Story

As a macOS user who wants to try GoPomodoro, I want to see the app's tray icon and open its dropdown panel, so that I can verify the app is running and see what controls will be available.

## Acceptance Criteria

1. **Tray icon is visible** — When the app launches, a tray icon appears in the macOS menu bar
2. **Tray icon shows idle state** — The icon displays in gray/neutral color indicating the app is ready but not running
3. **Dropdown opens on click** — Left-clicking the tray icon opens a dropdown panel below the icon
4. **Dropdown shows structure** — The panel displays:
   - A header area (where timer status will go)
   - A display area (where timer countdown will go)
   - Control buttons labeled "Start", "Pause", "Reset" (non-functional placeholders)
5. **Dropdown closes properly** — Clicking outside the dropdown or clicking the tray icon again closes the panel
6. **App runs in background** — The app has no dock icon, only exists as a tray application

## Use Case

### Actors
- **User**: A macOS user who has installed and launched GoPomodoro

### Preconditions
- GoPomodoro application is installed on macOS
- User has necessary permissions for the app to create a menu bar icon
- macOS system tray/menu bar is accessible and visible

### Main Flow

1. User launches the GoPomodoro application (e.g., from Applications folder or via command line)
2. System displays a gray/neutral tray icon in the macOS menu bar
3. Tray icon appears in the idle state (○ or similar visual indicator)
4. User left-clicks on the tray icon
5. System opens a dropdown panel directly below the tray icon
6. Dropdown panel displays:
   - Header section showing "Ready" or "Idle" state
   - Large central area reserved for timer display (showing placeholder like "25:00")
   - Three buttons in vertical layout:
     - "Start" button
     - "Pause" button (may be grayed out/disabled)
     - "Reset" button (may be grayed out/disabled)
7. User reviews the UI structure
8. User clicks outside the dropdown panel
9. System closes the dropdown panel
10. Tray icon remains visible in the menu bar

### Alternate / Exception Flows

**A1: User clicks tray icon to close dropdown**
- At step 7, instead of clicking outside:
  - User clicks the tray icon again
  - System closes the dropdown panel
  - Flow returns to step 10

**A2: User switches away from app**
- At step 7, user switches to another application or desktop space
- System automatically closes the dropdown panel
- Tray icon remains visible
- User can reopen dropdown by clicking the tray icon again

**A3: User attempts to interact with placeholder buttons**
- At step 7, user clicks "Start", "Pause", or "Reset" button
- System does nothing (buttons are non-functional in this increment)
- Dropdown remains open
- User can close dropdown as in main flow

**E1: App launch fails**
- At step 1, if the app cannot access menu bar permissions:
  - System may show an error or fail silently
  - No tray icon appears
  - (This increment assumes basic macOS permissions are granted; full error handling deferred)

## Context

GoPomodoro is a minimal macOS native pomodoro timer application that will live in the system tray. This is the **very first increment**, starting from zero code.

### Current Situation
- The project has a detailed PRD and constitution but no codebase yet
- We need to establish the foundational architecture for a macOS menu bar application
- The PRD envisions a full pomodoro cycle with timer logic, notifications, and persistence, but we're starting with just the UI shell

### Why This Matters
- Users need to be able to launch the app and see it's running before any timer functionality exists
- This increment validates our choice of Go UI framework and macOS tray integration approach
- It creates a testable foundation that later increments can build upon
- Following the constitution's "small, safe steps" principle, we separate UI structure from timer logic

### Key Constraints
- macOS-only for this increment
- Must run as a menu bar app (no dock icon)
- No timer functionality yet — purely structural UI
- Manual testing is acceptable per the lite constitution

## Goal

### Outcome
After this increment, users will be able to:
- Launch GoPomodoro and see it appear in their macOS menu bar
- Click the tray icon to open a dropdown panel
- See the structure of the UI (header, timer area, control buttons)
- Close the dropdown and verify the app remains running in the background

### Scope
This increment delivers:
- macOS menu bar / system tray integration
- A clickable tray icon with idle state visual
- A dropdown panel that opens/closes properly
- A structured layout with placeholder controls

### Non-Goals
This increment explicitly does **not** include:
- Timer countdown logic or state management
- Functional Start/Pause/Reset buttons
- Any pomodoro cycle logic
- Notifications
- Settings or persistence (no ~/.gopomodoro/ directory yet)
- Session tracking or statistics
- Color-coded states for work/break modes
- Multiple platform support

### Why This Is a Good Increment
- **Small and self-contained**: Focuses only on UI structure without business logic
- **Quickly testable**: Launch the app and verify the UI appears and behaves correctly
- **Low risk**: No data, no complex state, easy to iterate or change
- **Validates architecture early**: Confirms our Go framework choice works for macOS tray apps before investing in timer logic
- **Delivers visible progress**: Users (and stakeholders) can see the app taking shape
- **Aligns with constitution**: Follows "make it work, make it right" and "simple is better than complex"

## Tasks

### Task 1: App Launches as Menu Bar Application
- **User/Stakeholder Impact**: Users can start GoPomodoro and it appears in their macOS menu bar without cluttering the dock
- **Acceptance Clues**: 
  - Double-clicking the app or running it from terminal causes a tray icon to appear
  - No dock icon is visible when the app is running
  - App process is visible in Activity Monitor

### Task 2: Tray Icon Displays Idle State
- **User/Stakeholder Impact**: Users can see at a glance that the app is installed and running but not actively timing
- **Acceptance Clues**:
  - The menu bar shows a gray or neutral-colored icon (○ or similar)
  - Icon is clearly visible against both light and dark menu bar backgrounds
  - Icon is the correct size for macOS menu bar standards

### Task 3: Dropdown Panel Opens on Left-Click
- **User/Stakeholder Impact**: Users can access the app's controls with a single left-click
- **Acceptance Clues**:
  - Left-clicking the tray icon causes a panel to appear directly below the icon
  - The panel opens quickly (< 200ms perceived delay)
  - Only left-click triggers the dropdown (right-click does nothing in this increment)

### Task 4: Dropdown Shows Structured Layout with Placeholders
- **User/Stakeholder Impact**: Users can preview what the final timer interface will look like
- **Acceptance Clues**:
  - Panel displays a header area (may say "Ready" or "Idle")
  - Panel shows a large central area with placeholder timer display (e.g., "25:00")
  - Panel shows three clearly labeled buttons: "Start", "Pause", "Reset"
  - Layout is clean, readable, and matches the visual direction from the PRD
  - Colors and spacing are reasonable placeholders (don't need to be final)

### Task 5: Dropdown Closes Properly
- **User/Stakeholder Impact**: Users can dismiss the dropdown in natural ways
- **Acceptance Clues**:
  - Clicking anywhere outside the dropdown panel causes it to close
  - Clicking the tray icon again toggles the dropdown closed
  - Switching to another app or desktop space closes the dropdown
  - The dropdown doesn't get "stuck" open or require special actions to close

### Task 6: App Can Be Quit Cleanly
- **User/Stakeholder Impact**: Users can stop the app when they're done testing or using it
- **Acceptance Clues**:
  - Force-quitting from Activity Monitor works
  - Eventually, a "Quit" option or standard macOS quit mechanism should work
  - No orphaned processes remain after quit
  - App doesn't crash or hang on quit

### Task 7: Basic Visual Polish for Placeholder UI
- **User/Stakeholder Impact**: The app feels credible and not broken, even though functionality is limited
- **Acceptance Clues**:
  - Buttons look like buttons (not just text)
  - Layout doesn't overlap or appear broken
  - Font sizes are readable
  - Overall appearance is clean enough to share a screenshot

## Risks and Assumptions

### Risks

1. **Go UI Framework Choice**
   - **Risk**: We might choose a Go framework for macOS tray integration that has limitations we discover later
   - **Mitigation**: This increment is small enough to rewrite if we need to switch frameworks early

2. **macOS Permissions and Sandboxing**
   - **Risk**: macOS security features might require additional permissions or entitlements we haven't anticipated
   - **Mitigation**: Start with simplest approach (non-sandboxed, local development); address distribution concerns in later increments

3. **Visual Expectations**
   - **Risk**: The placeholder UI might set incorrect expectations about final visual polish
   - **Mitigation**: Keep placeholders clearly "in progress" looking; iterate on design in future increments

4. **Cross-Platform Complexity**
   - **Risk**: Choosing a very macOS-specific approach might make Linux/Windows support harder later
   - **Mitigation**: The PRD deprioritizes other platforms; macOS-first is acceptable for now

### Assumptions

1. **Development Environment**: Developer has a macOS machine available for building and testing
2. **Go Ecosystem**: A suitable Go library exists for macOS menu bar integration (e.g., systray, fyne, wails, or native bindings)
3. **Permissions**: Basic app launch permissions will be granted by the OS during development
4. **No Backend Required**: This increment needs no server, database, or external services
5. **Manual Testing Sufficient**: Per the lite constitution, automated UI tests are not required for this increment

## Success Criteria and Observability

### Success Criteria

This increment is successful when:

1. **Functional Success**:
   - The app launches and a tray icon appears in the macOS menu bar
   - Left-clicking the icon opens a dropdown panel
   - The panel shows placeholder UI elements (header, timer area, buttons)
   - Clicking outside or re-clicking the icon closes the dropdown
   - The app can be quit without crashing

2. **User Experience Success**:
   - A colleague or test user can launch the app and understand it's a timer app (even without functionality)
   - The UI doesn't look broken or confusing
   - Opening and closing the dropdown feels responsive

3. **Technical Success**:
   - Code compiles and runs on macOS (developer's machine at minimum)
   - No crashes or hangs during basic interaction
   - Clean app shutdown with no orphaned processes

### Observability

**How we will observe success after this increment:**

1. **Manual Testing Checklist**:
   - Launch app → tray icon visible? ✓
   - Click icon → dropdown opens? ✓
   - Click outside → dropdown closes? ✓
   - UI shows header, timer placeholder, buttons? ✓
   - Quit app → clean shutdown? ✓

2. **Visual Confirmation**:
   - Take a screenshot of the tray icon and dropdown
   - Compare against PRD mockup for basic structural match
   - Share screenshot with stakeholders for quick validation

3. **Developer Experience**:
   - Note which Go framework/library was chosen and how straightforward the integration was
   - Document any gotchas or unexpected macOS behaviors for future reference

4. **Process Check**:
   - Did this increment take hours/days (good) or weeks (too big)?
   - Was it easy to understand and implement without extensive design work?
   - Can we confidently move to the next increment (timer logic)?

**What to look for:**
- App appears in Activity Monitor when running
- Memory usage is reasonable (< 50MB per PRD)
- No error logs or crashes in Console.app
- Tray icon is clearly visible in both light and dark macOS themes

## Process Notes

### How This Increment Should Move Through the Workflow

1. **Small, Safe Changes**:
   - Start by getting the simplest possible tray icon to appear
   - Then add the dropdown panel
   - Then add placeholder UI elements
   - Commit working increments along the way

2. **Testing Approach**:
   - Manual testing is sufficient for this increment (per lite constitution)
   - Test on the developer's macOS machine
   - Verify basic interactions work before considering it done

3. **No Special Deployment**:
   - This is a local development increment
   - No need for CI/CD, app signing, or distribution yet
   - Running from `go run` or a local binary is fine

4. **Quick Feedback Loop**:
   - Aim to have something visible within the first day of work
   - Iterate quickly on the UI layout based on what looks reasonable
   - Don't over-engineer — placeholders are acceptable

5. **Constitution Alignment**:
   - This follows "make it work, make it right" — get the UI working first
   - This follows "simple is better than complex" — no fancy frameworks needed unless they're clearly simpler
   - This follows "small, safe steps" — UI structure before timer logic

## Follow-up Increments

### Immediate Next Increment
**Timer Countdown Logic and Display**
- Add actual 25-minute countdown timer
- Update the timer display in real-time (every second or every 30 seconds per PRD)
- Make timer state (idle, running, paused) affect what's displayed

### Subsequent Increments
**Functional Control Buttons**
- Wire up Start button to begin countdown
- Wire up Pause button to pause/resume
- Wire up Reset button to reset to 25:00

**Session Complete Detection**
- Detect when timer reaches 00:00
- Show some indication that the session is complete
- Reset state back to idle

**Config Directory and Settings**
- Create `~/.gopomodoro/` directory on first run
- Add basic settings.json support for timer durations

**Full Pomodoro Cycle**
- Implement work → short break → long break cycle
- Track which session in the cycle the user is on

**Notifications**
- Add native macOS notifications for session complete
- Support notification sounds

**Visual Polish and Color Coding**
- Implement the full color scheme from the PRD
- Add progress bar visualization
- Add cycle indicator (🍅🍅○○)

## PRD Entry (for docs/PRD.md)

```markdown
### Increment: Tray Icon and Dropdown UI

- **Increment ID**: `tray-icon-and-dropdown-ui`
- **Title**: Tray Icon and Dropdown UI
- **Status**: Proposed
- **Increment Folder**: `increments/tray-icon-and-dropdown-ui/`

**User Story**:  
As a macOS user who wants to try GoPomodoro, I want to see the app's tray icon and open its dropdown panel, so that I can verify the app is running and see what controls will be available.

**Acceptance Criteria**:
- Tray icon is visible in macOS menu bar when app launches
- Tray icon shows idle state (gray/neutral color)
- Left-clicking the tray icon opens a dropdown panel
- Dropdown displays header area, timer placeholder, and control buttons (Start, Pause, Reset)
- Dropdown closes when clicking outside or re-clicking the icon
- App runs as background-only application (no dock icon)

**Use Case Summary**:
User launches app → tray icon appears in menu bar → user clicks icon → dropdown panel opens showing placeholder UI structure → user clicks outside → dropdown closes. App remains running in background. Placeholder buttons are non-functional in this increment.
```
