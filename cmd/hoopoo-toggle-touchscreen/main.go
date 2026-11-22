// hoopoo:summary=Toggle the touch functionality of the screen on/off

package main

import (
  "os/exec"

  "hoopoo/pkg/niri"
)

const toggleFile = "touchscreen-disabled.kdl"

const kdlContent = `input {
    touch {
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
    exec.Command("hoopoo", "swayosd", "client", "--custom-icon", "device-support-touch-symbolic", "--custom-message", "Touchscreen enabled").Run()
  } else {
    // Disable by creating toggle file
    niri.WriteState(toggleFile, kdlContent)
    exec.Command("hoopoo", "restart", "niri").Run()
    exec.Command("hoopoo", "swayosd", "client", "--custom-icon", "touch-disabled-symbolic", "--custom-message", "Touchscreen disabled").Run()
  }
}
