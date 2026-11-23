package main

import (
  "os"
  "os/exec"
  "strconv"
  "strings"

  "hoopoo/pkg/hardware"
  "hoopoo/pkg/menu"
)

// MainMenu returns the main Hoopoo menu items (10 items - matching bash)
func MainMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰀻", Label: "Apps", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "walker"}}},
    {Icon: "󰧑", Label: "Learn", Action: &menu.SubmenuAction{MenuFunc: LearnMenu}},
    {Icon: "󱓞", Label: "Trigger", Action: &menu.SubmenuAction{MenuFunc: TriggerMenu}},
    {Icon: "", Label: "Appearance", Action: &menu.SubmenuAction{MenuFunc: AppearanceMenu}},
    {Icon: "", Label: "Setup", Action: &menu.SubmenuAction{MenuFunc: SetupMenu}},
    {Icon: "󰉉", Label: "Install", Action: &menu.SubmenuAction{MenuFunc: InstallMenu}},
    {Icon: "󰭌", Label: "Remove", Action: &menu.SubmenuAction{MenuFunc: RemoveMenu}},
    {Icon: "", Label: "Update", Action: &menu.SubmenuAction{MenuFunc: UpdateMenu}},
    {Icon: "", Label: "About", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "about"}}},
    {Icon: "", Label: "System", Action: &menu.SubmenuAction{MenuFunc: SystemMenu}},
  }
}

// LearnMenu - documentation and resources
func LearnMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "Keybindings", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"menu", "keybindings"}}},
    {Icon: "", Label: "Hoopoo", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "webapp", "https://github.com/rawalrauf/hoopoo"}}},
    {Icon: "", Label: "Niri", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "webapp", "https://niri-wm.github.io/niri/index.html"}}},
    {Icon: "󰣇", Label: "Arch", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "webapp", "https://wiki.archlinux.org/title/Main_page"}}},
    {Icon: "", Label: "Neovim", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "webapp", "https://www.lazyvim.org/keymaps"}}},
    {Icon: "󱆃", Label: "Bash", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "webapp", "https://devhints.io/bash"}}},
  }
}

// TriggerMenu - actions and utilities
func TriggerMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰔛", Label: "Reminder", Action: &menu.SubmenuAction{MenuFunc: ReminderMenu}},
    {Icon: "", Label: "Capture", Action: &menu.SubmenuAction{MenuFunc: CaptureMenu}},
    {Icon: "󰧸", Label: "Transcode", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"transcode"}}},
    {Icon: "", Label: "Share", Action: &menu.SubmenuAction{MenuFunc: ShareMenu}},
    {Icon: "", Label: "Toggle", Action: &menu.SubmenuAction{MenuFunc: TogglesMenu}},
    {Icon: "", Label: "Hardware", Action: &menu.SubmenuAction{MenuFunc: HardwareMenu}},
  }
}

// ReminderMenu - set reminders
func ReminderMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰔛", Label: "Set one", Action: &menu.FuncAction{Fn: setReminder}},
    {Icon: "󰔛", Label: "Show all", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"reminder", "show"}}},
    {Icon: "󰔛", Label: "Clear all", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"reminder", "clear"}}},
  }
}

// setReminder handles the multi-step reminder input process
func setReminder() error {
  // Step 1: Get minutes
  minutesCmd := exec.Command("hoopoo", "menu", "input", "Remind in minutes")
  minutesOutput, err := minutesCmd.Output()
  if err != nil {
    return nil // User cancelled
  }

  minutes := strings.TrimSpace(string(minutesOutput))
  if minutes == "" {
    return nil // User cancelled
  }

  // Validate minutes is a number
  if _, err := strconv.Atoi(minutes); err != nil {
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰔛",
      "Invalid reminder",
      "Enter the number of minutes",
      "-u", "critical",
      "-t", "5000",
    ).Run()
    return setReminder() // Retry
  }

  // Step 2: Get message
  messageCmd := exec.Command("hoopoo", "menu", "input", "Reminder message")
  messageOutput, err := messageCmd.Output()
  if err != nil {
    return nil // User cancelled
  }

  message := strings.TrimSpace(string(messageOutput))

  // Step 3: Set reminder
  if message != "" {
    return exec.Command("hoopoo", "reminder", minutes, message).Run()
  } else {
    return exec.Command("hoopoo", "reminder", minutes).Run()
  }
}

// CaptureMenu - screenshots and recordings
func CaptureMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "Screenshot", Action: &menu.SubmenuAction{MenuFunc: ScreenshotMenu}},
    {Icon: "", Label: "Screenrecord", Action: &menu.SubmenuAction{MenuFunc: ScreenrecordMenu}},
    {Icon: "󰴑", Label: "Text Extraction", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "text", "extraction"}}},
    {Icon: "󰃉", Label: "Color", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "color"}}},
  }
}

// ScreenshotMenu - screenshot options
func ScreenshotMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰩭", Label: "Select Region", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenshot"}}},
    {Icon: "", Label: "Full Screen", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenshot", "fullscreen"}}},
    {Icon: "", Label: "Active Window", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenshot", "window"}}},
  }
}

// ScreenrecordMenu - screenrecording options
func ScreenrecordMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "With no audio", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenrecording"}}},
    {Icon: "", Label: "With desktop audio", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenrecording", "--with-desktop-audio"}}},
    {Icon: "", Label: "With desktop + microphone audio", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenrecording", "--with-desktop-audio", "--with-microphone-audio"}}},
    {Icon: "", Label: "With desktop + microphone + webcam", Action: &menu.SubmenuAction{MenuFunc: WebcamMenu}},
  }
}

// WebcamMenu - webcam device selection (or auto-start if only one webcam)
func WebcamMenu() []menu.MenuItem {
  devices, err := hardware.ListWebcams()

  // If no webcams or error, use default
  if err != nil || len(devices) == 0 {
    return []menu.MenuItem{
      {Icon: "󰖠", Label: "Default webcam", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"capture", "screenrecording", "--with-desktop-audio", "--with-microphone-audio", "--with-webcam"}}},
    }
  }

  // If exactly one webcam, auto-start recording with it (no menu)
  if len(devices) == 1 {
    exec.Command("hoopoo", "capture", "screenrecording",
      "--with-desktop-audio",
      "--with-microphone-audio",
      "--with-webcam",
      "--webcam-device="+devices[0].Path).Start()
    os.Exit(0)
  }

  // Multiple webcams - build menu for selection
  var items []menu.MenuItem
  for _, device := range devices {
    devicePath := device.Path
    deviceName := device.Name

    // Clean up device name for display
    if deviceName == "" {
      deviceName = devicePath
    }

    items = append(items, menu.MenuItem{
      Icon:   "󰖠",
      Label:  deviceName,
      Action: &menu.CommandAction{
        Command: "hoopoo",
        Args:    []string{"capture", "screenrecording", "--with-desktop-audio", "--with-microphone-audio", "--with-webcam", "--webcam-device=" + devicePath},
      },
    })
  }

  return items
}

// ShareMenu - file sharing
func ShareMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "Clipboard", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"share", "clipboard"}}},
    {Icon: "", Label: "File", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"share", "file"}}},
    {Icon: "", Label: "Folder", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"share", "folder"}}},
  }
}

// TogglesMenu - system toggles (cleaned - visual toggles moved to Appearance)
func TogglesMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󱫖", Label: menu.Checkmark("Idle Lock", menu.IsIdleLockEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "idle"}}},
    {Icon: "󰂛", Label: menu.Checkmark("Notifications", menu.IsNotificationSilencingEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "notification", "silencing"}}},
    {Icon: "", Label: menu.Checkmark("Direct Boot", menu.IsDirectBootEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "floating", "terminal", "with", "presentation", "hoopoo config direct boot"}}},
    {Icon: "󰟵", Label: menu.Checkmark("Passwordless Sudo", menu.IsPasswordlessSudoEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "floating", "terminal", "with", "presentation", "hoopoo sudo passwordless"}}},
  }
}

// HardwareMenu - hardware toggles (conditional based on detected hardware)
func HardwareMenu() []menu.MenuItem {
  items := []menu.MenuItem{
    {Icon: "󰛧", Label: "Laptop Display", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"niri", "monitor", "internal", "toggle"}}},
    {Icon: "󰍹", Label: "Mirror Display", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"niri", "monitor", "internal", "mirror", "toggle"}}},
  }

  // Hybrid GPU (Asus ROG, some gaming laptops)
  if exec.Command("hoopoo", "hw", "hybrid", "gpu").Run() == nil {
    items = append(items, menu.MenuItem{
      Icon: "",
      Label: "Hybrid GPU",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "floating", "terminal", "with", "presentation", "hoopoo toggle hybrid gpu"}},
    })
  }

  // Touchpad
  if exec.Command("hoopoo", "hw", "touchpad").Run() == nil {
    items = append(items, menu.MenuItem{
      Icon: "󰟸",
      Label: "Touchpad",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "touchpad"}},
    })
  }

  // Touchscreen
  if exec.Command("hoopoo", "hw", "touchscreen").Run() == nil {
    items = append(items, menu.MenuItem{
      Icon: "󰆽",
      Label: "Touchscreen",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "touchscreen"}},
    })
  }

  return items
}

// AppearanceMenu - visual customization (organized by visual impact)
func AppearanceMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰸌", Label: "Theme", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "walker", "-m", "menus:hoopoothemes", "--width", "800", "--minheight", "400"}}},
    {Icon: "󰟵", Label: "Unlock Screen", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "walker", "-m", "menus:hoopoounlocks", "--width", "800", "--minheight", "400"}}},
    {Icon: "", Label: "Background", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "walker", "-m", "menus:hoopooBackgroundSelector", "--width", "800", "--minheight", "400"}}},
    {Icon: "", Label: "Font", Action: &menu.SubmenuAction{MenuFunc: FontMenu}},
    {Icon: "󰖯", Label: "Windows", Action: &menu.SubmenuAction{MenuFunc: WindowsMenu}},
    {Icon: "󰍜", Label: "Top Bar", Action: &menu.SubmenuAction{MenuFunc: TopBarMenu}},
  }
}

// WindowsMenu - window appearance settings
func WindowsMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰘇", Label: menu.Checkmark("Round Corners", menu.IsRoundCornersEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"niri", "window", "roundcorners", "toggle"}}},
    {Icon: "", Label: menu.Checkmark("Gaps", menu.IsGapsEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"niri", "window", "gaps", "toggle"}}},
    {Icon: "󰂵", Label: menu.Checkmark("Transparency", menu.IsTransparencyEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"niri", "window", "transparency", "toggle"}}},
  }
}

// TopBarMenu - waybar appearance and position settings
func TopBarMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰍜", Label: menu.Checkmark("Visibility", menu.IsWaybarEnabled()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "waybar"}}},
    {Icon: "󰂵", Label: menu.Checkmark("Transparency", menu.IsWaybarTransparent()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "waybar", "transparency"}}},
    {Icon: "󰡏", Label: "Position", Action: &menu.SubmenuAction{MenuFunc: WaybarPositionMenu}},
  }
}

// WaybarPositionMenu - waybar position selection (top or bottom)
func WaybarPositionMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰘣", Label: menu.Checkmark("Top", menu.IsWaybarPositionTop()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "waybar", "position"}}},
    {Icon: "󰘡", Label: menu.Checkmark("Bottom", menu.IsWaybarPositionBottom()), Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "waybar", "position"}}},
  }
}

// FontMenu - font selection
func FontMenu() []menu.MenuItem {
  currentFont := menu.GetCurrentFont()

  // Get available fonts
  cmd := exec.Command("hoopoo", "font", "list")
  output, err := cmd.Output()
  if err != nil {
    return []menu.MenuItem{
      {Icon: "󰀦", Label: "Error loading fonts", Action: &menu.CommandAction{Command: "echo", Args: []string{"Failed to load fonts"}}},
    }
  }

  fonts := strings.Split(strings.TrimSpace(string(output)), "\n")
  var items []menu.MenuItem

  for _, font := range fonts {
    font = strings.TrimSpace(font)
    if font == "" {
      continue
    }

    label := menu.Checkmark(font, font == currentFont)
    items = append(items, menu.MenuItem{
      Icon:   "",
      Label:  label,
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"font", "set", font}},
    })
  }

  return items
}

// SetupMenu - system configuration (now includes monitor scaling)
func SetupMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "Audio", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "audio"}}},
    {Icon: "", Label: "Wifi", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "wifi"}}},
    {Icon: "󰂯", Label: "Bluetooth", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "bluetooth"}}},
    {Icon: "󱐋", Label: "Power Profile", Action: &menu.SubmenuAction{MenuFunc: PowerProfileMenu}},
    {Icon: "", Label: "System Sleep", Action: &menu.SubmenuAction{MenuFunc: SystemSleepMenu}},
    {Icon: "󰍹", Label: "Monitors", Action: &menu.SubmenuAction{MenuFunc: MonitorsMenu}},
    {Icon: "", Label: "Defaults", Action: &menu.SubmenuAction{MenuFunc: DefaultsMenu}},
    {Icon: "󰱔", Label: "DNS", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"setup", "dns"}}},
    {Icon: "", Label: "Security", Action: &menu.SubmenuAction{MenuFunc: SecurityMenu}},
    {Icon: "", Label: "Config", Action: &menu.CommandAction{Command: "echo", Args: []string{"TODO: Config menu"}}},
  }
}

// MonitorsMenu - monitor configuration (now includes scaling from old Toggle menu)
func MonitorsMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰍹", Label: "Scale", Action: &menu.SubmenuAction{MenuFunc: MonitorScaleMenu}},
    {Icon: "", Label: "Edit Config", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "editor", "~/.config/niri/monitors.kdl"}}},
  }
}

// MonitorScaleMenu - smart routing: show monitor selection if multiple monitors, else show scale values directly
func MonitorScaleMenu() []menu.MenuItem {
  // Check if external monitors are connected
  if !menu.HasExternalMonitors() {
    // Single monitor: go directly to scale values
    return MonitorScaleValuesMenu("")
  }
  
  // Multiple monitors: show monitor selection first
  return MonitorSelectionMenu("scale")
}

// MonitorSelectionMenu - select which monitor to configure (for scale/resolution/transform)
func MonitorSelectionMenu(action string) []menu.MenuItem {
  outputs, err := menu.GetAllOutputNames()
  if err != nil || len(outputs) == 0 {
    // Fallback to direct action if can't get outputs
    switch action {
    case "scale":
      return MonitorScaleValuesMenu("")
    default:
      return []menu.MenuItem{}
    }
  }
  
  var items []menu.MenuItem
  for _, output := range outputs {
    // Create a closure to capture the output name
    outputName := output
    var submenuFunc func() []menu.MenuItem
    
    switch action {
    case "scale":
      submenuFunc = func() []menu.MenuItem {
        return MonitorScaleValuesMenu(outputName)
      }
    }
    
    items = append(items, menu.MenuItem{
      Icon:   "󰍹",
      Label:  outputName,
      Action: &menu.SubmenuAction{MenuFunc: submenuFunc},
    })
  }
  
  return items
}

// MonitorScaleValuesMenu - show scale values for a specific monitor (or default if empty)
func MonitorScaleValuesMenu(output string) []menu.MenuItem {
  current := menu.GetCurrentMonitorScale(output)
  
  // Build base command args - add --output flag if monitor specified
  var baseArgs []string
  if output != "" {
    baseArgs = []string{"niri", "monitor", "scaling", "cycle", "--output", output, "--set"}
  } else {
    baseArgs = []string{"niri", "monitor", "scaling", "cycle", "--set"}
  }
  
  return []menu.MenuItem{
    {Icon: "", Label: menu.Checkmark("1.0", current == "1"), Action: &menu.CommandAction{Command: "hoopoo", Args: append(baseArgs, "1")}},
    {Icon: "", Label: menu.Checkmark("1.25", current == "1.25"), Action: &menu.CommandAction{Command: "hoopoo", Args: append(baseArgs, "1.25")}},
    {Icon: "", Label: menu.Checkmark("1.5", current == "1.5"), Action: &menu.CommandAction{Command: "hoopoo", Args: append(baseArgs, "1.5")}},
    {Icon: "", Label: menu.Checkmark("2.0", current == "2"), Action: &menu.CommandAction{Command: "hoopoo", Args: append(baseArgs, "2")}},
  }
}

// PowerProfileMenu - power profile selection
func PowerProfileMenu() []menu.MenuItem {
  currentProfile := menu.GetCurrentPowerProfile()
  
  return []menu.MenuItem{
    {Icon: "󰓅", Label: menu.Checkmark("Performance", currentProfile == "performance"), Action: &menu.CommandAction{Command: "powerprofilesctl", Args: []string{"set", "performance"}}},
    {Icon: "󰾅", Label: menu.Checkmark("Balanced", currentProfile == "balanced"), Action: &menu.CommandAction{Command: "powerprofilesctl", Args: []string{"set", "balanced"}}},
    {Icon: "󰾆", Label: menu.Checkmark("Power Saver", currentProfile == "power-saver"), Action: &menu.CommandAction{Command: "powerprofilesctl", Args: []string{"set", "power-saver"}}},
  }
}

// SystemSleepMenu - suspend and hibernation configuration
func SystemSleepMenu() []menu.MenuItem {
  items := []menu.MenuItem{}
  
  // Suspend toggle
  if menu.IsSuspendDisabled() {
    items = append(items, menu.MenuItem{
      Icon:   "󰒲",
      Label:  "Enable Suspend",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "suspend"}},
    })
  } else {
    items = append(items, menu.MenuItem{
      Icon:   "󰒲",
      Label:  "Disable Suspend",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"toggle", "suspend"}},
    })
  }
  
  // Hibernation setup/remove
  if menu.IsHibernationAvailable() {
    items = append(items, menu.MenuItem{
      Icon:   "󰤁",
      Label:  "Disable Hibernate",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "floating", "terminal", "with", "presentation", "hoopoo hibernation remove"}},
    })
  } else {
    items = append(items, menu.MenuItem{
      Icon:   "󰤁",
      Label:  "Enable Hibernate",
      Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"launch", "floating", "terminal", "with", "presentation", "hoopoo hibernation setup"}},
    })
  }
  
  return items
}

// SecurityMenu - security settings
func SecurityMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰈷", Label: "Fingerprint", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"setup", "security", "fingerprint"}}},
    {Icon: "", Label: "Fido2", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"setup", "security", "fido2"}}},
  }
}

// DefaultsMenu - default applications
func DefaultsMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "Browser", Action: &menu.CommandAction{Command: "echo", Args: []string{"TODO: Browser selection"}}},
    {Icon: "", Label: "Terminal", Action: &menu.CommandAction{Command: "echo", Args: []string{"TODO: Terminal selection"}}},
    {Icon: "", Label: "Editor", Action: &menu.CommandAction{Command: "echo", Args: []string{"TODO: Editor selection"}}},
  }
}

// InstallMenu - placeholder
func InstallMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰏔", Label: "TODO: Install menu", Action: &menu.CommandAction{Command: "echo", Args: []string{"Not implemented yet"}}},
  }
}

// RemoveMenu - placeholder
func RemoveMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "󰆴", Label: "TODO: Remove menu", Action: &menu.CommandAction{Command: "echo", Args: []string{"Not implemented yet"}}},
  }
}

// UpdateMenu - placeholder
func UpdateMenu() []menu.MenuItem {
  return []menu.MenuItem{
    {Icon: "", Label: "TODO: Update menu", Action: &menu.CommandAction{Command: "echo", Args: []string{"Not implemented yet"}}},
  }
}

// AboutMenu - placeholder
// AboutMenu - about Hoopoo
// SystemMenu - placeholder
func SystemMenu() []menu.MenuItem {
  items := []menu.MenuItem{
    {Icon: "", Label: "Lock", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"system", "lock"}}},
  }

  // Add Suspend only if not disabled
  if !menu.IsSuspendDisabled() {
    items = append(items, menu.MenuItem{Icon: "󰒲", Label: "Suspend", Action: &menu.CommandAction{Command: "systemctl", Args: []string{"suspend"}}})
  }

  // Add Hibernate if available
  if menu.IsHibernationAvailable() {
    items = append(items, menu.MenuItem{Icon: "󰤁", Label: "Hibernate", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"system", "hibernate"}}})
  }

  items = append(items,
    menu.MenuItem{Icon: "󰍃", Label: "Logout", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"system", "logout"}}},
    menu.MenuItem{Icon: "󰜉", Label: "Restart", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"system", "reboot"}}},
    menu.MenuItem{Icon: "󰐥", Label: "Shutdown", Action: &menu.CommandAction{Command: "hoopoo", Args: []string{"system", "shutdown"}}},
  )

  return items
}
