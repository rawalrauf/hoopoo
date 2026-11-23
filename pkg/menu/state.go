// State detection helpers for menu checkmarks

package menu

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "hoopoo/pkg/font"
  "hoopoo/pkg/hardware"
  "hoopoo/pkg/niri"
  "hoopoo/pkg/system/process"
  "hoopoo/pkg/toggle"
)

// Checkmark appends ✓ to label if condition is true
func Checkmark(label string, enabled bool) string {
  if enabled {
    return label + " ✓"
  }
  return label
}

// IsIdleLockEnabled checks if swayidle is running
func IsIdleLockEnabled() bool {
  return process.IsRunning("swayidle")
}

// IsNotificationSilencingEnabled checks if do-not-disturb is OFF (notifications enabled)
func IsNotificationSilencingEnabled() bool {
  cmd := exec.Command("makoctl", "mode")
  output, err := cmd.Output()
  if err != nil {
    return true // Default: notifications enabled
  }
  return !strings.Contains(string(output), "do-not-disturb")
}

// IsWaybarEnabled checks if waybar is running
func IsWaybarEnabled() bool {
  return process.IsRunning("waybar")
}

// IsWaybarTransparent checks waybar transparency toggle state
func IsWaybarTransparent() bool {
  return toggle.Enabled("waybar-transparent")
}

// IsWaybarPositionTop checks if waybar is at top (default if no toggle file)
func IsWaybarPositionTop() bool {
  return toggle.Disabled("waybar-position-bottom")
}

// IsWaybarPositionBottom checks if waybar is at bottom
func IsWaybarPositionBottom() bool {
  return toggle.Enabled("waybar-position-bottom")
}

// IsGapsEnabled checks niri gaps toggle
func IsGapsEnabled() bool {
  return niri.StateEnabled("gaps-enabled.kdl")
}

// IsTransparencyEnabled checks niri transparency toggle (inverted)
func IsTransparencyEnabled() bool {
  return !niri.StateEnabled("transparency-disabled.kdl")
}

// IsRoundCornersEnabled checks niri round corners toggle
func IsRoundCornersEnabled() bool {
  return niri.StateEnabled("roundcorners-enabled.kdl")
}

// IsDirectBootEnabled checks direct boot toggle
func IsDirectBootEnabled() bool {
  return toggle.Enabled("direct-boot-enabled")
}

// IsPasswordlessSudoEnabled checks temporary sudo file
func IsPasswordlessSudoEnabled() bool {
  user := os.Getenv("USER")
  _, err := os.Stat("/tmp/hoopoo-sudo-passwordless-" + user)
  return err == nil
}

// GetCurrentTheme returns the current theme name
func GetCurrentTheme() string {
  // TODO: Implement theme detection
  return ""
}

// GetCurrentFont returns the current font name
func GetCurrentFont() string {
  return font.Current()
}

// GetCurrentPowerProfile returns current power profile
func GetCurrentPowerProfile() string {
  cmd := exec.Command("powerprofilesctl", "get")
  output, err := cmd.Output()
  if err != nil {
    return ""
  }
  return strings.TrimSpace(string(output))
}

// GetCurrentBrowser returns the current default browser
func GetCurrentBrowser() string {
  cmd := exec.Command("xdg-settings", "get", "default-web-browser")
  output, err := cmd.Output()
  if err != nil {
    return ""
  }
  browser := strings.TrimSpace(string(output))
  // Strip .desktop extension
  return strings.TrimSuffix(browser, ".desktop")
}

// GetCurrentTerminal returns the current default terminal
func GetCurrentTerminal() string {
  // Check hoopoo's terminal state
  cmd := exec.Command("hoopoo", "default", "terminal")
  output, err := cmd.Output()
  if err != nil {
    return ""
  }
  return strings.TrimSpace(string(output))
}

// GetCurrentEditor returns the current default editor
func GetCurrentEditor() string {
  // Check EDITOR env var first
  if editor := os.Getenv("EDITOR"); editor != "" {
    return editor
  }
  return ""
}

// GetCurrentMonitorScale returns the current monitor scale for a specific output
func GetCurrentMonitorScale(output string) string {
  // If no output specified, try eDP first, then focused
  if output == "" {
    var err error
    output, err = niri.GetInternalDisplay()
    if err != nil {
      output, err = niri.GetFocusedOutput()
      if err != nil {
        return "1"
      }
    }
  }
  
  scale, err := niri.GetCurrentScale(output)
  if err != nil {
    return "1"
  }
  
  // Format scale value (remove trailing zeros)
  return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", scale), "0"), ".")
}

// GetCurrentMonitorResolution returns the current resolution for a specific output
func GetCurrentMonitorResolution(output string) string {
  if output == "" {
    var err error
    output, err = niri.GetInternalDisplay()
    if err != nil {
      output, err = niri.GetFocusedOutput()
      if err != nil {
        return ""
      }
    }
  }
  
  resolution, err := niri.GetCurrentResolution(output)
  if err != nil {
    return ""
  }
  
  return resolution
}

// GetCurrentMonitorTransform returns the current transform for a specific output
func GetCurrentMonitorTransform(output string) string {
  if output == "" {
    var err error
    output, err = niri.GetInternalDisplay()
    if err != nil {
      output, err = niri.GetFocusedOutput()
      if err != nil {
        return "normal"
      }
    }
  }
  
  transform, err := niri.GetCurrentTransform(output)
  if err != nil {
    return "normal"
  }
  
  return transform
}

// GetCurrentUpdateChannel returns the update channel (stable/dev)
func GetCurrentUpdateChannel() string {
  cmd := exec.Command("hoopoo", "channel")
  output, err := cmd.Output()
  if err != nil {
    return "stable"
  }
  return strings.TrimSpace(string(output))
}

// IsSuspendDisabled checks if suspend is disabled
func IsSuspendDisabled() bool {
  return toggle.Enabled("suspend-off")
}

// IsHibernationAvailable checks if hibernation is available
func IsHibernationAvailable() bool {
  cmd := exec.Command("hoopoo", "hibernation", "available")
  return cmd.Run() == nil
}

// HasExternalMonitors checks if external monitors are connected
func HasExternalMonitors() bool {
  return hardware.HasExternalMonitors()
}

// GetAllOutputNames returns all connected output names
func GetAllOutputNames() ([]string, error) {
  return niri.GetAllOutputs()
}
