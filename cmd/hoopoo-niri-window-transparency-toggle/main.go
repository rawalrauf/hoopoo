// hoopoo:summary=Toggle window transparency on/off in niri

package main

import (
  "os/exec"
  "syscall"

  "hoopoo/pkg/niri"
)

const toggleFile = "transparency-disabled.kdl"

const kdlContent = `window-rule {
    opacity 1.0
}`

func main() {
  // Check if currently disabled (inverted logic)
  isDisabled := exec.Command("hoopoo", "niri", "toggle", "enabled", toggleFile).Run() == nil

  if isDisabled {
    // Enable transparency by removing toggle file
    exec.Command("hoopoo", "niri", "toggle", "clear", toggleFile).Run()
    exec.Command("hoopoo", "notification", "send", "  Window transparency enabled").Run()
  } else {
    // Disable transparency by creating toggle file
    niri.WriteState(toggleFile, kdlContent)
    exec.Command("hoopoo", "notification", "send", "  Window transparency disabled").Run()
  }

  // Restart niri
  cmd := exec.Command("hoopoo", "restart", "niri")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
