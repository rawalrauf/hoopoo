// hoopoo:summary=Toggle window round corners on/off in niri

package main

import (
  "os/exec"
  "syscall"

  "hoopoo/pkg/niri"
)

const toggleFile = "roundcorners-enabled.kdl"

const kdlContent = `window-rule {
    clip-to-geometry true
    geometry-corner-radius 12
}`

func main() {
  // Check if currently enabled
  isEnabled := exec.Command("hoopoo", "niri", "toggle", "enabled", toggleFile).Run() == nil

  if isEnabled {
    // Disable round corners by removing toggle file (back to default 0)
    exec.Command("hoopoo", "niri", "toggle", "clear", toggleFile).Run()
    exec.Command("hoopoo", "notification", "send", "  Window round corners disabled").Run()
  } else {
    // Enable round corners by creating toggle file with radius 8
    niri.WriteState(toggleFile, kdlContent)
    exec.Command("hoopoo", "notification", "send", "  Window round corners enabled").Run()
  }

  // Restart niri
  cmd := exec.Command("hoopoo", "restart", "niri")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
