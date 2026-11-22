// Input device detection (touchpad, touchscreen, etc.)

package input

import (
  "os"
  "strings"
)

// HasTouchpad checks if a touchpad is present
func HasTouchpad() bool {
  data, err := os.ReadFile("/proc/bus/input/devices")
  if err != nil {
    return false
  }
  
  lower := strings.ToLower(string(data))
  return strings.Contains(lower, "touchpad") || 
         strings.Contains(lower, "trackpad") ||
         strings.Contains(lower, "synaptics")
}

// HasTouchscreen checks if a touchscreen is present
func HasTouchscreen() bool {
  data, err := os.ReadFile("/proc/bus/input/devices")
  if err != nil {
    return false
  }
  
  lower := strings.ToLower(string(data))
  return strings.Contains(lower, "touchscreen") || 
         strings.Contains(lower, "multi-touch") ||
         strings.Contains(lower, "wacom")
}
