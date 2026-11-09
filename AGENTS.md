# Style

- Two spaces for indentation, no tabs
- Use bash 5 conditionals: use `[[ ]]` for string/file tests and `(( ))` for numeric tests
- In `[[ ]]`, don't quote variables, but do quote string literals when comparing values (e.g., `[[ $branch == "dev" ]]`)
- Prefer `(( ))` over numeric operators inside `[[ ]]` (e.g., `(( count < 50 ))`, not `[[ $count -lt 50 ]]`)
- For strings/paths with spaces, quote them instead of escaping spaces with `\ ` (e.g., `"$APP_DIR/Disk Usage.desktop"`, not `$APP_DIR/Disk\ Usage.desktop`)
- Shebangs must use `#!/bin/bash` consistently (never `#!/usr/bin/env bash`)
- Scripts under `install/` and `migrations/` may be sourced and intentionally omit shebangs

# Command Naming

All commands start with `hoopoo-`. Prefixes indicate purpose.

The authoritative command group list lives in `bin/hoopoo` in `GROUP_DESCRIPTIONS`. Keep `GROUP_DESCRIPTIONS` updated when adding a new command prefix.

Common prefixes include:

- `cmd-` - check if commands exist, misc utility commands
- `capture-` - screenshots, screen recordings, and other capture tools
- `pkg-` - package management helpers
- `hw-` - hardware detection (return exit codes for use in conditionals)
- `refresh-` - copy default config to user's `~/.config/`
- `restart-` - restart a component
- `launch-` - open applications
- `install-` - install optional software
- `setup-` - interactive setup wizards
- `toggle-` - toggle features on/off
- `theme-` - theme management
- `update-` - update components

Other current prefixes include:

- `ac-`, `audio-`, `battery-`, `branch-`, `brightness-`, `channel-`, `config-`, `debug-`, `dev-`, `drive-`, `first-`, `font-`, `haptic-`, `hibernation-`, `hook-`, `hyprland-`, `menu-`, `migrate-`, `notification-`, `npx-`, `plymouth-`, `powerprofiles-`, `reinstall-`, `remove-`, `screensaver-`, `show-`, `snapshot-`, `state-`, `sudo-`, `swayosd-`, `system-`, `transcode-`, `tui-`, `tz-`, `upload-`, `version-`, `voxtype-`, `webapp-`, `wifi-`, `windows-`

# Command Metadata

Commands in `bin/` can declare CLI metadata in comments near the top of the file. `bin/hoopoo` scans the first 80 lines, and tests expect command metadata to remain valid.

Supported metadata keys:

- `# hoopoo:summary=...` - short help text
- `# hoopoo:group=...` - command group when it differs from the filename-derived prefix
- `# hoopoo:name=...` - command name within the group
- `# hoopoo:args=...` - usage arguments
- `# hoopoo:examples=...` - examples separated with ` | `
- `# hoopoo:alias=...` / `# hoopoo:aliases=...` - alternate routes
- `# hoopoo:hidden=true` - hide from default command listings
- `# hoopoo:requires-sudo=true` - mark commands that require sudo

Prefer explicit metadata for user-facing commands. Keep routes consistent with the filename unless there is a deliberate alias or compatibility route.

Example:

```bash
# hoopoo:summary=Take a screenshot
# hoopoo:group=capture
# hoopoo:args=[smart|region|windows|fullscreen] [slurp|copy]
# hoopoo:examples=hoopoo screenshot | hoopoo capture screenshot region
# hoopoo:aliases=hoopoo screenshot
```

# Install Scripts

Install entry points (`install.sh`, `boot.sh`) use `#!/bin/bash`. Many scripts under `install/` are sourced via `run_logged` and intentionally do not have shebangs.

Install stage files follow this pattern:

- `install/*/all.sh` lists scripts in execution order
- leaf scripts are sourced by `run_logged $HOOPOO_INSTALL/path/to/script.sh`
- avoid `exit` in sourced install scripts unless intentionally aborting the install
- use `$HOOPOO_INSTALL` and `$HOOPOO_PATH` instead of hard-coded Hoopoo paths
- keep hardware-specific logic under `install/config/hardware/`
- prefer helper commands for package and command checks where available

Raw `command -v`, `pacman`, and `pacman-key` are acceptable in bootstrap/preflight/package-helper contexts where the helper commands may not be available yet or where direct package-manager behavior is the point of the script.

# Helper Commands

Use these instead of raw shell commands:

- `hoopoo-cmd-missing` / `hoopoo-cmd-present` - check for commands
- `hoopoo-pkg-missing` / `hoopoo-pkg-present` - check for packages
- `hoopoo-pkg-add` - install packages (handles both pacman and AUR)
- `hoopoo-hw-asus-rog` - detect ASUS ROG hardware (and similar `hw-*` commands)

Exceptions are allowed for bootstrap, preflight, migration, and package-helper scripts where the helper may not be available yet, where the helper itself is being implemented, or where direct package-manager behavior is required.

# Config Structure

- `config/` - default configs copied to `~/.config/`
- `default/themed/*.tpl` - templates with `{{ variable }}` placeholders for theme colors
- `themes/*/colors.toml` - theme color definitions (accent, background, foreground, color0-15)

# Visual Changes

When making visual changes, such as Waybar styles or desktop appearance, always take and analyze a screenshot after applying the change to verify the result. Use `hoopoo capture screenshot fullscreen save` for fullscreen screenshots.

# Refresh Pattern

To copy a default config to user config with automatic backup:

```bash
hoopoo-refresh-config hypr/hyprlock.conf
```

This copies `~/.local/share/hoopoo/config/hypr/hyprlock.conf` to `~/.config/hypr/hyprlock.conf`.

# Migrations

To create a new migration, run `hoopoo-dev-add-migration --no-edit`. This creates a migration file named after the unix timestamp of the last commit.

New migration format:
- File permissions must be `0644` (`-rw-r--r--`); migrations are sourced, not executed directly
- No shebang line
- Start with an `echo` describing what the migration does
- Use `$HOOPOO_PATH` to reference the hoopoo directory
- Prefer helper commands such as `hoopoo-cmd-present`, `hoopoo-cmd-missing`, `hoopoo-pkg-present`, and `hoopoo-pkg-missing`

Some older migrations predate these rules. Do not copy older migrations that start with shebangs, omit the leading `echo`, or hard-code `~/.local/share/hoopoo`.

Migrations may use raw `pacman`, `command -v`, or direct config edits when needed for historical compatibility or one-off repair work.

Example:
```bash
echo "Disable fingerprint in hyprlock if fingerprint auth is not configured"

if hoopoo-cmd-missing fprintd-list || ! fprintd-list "$USER" 2>/dev/null | grep -q "finger"; then
  sed -i 's/fingerprint:enabled = .*/fingerprint:enabled = false/' ~/.config/hypr/hyprlock.conf
fi
```

# Go Migration Workflow (Bash → Go)

When porting bash commands to Go, follow this mandatory workflow to ensure correctness and performance:

## 1. Analysis Phase

**For each command to port:**

1. Read the hoopoo bash implementation in `bin/hoopoo-<command>`
2. Read the omarchy bash implementation in `~/.local/share/omarchy-repos/omarchy/bin/omarchy-<command>` (if exists)
3. Compare approaches and identify:
   - Edge cases handled by omarchy but missing in hoopoo
   - Better algorithms or optimizations
   - Hardware compatibility issues (Apple Silicon, Chromebooks, Framework)
   - Multi-device scenarios (multiple batteries, monitors, audio sinks)
4. **Choose the better approach** - prefer omarchy's robustness if it handles more edge cases

## 2. Implementation Phase

**Design the Go version:**

1. Create shared library code in `internal/` packages for reusable logic
2. Keep command entry points in `cmd/hoopoo-<command>/main.go` minimal
3. **Use goroutines where beneficial:**
   - Parallel file reads (multiple sysfs files)
   - Concurrent API calls (D-Bus, IPC)
   - Multi-device queries (all batteries, all audio sinks)
   - Background tasks (monitoring, watching for changes)
4. Avoid goroutines for single operations (overhead > benefit)
5. Apply performance optimizations:
   - Cache expensive lookups (device discovery, file paths)
   - Minimize allocations (buffer reuse, manual parsing)
   - Direct system calls where possible (sysfs > upower, PipeWire > pulseaudio)

## 3. Verification Phase

**Test correctness before benchmarking:**

1. Build the Go binary: `go build -ldflags="-s -w" -o bin/<command>.go ./cmd/<command>`
2. Run both bash and Go versions side-by-side with same inputs
3. Compare outputs - they must match exactly (unless fixing a bug)
4. Test edge cases:
   - No hardware present (no battery, no audio)
   - Multiple devices (BAT0 + BAT1, multiple monitors)
   - Missing sysfs files (older kernels, different hardware)
   - Invalid inputs (empty strings, out of range values)
5. **Only proceed to benchmarks when output is identical**

## 4. Benchmark Phase

**Measure performance improvements:**

1. Use `hyperfine --warmup 3 --runs 50` for accurate measurements
2. Compare bash vs Go versions
3. Document speedup ratio and memory usage
4. Verify the Go version is actually faster (no regressions)

## 5. Final Review Phase

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

## 6. Integration Phase

**After all checks pass:**

1. Keep both versions temporarily (`hoopoo-command` bash + `hoopoo-command.go` binary)
2. Update `Makefile` to build the new Go command
3. Document the migration in comments
4. Once stable, remove bash version and rename `.go` binary to final name

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

## Success Criteria

A command port is complete when:
- ✅ Output matches bash version exactly (or fixes a known bug)
- ✅ All edge cases tested and handled
- ✅ Benchmarks show significant speedup (>2x minimum)
- ✅ Final code review completed with no issues found
- ✅ Goroutines used where beneficial for performance
