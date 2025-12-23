# Project Constitution

constitution-mode: lite

## 1. Purpose and Scope

This constitution guides the development of **GoPomodoro**, a minimal macOS native pomodoro timer application written in Go.

GoPomodoro is a single-developer tool focused on simplicity and effectiveness. This constitution defines how we approach building and evolving this software—keeping increments small, designs lightweight, and code clean.

The constitution applies to:
- The codebase under this repository
- Documentation and increment artifacts we create
- Design and implementation decisions we make

## 2. Implementation & Doc Layout

### Increment Artifacts
- **Location:** `docs/increments/<slug>/`
- **Files:**
  - `increment.md` — What we're building and why
  - `design.md` — How we'll build it (lightweight, skip if obvious)
  - `implement.md` — Step-by-step task breakdown

### Improve Artifacts (Lightweight)
- **Location:** `docs/improve/`
- **Filename pattern:** `YYYY-MM-DD-improve.md`
- **Usage:** Sparingly—only when stepping back to reflect adds clear value

### ADR Artifacts (Minimal)
- **Location:** `docs/adr/`
- **Filename pattern:** `ADR-YYYY-MM-DD-<slug>.md`
- **Usage:** Only for significant architectural decisions
  - Examples: GUI framework choice, state persistence approach, notification system design

### Other Documentation
- **Architecture notes:** `docs/architecture/` (if complexity grows beyond the PRD)
- **Runbooks/ops:** `docs/ops/` (if deployment/distribution needs documenting)

## 3. Design & Delivery Principles

### Small, Safe Steps
Build in tiny increments that can be tested and committed independently. Each increment should take hours or days, not weeks.

**In practice:**
- Add the tray icon before building the dropdown menu
- Get the timer counting before adding pause/resume
- Hardcode one theme before implementing theme switching

### Simple Is Better Than Complex
Prefer straightforward solutions over clever abstractions. Don't build for imaginary future requirements.

**In practice:**
- Start with JSON file persistence, not a database
- Use a simple state machine for timer states, not an event bus
- Keep the UI minimal—resist feature creep

### Make It Work, Make It Right
Get something working first, then refactor to clean it up. Refactoring is normal work, not a separate phase.

**In practice:**
- Hardcode timer durations (25/5/15 minutes) initially, then make them configurable when needed
- Build the UI with basic colors first, then refactor for light/dark mode support
- If code gets messy while figuring out macOS APIs, clean it up once it works

## 4. Testing and CI

### Testing (Lite Expectations)
- Manual testing is fine for UI and integration
- Add automated tests for core timer logic and state machine when they stabilize
- Don't over-test trivial code—focus tests where bugs would hurt

### CI/CD
- No CI requirement for this single-developer project
- If you add CI later (GitHub Actions), keep it simple: build and run tests

## 5. When to Use ADRs and Improve

### ADRs
Use ADRs only for **significant architectural choices** that have lasting impact:
- Which Go GUI library to use (fyne, wails, native bindings, etc.)
- How to handle macOS tray integration
- Persistence format and location
- Notification system design

Skip ADRs for routine decisions that can be changed easily.

### Improve
Use Improve docs when you need to **step back and reflect** on what's working and what isn't:
- After completing a major milestone (e.g., "first working version")
- When you notice code or design smells accumulating
- Before planning a significant refactor

For a lite project, you might only write 2-3 Improve docs across the entire development lifecycle.

---

## Acceptance

This constitution succeeds when:
- Increments feel natural and quick to create
- Designs are lightweight but clear enough to guide implementation
- Code stays simple and maintainable
- You're building steadily without getting bogged down in process
