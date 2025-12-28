# ADR: Hexagonal Architecture for Testability

**Date:** 2025-12-28  
**Status:** Accepted  
**Context:** Hexagonal-for-testing increment

---

## Context

After implementing the complete pomodoro cycle with UI integration, the application faced critical architectural issues that prevented effective testing and caused runtime problems:

### Problems Identified

1. **Deadlocks in Production**
   - Timer callbacks re-entered timer methods while holding locks, causing UI freezes
   - Circular dependency: UI → Timer → UI callback created fragile execution paths
   - Lock contention between `timer.Start()` and subsequent UI state queries

2. **Untestable Business Logic**
   - Cannot test complete pomodoro cycles (4 work sessions + breaks) without starting the systray GUI
   - Business logic was intertwined with `systray` library implementation details
   - No way to verify cycle progression, session transitions, or timing behavior through automated tests

3. **Tight Coupling**
   - UI package (`internal/ui/window.go`) acted as both orchestrator AND view layer
   - Business rules scattered across timer callbacks and UI event handlers
   - No clear separation between "what should happen" (business logic) and "how to display it" (infrastructure)

4. **Circular Dependencies**
   - Window created Timer
   - Timer called back to Window methods
   - Window queried Timer state
   - Result: fragile code with deadlock potential

### Why Hexagonal Architecture

The decision to adopt hexagonal (ports and adapters) architecture was driven by the need to:

- **Enable automated testing** of complete pomodoro cycles without requiring GUI interaction
- **Eliminate deadlocks** by removing circular callback dependencies
- **Separate business logic** from infrastructure (systray, file system, time)
- **Support future drivers** (CLI, web interface, different test frameworks) without changing core logic
- **Align with project principles** from CONSTITUTION.md: simplicity, testability, clear boundaries

## Decision

Refactor the application to use **Hexagonal Architecture (Ports and Adapters pattern)** with a central `PomodoroService` that:

1. **Owns all business logic** - Session state machine, cycle progression, timing rules
2. **Exposes command methods** (inbound ports) - `StartSession()`, `PauseSession()`, `SkipSession()`, etc.
3. **Publishes events through interfaces** (outbound ports) - `Notifier`, `Storage`, `Clock`
4. **Never calls back to drivers** - Service emits events; adapters choose how to react
5. **Is fully testable in isolation** - Using mock adapters and deterministic clock

### Architecture Structure

```
Primary Adapters (Drivers)          Inbound Ports (Commands)
┌─────────────────────┐             ┌────────────────────────┐
│  SystrayAdapter     │────────────▶│  StartSession()        │
│  - Button handlers  │             │  PauseSession()        │
│  - Implements       │             │  ResumeSession()       │
│    Notifier         │             │  SkipSession()         │
└─────────────────────┘             │  ResetCycle()          │
                                    │  GetState()            │
┌─────────────────────┐             └────────────────────────┘
│  TestDriver         │                       │
│  - Programmatic     │                       ▼
│    commands         │             ┌────────────────────────┐
│  - Records events   │             │   PomodoroService      │
│  - Asserts state    │             │  - Business logic      │
└─────────────────────┘             │  - State machine       │
                                    │  - Cycle rules         │
                                    └────────────────────────┘
                                              │
                                              ▼
                                    Outbound Ports (Events)
                                    ┌────────────────────────┐
                                    │  Notifier interface:   │
                                    │    SessionStarted()    │
                                    │    SessionTick()       │
                                    │    SessionCompleted()  │
                                    │    StateChanged()      │
                                    │                        │
Secondary Adapters (Driven)        │  Clock interface       │
┌─────────────────────┐             │  Storage interface     │
│  SystrayNotifier    │◀────────────└────────────────────────┘
│  FileStorage        │
│  MockClock/RealClock│
└─────────────────────┘
```

The architecture follows the hexagonal pattern with:

- **Primary Adapters (Drivers):** Handle user interactions via UI and enable programmatic testing
- **Inbound Ports:** Command methods that control the pomodoro cycle
- **Core Domain:** Contains all business logic, state machine, and cycle rules
- **Outbound Ports:** Event interfaces that the service uses to communicate with infrastructure
- **Secondary Adapters (Driven):** Concrete implementations for notifications, storage, and time

### Interface Contracts

**Inbound Ports (Commands):** The service exposes command methods to control session lifecycle and query methods to inspect current state.

**Outbound Ports (Events):** The service depends on interfaces for notifications, storage, and time abstraction, enabling dependency injection and testability.

## Consequences

### Benefits

1. **Automated Testing Enabled**
   - Complete pomodoro cycles testable in milliseconds using mock time
   - No UI or manual interaction required
   - Deterministic, repeatable test execution

2. **Deadlocks Eliminated**
   - Service never calls back to UI
   - One-way event flow: UI → Service (commands), Service → UI (events)
   - No lock contention between layers

3. **Clear Separation of Concerns**
   - Business logic isolated in core domain
   - UI adapter only handles rendering and user input
   - Infrastructure (time, storage) abstracted behind interfaces

4. **Future Flexibility**
   - Easy to add new drivers (CLI, web UI, different test frameworks)
   - Can swap implementations of time, storage, notification
   - Business logic unchanged when infrastructure changes

5. **Better Maintainability**
   - Single source of truth for state machine
   - Business rules explicit and testable
   - Dependencies clearly defined through interfaces

### Drawbacks

1. **Increased Complexity**
   - More interfaces and indirection
   - Additional modules for core domain and adapters
   - Adapter pattern adds conceptual overhead

2. **More Files**
   - Split logic across more modules
   - New test harness infrastructure
   - May be harder to navigate initially

3. **Late Binding Required**
   - Service and adapters have circular awareness
   - Application bootstrap must wire components carefully
   - Notifier interface requires late binding pattern

### Mitigation

- Keep interfaces minimal and focused
- Document architecture clearly
- Provide example tests showing pattern usage
- Only add abstraction where it enables testing

## Alternatives Considered

### 1. Keep Callback-Driven Architecture, Add Test Hooks

**Approach:** Keep current structure but add test-specific hooks or flags to bypass UI initialization.

**Rejected because:**
- Doesn't solve deadlock problem
- Test hooks pollute production code
- Still tightly couples business logic to systray
- Doesn't enable clean, readable tests

### 2. Extract Business Logic to Helper Functions

**Approach:** Move logic from UI event handlers to pure functions, keep callback structure.

**Rejected because:**
- State management still scattered
- Doesn't enable full cycle testing
- Callbacks still create circular dependencies
- Partial solution that doesn't address root causes

### 3. Use Event Bus / Message Queue

**Approach:** Introduce pub/sub system for decoupling components.

**Rejected because:**
- Over-engineering for a lite-mode desktop app
- Adds runtime complexity and debugging difficulty
- Hexagonal achieves same decoupling with simpler, standard patterns

### 4. Model-View-Presenter (MVP) or MVVM

**Approach:** Use UI pattern to separate view from logic.

**Rejected because:**
- Still tied to UI framework lifecycle
- Doesn't abstract infrastructure (time, storage)
- Hexagonal more clearly separates testable core from all infrastructure

---

**References:**

- Hexagonal Architecture: Alistair Cockburn, "Hexagonal Architecture" (2005)
- Dependency Inversion Principle: Robert C. Martin, "Agile Software Development" (2003)
