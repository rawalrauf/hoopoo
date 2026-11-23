// Niri window management

package niri

import (
  "encoding/json"
  "os/exec"
  "strconv"
)

// GetWindows returns all open windows
func GetWindows() ([]Window, error) {
  out, err := exec.Command("niri", "msg", "-j", "windows").Output()
  if err != nil {
    return nil, err
  }

  var windows []Window
  if err := json.Unmarshal(out, &windows); err != nil {
    return nil, err
  }

  return windows, nil
}

// GetFocusedWindow returns the currently focused window
func GetFocusedWindow() (*Window, error) {
  out, err := exec.Command("niri", "msg", "--json", "focused-window").Output()
  if err != nil {
    return nil, err
  }

  var window Window
  if err := json.Unmarshal(out, &window); err != nil {
    return nil, err
  }

  return &window, nil
}

// CloseWindow closes a window by ID
func CloseWindow(id int) error {
  return exec.Command("niri", "msg", "action", "close-window", "--id", strconv.Itoa(id)).Run()
}

// FocusWindow focuses a window by ID
func FocusWindow(id int) error {
  return exec.Command("niri", "msg", "action", "focus-window", "--id", strconv.Itoa(id)).Run()
}
