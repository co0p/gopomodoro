GoPomodoro
==========

A minimal Pomodoro timer that lives in your taskbar, helping you maintain focus through structured work intervals.

## What is the Pomodoro Technique?

The Pomodoro Technique is a time management method developed by Francesco Cirillo in the late 1980s. It uses a timer to break work into focused intervals, traditionally 25 minutes in length, separated by short breaks.

## The Core Rules

### The Basic Cycle

1. **Work (Pomodoro)**: 25 minutes of focused, uninterrupted work
2. **Short Break**: 5 minutes to rest and recharge
3. **Repeat**: Complete 4 pomodoros
4. **Long Break**: 15 minutes after completing 4 pomodoros

### The Essential Principles

- **Indivisible**: A pomodoro cannot be split. If interrupted, it must be abandoned and restarted.
- **No Multitasking**: One task per pomodoro. Full focus, no distractions.
- **No Exceptions**: If a pomodoro begins, it must complete. No checking off early.
- **Protect Your Time**: If an interruption is unavoidable, the pomodoro is void.

## Interface

### Taskbar (Collapsed)
- **Icon**: Shows current state (🍅 Pomodoro / ☕ Short Break / 🌴 Long Break)
- **Timer**: Minutes remaining displayed next to icon

Example: `🍅 23` or `☕ 4`

### Tray (Opened)
- **Current State**: Visual indicator of phase
- **Timer**: Minutes remaining (if applicable)
- **Actions**: Three buttons
  - **Start**: Begin a pomodoro
  - **Pause**: Pause the current timer
  - **Reset**: Abandon current pomodoro and restart

### Notifications
- **Sound alerts**: A brief tone plays automatically when phases transition
  - Pomodoro → Break (short or long)
  - Break → Pomodoro
  - Long break → Idle
- **Purpose**: Stay focused without watching the timer constantly

## Timer Intervals

| State | Duration | Icon |
|-------|----------|------|
| Pomodoro | 25 minutes | 🍅 |
| Short Break | 5 minutes | ☕ |
| Long Break | 15 minutes | 🌴 |

After 4 completed pomodoros, the timer automatically moves to a long break.

## Controls

### Start
- Begins a new 25-minute pomodoro
- Resumes a paused pomodoro
- Starts the next break after a pomodoro completes

### Pause
- Temporarily stops the current timer
- Time freezes until resumed
- Use sparingly—pausing defeats the purpose of timeboxing

### Reset
- Abandons the current pomodoro or break
- Returns to idle state, ready to start fresh
- Use when interruptions make the current pomodoro invalid

## Flags

### --silent
- Disables all sound notifications
- Visual timer updates continue normally
- Usage: `gopomodoro --silent`

## The Philosophy

> "The Pomodoro Technique isn't about the time you have, it's about the focus you bring."

### Why It Works

- **Timeboxing**: Finite intervals create urgency and prevent procrastination
- **Breaks Are Mandatory**: Regular rest prevents burnout and maintains mental freshness
- **Rhythm**: Predictable cycles create a sustainable work rhythm
- **Simplicity**: No configuration, no complexity—just focus

### Best Practices

✅ **Do:**
- Plan your tasks before starting
- Eliminate distractions during pomodoros
- Actually take your breaks
- Reset if genuinely interrupted
- Trust the fixed intervals

❌ **Don't:**
- Check email during a pomodoro
- Skip breaks to "power through"
- Split a pomodoro across multiple tasks
- Pause for trivial interruptions
- Obsess over the timer—it's just a guide

## Design Principles

This timer embraces minimalism:

- **Subtle notifications**: Optional sound alerts at phase transitions (disable with --silent)
- **No customization**: The traditional intervals work. Trust the method.
- **No statistics**: Focus on the present work, not the metrics
- **No complexity**: Three actions. One purpose. Pure focus.

## Privacy

This timer:
- ✅ Runs locally on your machine
- ❌ Does not track what you're working on
- ❌ Does not collect or send any data

## Credits

Based on the Pomodoro Technique® by Francesco Cirillo. The Pomodoro Technique® and Pomodoro™ are registered trademarks of Francesco Cirillo.

## License

[Your License Here]

---

**Remember**: The timer is just a tool. Your focus is the real power. 🍅