// Niri compositor core utilities

package niri

import (
  "os/exec"
  "strconv"
)

// Reload reloads the niri configuration
func Reload() error {
  return exec.Command("niri", "msg", "action", "load-config-file").Run()
}

// FocusWorkspace focuses a workspace by number
func FocusWorkspace(num int) error {
  return exec.Command("niri", "msg", "action", "focus-workspace", strconv.Itoa(num)).Run()
}
