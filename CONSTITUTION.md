# Project Constitution

## About This Project

GoPomodoro is a system tray tool for Pomodoro time management built with Go and Fyne.io. It provides a lightweight, desktop-native way to manage focus sessions with minimal friction.

---

## Core Principles

### 1. Desktop-Native Experience _(Pillar: Design Integrity)_

**Statement:** System tray integration is our primary interface—the application lives where users work, not in another window competing for attention.

**Rationale:** Pomodoro tools should reduce cognitive overhead, not add to it. A native system tray presence ensures users can control timers without breaking flow or switching contexts.

**In Practice:**
- All core actions (start/pause/stop, session selection) accessible from tray menu
- Visual feedback through tray icon state changes (running, paused, break)
- No persistent windows—the app exists in the background until needed
- Platform conventions respected (macOS menu bar, Windows system tray, Linux status bar)

### 2. Comprehensive Test Coverage _(Pillar: Test Strategy)_

**Statement:** We test timer logic, state transitions, and UI behavior comprehensively—confidence in correctness matters more than shipping speed.

**Rationale:** Timer accuracy and state consistency are non-negotiable. Users rely on Pomodoro intervals to structure their work; a broken timer breaks trust.

**In Practice:**
- Unit tests for all timer state machine transitions (running → paused → stopped)
- Test time calculations and remaining time display logic
- Test session type switching (work → short break → long break cycles)
- Integration tests for tray interactions where feasible with Fyne testing tools
- Fast iteration on experiments, careful testing before merging to main

### 3. Minimal Dependencies _(Pillar: Dependency Discipline)_

**Statement:** We prefer Go's standard library over third-party packages—dependencies are added only when they provide essential functionality we cannot reasonably build ourselves.

**Rationale:** Every dependency increases binary size, maintenance burden, and attack surface. For a desktop tool that runs continuously, staying lean matters.

**In Practice:**
- Fyne.io is our only UI framework dependency (essential for cross-platform GUI)
- Timer logic uses `time` package from standard library
- No logging frameworks—standard `log` package is sufficient
- No configuration libraries—simple file I/O or JSON unmarshaling
- Before adding a package, ask: "Can we implement this in <100 lines?"

### 4. Single Responsibility Components _(Pillar: Simplicity First)_

**Statement:** Each Go package does one thing—timer logic, UI rendering, and state management remain clearly separated.

**Rationale:** Clear boundaries make testing easier and prevent entangled code. When timer logic and UI are coupled, both become harder to change.

**In Practice:**
- `timer` package handles countdown logic and state transitions (no UI imports)
- `ui` package manages Fyne widgets and tray menu (no timer implementation)
- `state` package coordinates between timer and UI (thin orchestration layer)
- Duplication is acceptable until patterns emerge naturally (pragmatic refactoring)

### 5. Resource Efficiency _(Pillar: Technical Debt Boundaries)_

**Statement:** The application must remain lightweight—memory footprint under 50MB, CPU usage negligible when idle, battery impact minimal.

**Rationale:** Users run GoPomodoro continuously throughout their workday. A resource-hungry app would be uninstalled quickly.

**In Practice:**
- Use efficient time.Ticker for countdown updates (not polling loops)
- Profile memory allocations during development
- Test on low-power devices (older laptops, ARM-based machines)
- Quick hacks allowed in experimental features (labeled with TODO comments)
- Performance regressions block merges—measure before and after

---

### Pillar Coverage

_This constitution addresses the following pillars:_
- ✓ **Design Integrity** (Principle #1 - Desktop-Native Experience)
- ✓ **Test Strategy** (Principle #2 - Comprehensive Test Coverage)
- ✓ **Dependency Discipline** (Principle #3 - Minimal Dependencies)
- ✓ **Simplicity First** (Principle #4 - Single Responsibility Components)
- ✓ **Technical Debt Boundaries** (Principle #5 - Resource Efficiency)

## Technical Decisions

### Languages
- **Go:** Primary language for its simplicity, excellent standard library, and efficient concurrency primitives ideal for timer management.

### Frameworks
- **Fyne.io:** Cross-platform GUI framework chosen for native system tray support and declarative UI approach that minimizes boilerplate.

### Deployment
- **Static Binaries:** Single-file executables compiled for macOS, Linux, and Windows—no runtime dependencies or installation wizards required.

---

**Last Updated:** 2025-11-27
