# Implement: Tray Icon and Dropdown UI

## Context

- **Increment Goal**: Enable users to launch GoPomodoro, see a tray icon in the macOS menu bar, click it to open a dropdown panel with placeholder UI, and close it naturally
- **Key Non-Goals**: No timer logic, no functional buttons, no persistence, no notifications
- **Design Approach**: Use `getlantern/systray` for tray integration, create three components (main app, tray manager, ui window), build with Makefile
- **Constitution Mode**: `lite` — keep tests minimal but meaningful, manual testing acceptable for UI, simple logging

**Related Documents:**
- [increment.md](increment.md)
- [design.md](design.md)
- [CONSTITUTION.md](../../CONSTITUTION.md)

**Status**: Not started  
**Next step**: Step 1 — Initialize Go module and project structure

---

## 1. Workstreams

- **Workstream A** — Project Foundation & Build (Go module, Makefile, dependencies)
- **Workstream B** — Tray Icon Integration (systray library, icon management, click handling)
- **Workstream C** — UI Window/Panel (dropdown creation, show/hide, placeholder layout)
- **Workstream D** — Main Application Wiring (bootstrap, component integration, smoke test flag, logging)

---

## 2. Steps

### Step 1: Initialize Go module and project structure

- [ ] **Step 1: Initialize Go module and project structure**

**Workstream**: A  
**Based on Design**: §4 Architecture and Boundaries — Component Structure

**Files**:
- `go.mod`
- `cmd/gopomodoro/`
- `internal/tray/`
- `internal/ui/`
- `assets/`

**TDD Cycle**:

- **Red — Failing test first**:
  - Create `cmd/gopomodoro/main_test.go` with a simple test:
    ```go
    package main
    
    import "testing"
    
    func TestModuleExists(t *testing.T) {
        // This test will fail until go.mod exists
        // Validates module name when it does exist
    }
    ```
  - Run `go test ./cmd/gopomodoro` → fails (no go.mod, module not initialized)

- **Green — Make the test(s) pass**:
  - Run `go mod init github.com/juliangodesa/gopomodoro`
  - Create directory structure:
    ```bash
    mkdir -p cmd/gopomodoro
    mkdir -p internal/tray
    mkdir -p internal/ui
    mkdir -p assets
    mkdir -p bin
    ```
  - Update test to verify module name if needed
  - Run `go test ./...` → passes (module initialized)

- **Refactor — Clean up with tests green**:
  - Verify directory structure is correct
  - Ensure go.mod has correct module path
  - Add `.gitignore` for `bin/` directory

**CI / Checks**:
- `go mod verify`
- `go test ./...`

---

### Step 2: Create minimal Makefile

- [ ] **Step 2: Create minimal Makefile**

**Workstream**: A  
**Based on Design**: §8 CI/CD and Rollout — Makefile targets, Machine-Readable Artifacts

**Files**:
- `Makefile`

**TDD Cycle**:

- **Red — Failing test first**:
  - Try running `make build` in terminal → fails (no Makefile)
  - Observe error: "make: *** No targets specified and no makefile found"

- **Green — Make the test(s) pass**:
  - Create `Makefile` with minimal targets:
    ```makefile
    .PHONY: build run test clean
    
    BINARY_NAME=gopomodoro
    BINARY_PATH=bin/$(BINARY_NAME)
    
    build:
    	@mkdir -p bin
    	@go build -o $(BINARY_PATH) ./cmd/gopomodoro
    
    run: build
    	@$(BINARY_PATH)
    
    test:
    	@go test ./...
    
    clean:
    	@rm -rf bin/
    ```
  - Run `make build` → will fail until main.go exists, but Makefile itself works

- **Refactor — Clean up with tests green**:
  - Ensure targets use tabs (not spaces)
  - Verify `@` suppresses command echo for clean output
  - Test each target individually

**CI / Checks**:
- `make clean` succeeds
- `make test` runs (even if no tests yet)
- Makefile syntax is valid

---

### Step 3: Add systray dependency

- [ ] **Step 3: Add systray dependency**

**Workstream**: A  
**Based on Design**: §6 Contracts and Data — External Dependency (getlantern/systray)

**Files**:
- `go.mod`
- `go.sum`

**TDD Cycle**:

- **Red — Failing test first**:
  - Create a simple test that imports systray:
    ```go
    // internal/tray/tray_test.go
    package tray
    
    import (
        "testing"
        "github.com/getlantern/systray"
    )
    
    func TestSystrayImport(t *testing.T) {
        // Just validate we can import systray
        _ = systray.Run
    }
    ```
  - Run `go test ./internal/tray` → fails (dependency not found)

- **Green — Make the test(s) pass**:
  - Run `go get github.com/getlantern/systray`
  - Run `go mod tidy`
  - Run `go test ./internal/tray` → passes

- **Refactor — Clean up with tests green**:
  - Verify `go.mod` has the dependency listed
  - Verify `go.sum` is generated
  - Run `go mod verify` to ensure checksums are valid

**CI / Checks**:
- `go mod verify`
- `go build ./...`
- `go test ./...`

---

### Step 4: Create tray package with Initialize function

- [ ] **Step 4: Create tray package with Initialize function**

**Workstream**: B  
**Based on Design**: §6 Contracts and Data — internal/tray Package interface

**Files**:
- `internal/tray/tray.go`
- `internal/tray/tray_test.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Create `internal/tray/tray_test.go`:
    ```go
    package tray
    
    import "testing"
    
    func TestInitialize(t *testing.T) {
        err := Initialize()
        if err != nil {
            t.Fatalf("Initialize() failed: %v", err)
        }
    }
    ```
  - Run `go test ./internal/tray` → fails (no tray.go file, no Initialize function)

- **Green — Make the test(s) pass**:
  - Create `internal/tray/tray.go`:
    ```go
    // Package tray manages the system tray icon
    package tray
    
    // Initialize sets up the system tray icon
    func Initialize() error {
        // Stub for now - will integrate with systray later
        return nil
    }
    ```
  - Run `go test ./internal/tray` → passes

- **Refactor — Clean up with tests green**:
  - Add package documentation comment
  - Consider adding internal state tracking if needed
  - Keep it simple for now

**CI / Checks**:
- `go test ./internal/tray`
- `make test`

---

### Step 5: Implement icon loading and SetIcon

- [ ] **Step 5: Implement icon loading and SetIcon**

**Workstream**: B  
**Based on Design**: §6 Contracts and Data — SetIcon(iconData []byte)

**Files**:
- `internal/tray/tray.go`
- `internal/tray/tray_test.go`
- `assets/icon-idle.png`

**TDD Cycle**:

- **Red — Failing test first**:
  - Create a simple gray circle PNG icon (16x16 or 32x32) in `assets/icon-idle.png`
  - Add test for SetIcon:
    ```go
    func TestSetIcon(t *testing.T) {
        iconData := []byte{/* minimal PNG bytes or load from file */}
        err := SetIcon(iconData)
        if err != nil {
            t.Fatalf("SetIcon() failed: %v", err)
        }
    }
    ```
  - Run `go test ./internal/tray` → fails (no SetIcon function)

- **Green — Make the test(s) pass**:
  - Implement `SetIcon` in `tray.go`:
    ```go
    import "github.com/getlantern/systray"
    
    // SetIcon updates the tray icon image
    func SetIcon(iconData []byte) error {
        systray.SetIcon(iconData)
        return nil
    }
    ```
  - Note: systray.SetIcon() must be called from within systray.Run(), so this test may need adjustment or we accept it as a stub for now
  - Run `go test ./internal/tray` → passes

- **Refactor — Clean up with tests green**:
  - Add helper function to load icon from assets:
    ```go
    import (
        "os"
        "path/filepath"
    )
    
    func LoadIconFromAssets() ([]byte, error) {
        return os.ReadFile(filepath.Join("assets", "icon-idle.png"))
    }
    ```
  - Update test to use actual icon file if present
  - Handle errors properly

**CI / Checks**:
- `go test ./internal/tray`
- Verify icon file exists in `assets/`

---

### Step 6: Add click handler registration

- [ ] **Step 6: Add click handler registration**

**Workstream**: B  
**Based on Design**: §6 Contracts and Data — OnClick(handler func())

**Files**:
- `internal/tray/tray.go`
- `internal/tray/tray_test.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Add test for OnClick:
    ```go
    func TestOnClick(t *testing.T) {
        called := false
        handler := func() { called = true }
        
        OnClick(handler)
        
        // Verify handler was registered
        // (May need to expose internal state or accept simple registration test)
    }
    ```
  - Run `go test ./internal/tray` → fails (no OnClick function)

- **Green — Make the test(s) pass**:
  - Add package-level variable to store handler:
    ```go
    var clickHandler func()
    
    // OnClick registers a callback for tray icon clicks
    func OnClick(handler func()) {
        clickHandler = handler
    }
    ```
  - Run `go test ./internal/tray` → passes

- **Refactor — Clean up with tests green**:
  - Ensure handler is nil-safe when called
  - Consider adding getter for testing purposes
  - Document that handler is called from systray event loop

**CI / Checks**:
- `go test ./internal/tray`
- `make test`

---

### Step 7: Create ui package with CreateWindow function

- [ ] **Step 7: Create ui package with CreateWindow function**

**Workstream**: C  
**Based on Design**: §6 Contracts and Data — internal/ui Package interface

**Files**:
- `internal/ui/window.go`
- `internal/ui/window_test.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Create `internal/ui/window_test.go`:
    ```go
    package ui
    
    import "testing"
    
    func TestCreateWindow(t *testing.T) {
        win, err := CreateWindow()
        if err != nil {
            t.Fatalf("CreateWindow() failed: %v", err)
        }
        if win == nil {
            t.Fatal("CreateWindow() returned nil window")
        }
    }
    ```
  - Run `go test ./internal/ui` → fails (no ui package)

- **Green — Make the test(s) pass**:
  - Create `internal/ui/window.go`:
    ```go
    // Package ui manages the dropdown window/panel
    package ui
    
    // Window represents the dropdown panel
    type Window struct {
        visible bool
    }
    
    // CreateWindow initializes the dropdown window with placeholder UI
    func CreateWindow() (*Window, error) {
        return &Window{visible: false}, nil
    }
    ```
  - Run `go test ./internal/ui` → passes

- **Refactor — Clean up with tests green**:
  - Add package documentation
  - Consider what fields Window will need (position, content, etc.)
  - Keep minimal for now

**CI / Checks**:
- `go test ./internal/ui`
- `make test`

---

### Step 8: Implement Show and Hide methods

- [ ] **Step 8: Implement Show and Hide methods**

**Workstream**: C  
**Based on Design**: §6 Contracts and Data — window.Show(x, y), window.Hide()

**Files**:
- `internal/ui/window.go`
- `internal/ui/window_test.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Add tests for Show and Hide:
    ```go
    func TestShowHide(t *testing.T) {
        win, _ := CreateWindow()
        
        err := win.Show(100, 100)
        if err != nil {
            t.Fatalf("Show() failed: %v", err)
        }
        if !win.IsVisible() {
            t.Error("Window should be visible after Show()")
        }
        
        err = win.Hide()
        if err != nil {
            t.Fatalf("Hide() failed: %v", err)
        }
        if win.IsVisible() {
            t.Error("Window should not be visible after Hide()")
        }
    }
    ```
  - Run `go test ./internal/ui` → fails (no Show, Hide, or IsVisible methods)

- **Green — Make the test(s) pass**:
  - Implement methods:
    ```go
    // Show displays the window at the specified screen coordinates
    func (w *Window) Show(x, y int) error {
        w.visible = true
        // Actual window display logic will be added when integrating with UI framework
        return nil
    }
    
    // Hide conceals the window
    func (w *Window) Hide() error {
        w.visible = false
        return nil
    }
    
    // IsVisible returns whether the window is currently displayed
    func (w *Window) IsVisible() bool {
        return w.visible
    }
    ```
  - Run `go test ./internal/ui` → passes

- **Refactor — Clean up with tests green**:
  - Add position tracking if needed
  - Consider adding window content initialization
  - Document that actual UI rendering happens in integration

**CI / Checks**:
- `go test ./internal/ui`
- `make test`

---

### Step 9: Add placeholder UI layout

- [ ] **Step 9: Add placeholder UI layout**

**Workstream**: C  
**Based on Design**: §5 Scope — Placeholder UI Elements (header, timer, buttons)

**Files**:
- `internal/ui/window.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Manual test expectation: When `CreateWindow()` is called, window should have structure for header, timer, and buttons
  - Since this is UI layout, automated testing is minimal (lite mode)
  - Expectation: Code should initialize placeholder content

- **Green — Make the test(s) pass**:
  - Update `CreateWindow()` to initialize placeholder layout:
    ```go
    import "github.com/getlantern/systray"
    
    type Window struct {
        visible bool
        // Placeholder UI elements (using systray menu items as fallback)
        header  *systray.MenuItem
        timer   *systray.MenuItem
        btnStart *systray.MenuItem
        btnPause *systray.MenuItem
        btnReset *systray.MenuItem
    }
    
    func CreateWindow() (*Window, error) {
        w := &Window{visible: false}
        // Note: Actual menu item creation must happen inside systray.Run()
        // This is a structural placeholder
        return w, nil
    }
    
    // InitializeMenu creates the actual menu items (called from systray.Run)
    func (w *Window) InitializeMenu() {
        w.header = systray.AddMenuItem("Ready", "Current state")
        w.header.Disable()
        
        w.timer = systray.AddMenuItem("25:00", "Timer display")
        w.timer.Disable()
        
        systray.AddSeparator()
        
        w.btnStart = systray.AddMenuItem("Start", "Start timer")
        w.btnStart.Disable() // Non-functional for this increment
        
        w.btnPause = systray.AddMenuItem("Pause", "Pause timer")
        w.btnPause.Disable()
        
        w.btnReset = systray.AddMenuItem("Reset", "Reset timer")
        w.btnReset.Disable()
    }
    ```

- **Refactor — Clean up with tests green**:
  - Extract menu creation to separate method
  - Document that this uses systray menu as fallback (simpler than native window for lite mode)
  - Consider if native window is needed or if systray menu is sufficient

**CI / Checks**:
- `go build ./internal/ui` compiles
- Visual inspection via `make run` (after integration)

---

### Step 10: Create main.go with basic structure and -smoke flag

- [ ] **Step 10: Create main.go with basic structure and -smoke flag**

**Workstream**: D  
**Based on Design**: §4 Architecture — cmd/gopomodoro/main.go

**Files**:
- `cmd/gopomodoro/main.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Try `go build ./cmd/gopomodoro` → fails (no main.go)
  - Try `make build` → fails with no main package

- **Green — Make the test(s) pass**:
  - Create `cmd/gopomodoro/main.go`:
    ```go
    package main
    
    import (
        "flag"
        "log"
    )
    
    var smokeTest = flag.Bool("smoke", false, "Run smoke test (start and immediately exit)")
    
    func main() {
        flag.Parse()
        
        log.Println("[INFO] GoPomodoro starting...")
        
        if *smokeTest {
            log.Println("[INFO] Smoke test mode - exiting immediately")
            return
        }
        
        // Actual app initialization will be added in next steps
        log.Println("[INFO] Ready (not yet integrated with systray)")
    }
    ```
  - Run `make build` → succeeds
  - Run `./bin/gopomodoro -smoke` → exits immediately with logs

- **Refactor — Clean up with tests green**:
  - Ensure logging format is consistent
  - Add basic error handling structure
  - Document flags

**CI / Checks**:
- `make build` succeeds
- `./bin/gopomodoro -smoke` exits with code 0
- `make test` still passes

---

### Step 11: Wire tray initialization in main.go

- [ ] **Step 11: Wire tray initialization in main.go**

**Workstream**: D  
**Based on Design**: §3 Proposed Solution — Data Flow (systray.Run with onReady callback)

**Files**:
- `cmd/gopomodoro/main.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Run `make run` → app starts but no tray icon appears
  - Manual test: Look for tray icon in menu bar → not there

- **Green — Make the test(s) pass**:
  - Update `main.go` to call systray.Run:
    ```go
    import (
        "github.com/getlantern/systray"
        "github.com/juliangodesa/gopomodoro/internal/tray"
    )
    
    func main() {
        flag.Parse()
        
        log.Println("[INFO] GoPomodoro starting...")
        
        if *smokeTest {
            log.Println("[INFO] Smoke test mode - running minimal initialization")
            systray.Run(onReady, onExit)
            return
        }
        
        systray.Run(onReady, onExit)
    }
    
    func onReady() {
        log.Println("[INFO] Initializing tray icon...")
        
        if err := tray.Initialize(); err != nil {
            log.Fatalf("[ERROR] Failed to initialize tray: %v", err)
        }
        
        // Load and set icon
        iconData, err := tray.LoadIconFromAssets()
        if err != nil {
            log.Fatalf("[ERROR] Failed to load icon: %v", err)
        }
        
        if err := tray.SetIcon(iconData); err != nil {
            log.Fatalf("[ERROR] Failed to set icon: %v", err)
        }
        
        systray.SetTooltip("GoPomodoro")
        log.Println("[INFO] Tray icon initialized successfully")
        
        if *smokeTest {
            log.Println("[INFO] Smoke test - quitting after initialization")
            systray.Quit()
        }
    }
    
    func onExit() {
        log.Println("[INFO] Application shutting down")
    }
    ```
  - Run `make run` → tray icon appears in menu bar
  - Run `./bin/gopomodoro -smoke` → tray icon briefly appears, then exits

- **Refactor — Clean up with tests green**:
  - Extract onReady and onExit to separate functions
  - Add error handling
  - Ensure smoke test path works correctly

**CI / Checks**:
- `make run` → tray icon visible (manual check)
- `./bin/gopomodoro -smoke` → exits cleanly after brief init
- No crashes or errors in logs

---

### Step 12: Wire UI window creation and click handling

- [ ] **Step 12: Wire UI window creation and click handling**

**Workstream**: D  
**Based on Design**: §3 Proposed Solution — Data Flow (click event → window toggle)

**Files**:
- `cmd/gopomodoro/main.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Run `make run` → click tray icon → nothing happens
  - Manual test: No menu or window appears on click

- **Green — Make the test(s) pass**:
  - Update `onReady()` to create window and register click handler:
    ```go
    import (
        "github.com/juliangodesa/gopomodoro/internal/ui"
    )
    
    func onReady() {
        // ... existing tray initialization code ...
        
        // Create UI window
        window, err := ui.CreateWindow()
        if err != nil {
            log.Fatalf("[ERROR] Failed to create window: %v", err)
        }
        
        // Initialize menu items (systray menu approach)
        window.InitializeMenu()
        log.Println("[INFO] Dropdown window created")
        
        // For systray menu approach, menu is always "ready" to show
        // No explicit click handler needed - systray manages menu display
        
        if *smokeTest {
            log.Println("[INFO] Smoke test - quitting after initialization")
            systray.Quit()
        }
    }
    ```
  - Run `make run` → click icon → menu appears with placeholder items

- **Refactor — Clean up with tests green**:
  - If using custom window approach (not systray menu), add click handler:
    ```go
    // Alternative: Custom window toggle (if not using systray menu)
    tray.OnClick(func() {
        log.Println("[INFO] Tray icon clicked")
        if window.IsVisible() {
            window.Hide()
            log.Println("[INFO] Window hidden")
        } else {
            window.Show(100, 100) // Position would be calculated based on tray icon
            log.Println("[INFO] Window shown")
        }
    })
    ```
  - Document which approach is being used
  - Keep it simple for lite mode

**CI / Checks**:
- `make run` → click icon → menu/window appears (manual check)
- Menu shows: Ready, 25:00, Start, Pause, Reset
- All items are disabled (non-functional as expected)

---

### Step 13: Add LSUIElement for dock icon suppression

- [ ] **Step 13: Add LSUIElement for dock icon suppression**

**Workstream**: D  
**Based on Design**: §5 Scope — Suppress dock icon (LSUIElement)

**Files**:
- `cmd/gopomodoro/main.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Run `make run` → check Activity Monitor and Dock → dock icon may appear
  - Manual test: App should have no dock icon, only tray icon

- **Green — Make the test(s) pass**:
  - Research macOS LSUIElement approach for Go applications
  - For lite mode, acceptable approaches:
    - Runtime configuration (if available via systray or macOS APIs)
    - Document limitation if dock icon appears
    - For distribution, would need Info.plist with LSUIElement=true
  - Add code or documentation:
    ```go
    // Note: To fully suppress dock icon, this app should be built as a .app bundle
    // with Info.plist containing LSUIElement=true
    // For development builds, dock icon may appear
    // This is acceptable for lite mode
    ```
  - If there's a runtime approach available, implement it

- **Refactor — Clean up with tests green**:
  - Document the approach taken
  - Note in code comments what's needed for distribution
  - Accept limitation for development builds if necessary

**CI / Checks**:
- `make run` → verify dock icon behavior (manual)
- Document findings in code comments

---

### Step 14: Add logging for key events

- [ ] **Step 14: Add logging for key events**

**Workstream**: D  
**Based on Design**: §9 Observability — Logging strategy

**Files**:
- `cmd/gopomodoro/main.go`
- `internal/tray/tray.go`
- `internal/ui/window.go`

**TDD Cycle**:

- **Red — Failing test first**:
  - Run `make run` → observe logs → may be minimal or inconsistent
  - Expected: Structured logs for all key events

- **Green — Make the test(s) pass**:
  - Ensure consistent logging throughout:
    ```go
    // In main.go - already added in previous steps
    log.Println("[INFO] GoPomodoro starting...")
    log.Println("[INFO] Tray icon initialized successfully")
    log.Println("[INFO] Dropdown window created")
    log.Println("[INFO] Application shutting down")
    
    // In tray/tray.go
    func Initialize() error {
        log.Println("[INFO] Tray initialization called")
        return nil
    }
    
    func SetIcon(iconData []byte) error {
        systray.SetIcon(iconData)
        log.Printf("[INFO] Tray icon set (%d bytes)", len(iconData))
        return nil
    }
    
    // In ui/window.go
    func (w *Window) Show(x, y int) error {
        w.visible = true
        log.Printf("[INFO] Window shown at position (x: %d, y: %d)", x, y)
        return nil
    }
    
    func (w *Window) Hide() error {
        w.visible = false
        log.Println("[INFO] Window hidden")
        return nil
    }
    ```
  - Run `make run` → see structured logs in terminal

- **Refactor — Clean up with tests green**:
  - Ensure all log statements use consistent format: `[INFO]`, `[ERROR]`
  - Add timestamps if desired (log package includes them by default)
  - Keep logging simple and readable

**CI / Checks**:
- `make run` → verify logs appear in stdout
- Check log format is consistent
- Ensure no excessive logging that clutters output

---

## 3. Rollout & Validation Notes

### Suggested Grouping into PRs

For lite mode, the entire increment could be a single PR, but if desired to split:

- **PR 1: Foundation** (Steps 1-3)
  - Go module, Makefile, dependencies
  - Validates: Project structure is set up, builds work

- **PR 2: Core Packages** (Steps 4-9)
  - Tray and UI packages with tests
  - Validates: Components can be built and tested independently

- **PR 3: Integration** (Steps 10-14)
  - Main application wiring, smoke test, logging
  - Validates: End-to-end functionality, tray icon appears, menu works

### Validation Checkpoints

**After Step 3:**
- `make build` succeeds
- `make test` runs (even with minimal tests)
- Project structure is clean

**After Step 9:**
- All package tests pass: `go test ./internal/...`
- Packages compile: `go build ./internal/...`

**After Step 12:**
- Manual test: Run `make run`, click tray icon, see menu with placeholder items
- All menu items are visible and disabled
- App runs without crashes

**After Step 14:**
- Smoke test: `./bin/gopomodoro -smoke` exits cleanly with logs
- Full run: `make run` shows structured logs and working tray functionality
- All acceptance criteria from increment.md are met

### Final Manual Testing Checklist

From increment.md acceptance criteria:

1. ✓ Tray icon is visible in macOS menu bar
2. ✓ Icon displays gray/idle state
3. ✓ Left-click opens dropdown (menu)
4. ✓ Dropdown shows structured layout (header, timer, buttons)
5. ✓ Dropdown closes properly (systray manages this)
6. ✓ App runs in background (verify no dock icon or document limitation)

### Smoke Test Usage

The `-smoke` flag enables automated testing:
- `./bin/gopomodoro -smoke` — Start app, initialize tray, quit immediately
- Useful for CI/CD validation (when CI is added)
- Verifies app can start and shut down cleanly
- Can be incorporated into `make test` if desired
