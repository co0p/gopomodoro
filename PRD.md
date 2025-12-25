# Product Requirements Document:  GoPomodoro

## Overview

**Product Name:** GoPomodoro  
**Version:** 1.0  
**Last Updated:** 2025-12-23  
**Document Owner:** co0p  

### Executive Summary

GoPomodoro is a minimal, native traybar pomodoro timer application focused on helping users maintain focus through the pomodoro technique. The app prioritizes simplicity and unobtrusiveness while providing just enough engagement through color coding, notifications, and simple streak tracking.

### Product Vision

A pomodoro timer that lives in your system tray, stays out of your way, and provides subtle encouragement to maintain focus. No bloat, no complexity—just an effective focus tool. 

---

## Core Principles

1. **Minimal by Default** - Single-click access, no unnecessary features
2. **Visual Clarity** - Status communicated through color and icons
3. **Non-Intrusive** - Lives in tray, gentle notifications
4. **Simple Persistence** - Human-readable files for settings and tracking
5. **No Settings UI** - Manual JSON editing for advanced customization

---

## User Personas

### Primary:  Focused Developer
- **Name:** Alex
- **Context:** Works from home, easily distracted
- **Needs:** Simple focus tool, doesn't want complexity
- **Behavior:** Wants to start/stop quickly, occasional stats check

### Secondary: Remote Worker
- **Name:** Jordan  
- **Context:** Back-to-back meetings, needs structured breaks
- **Needs:** Reminders to take breaks, visual time tracking
- **Behavior:** Relies on notifications, checks streaks for motivation

---

## Product Specifications

### 1. User Interface

#### Tray Icon States
The tray icon is the primary UI element, always visible: 

| State | Icon | Color | Display |
|-------|------|-------|---------|
| Idle | ○ | Gray | Ready |
| Work Session | 🍅 | Red | 25 |
| Short Break | ☕ | Green | 5 |
| Long Break | 🌟 | Blue | 15 |
| Paused | ⏸️ | Gray | 12: 34 |

#### Dropdown Panel (Left Click Only)
```
┌─────────────────────────────┐
│  🍅 Work Session            │  ← Color-coded header
│     24: 37                   │  ← Large timer
│  ━━━━━━━━━━━━━━━━━━━━━━━  │  ← Progress bar
├─────────────────────────────┤
│  Session 2/4    🍅🍅○○      │  ← Cycle indicator
├─────────────────────────────┤
│  [  ▶️  Start  ]            │  ← Context-aware buttons
│  [  ⏸️  Pause  ]            │
│  [  ⏭️  Skip   ]            │
│  [  🔄  Reset  ]            │
├─────────────────────────────┤
│  Today:  6🍅  Streak: 🔥3    │  ← Stats footer
└─────────────────────────────┘
```

#### Color Scheme

**Light Mode:**
- **Work:** Background `#FFE5E5`, Text `#C0392B`, Accent `#E74C3C`
- **Short Break:** Background `#E5F9E5`, Text `#27AE60`, Accent `#2ECC71`
- **Long Break:** Background `#E5F2FF`, Text `#2980B9`, Accent `#3498DB`
- **Paused:** Background `#F0F0F0`, Text `#7F8C8D`, Accent `#95A5A6`

**Dark Mode:**
- **Work:** Background `#3D1F1F`, Text `#FF6B6B`, Accent `#FF6B6B`
- **Short Break:** Background `#1F3D1F`, Text `#51CF66`, Accent `#51CF66`
- **Long Break:** Background `#1F2D3D`, Text `#4DABF7`, Accent `#4DABF7`
- **Paused:** Background `#2C2C2C`, Text `#ADB5BD`, Accent `#868E96`

#### Button States

**Idle State:**
- Start Work
- (Reset Stats - optional)

**Active Work Session:**
- Pause
- Skip
- Reset

**Active Break:**
- Start Work Now
- Pause Break
- Reset Cycle

**Paused State:**
- Resume
- Skip
- Reset

### 2. Pomodoro Logic

#### Session Flow
```
Start → Work (25m) → Short Break (5m) → Work (25m) → Short Break (5m) 
     → Work (25m) → Short Break (5m) → Work (25m) → Long Break (15m) 
     → [Cycle Repeats]
```

#### Default Timing
- **Work Session:** 25 minutes
- **Short Break:** 5 minutes
- **Long Break:** 15 minutes
- **Cycles Before Long Break:** 4

#### Session Rules
1. Work sessions must be completed to count toward stats
2. Skipped sessions are logged but don't count as completed
3. Long break unlocks after 4 completed work sessions
4. Pausing stops the timer but maintains session state
5. Reset clears current cycle back to session 1

### 3. Notifications

#### Session Complete - Work
```
Title: 🍅 Pomodoro Complete!
Body: Great focus!  Time for a 5 minute break.
Actions: [Take Break] [Keep Working]
Sound: Gentle chime (ding. wav)
```

#### Session Complete - Short Break
```
Title: ☕ Break Over!
Body: Ready to dive back in? 
Actions: [Start Work] [Extend Break]
Sound: Soft bell
```

#### Session Complete - Long Break
```
Title: 🌟 Long Break Over!
Body: You completed 4 pomodoros!  Feeling refreshed? 
Actions: [Start Fresh Cycle]
Sound: Success chime
```

#### Milestone - Streak
```
Title: 🔥 3 Day Streak!
Body: You're building a great focus habit!
Actions: [Dismiss]
Sound: None (silent)
Trigger: Every 3, 7, 14, 30, 60, 90 day milestones
```

#### Notification Settings
- **Enabled/Disabled:** Global toggle
- **Sound:** Separate toggle
- **Types:** Session complete, break complete, milestones

### 4. Data Persistence

#### Directory Structure
```
$HOME/. gopomodoro/
├── settings.json       # User configuration
├── sessions.log        # Append-only session log
└── stats.json          # Computed statistics
```

#### settings.json
```json
{
  "workDuration": 25,
  "shortBreakDuration": 5,
  "longBreakDuration": 15,
  "cyclesBeforeLongBreak": 4,
  "notifications": {
    "enabled": true,
    "sound":  true,
    "sessionComplete": true,
    "breakComplete": true,
    "milestones": true
  },
  "autoStartBreaks": true,
  "autoStartWork": false,
  "theme": "auto"
}
```

**Field Descriptions:**
- `workDuration`: Minutes for work session (default:  25)
- `shortBreakDuration`: Minutes for short break (default: 5)
- `longBreakDuration`: Minutes for long break (default: 15)
- `cyclesBeforeLongBreak`: Work sessions before long break (default: 4)
- `notifications.enabled`: Master toggle for all notifications
- `notifications.sound`: Enable/disable notification sounds
- `notifications.sessionComplete`: Show work session complete notifications
- `notifications.breakComplete`: Show break complete notifications
- `notifications.milestones`: Show streak milestone notifications
- `autoStartBreaks`: Automatically start break timer when work completes
- `autoStartWork`: Automatically start work timer when break completes
- `theme`: UI theme (`"light"`, `"dark"`, `"auto"`)

#### sessions.log
Append-only CSV format:
```
timestamp,session_type,event,duration_minutes
2025-12-23T09:00:00Z,work,started,0
2025-12-23T09:25:00Z,work,completed,25
2025-12-23T09:25:01Z,short_break,started,0
2025-12-23T09:30:01Z,short_break,completed,5
2025-12-23T09:30:02Z,work,started,0
2025-12-23T09:35:02Z,work,skipped,5
2025-12-23T09:35:03Z,short_break,started,0
```

**Event Types:**
- `started`: Session began
- `completed`: Session finished naturally
- `skipped`: User skipped before completion
- `paused`: Timer paused
- `resumed`: Timer resumed from pause

#### stats.json
```json
{
  "totalPomodoros": 247,
  "totalFocusMinutes": 6175,
  "currentStreak": 3,
  "longestStreak": 12,
  "lastSessionDate": "2025-12-23",
  "dailyStats": {
    "2025-12-23": {
      "completed": 6,
      "skipped": 1,
      "focusMinutes": 150,
      "breakMinutes": 25
    },
    "2025-12-22": {
      "completed": 8,
      "skipped": 0,
      "focusMinutes": 200,
      "breakMinutes":  35
    }
  }
}
```

**Calculation Rules:**
- `totalPomodoros`: Count of all completed work sessions
- `currentStreak`: Consecutive days with ≥1 completed pomodoro
- `longestStreak`: Historical longest streak
- `dailyStats`: Rolling 30-day window
- Stats regenerated from `sessions.log` on app start

### 5. Progress Visualization

#### Progress Bar
Visual indicator showing elapsed time in current session:
```
Empty:    ○○○○○○○○○○○○○○○○○○○○  (0%)
25%:      ●●●●●○○○○○○○○○○○○○○○  
50%:      ●●●●●●●●●●○○○○○○○○○○  
75%:      ●●●●●●●●●●●●●●●○○○○○  
Complete: ●●●●●●●●●●●●●●●●●●●●  (100%)
```
- Updates every 30 seconds (not real-time)
- Color matches current session type
- Smooth animation on updates

#### Cycle Indicator
Shows progress through 4-session cycle:
```
Session 1/4:  🍅○○○
Session 2/4:  🍅🍅○○
Session 3/4:  🍅🍅🍅○
Session 4/4:  🍅🍅🍅🍅
```

### 6. System Integration

#### Platform Support
- **Primary:** macOS (native tray support)
- **Secondary:** Linux (system tray)
- **Future:** Windows

#### System Tray Behavior
- App runs in background only (no dock/taskbar icon)
- Tray icon always visible
- Single left-click opens dropdown
- Dropdown closes when clicking outside
- App survives system sleep/wake

#### Startup Behavior
- Does not auto-start on system boot (user configurable via OS)
- Restores state if app crashes (reads from sessions.log)
- Creates `~/.gopomodoro/` if missing

---

## Technical Requirements

### Technology Stack
- **Language:** Go
- **UI Framework:** 
  - macOS: Native Cocoa (via cgo or go-astilectron)
  - Linux: systray library
- **Storage:** JSON files, CSV logs
- **Notifications:** Native OS notifications

### Performance Requirements
- **Memory Usage:** < 50 MB idle
- **CPU Usage:** < 1% idle, < 5% active
- **Startup Time:** < 1 second
- **File I/O:** Async writes, no blocking

### Data Requirements
- **Settings:** Hot-reload on file change (optional)
- **Stats:** Regenerate from log on startup
- **Log Retention:** Keep 90 days, auto-prune older

---

## Feature Development Roadmap

### Phase 1: MVP Core (v0.1)
**Goal:** Basic functional pomodoro timer

- [x] **F1.1** Tray icon with basic states (idle, work, break) — See [Tray Icon and Dropdown UI Increment](docs/increments/tray-icon-and-dropdown-ui/increment.md)
- [x] **F1.2** Left-click dropdown UI — See [Tray Icon and Dropdown UI Increment](docs/increments/tray-icon-and-dropdown-ui/increment.md)
- [x] **F1.3** Timer logic (25/5/15 hardcoded)
- [x] **F1.4** Start/Pause/Reset buttons
- [x] **F1.5** Basic timer display (MM:SS)
- [x] **F1.6** Session complete detection
- [x] **F1.7** Create `~/.gopomodoro/` directory on first run

**Deliverable:** Can start/pause/reset a single 25-minute timer

---

### Phase 2: Pomodoro Cycle (v0.2)
**Goal:** Complete pomodoro workflow

- [ ] **F2.1** Implement 4-session cycle logic
- [ ] **F2.2** Short break (5 min) after work session
- [ ] **F2.3** Long break (15 min) after 4 cycles
- [ ] **F2.4** Cycle indicator UI (🍅🍅○○)
- [ ] **F2.5** Skip button functionality
- [ ] **F2.6** Auto-transition between sessions (optional)
- [ ] **F2.7** Tray icon updates per session type

**Deliverable:** Full pomodoro cycle works automatically

---

### Phase 3: Visual Feedback (v0.3)
**Goal:** Color coding and progress visualization

- [ ] **F3.1** Color-coded header backgrounds
- [ ] **F3.2** Progress bar implementation
- [ ] **F3.3** Tray icon color changes
- [ ] **F3.4** Dark mode support
- [ ] **F3.5** Smooth transitions/animations
- [ ] **F3.6** Timer updates every 30 seconds (not real-time)

**Deliverable:** Visually polished UI with clear state indication

---

### Phase 4: Notifications (v0.4)
**Goal:** User alerts and engagement

- [ ] **F4.1** Native OS notification integration
- [ ] **F4.2** Work session complete notification
- [ ] **F4.3** Break complete notification
- [ ] **F4.4** Long break complete notification
- [ ] **F4.5** Notification sound support
- [ ] **F4.6** Notification action buttons (Take Break, Start Work)

**Deliverable:** Users are notified at key moments

---

### Phase 5: Data Persistence (v0.5)
**Goal:** Settings and session logging

- [ ] **F5.1** Create `settings.json` with defaults
- [ ] **F5.2** Read settings on app start
- [ ] **F5.3** Apply custom durations from settings
- [x] **F5.4** Create `sessions.log` file
- [x] **F5.5** Log session events (start/complete/skip)
- [x] **F5.6** Append-only logging with proper timestamps
- [ ] **F5.7** Handle file I/O errors gracefully

**Deliverable:** Settings persist, sessions are logged

---

### Phase 6: Statistics & Tracking (v0.6)
**Goal:** Simple progress tracking

- [ ] **F6.1** Generate `stats.json` from sessions.log
- [ ] **F6.2** Calculate today's completed pomodoros
- [ ] **F6.3** Calculate current streak
- [ ] **F6.4** Display stats in dropdown footer
- [ ] **F6.5** Daily stats aggregation
- [ ] **F6.6** Total focus time calculation
- [ ] **F6.7** Auto-prune logs older than 90 days

**Deliverable:** Users can see basic stats and streaks

---

### Phase 7: Polish & Edge Cases (v0.7)
**Goal:** Production-ready stability

- [ ] **F7.1** Handle system sleep/wake
- [ ] **F7.2** Restore state after crash
- [ ] **F7.3** Proper cleanup on app quit
- [ ] **F7.4** Settings hot-reload (optional)
- [ ] **F7.5** Error logging to `~/.gopomodoro/error.log`
- [ ] **F7.6** Validate settings. json on load
- [ ] **F7.7** Graceful degradation if files corrupted

**Deliverable:** Stable, reliable app for daily use

---

### Phase 8: Enhanced Engagement (v0.8)
**Goal:** Subtle motivation features

- [ ] **F8.1** Streak milestone notifications (3, 7, 14, 30 days)
- [ ] **F8.2** Extended stats view (hover or expand)
- [ ] **F8.3** Weekly summary
- [ ] **F8.4** Longest streak tracking
- [ ] **F8.5** Best day tracking
- [ ] **F8.6** Custom notification messages

**Deliverable:** Users feel motivated by progress

---

### Phase 9: Platform Expansion (v0.9)
**Goal:** Cross-platform support

- [ ] **F9.1** Linux system tray integration
- [ ] **F9.2** Linux notification support
- [ ] **F9.3** Windows system tray (future)
- [ ] **F9.4** Platform-specific build scripts
- [ ] **F9.5** CI/CD for multi-platform builds

**Deliverable:** Works on macOS and Linux

---

### Phase 10: Advanced Features (v1.0)
**Goal:** Power user enhancements

- [ ] **F10.1** Global keyboard shortcuts (start/pause/skip)
- [ ] **F10.2** Idle detection and auto-pause
- [ ] **F10.3** Do Not Disturb integration
- [ ] **F10.4** Export stats to CSV
- [ ] **F10.5** Manual session entry (backfill)
- [ ] **F10.6** Custom themes via settings
- [ ] **F10.7** Sound customization (bring your own chime)

**Deliverable:** Feature-complete v1.0 release

---

## Success Metrics

### User Engagement
- **Daily Active Users:** Users opening app daily
- **Average Sessions/Day:** Target 6-8 pomodoros
- **Streak Retention:** Users maintaining 3+ day streaks

### Performance
- **Crash Rate:** < 0.1% of sessions
- **Notification Delivery:** > 99% success rate
- **File I/O Errors:** < 0.01% of operations

### Adoption
- **User Retention (7-day):** > 60%
- **User Retention (30-day):** > 40%

---

## Non-Goals

### Out of Scope (v1.0)
- ❌ Cloud sync across devices
- ❌ Mobile apps
- ❌ Team/collaborative features
- ❌ Browser extension
- ❌ Task/project management integration
- ❌ In-app settings UI (edit JSON manually)
- ❌ Charts and graphs
- ❌ Pomodoro customization beyond timing
- ❌ Social features or leaderboards
- ❌ Paid features or monetization

---

## Appendix

### A. File Format Examples

#### Minimal settings.json
```json
{
  "workDuration": 25,
  "shortBreakDuration": 5,
  "longBreakDuration": 15,
  "cyclesBeforeLongBreak": 4
}
```

#### Custom settings.json
```json
{
  "workDuration": 50,
  "shortBreakDuration": 10,
  "longBreakDuration": 30,
  "cyclesBeforeLongBreak": 3,
  "notifications": {
    "enabled": true,
    "sound": false,
    "sessionComplete": true,
    "breakComplete": false,
    "milestones":  true
  },
  "autoStartBreaks":  false,
  "autoStartWork": false,
  "theme": "dark"
}
```

### B. Notification Sound Recommendations
- **Work Complete:** Single bell chime (0.5s)
- **Break Complete:** Soft gong (0.8s)
- **Long Break Complete:** Triple chime (1.2s)
- **Milestone:** None (silent notification)

### C. Edge Case Handling

| Scenario | Behavior |
|----------|----------|
| System sleep during session | Pause timer, resume on wake |
| App crash mid-session | Restore from last log entry |
| Missing settings.json | Create with defaults |
| Corrupted stats.json | Regenerate from sessions.log |
| Invalid settings values | Use defaults, log warning |
| Disk full | Graceful failure, notify user |
| Multiple instances | Prevent second instance from starting |

### D.  Accessibility Considerations
- High contrast mode support
- Screen reader compatibility for notifications
- Keyboard navigation in dropdown
- Configurable notification duration

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-12-23 | Initial PRD |

---

## Approval

**Product Owner:** co0p  
**Status:** Draft  
**Next Review:** After MVP (Phase 1)
