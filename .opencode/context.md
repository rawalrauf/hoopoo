# Hoopoo Project Context

## Project Overview

Hoopoo is an Arch Linux-based distribution focused on performance and user experience. This repository contains the core system commands and utilities.

## Repository Structure

### Development Repositories (`~/.local/share/hoopoo-repos/`)

**hoopoo/** (this repo)
- Main development repository for core commands and utilities
- Currently on `go` branch doing bash → Go migration
- 282 total commands, migrating and testing one by one
- Commands live in `bin/hoopoo-*` (bash) and `cmd/hoopoo-*/main.go` (Go)
- Shared Go libraries in `pkg/` with professional sub-package architecture:
  - `pkg/notification/` - desktop notifications (mako integration)
    - `pkg/notification/progress/` - progress bar formatting
  - `pkg/battery/` - battery-related operations
  - `pkg/hardware/` - hardware detection (restructured with sub-packages)
    - `pkg/hardware/cpu/` - CPU detection
    - `pkg/hardware/dmi/` - DMI/SMBIOS queries
    - `pkg/hardware/gpu/` - GPU detection
    - `pkg/hardware/input/` - input device detection
    - `pkg/hardware/audio/` - audio device detection
    - `pkg/hardware/webcam/` - webcam device detection
  - `pkg/ytdlp/` - yt-dlp download operations
  - `pkg/media/` - media processing
    - `pkg/media/video/` - video thumbnail generation and processing
    - `pkg/media/transcode/` - picture and video transcoding
    - `pkg/media/recorder/` - screen recording with webcam overlay
  - `pkg/portal/` - XDG Desktop Portal integration via D-Bus
    - `pkg/portal/fileselect.go` - native file/folder picker (replaces fzf/walker for file selection)
  - `pkg/niri/` - niri compositor operations (GetWindows, GetInternalDisplay, FocusWindow, EventStream, etc.)
    - `pkg/niri/config/` - KDL config editor with comment preservation
  - `pkg/toggle/` - generic toggle state management (Enabled, Disabled, Clear, Write)
  - `pkg/font/` - font detection and management (state file-based)
  - `pkg/walker/` - font-aware walker dimensions and launcher utilities
  - `pkg/system/paths/` - path utilities (UserConfig, State, Hoopoo, etc.)
  - `pkg/system/process/` - process management (IsRunning, Kill, KillForce, KillByPattern, Start, SendSignal)
  - `pkg/system/hooks/` - hook execution
  - `pkg/system/systemd/` - systemd service management

**hoopoo-bash/**
- Reference copy of bash scripts (on `dev` branch)
- Used for benchmarking Go binaries against original bash implementations
- Reference for fixing/debugging migrated commands
- Keep this on `dev` or `master` branch, never switch to `go` branch

**hoopoo-pkgs/**
- Package definitions for Hoopoo-specific packages
- Build system for AUR and custom packages

**hoopoo-iso/**
- ISO builder and installer
- Live environment configuration

### Installed Version (`~/.local/share/hoopoo/`)

- Production installation (currently v1.0.0)
- Live system running bash commands
- Used for local testing: after building Go binaries in `hoopoo-repos/hoopoo/`, we copy them here to test
- Testing workflow: `make all` in dev repo → copy to `bin/` → test commands

### Upstream Reference (`~/.local/share/omarchy-repos/`)

Hoopoo was originally forked from an older version of Omarchy. Upstream has since grown significantly.

**omarchy/** - Current upstream Omarchy (reference for better approaches)
**old-omarchy/** - The version Hoopoo was based on
**omarchy-iso/** and **omarchy-pkgs/** - Their build repositories

## Migration Goals

**Why Go?**
1. **Performance** - 10-100x faster than bash for system queries
2. **Optimization** - Parallel execution, efficient sysfs reads, minimal overhead
3. **Full-scale distro** - More robust, maintainable codebase for production use
4. **Type safety** - Catch errors at compile time
5. **Single binary** - No shell interpreter dependencies

## Go Migration Workflow (Bash → Go)

When porting bash commands to Go, follow this workflow:

### 1. Analysis Phase

**For each command to port:**

1. Read the hoopoo bash implementation in `bin/hoopoo-<command>`
2. Read the omarchy bash implementation in `~/.local/share/omarchy-repos/omarchy/bin/omarchy-<command>` (if exists)
3. Compare approaches and identify:
   - Edge cases handled by omarchy but missing in hoopoo
   - Better algorithms or optimizations
   - Hardware compatibility issues (Apple Silicon, Chromebooks, Framework)
   - Multi-device scenarios (multiple batteries, monitors, audio sinks)
4. **Choose the better approach** - prefer omarchy's robustness if it handles more edge cases

### 2. Implementation Phase

**Design the Go version:**

1. Create shared library code in `pkg/` packages for reusable logic, using sub-packages for organization
2. Keep command entry points in `cmd/hoopoo-<command>/main.go` minimal (~15-30 lines)
3. **Package organization:**
   - Use professional `pkg/` structure with sub-packages (e.g., `pkg/notification/progress/`)
   - Group by domain functionality (notification, battery, hardware, media, system utilities)
   - Use short function names (`battery.IsPresent()` not `battery.BatteryIsPresent()`)
   - Commands call each other via dispatcher: `exec.Command("hoopoo", "group", "command")`
4. **Use goroutines where beneficial:**
   - Parallel file reads (multiple sysfs files)
   - Concurrent API calls (D-Bus, IPC)
   - Multi-device queries (all batteries, all audio sinks)
   - Background tasks (monitoring, watching for changes)
5. Avoid goroutines for single operations (overhead > benefit)
6. Apply performance optimizations:
   - Cache expensive lookups (device discovery, file paths)
   - Minimize allocations (buffer reuse, manual parsing)
   - Direct system calls where possible (sysfs > upower, PipeWire > pulseaudio)

### 3. Verification Phase

**Test correctness before benchmarking:**

1. Build: `make all` (builds all Go commands)
2. Copy to test environment: `cp bin/hoopoo-* ~/.local/share/hoopoo/bin/`
3. Copy metadata sources: `cp -r cmd/hoopoo-* ~/.local/share/hoopoo/cmd/`
4. Run both bash and Go versions side-by-side with same inputs
5. Compare outputs - they must match exactly (unless fixing a bug)
6. Test edge cases:
   - No hardware present (no battery, no audio)
   - Multiple devices (BAT0 + BAT1, multiple monitors)
   - Missing sysfs files (older kernels, different hardware)
   - Invalid inputs (empty strings, out of range values)
7. **Only proceed to benchmarks when output is identical**

### 4. Benchmark Phase

**Measure performance improvements:**

1. Use `hyperfine --warmup 3 --runs 50` for accurate measurements
2. Compare bash vs Go versions
3. Document speedup ratio and memory usage
4. Verify the Go version is actually faster (no regressions)

### 5. Final Review Phase

**Critical - analyze the code you just wrote:**

1. Read all Go files you created (`cmd/` and `internal/`)
2. Look for:
   - **Concurrency opportunities missed** - can goroutines speed this up?
   - Logic errors or off-by-one bugs
   - Missing error handling
   - Unnecessary allocations or string conversions
   - Hardcoded paths that should be configurable
   - Edge cases not handled
   - Code duplication that should be extracted
3. Apply any improvements found
4. Rebuild and re-benchmark if changes are significant

### 6. Integration Phase

**After all checks pass:**

1. Update `Makefile` to build the new Go command
2. Build and copy to test environment
3. Verify with `hoopoo commands --check` (should show 282 commands)
4. Test the command through dispatcher: `hoopoo <group> <command>`
5. Commit to `go` branch with descriptive message

## Goroutine Guidelines

**Use goroutines for:**
- Reading multiple independent sysfs files in parallel
- Querying multiple devices concurrently (all batteries, all audio sinks)
- D-Bus calls that can run in parallel
- Watch/monitor commands that need event loops

**Avoid goroutines for:**
- Single file reads (overhead > benefit)
- Sequential operations with dependencies
- Commands that run once and exit immediately
- Simple data transformations

**Pattern for parallel file reads:**
```go
type result struct {
  value int
  err   error
}

ch := make(chan result, len(files))
for _, file := range files {
  go func(f string) {
    val, err := readInt(f)
    ch <- result{val, err}
  }(file)
}

for range files {
  r := <-ch
  // handle result
}
```

## Command Metadata System

Commands use `// hoopoo:` comments for metadata (both bash `#` and Go `//` styles supported):

```go
// hoopoo:summary=Short description ending with period.
// hoopoo:args=<required> [optional]
// hoopoo:examples=hoopoo cmd example1 | hoopoo cmd example2
// hoopoo:hidden=true
```

The dispatcher (`bin/hoopoo`) reads metadata directly from source files in `cmd/*/main.go` to build help and routing.

## Command Format Rules

**CRITICAL: All commands MUST use fully space-separated format with NO dashes in subcommands:**

✅ **Correct:**
```go
exec.Command("hoopoo", "notification", "send", ...)
exec.Command("hoopoo", "launch", "or", "focus", ...)
exec.Command("hoopoo", "niri", "toggle", "disabled", ...)
exec.Command("hoopoo", "swayosd", "kbd", "brightness", ...)
```

❌ **Wrong:**
```go
exec.Command("hoopoo", "notification-send", ...)           // NO dashes in subcommands
exec.Command("hoopoo", "launch-or-focus", ...)            // NO dashes in subcommands
exec.Command("hoopoo", "niri-toggle-disabled", ...)       // NO dashes in subcommands
exec.Command("hoopoo", "swayosd", "kbd-brightness", ...)  // NO dashes in subcommands
```

**Why:** The bash dispatcher automatically converts spaces to dashes when finding binaries:
- User types: `hoopoo notification send "Title"`
- Dispatcher finds: `hoopoo-notification-send` binary
- This prepares for single-binary migration where subcommands are naturally space-separated

**Usage messages:** Always show space-separated format to match what users type:
```go
fmt.Fprintln(os.Stderr, "Usage: hoopoo notification send <headline> [description]")
```

**Binary names:** Keep dashes in actual binary names (`hoopoo-notification-send`) until single-binary migration.

## Single Binary Migration Plan (Final Phase)

**After all bash scripts are migrated and tested, consolidate into a single `hoopoo` binary.**

### Architecture:

**Current (multi-binary):**
```
bin/hoopoo-notification-send    (standalone binary, 2MB)
bin/hoopoo-launch-or-focus      (standalone binary, 2MB)
bin/hoopoo-brightness-keyboard  (standalone binary, 2MB)
... (38 binaries × 2-5MB = 100-200MB total)
```

**Future (single binary):**
```
bin/hoopoo                      (single binary, 5-10MB)
```

### Migration Steps:

1. **Rename all cmd directories** (remove `hoopoo-` prefix):
   ```bash
   cmd/hoopoo-notification-send/ → cmd/notification-send/
   cmd/hoopoo-launch-or-focus/   → cmd/launch-or-focus/
   cmd/hoopoo-brightness-keyboard/ → cmd/brightness-keyboard/
   ```

2. **Convert all commands from `main()` to `Run(args []string)`:**
   ```go
   // Before:
   package main
   func main() {
     args := os.Args[1:]
     // ... logic
   }
   
   // After:
   package notification_send
   func Run(args []string) {
     // ... same logic, use args parameter
   }
   ```

3. **Create smart dispatcher** (`main.go` at repo root):
   ```go
   package main
   
   import (
     "os"
     "strings"
     // Auto-import all cmd packages
   )
   
   func main() {
     args := os.Args[1:]
     
     // Auto-discover and route to commands based on directory scan
     // Convert space-separated args to package names
     // Call appropriate Run() function
   }
   ```

4. **Smart dispatcher features:**
   - Auto-discovers all packages in `cmd/` directory
   - Converts space-separated args to command routing (e.g., `notification send` → `notification_send.Run()`)
   - No manual route registration needed - reads from `cmd/` structure
   - Supports multi-level commands (e.g., `niri toggle disabled` → `niri_toggle_disabled.Run()`)
   - Parses metadata from command files for help generation

5. **Build single binary:**
   ```bash
   go build -ldflags="-s -w" -o bin/hoopoo main.go
   ```

### Benefits:

- ✅ **Size:** 5-10MB instead of 100-200MB (90-95% reduction)
- ✅ **Speed:** No exec overhead between commands (10-100x faster for chained operations)
- ✅ **Distribution:** Single file to deploy
- ✅ **Maintenance:** One binary to build, test, and version
- ✅ **Natural routing:** Space-separated commands work natively without bash dispatcher

### Timeline:

- **Phase 1 (current):** Migrate all bash scripts to Go, test individually
- **Phase 2 (testing):** Verify all 282 commands work perfectly with multi-binary approach
- **Phase 3 (consolidation):** Execute single-binary migration (1-2 hours, mostly mechanical)

### Notes:

- Keep bash dispatcher during migration for compatibility
- All space-separated command format already matches single-binary architecture
- Package names use underscores (Go requirement): `launch_or_focus`, `notification_send`
- Directory names keep dashes for readability: `cmd/launch-or-focus/`

## Module Information

- **Module path:** `hoopoo` (changed from `hoopoo.net`)
- **Architecture:** Professional pkg/ structure with sub-packages for organization
- **Command execution:** Commands call each other via dispatcher: `exec.Command("hoopoo", "notification", "send", ...)`
- **Build:** `make all` compiles all Go commands to `bin/hoopoo-*`
- **Binary stripping:** `-ldflags="-s -w"` for smaller binaries

## Current Migration Status

**Completed command groups:**
- `battery-*` (3 commands: present, remaining, monitor)
- `hw-*` (21 commands: hardware detection, GPU checks, input devices, monitor recovery, cleanup)
- `niri-*` (15 commands: monitor control, window management, toggle utilities, event stream)
- `toggle-*` (5 commands: touchpad, touchscreen, generic toggle utilities)
- `notification-*` (2 commands: send, dismiss - desktop notifications with progress bars)
- `version-*` (3 commands: version, branch, pkgs)
- `restart-*` (17 commands: mako, niri, waybar, swayosd, etc.)
- `refresh-*` (14 commands: config file refresh utilities)
- `launch-*` (15 commands: application launchers with focus logic, walker integration)
- `font-*` (3 commands: current, list, set - state file-based font management)
- `menu-*` (2 commands: keybindings, tmux-keybindings - interactive search menus)
- `sudo-reset` (1 command: reset faillock)
- `capture-*` (5 commands: screenshot with satty, screenrecording with webcam, text extraction OCR, color picker, transcode)
- `share-*` (3 commands: clipboard, file, folder - LocalSend integration with portal file picker)
- `transcode` (1 command: picture/video transcoding with portal file picker)
- `file-select` (1 command: native file/folder picker using XDG Desktop Portal)

**Total:** 111 Go commands migrated

## Architecture Principles

**Domain-focused organization:**
- Organize by what (battery, hw, notification) not how (sysfs, dbus)
- Commands call binaries instead of internal functions when possible
- Internal packages provide reusable functionality
- Clean separation between infrastructure and domain logic

**Binary vs Package usage (CRITICAL):**
- **NEVER call hoopoo binaries from pkg/** - pkg/ is pure logic only
- **In cmd/, prefer binaries over pkg when available:**
  - Use `exec.Command("hoopoo", "notification", "send", ...)` instead of `notification.Send()`
  - Use `exec.Command("hoopoo", "restart", "waybar")` instead of direct systemd calls
  - Use `exec.Command("hoopoo", "file", "select", ...)` for interactive file picking instead of walker/fzf
  - Exception: When pkg provides the ONLY interface (e.g., `niri.GetWindows()`, `font.Current()`, `portal.SelectFiles()`)
- **Use packages directly for:**
  - System infrastructure (paths, toggle, process, hooks)
  - Complex operations with no binary equivalent (niri.GetWindows, niri.EventStream)
  - State file operations (toggle.Write, font.StateFile)
  - Reading configuration (font.Current, walker dimensions)
  - Portal file selection when calling from pkg context (rare, prefer binary in cmd/)
- **Rule:** pkg/ should only call system tools (fc-list, niri msg, notify-send), never hoopoo binaries
- **Exception:** Low-level operations that form the foundation for binaries use pkg directly

**Examples of correct usage:**
```go
// ✅ CORRECT in cmd/
exec.Command("hoopoo", "notification", "send", "Title").Run()
exec.Command("hoopoo", "restart", "waybar").Run()
exec.Command("hoopoo", "file", "select", "--multiple").Output() // Interactive file picking
fontName := font.Current() // No binary equivalent, pkg is source
dims := walker.Get() // Reading config, not executing action
paths, _ := portal.SelectFiles(...) // Direct pkg usage when binary call is awkward

// ✅ CORRECT in pkg/
exec.Command("fc-list", ":spacing=100").Output() // System tool
exec.Command("niri", "msg", "action", "quit").Run() // System tool
exec.Command("notify-send", "Title").Run() // System tool

// ❌ WRONG in pkg/
exec.Command("hoopoo", "restart", "waybar").Run() // Never call hoopoo from pkg
exec.Command("hoopoo", "notification", "send").Run() // Never call hoopoo from pkg
exec.Command("hoopoo", "file", "select").Run() // Never call hoopoo from pkg
```

## Code Style

**Go conventions:**
- Two spaces for indentation
- Short, clear function names without package stuttering
- One-line package comments only (no function docs in internal packages)
- Minimal comments in cmd files (just metadata)
- No "Get" prefixes (`Level()` not `GetLevel()`)

**Toggle API (generic package for all toggle operations):**
```go
// Generic toggles (works with any path relative to toggles/)
toggle.Enabled("waybar-off")                    → ~/.local/state/hoopoo/toggles/waybar-off
toggle.Enabled("niri/gaps-enabled.kdl")         → ~/.local/state/hoopoo/toggles/niri/gaps-enabled.kdl
toggle.Disabled("waybar-off")                   → inverse of Enabled
toggle.Clear("waybar-off")                      → removes toggle file
toggle.Write("waybar-off", "content")           → writes toggle file

// Niri wrappers (convenience functions for niri toggles)
niri.StateEnabled("gaps-enabled.kdl")           → toggle.Enabled("niri/gaps-enabled.kdl")
niri.StateDisabled("gaps-enabled.kdl")          → toggle.Disabled("niri/gaps-enabled.kdl")
niri.ClearState("gaps-enabled.kdl")             → toggle.Clear("niri/gaps-enabled.kdl")
niri.WriteState("gaps-enabled.kdl", "content")  → toggle.Write("niri/gaps-enabled.kdl", "content")
```

**Architecture principle:**
- `pkg/toggle/` is generic infrastructure (system-level)
- `pkg/niri/` wraps toggle for niri-specific convenience
- Commands use niri wrappers when working with niri toggles
- Commands use toggle directly for non-niri toggles

**Notification API:**
```go
// Command-line usage (NEW format with --glyph flag):
exec.Command("hoopoo", "notification", "send", "--glyph=󰔛", "Headline", "Description", "-u", "critical").Run()

// Package usage:
notification.Send(notification.Options{
  Glyph:       "🔋",
  Headline:    "Battery Low",
  Description: "15% remaining",
  Urgency:     "critical",
  Icon:        "battery-caution",
  Timeout:     5000,
  SoundFile:   "/usr/share/sounds/freedesktop/stereo/suspend-error.oga",
})
```

Progress bars use `pkg/notification/progress/Format()` with `■□` blocks.

All fields are optional.

**Portal file selection (replaces walker/fzf for file picking):**

**ALWAYS use portal for interactive file selection** - provides native file manager dialog across all desktop environments.

```go
// Binary usage (from any command):
cmd := exec.Command("hoopoo", "file", "select",
  "--title", "Select files to transcode",
  "--extensions", "mp4 mov mkv jpg png",
  "--multiple",  // Allow multiple selection
)
output, err := cmd.Output()
// Returns newline-separated file paths

// For directories:
cmd := exec.Command("hoopoo", "file", "select",
  "--title", "Select folder to share",
  "--directory",
)

// Package usage (direct API):
import "hoopoo/pkg/portal"

paths, err := portal.SelectFiles(portal.FileSelectOptions{
  Title:      "Select media files",
  Extensions: []string{"mp4", "mov", "mkv", "jpg", "png"},
  Multiple:   true,     // Allow multiple files
  Directory:  false,    // false for files, true for folders
})
```

**When to use portal:**
- ✅ Any interactive file/folder selection in commands
- ✅ Transcode, share, file uploads, media selection
- ✅ Replaces walker file menus and fzf file pickers
- ❌ Don't use for config file paths (use paths.UserConfig)
- ❌ Don't use for programmatic file discovery (use filepath.Walk)

**Portal error handling:**
```go
if err != nil {
  if errors.Is(err, portal.ErrNoFilesSelected) || errors.Is(err, portal.ErrSelectionCancelled) {
    os.Exit(0) // User cancelled - exit silently
  }
  fmt.Fprintf(os.Stderr, "error: %v\n", err)
  os.Exit(1)
}
```

**Benefits:**
- Native file manager integration (GNOME Files, Dolphin, etc.)
- Works in all desktop environments (GNOME, KDE, wlroots)
- Sandboxed app support (Flatpak-safe)
- Better UX than terminal-based pickers
- File type filtering built-in

**Common patterns:**
- Use `paths.State()`, `paths.Hoopoo()`, `paths.Runtime()` for path resolution
- Use `portal.SelectFiles()` for interactive file/folder selection (replaces walker/fzf)
- Call external binaries with `exec.Command("hoopoo", "group", "command")` via dispatcher
- Use pkg/ packages for reusable logic, organize with sub-packages when needed
- Commands are modular and call each other, designed for eventual single-binary consolidation

## Success Criteria

A command port is complete when:
- ✅ Output matches bash version exactly (or fixes a known bug)
- ✅ All edge cases tested and handled
- ✅ Benchmarks show significant speedup (>2x minimum)
- ✅ Final code review completed with no issues found
- ✅ Goroutines used where beneficial for performance
- ✅ Command works through dispatcher in test environment
