# Design (Emergent)

> This document reflects architecture that emerged through TDD.
> Updated during Promote phase. Not a planning document.

## Current Structure

The application follows a ports-and-adapters pattern with the domain (`pkg/`) at the center, adapters in subpackages (`pkg/tray/`, `pkg/ticker/`), and wiring in `cmd/`.

### pkg/
- **Purpose**: Core domain logic for pomodoro cycle state machine
- **Key patterns**: State machine for timer transitions, observer pattern for state changes

### pkg/tray/
- **Purpose**: System tray UI adapter
- **Key patterns**: Formatter pattern separates display logic from rendering

## Patterns Discovered

### Formatter Pattern for Testable Display
- **What**: Separate display string formatting from UI rendering. `Formatter.Format(CycleState, time.Duration) → string`
- **Why it emerged**: Testing tray display required avoiding systray dependency. Tests needed to verify display strings without launching UI.
- **Where used**: `pkg/tray/formatter.go`, `pkg/tray/tray.go`

## History

| Date | Increment | Changes |
|------|-----------|---------|
| 2026-01-29 | Short Break and Time Display | Formatter pattern extracted for testable display logic |
