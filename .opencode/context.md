# Hoopoo Project Context

## Project Overview

Hoopoo is an Arch Linux-based distribution focused on performance and user experience. This repository contains the core system commands and utilities.

## Repository Structure

### Development Repositories (`~/.local/share/hoopoo-repos/`)

**hoopoo/** (this repo)
- Main development repository for core commands and utilities
- Currently on `go` branch doing bash → Go migration
- 282 total commands/script, we are migrating and testing one by one
- Commands live in `bin/hoopoo-*` (bash) and `cmd/hoopoo-*/main.go` (Go)
- Shared Go libraries in `internal/` organized by domain:
  - `internal/sysfs/battery/` - battery-related sysfs operations
  - `internal/sysfs/hw/` - hardware detection
  - `internal/common/notify/` - desktop notifications
  - `internal/common/hooks/` - hook execution
  - `internal/common/paths/` - path utilities

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

1. Create shared library code in `internal/` packages for reusable logic
2. Keep command entry points in `cmd/hoopoo-<command>/main.go` minimal (~15-30 lines)
3. **Package organization:**
   - Group by infrastructure domain (`sysfs/battery/`, `sysfs/hw/`, `dbus/`, etc.)
   - Keep cross-cutting concerns in `common/` (notify, hooks, paths)
   - Use short function names (`battery.IsPresent()` not `battery.BatteryIsPresent()`)
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

## Module Information

- **Module path:** `hoopoo.net` (no `/hoopoo` suffix)
- **Architecture:** Separate binaries per command (not dispatcher-based)
- **Build:** `make all` compiles all Go commands to `bin/hoopoo-*`
- **Binary stripping:** `-ldflags="-s -w"` for smaller binaries

## Current Migration Status

**Completed command groups:**
- `battery-*` (3 commands: present, remaining, monitor)
- `hw-*` (19 commands: hardware detection, GPU checks, input devices)
- `launch-or-focus` (1 command: window management)

**Total:** 23 Go commands migrated out of 282

## Code Style

**Go conventions:**
- Two spaces for indentation
- Short, clear function names without package stuttering
- One-line package comments only (no function docs in internal packages)
- Minimal comments in cmd files (just metadata)
- No "Get" prefixes (`Level()` not `GetLevel()`)

**Notification API:**
```go
notify.Send(notify.Options{
  Glyph:       "🔋",
  Headline:    "Battery Low",
  Description: "15% remaining",
  Urgency:     notify.Critical,
  SoundFile:   "battery-low",
})
```

All fields are optional.

## Success Criteria

A command port is complete when:
- ✅ Output matches bash version exactly (or fixes a known bug)
- ✅ All edge cases tested and handled
- ✅ Benchmarks show significant speedup (>2x minimum)
- ✅ Final code review completed with no issues found
- ✅ Goroutines used where beneficial for performance
- ✅ Command works through dispatcher in test environment
