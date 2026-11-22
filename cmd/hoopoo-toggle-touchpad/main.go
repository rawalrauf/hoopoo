// hoopoo:summary=Toggle the touchpad on/off

package main

import (
  "os/exec"

  "hoopoo/pkg/niri"
)

const toggleFile = "touchpad-disabled.kdl"

const kdlContent = `input {
    touchpad {
        off
    }
}`

func main() {
  // Check if currently disabled
  isDisabled := exec.Command("hoopoo", "niri", "toggle", "enabled", toggleFile).Run() == nil

  if isDisabled {
    // Enable by removing toggle file
    exec.Command("hoopoo", "niri", "toggle", "clear", toggleFile).Run()
    exec.Command("hoopoo", "restart", "niri").Run()
    exec.Command("hoopoo", "swayosd", "client", "--custom-icon", "input-touchpad-symbolic", "--custom-message", "Touchpad enabled").Run()
  } else {
    // Disable by creating toggle file
    niri.WriteState(toggleFile, kdlContent)
    exec.Command("hoopoo", "restart", "niri").Run()
    exec.Command("hoopoo", "swayosd", "client", "--custom-icon", "touchpad-disabled-symbolic", "--custom-message", "Touchpad disabled").Run()
  }
}
