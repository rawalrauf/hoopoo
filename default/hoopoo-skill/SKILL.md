---
name: hoopoo
description: >
  REQUIRED for end-user customization of Linux desktop, window manager, or system config.
  Use when editing ~/.config/niri/, ~/.config/waybar/, ~/.config/walker/,
  ~/.config/alacritty/, ~/.config/foot/, ~/.config/kitty/, ~/.config/ghostty/, ~/.config/mako/,
  or ~/.config/hoopoo/. Triggers: niri, window rules, keybindings, monitors, gaps, borders,
  opacity, waybar, walker, terminal config, themes, wallpaper, night light, idle, lock screen,
  screenshots, layer rules, workspace settings, display config, and user-facing hoopoo commands.
  Excludes Hoopoo source development in ~/.local/share/hoopoo/ and `hoopoo dev` workflows.
---

# Hoopoo Skill

Manage a custom niri-based Arch Linux system built on top of [Hoopoo](https://hoopoo.dev/).

This is a fork of DHH's hoopoo — migrated from Hyprland to **niri** as the Wayland compositor.
This skill is for end-user customization on installed systems only.
It is NOT for contributing to Hoopoo source code.

## When This Skill MUST Be Used

**ALWAYS invoke this skill for end-user requests involving ANY of these:**

- Editing ANY file in `~/.config/niri/` (window rules, keybindings, monitors, looknfeel, etc.)
- Editing ANY file in `~/.config/waybar/`, `~/.config/walker/`, `~/.config/mako/`
- Editing terminal configs (alacritty, foot, kitty, ghostty)
- Editing ANY file in `~/.config/hoopoo/`
- Window behavior, opacity, gaps, borders, layer rules, workspace settings
- Themes, wallpapers, fonts, appearance changes
- Display/monitor configuration, mirror mode
- Screenshots, screen recording, reminders, night light, idle behavior, lock screen
- User-facing `hoopoo` commands

**If you're about to edit a config file in ~/.config/ on this system, STOP and use this skill first.**

**Do NOT use this skill for Hoopoo source development** (editing files in `~/.local/share/hoopoo/`).

## Critical Safety Rules

**NEVER modify anything in `~/.local/share/hoopoo/`** — reading is safe and encouraged.

This directory is managed by git. Changes will be lost on `hoopoo update` and cause conflicts.

```
~/.local/share/hoopoo/     # READ-ONLY - NEVER EDIT (reading is OK)
├── bin/                    # Source scripts (symlinked to PATH)
├── config/                 # Default config templates
├── default/niri/           # Niri defaults (sourced by user config via include)
├── themes/                 # Stock themes
├── migrations/             # Update migrations
└── install/                # Installation scripts
```

**Always use these safe locations instead:**
- `~/.config/` — user configuration (safe to edit)
- `~/.config/hoopoo/themes/<custom-name>/` — custom themes
- `~/.config/hoopoo/hooks/` — custom automation hooks

## System Architecture

| Component | Purpose | Config Location |
|-----------|---------|-----------------|
| **Arch Linux** | Base OS | `/etc/`, `~/.config/` |
| **niri** | Wayland compositor/WM | `~/.config/niri/` |
| **Waybar** | Status bar | `~/.config/waybar/` |
| **Walker** | App launcher | `~/.config/walker/` |
| **Alacritty/Foot/Kitty/Ghostty** | Terminals | `~/.config/<terminal>/` |
| **Mako** | Notifications | `~/.config/mako/` |
| **SwayOSD** | On-screen display | `~/.config/swayosd/` |
| **swayidle** | Idle management | `~/.config/swayidle/` |
| **hyprlock** | Lock screen | `~/.config/hypr/hyprlock.conf` |

## Niri Config Architecture

Niri config uses KDL format. The user's `~/.config/niri/config.kdl` is the entry point — it sources everything via `include` directives:

```
~/.config/niri/
├── config.kdl          # Entry point — includes everything below
├── autostart.kdl       # Startup applications
├── monitors.kdl        # Display/monitor configuration
├── input.kdl           # Keyboard, mouse, touchpad settings
├── bindings.kdl        # Keybindings (sourced from default + user overrides)
├── looknfeel.kdl       # Appearance (gaps, borders, animations, opacity)
├── flags.kdl           # Optional toggle includes (gaps, transparency, mirror, etc.)
└── system.kdl          # Window rules (floating, fullscreen, etc.)
```

**Key architecture detail:** Hoopoo defaults live in `~/.local/share/hoopoo/default/niri/` and are included by `config.kdl` via:
```kdl
include "~/.local/share/hoopoo/default/niri/default.kdl"
```

User files in `~/.config/niri/` override/extend the defaults. Edit those — not the defaults.

### Niri Reload

Niri does NOT auto-reload on save. Always reload after config changes:

```bash
hoopoo restart niri         # or: niri msg action load-config-file
```

Always validate config before reloading:
```bash
niri validate --config ~/.config/niri/config.kdl
```

If validation fails, fix errors before reloading.

### Niri Toggles

Toggles are managed via flag files in `~/.local/state/hoopoo/toggles/` and KDL snippets in `~/.local/share/hoopoo/default/niri/toggles/`. The `flags.kdl` optionally includes them.

Current toggles and their binds:
- `Mod+Backspace` — window transparency
- `Mod+Alt+Backspace` — window gaps
- `Mod+Alt+Delete` — laptop display toggle
- `Mod+Alt+Shift+Delete` — mirror display
- `Mod+Alt+W` — top bar (waybar)
- `Mod+Alt+Shift+W` — top bar transparency
- `Mod+Alt+I` — idle lock
- `Mod+Alt+Period` — toggle menu

## Keybinding Conventions

```
Mod = Super key

Mod+Space               — app launcher (walker)
Mod+Alt+Space           — hoopoo menu
Mod+Alt+<letter>        — utilities, toggles, menus
Mod+Shift+<letter>      — launch apps
Mod+Shift+Ctrl+<letter> — TUI controls
Mod+Ctrl+<letter>       — tiling controls
```

**When adding keybindings**, edit `~/.config/niri/bindings.kdl`. Check for conflicts first:
```bash
grep -r "Mod" ~/.config/niri/bindings.kdl ~/.local/share/hoopoo/default/niri/bindings/
```

Niri keybind format:
```kdl
binds {
  Mod+E { spawn "nautilus"; }
  Mod+Shift+Q { close-window; }
}
```

## Waybar

```
~/.config/waybar/
├── config.jsonc        # Bar layout and modules
└── style.css           # Styling (transparent.css dynamically imported by toggle)
```

Waybar does NOT auto-reload. After any config change:
```bash
hoopoo restart waybar
```

The waybar transparency toggle dynamically adds/removes a `transparent.css` import and sends `SIGUSR2` for CSS reload — no full restart needed for that specific toggle.

## Command Discovery

```bash
hoopoo commands                  # List all commands
hoopoo <group> --help            # Commands in a group
hoopoo commands --json           # Machine-readable listing
cat $(which hoopoo-theme-set)    # Read a command's source
```

### Command Groups

| Group | Purpose | Example |
|-------|---------|---------|
| `hoopoo refresh` | Reset config to defaults (backs up first) | `hoopoo refresh niri` |
| `hoopoo restart` | Restart a service/app | `hoopoo restart waybar` |
| `hoopoo toggle` | Toggle feature on/off | `hoopoo toggle nightlight` |
| `hoopoo theme` | Theme management | `hoopoo theme set <name>` |
| `hoopoo install` | Install optional software | `hoopoo install docker` |
| `hoopoo launch` | Launch apps | `hoopoo launch browser` |
| `hoopoo capture` | Screenshots and recordings | `hoopoo capture screenshot` |
| `hoopoo reminder` | Desktop notification reminders | `hoopoo reminder 15 "Call back"` |
| `hoopoo pkg` | Package management | `hoopoo pkg add <pkg>` |
| `hoopoo setup` | Setup wizards | `hoopoo setup fingerprint` |
| `hoopoo update` | System updates | `hoopoo update` |
| `hoopoo niri` | Niri-specific controls | `hoopoo niri monitor mirror` |

## Configuration Locations

### Niri (Window Manager)

Edit files in `~/.config/niri/` — never in `~/.local/share/hoopoo/default/niri/`.

Key files:
- `monitors.kdl` — display configuration
- `input.kdl` — keyboard/touchpad/mouse
- `bindings.kdl` — keybindings
- `looknfeel.kdl` — gaps, borders, animations, opacity
- `autostart.kdl` — startup apps
- `window-rules.kdl` — window rules

### Other Configs

| App | Location |
|-----|----------|
| btop | `~/.config/btop/btop.conf` |
| fastfetch | `~/.config/fastfetch/config.jsonc` |
| lazygit | `~/.config/lazygit/config.yml` |
| starship | `~/.config/starship.toml` |
| git | `~/.config/git/config` |
| walker | `~/.config/walker/config.toml` |
| hyprlock | `~/.config/hypr/hyprlock.conf` |

## Safe Customization Patterns

### Pattern 1: Edit User Config Directly

```bash
# 1. Read current config
cat ~/.config/niri/bindings.kdl

# 2. Make changes with Edit tool

# 3. Validate
niri validate --config ~/.config/niri/config.kdl

# 4. Reload
hoopoo restart niri
```

### Pattern 2: Make a New Theme

1. Create a directory under `~/.config/hoopoo/themes/`
2. See how an existing theme is done: `ls ~/.local/share/hoopoo/themes/catppuccin/`
3. Add backgrounds and a `colors.toml` matching the theme format
4. Apply: `hoopoo theme set "Name of new theme"`

### Pattern 3: Use Hooks for Automation

```bash
~/.config/hoopoo/hooks/
├── theme-set        # Runs after theme change (receives theme name as $1)
├── font-set         # Runs after font change
└── post-update      # Runs after hoopoo update
```

### Pattern 4: Reset to Defaults — ALWAYS SEEK USER CONFIRMATION FIRST

```bash
hoopoo refresh niri          # Reset all niri configs (backs up first)
hoopoo refresh waybar        # Reset waybar config
hoopoo refresh walker        # Reset walker config
```

## Common Tasks

### Themes

```bash
hoopoo theme list              # Show available themes
hoopoo theme current           # Show current theme
hoopoo theme set <name>        # Apply theme (use "Tokyo Night" not "tokyo-night")
hoopoo theme bg next           # Cycle wallpaper
hoopoo theme install <url>     # Install from git repo
```

### Display/Monitors

Edit `~/.config/niri/monitors.kdl`. List monitors:
```bash
niri msg outputs
```

Niri monitor format:
```kdl
output "eDP-1" {
  mode "1920x1200@60"
  scale 1.5
  position x=0 y=0
}
```

Mirror mode: `Mod+Alt+Shift+Delete` or `hoopoo niri monitor mirror`

### Window Rules

Window rules go in `~/.config/niri/system.kdl`. Niri KDL format:
```kdl
window-rule {
  match app-id="org.gnome.Calculator"
  open-floating true
  default-floating-size { width 400; height 500; }
}
```

Always validate after adding rules:
```bash
niri validate --config ~/.config/niri/config.kdl
```

### Fonts

```bash
hoopoo font list
hoopoo font current
hoopoo font set <name>
```

### System

```bash
hoopoo update                  # Full system update
hoopoo version                 # Show version
hoopoo debug --no-sudo --print # Debug info (ALWAYS use these flags)
hoopoo system lock             # Lock screen
hoopoo system shutdown         # Shutdown
hoopoo system reboot           # Reboot
```

### Reminders

```bash
hoopoo reminder 15 "Pickup Jack"
hoopoo reminder 60 "Check laundry"
hoopoo reminder show
hoopoo reminder clear
```

## Troubleshooting

```bash
hoopoo debug --no-sudo --print   # Debug info (always use these flags)
hoopoo upload log                # Upload logs for support
hoopoo refresh niri              # Reset niri configs to defaults
hoopoo refresh <app>             # Reset any app config to defaults
hoopoo reinstall                 # Nuclear option — full config reinstall
```

## Decision Framework

1. **Is it a stock hoopoo command?** Use it directly
2. **Is it a config edit?** Edit in `~/.config/`, never `~/.local/share/hoopoo/`
3. **Is it a theme customization?** Create a NEW custom theme directory
4. **Is it automation?** Use hooks in `~/.config/hoopoo/hooks/`
5. **Is it a package install?** Use `hoopoo pkg add <pkgs...>` (or `hoopoo pkg aur add <pkgs...>` for AUR)
6. **Unsure if command exists?** Run `hoopoo commands` first

## Out of Scope

This skill does NOT cover Hoopoo source development. Do not use it for:
- Editing files in `~/.local/share/hoopoo/`
- Creating or editing migrations
- Running `hoopoo dev ...` commands

## Example Requests

- "Change my theme to catppuccin" → `hoopoo theme set catppuccin`
- "Add a keybinding for Super+E to open file manager" → check conflicts, add to `~/.config/niri/bindings.kdl`, validate, reload
- "Configure my external monitor" → edit `~/.config/niri/monitors.kdl`, validate, reload
- "Make window gaps smaller" → edit `~/.config/niri/looknfeel.kdl`, validate, reload
- "Set up night light" → `hoopoo toggle nightlight`
- "Set a reminder in 15 minutes" → `hoopoo reminder 15 "message"`
- "Mirror my display" → `Mod+Alt+Shift+Delete` or `hoopoo niri monitor mirror`
- "Reset waybar to defaults" → `hoopoo refresh waybar`
- "Run a script every time I change themes" → create `~/.config/hoopoo/hooks/theme-set`
