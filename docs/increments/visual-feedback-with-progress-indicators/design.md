# Design: Visual Feedback with Progress Indicators

## Context and Problem

GoPomodoro has completed its foundational phases. Phase 1 (MVP Core) delivered a functional timer with system tray integration. Phase 2 (Pomodoro Cycle) implemented the complete 4-session workflow with automatic transitions between work and break periods. The app works reliably, and users can start, pause, reset, and skip through pomodoro sessions.

### Current State

The UI is entirely text-based, relying on `systray` menu items:
- Generic text headers: "Work Session", "Ready for Break"
- Timer in MM:SS format: "24:37"
- Cycle indicator with emoji: "Session 2/4  🍅🍅○○"
- Timer updates every second (1-second tick interval)

### The Problem

Users must carefully read and interpret text to understand their current state. There's no immediate visual distinction between a work session and a break at a glance. Progress through a session isn't visible—users must mentally calculate what "24:37 remaining" means in terms of completion percentage. The MM:SS format provides more precision than needed for 25-minute sessions, and every-second updates consume CPU unnecessarily.

### Why Now

With core functionality stable and tested, this is the right moment to improve user experience. Visual feedback transforms the app from "functional but plain" to "functional and pleasant to use." This increment delivers clear value without touching the stable timer and session logic. The changes are purely presentational, making them low-risk and easily reversible.

### Components in Play

- **UI Package** ([internal/ui/window.go](../../../internal/ui/window.go)): Currently manages dropdown menu rendering, timer display formatting, button states, and UI update coordination
- **Tray Package** ([internal/tray/tray.go](../../../internal/tray/tray.go)): Already has complete icon infrastructure with `UpdateIcon()` method and state-based icon selection
- **Timer Package** ([internal/timer/timer.go](../../../internal/timer/timer.go)): Maintains 1-second tick interval and fires callbacks on each tick
- **Session Package** ([internal/session/session.go](../../../internal/session/session.go)): Already has cycle indicator formatting with emoji tomatoes

### References

- [increment.md](./increment.md) — Product-level requirements for this increment
- [CONSTITUTION.md](../../../CONSTITUTION.md) — Project values, principles, and lite-mode expectations
- [ARCHITECTURE.md](../../../ARCHITECTURE.md) — System architecture and component boundaries
- [PRD.md](../../../PRD.md) — Phase 3 Visual Enhancements section

---

## Proposed Solution (Technical Overview)

Enhance the visual presentation within `systray` library constraints using Unicode characters, emoji, and icon file swapping. The solution focuses on five key improvements:

1. **Progress Bar Visualization**: Add a new menu item that displays a 10-segment Unicode progress bar (●●●●●○○○○○) that fills proportionally as the session advances
2. **Minutes-Only Timer Display**: Modify the `formatTime()` function to return "Xmin" format instead of "MM:SS"
3. **10-Second Update Frequency**: Change the timer's tick interval from 1 second to 10 seconds
4. **Session-Specific Header Emoji**: Enhance existing header updates to include emoji for each state (🍅 Work Session, ☕ Short Break, 🌟 Long Break, ⏸️ Paused)
5. **Tray Icon Updates**: Verify and integrate existing tray icon system (already implemented in Phase 2)

### Components Involved

**UI Package** (primary changes):
- Add new `progressBar` menu item to display visual progress
- Modify `formatTime()` to return minutes-only format
- Add `formatProgressBar()` function to render Unicode progress segments
- Update header text in `handleTimerStarted()` and pause/resume handlers to include emoji
- Modify `handleTimerTick()` to update progress bar alongside timer display
- Add `sessionStartTime` and `sessionDuration` fields to track progress calculation

**Timer Package** (minor adjustment):
- Change `tickInterval` constant from `1 * time.Second` to `10 * time.Second`

**Tray Package** (no changes needed):
- Already complete with `UpdateIcon()` method
- Icon assets already exist for all states

**Session Package** (no changes needed):
- `FormatCycleIndicator()` already uses emoji tomatoes correctly

### Typical Flow After Changes

1. User clicks tray icon to open dropdown
2. Dropdown displays:
   - Header: "🍅 Work Session" (with emoji)
   - Timer: "12min" (minutes only)
   - Progress bar: "●●●●●○○○○○" (50% filled for halfway through session)
   - Cycle indicator: "Session 2/4  🍅🍅○○" (unchanged, already has emoji)
3. Every 10 seconds:
   - Timer updates: "12min" → "11min"
   - Progress bar updates: "●●●●●○○○○○" → "●●●●●●○○○○"
4. User pauses:
   - Header changes to "⏸️ Paused"
   - Tray icon changes to paused icon
   - Timer and progress bar freeze at current values
5. Session completes:
   - Progress bar shows fully filled: "●●●●●●●●●●"
   - Auto-transitions to next session
   - Header, timer, and icon all update for new session type

---

## Scope and Non-Scope (Technical)

### In-Scope

**Progress Bar Visualization:**
- New menu item positioned between timer display and cycle indicator
- Calculates fill percentage based on elapsed time vs total session duration
- Renders 10 Unicode circle segments (filled ● or empty ○)
- Updates every 10 seconds alongside timer
- Resets to empty on session transitions

**Minutes-Only Timer Display:**
- Modify `formatTime()` function signature unchanged, behavior changed
- Returns "0min" through "25min" format
- Handles edge case: when remaining < 60 seconds, displays "0min"

**10-Second Update Frequency:**
- Change `tickInterval` in timer package from 1s to 10s
- All tick callbacks fire every 10 seconds instead of every second
- No changes to callback signatures or state management

**Session-Specific Header Emoji:**
- Update header text in `handleTimerStarted()` for work/short break/long break
- Update header in pause handler to show "⏸️ Paused"
- Update header in resume handler to restore session-specific emoji
- Idle state shows "○ Idle" or "Ready"

**Tray Icon Verification:**
- Verify `updateTrayIcon()` is called in all appropriate state transitions
- Verify icon assets exist for all states (already confirmed)
- No new code needed—integration already complete from Phase 2

**Paused State Visual Distinction:**
- Header shows "⏸️ Paused" when `timer.StatePaused`
- Tray icon uses paused-specific image
- Timer value freezes (already works, no change needed)
- Progress bar maintains fill state (no reset on pause)

### Out-of-Scope

- Migrating to different UI framework (staying with `systray`)
- Color-coded backgrounds or rich text formatting (not possible with `systray`)
- Dark mode support (deferred to future increment)
- Smooth animations or transitions (not feasible with static menu items)
- Changes to timer state machine logic
- Changes to session cycle rules or duration logic
- Changes to data persistence or logging format
- Settings-based customization of update frequency or display format
- Accessibility features (screen reader support, high contrast)

---

## Architecture and Boundaries

### Component Architecture

No new components are introduced. Changes are confined to existing packages:

```
UI Package (internal/ui/window.go)
├─ Owns presentation logic
├─ Add: progressBar menu item
├─ Modify: formatTime() behavior
├─ Modify: header text updates with emoji
└─ Modify: tick handler to update progress bar

Timer Package (internal/timer/timer.go)
├─ Owns timing and state management
└─ Modify: tickInterval constant (implementation detail)

Tray Package (internal/tray/tray.go)
└─ No changes (already complete)

Session Package (internal/session/session.go)
└─ No changes (already complete)
```

### Boundaries Respected

**Separation of Concerns:**
- UI Package continues to own all presentation logic and formatting
- Timer Package continues to own timing state and tick generation
- Session Package remains pure business logic with no UI dependencies
- Tray Package maintains clean interface via `UpdateIcon(sessionType, state)`

**Layering:**
- No circular dependencies introduced
- UI depends on Timer and Session (already established)
- Timer and Session remain independent of UI
- Tray is called by UI but doesn't call back

**Data Flow:**
- Timer generates tick events → UI updates display
- Session provides state → UI renders appropriate visuals
- UI commands timer (start/pause/reset) → Timer updates state → UI reflects changes

### Guardrails from Constitution

**"Simple Is Better Than Complex":**
- Using Unicode characters rather than attempting custom rendering
- 10-segment progress bar (simple math, no complex interpolation)
- Minutes-only format (simpler than MM:SS for this use case)

**"Small, Safe Steps":**
- Each change can be implemented and tested independently
- Progress bar can be added without touching timer display
- Timer format can be changed without affecting progress bar
- Tick interval change is isolated to timer package

**"Make It Work, Make It Right":**
- Initial implementation can hardcode 10 segments
- Can refactor progress bar calculation later if needed
- Format functions can be extracted if they grow complex

---

## Contracts and Data

### No New Contracts

All changes are internal to UI and Timer packages. No new public APIs, events, or data structures are introduced to other packages.

### Modified Behaviors (Backward Compatible)

**Timer Tick Callback:**
- **Before**: Fires every 1 second with remaining time
- **After**: Fires every 10 seconds with remaining time
- **Compatibility**: Consumers (UI package) unaffected—they receive the same data, just less frequently
- **Contract preserved**: `onTick(remaining int)` signature unchanged

**Timer Display Format:**
- **Before**: `formatTime(1500)` returns "25:00"
- **After**: `formatTime(1500)` returns "25min"
- **Scope**: Internal to UI package, no external consumers
- **Contract**: Function signature `formatTime(seconds int) string` unchanged

### Progress Bar Logic

**Calculation:**
```
elapsed = sessionDuration - remaining
fillPercentage = elapsed / sessionDuration
filledSegments = round(fillPercentage * 10)
```

**Rendering:**
- For `filledSegments` from 0 to 10:
  - Append `filledSegments` × "●" (U+25CF Black Circle)
  - Append `(10 - filledSegments)` × "○" (U+25CB White Circle)
- Result: Always exactly 10 characters

**Edge Cases:**
- Session start (remaining = duration): 0 filled segments → "○○○○○○○○○○"
- Session end (remaining = 0): 10 filled segments → "●●●●●●●●●●"
- Paused state: Preserve current fill, no recalculation
- Resume state: Continue filling from current position

### Data Model Unchanged

- Session state (`CurrentType`, `CompletedWorkSessions`) unchanged
- Timer state (`state`, `remaining`, `sessionType`) unchanged
- Storage format (CSV logging) unchanged
- No migration or versioning needed

---

## Testing and Safety Net

### Testing Strategy (Aligned with Lite Constitution)

Per the constitution's lite mode, manual testing is appropriate for visual changes. Automated tests are focused on core logic only.

### Manual Visual Verification (Primary Method)

**Progress Bar Rendering:**
- Verify progress bar menu item appears in correct position
- Check that 10 segments render clearly (10 circles visible)
- Confirm filled circles (●) are visually distinct from empty (○)
- Verify progress bar fills proportionally during active session
- Check that progress bar doesn't overflow beyond 10 segments

**Timer Format:**
- Verify timer displays "25min" at session start
- Confirm minutes decrement correctly: "25min" → "24min" → "23min"
- Check edge case: final minute shows "0min" (not negative)
- Verify format is consistent across all session types

**Update Frequency:**
- Observe that timer updates occur approximately every 10 seconds
- Verify progress bar updates in sync with timer
- Confirm no perceived lag or "frozen" UI
- Check that updates feel responsive enough for 25-minute sessions

**Header Emoji:**
- Work session shows: "🍅 Work Session"
- Short break shows: "☕ Short Break"
- Long break shows: "🌟 Long Break"
- Paused state shows: "⏸️ Paused"
- Idle state shows: "Ready" or "○ Idle"

**Tray Icon Changes:**
- Verify different icons appear in macOS menu bar for each state
- Confirm icon changes occur in sync with session transitions
- Check that paused icon is visually distinct

**State Transitions:**
- Start → Running: Header + icon + progress bar all update
- Running → Paused: Header shows paused emoji, progress bar freezes
- Paused → Running: Header restores session emoji, progress bar continues
- Complete → Next Session: Progress bar resets to empty, timer resets

### Minimal Automated Tests

**Unit Tests for Format Functions:**

Test `formatTime()` behavior:
- `formatTime(1500)` returns `"25min"`
- `formatTime(720)` returns `"12min"`
- `formatTime(60)` returns `"1min"`
- `formatTime(59)` returns `"0min"`
- `formatTime(0)` returns `"0min"`

Test `formatProgressBar()` calculation:
- 0% progress (elapsed=0, duration=1500): "○○○○○○○○○○"
- 10% progress (elapsed=150, duration=1500): "●○○○○○○○○○"
- 50% progress (elapsed=750, duration=1500): "●●●●●○○○○○"
- 100% progress (elapsed=1500, duration=1500): "●●●●●●●●●●"
- Edge: slightly over 100% should cap at 10 segments

**No Integration Tests Required:**

For this lite-mode project, manual end-to-end testing through complete pomodoro cycles is sufficient and more appropriate than automated UI testing.

### Edge Cases to Verify Manually

- Timer showing "0min" during final 59 seconds
- Progress bar at 100% before transition (should show 10 filled segments, no overflow)
- Rapid session transitions via skip button (verify all visuals update correctly)
- Pause mid-session, wait, resume (verify progress continues from paused point)
- Multiple cycles in one sitting (verify visual consistency across cycles)

### No Regression Testing Needed

Since changes are purely presentational:
- Existing timer functionality (Start, Pause, Resume, Reset, Skip) continues working
- Session cycle logic unchanged (short breaks, long breaks, auto-transitions)
- Session logging still records events correctly
- Button enable/disable logic unchanged

If any behavioral regression appears during manual testing, it indicates a bug in implementation (not expected given scope).

---

## CI/CD and Rollout

### CI Implications

**No CI pipeline exists** per the constitution's lite mode for this single-developer project. If CI is added in the future, no special steps are required for this increment.

### Build Process

Standard build process unchanged:
```bash
make build
```

No new dependencies, no build configuration changes needed.

### Rollout Plan

**Development and Testing:**
1. Implement changes incrementally (progress bar, then timer format, then emoji, then tick interval)
2. Test each piece independently before integrating
3. Run through at least one complete 4-session pomodoro cycle manually
4. Verify all states: idle, work, short break, long break, paused
5. Test pause/resume and skip functionality
6. Check performance with Activity Monitor (verify CPU usage decrease)

**Deployment:**
- Build locally: `make build`
- Copy binary to `bin/gopomodoro`
- Replace running instance (quit and restart)
- No data migration needed
- No configuration changes required

**Staged Approach (Optional):**
- Can test during one work session before committing to full day
- Can run old version in parallel initially (different binary name)
- Low-risk changes allow immediate adoption

### Rollback Plan

**If visual issues appear:**
- Revert to previous version (simple binary replacement)
- Git revert the commit
- Rebuild and redeploy

**If rendering issues on specific macOS versions:**
- Can adjust Unicode characters (try different circles: ⬤⚪, ⚫⚪, 🔴⚪)
- Can fall back to ASCII characters (XX____) if Unicode fails
- Can disable progress bar by commenting out menu item

**If 10s updates feel too slow:**
- Can adjust `tickInterval` to 5s or 7s
- Simple one-line constant change
- Rebuild and test

**No Data Concerns:**
- No persistence format changes
- No backward compatibility issues with session logs
- No settings migration needed

### Manual Steps

**None beyond standard build:**
- No database migrations
- No configuration file updates
- No environment variable changes
- No cleanup scripts needed

---

## Observability and Operations

### Logging (Lite Expectations)

**No new logging required** for this increment:
- Visual changes don't generate events worth logging
- Existing session logging (started/completed/skipped) remains unchanged
- Timer state transitions already logged in previous phases

**Optional debug logging during development:**
- Can add temporary logs for progress bar segment calculation
- Can log tick interval changes for verification
- Remove or comment out before final release

### Performance Observability

**CPU Usage Verification:**
- Use macOS Activity Monitor to observe GoPomodoro process
- Compare CPU usage before/after tick interval change
- Expected result: Measurably lower CPU usage during active sessions (10x fewer ticks)
- No specific metrics collection needed—visual observation sufficient

**Memory:**
- No memory impact expected (one additional menu item, negligible)
- No monitoring needed

### Operational Considerations

**Unicode and Emoji Rendering:**
- Dependent on macOS system fonts
- Assumes macOS 10.12+ with full emoji support
- If rendering issues appear on older macOS: can substitute ASCII alternatives

**systray Library Behavior:**
- Assumes systray correctly displays Unicode in menu item text
- Prior phases successfully used emoji (🍅○○○), so this assumption is validated
- No known issues with current `getlantern/systray` version

**No Monitoring/Alerting Needed:**
- Single-user desktop application
- Visual issues are immediately apparent to user
- No distributed system concerns
- No SLOs or uptime requirements

### What "Good" Looks Like

After rollout, the following should be true:
- Timer updates are visible approximately every 10 seconds
- CPU usage during active sessions is lower than before
- Progress bar fills smoothly (no sudden jumps, though some granularity acceptable)
- All emoji render correctly in dropdown menu
- Tray icons change in sync with session transitions
- User can glance at dropdown and immediately understand state without reading text

---

## Risks, Trade-offs, and Alternatives

### Known Risks

**Unicode Rendering Inconsistency:**
- **Risk**: Different macOS versions or system font configurations might render ● and ○ inconsistently (different sizes, spacing issues)
- **Likelihood**: Low (standard Unicode, broad support)
- **Mitigation**: Use well-supported Unicode characters (U+25CF, U+25CB). If issues arise, can substitute with alternative characters or ASCII fallback
- **Fallback**: Switch to `[XXXXX     ]` style ASCII progress bar

**10-Second Updates Feel Too Slow:**
- **Risk**: Users might perceive 10-second intervals as unresponsive, especially if accustomed to second-by-second countdown
- **Likelihood**: Medium (subjective user preference)
- **Mitigation**: Start with 10s based on increment requirements. Can adjust to 5s or 7s if user feedback demands it
- **Fallback**: One-line constant change to reduce interval

**Progress Bar Segment Alignment:**
- **Risk**: With 10 segments and 25-minute sessions (1500 seconds), each segment represents 150 seconds (2.5 minutes). Progress might appear to "jump" rather than fill smoothly
- **Likelihood**: High (inherent to 10-segment design)
- **Impact**: Low (acceptable for lite UI, users unlikely to notice exact segment boundaries)
- **Mitigation**: Accept as design trade-off. Can increase segments to 20 if feedback indicates issue, but adds visual clutter

**Emoji Font Availability:**
- **Risk**: Older macOS versions (<10.12) might not render emoji correctly
- **Likelihood**: Very Low (target audience likely on recent macOS)
- **Mitigation**: Minimum macOS version can be documented. Can provide text fallbacks if needed
- **Fallback**: Use ASCII symbols (*, @, #) instead of emoji

**Icon File Requirements:**
- **Risk**: Implementation assumes icon files exist and are correctly formatted
- **Likelihood**: None (assets already confirmed to exist in `/assets/` directory)
- **Mitigation**: N/A (already validated)

### Trade-offs

**10 Seconds vs 1 Second Updates:**
- **Chosen**: 10 seconds
- **Trade-off**: Less responsive feel for significant CPU/battery savings
- **Rationale**: For 25-minute sessions, 10-second granularity provides sufficient feedback. CPU efficiency matters for background app. Matches increment requirements.
- **When to revisit**: If multiple users report feeling disconnected from timer progress

**Minutes-Only vs MM:SS Display:**
- **Chosen**: Minutes only ("25min", "12min", "0min")
- **Trade-off**: Less precision for improved scannability
- **Rationale**: Users don't need second-level precision for pomodoro technique. Simpler display is easier to parse at a glance. Matches increment requirements.
- **When to revisit**: If users request more precise countdown (can make configurable in future)

**10 Segments vs More Granular Progress Bar:**
- **Chosen**: 10 segments (●●●○○○○○○○)
- **Trade-off**: Some visual "jumping" between segments vs cleaner appearance
- **Rationale**: 10 segments fit well in menu width. Simple to calculate (multiples of 10%). More segments add visual clutter.
- **When to revisit**: If users report confusion about progress jumps

**Unicode Circles vs Other Progress Indicators:**
- **Chosen**: ● (filled circle) and ○ (empty circle)
- **Trade-off**: Potential rendering issues vs visual appeal
- **Rationale**: Circles are semantically clear (filled = done, empty = remaining). Widely supported Unicode. Matches common design patterns.
- **When to revisit**: If rendering issues appear on target macOS versions

### Alternatives Considered

**5-Second Tick Interval:**
- **Alternative**: Update every 5 seconds instead of 10
- **Rejected**: Unnecessary compromise. 10 seconds provides sufficient feedback for 25-minute sessions while maximizing CPU savings.
- **When to reconsider**: If user testing reveals 10s feels too slow

**Text-Based Percentage Display:**
- **Alternative**: Show "50% complete" or "12/25 minutes elapsed"
- **Rejected**: Less visually engaging than graphical progress bar. Doesn't provide at-a-glance understanding.
- **When to reconsider**: If Unicode rendering proves problematic

**Different Progress Bar Characters:**
- **Alternatives considered**:
  - ▰▱ (larger rectangles)
  - █░ (blocks)
  - ⬤⚪ (larger circles)
  - ASCII: `[XXXXX     ]`
- **Choice rationale**: ●○ provides good balance of visibility and compactness
- **Open to adjustment**: If rendering issues arise, can swap characters with one-line change

**Variable Segment Count:**
- **Alternative**: Use 15, 20, or 25 segments for smoother progress
- **Rejected**: Adds visual clutter. Calculation complexity not worth marginal smoothness gain.
- **When to reconsider**: If users specifically request finer granularity

**Configurable Update Frequency:**
- **Alternative**: Add setting to choose 5s/10s/15s update intervals
- **Rejected**: Over-engineering for lite project. Can add later if clear user demand.
- **When to reconsider**: If different users have strong opposing preferences

---

## Follow-up Work

### Future Increments (Explicitly Deferred)

**Dark Mode Support:**
- Detect macOS system theme (light/dark)
- Maintain separate icon sets for each theme
- Possibly adjust emoji choices for better contrast
- Estimated scope: Small increment, similar complexity to this one

**Settings-Based Customization:**
- Allow users to configure visual preferences via `settings.json`
- Options: update frequency (5s/10s/30s), timer format (min vs MM:SS), emoji on/off
- Requires settings loading in UI initialization
- Estimated scope: Medium increment (settings infrastructure)

**Alternative Progress Bar Styles:**
- Support multiple progress indicator types in settings
- Options: circles, blocks, ASCII, percentage text
- Estimated scope: Small increment once settings infrastructure exists

**Extended Progress Information:**
- Show estimated completion time ("Work ends at 3:45 PM")
- Show elapsed time in addition to remaining
- Daily progress summary in dropdown footer
- Estimated scope: Small increment, presentation logic only

**Icon Design Improvements:**
- Higher-resolution icon assets for Retina displays
- Subtle color variations (within menu bar conventions)
- Potentially animated icon sequences (if macOS supports)
- Estimated scope: Design work + small integration increment

**Accessibility Features:**
- Text-based alternatives for screen reader compatibility
- High-contrast icon variants
- Audio feedback for state transitions (subtle sounds)
- Estimated scope: Medium increment, requires accessibility testing

### Tech Debt or Clean-up

**None anticipated** for this increment:
- Changes are straightforward additions/modifications
- No temporary workarounds introduced
- No complex refactoring needed
- Code remains simple and maintainable

**Potential future refactoring** (not required now):
- If progress bar logic grows complex: extract to separate function/file
- If format functions proliferate: create formatting utility package
- If emoji mappings grow: extract to constants/configuration

### Validation Tasks After Rollout

**Immediate (Day 1):**
- Verify progress bar renders correctly during first work session
- Confirm 10-second updates feel responsive enough
- Check that all emoji appear correctly

**Short-term (Week 1):**
- Monitor for any macOS version compatibility reports
- Verify CPU usage improvement is noticeable
- Gather subjective feedback on visual improvements

**Long-term (Month 1):**
- Assess whether 10-second interval remains acceptable
- Determine if any emoji or Unicode characters need adjustment
- Evaluate whether any follow-up visual enhancements are desired

---

## References

- **Increment Definition**: [increment.md](./increment.md)
- **Project Constitution**: [CONSTITUTION.md](../../../CONSTITUTION.md)
- **Architecture Documentation**: [ARCHITECTURE.md](../../../ARCHITECTURE.md)
- **Product Requirements**: [PRD.md](../../../PRD.md) (Phase 3: Visual Enhancements)
- **Related Code**:
  - UI Package: [internal/ui/window.go](../../../internal/ui/window.go)
  - Tray Package: [internal/tray/tray.go](../../../internal/tray/tray.go)
  - Timer Package: [internal/timer/timer.go](../../../internal/timer/timer.go)
  - Session Package: [internal/session/session.go](../../../internal/session/session.go)
