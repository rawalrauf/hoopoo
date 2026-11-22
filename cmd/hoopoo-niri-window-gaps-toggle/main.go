// hoopoo:summary=Toggle window gaps on/off in niri

package main

import (
  "os/exec"
  "syscall"

  "hoopoo/pkg/niri"
)

const toggleFile = "gaps-enabled.kdl"

const kdlContent = `layout {
    gaps 12
}`

func main() {
  // Check if currently enabled
  isEnabled := exec.Command("hoopoo", "niri", "toggle", "enabled", toggleFile).Run() == nil

  if isEnabled {
    // Disable by removing toggle file
    exec.Command("hoopoo", "niri", "toggle", "clear", toggleFile).Run()
    exec.Command("hoopoo", "notification", "send", "Window gaps disabled").Run()
  } else {
    // Enable by creating toggle file
    niri.WriteState(toggleFile, kdlContent)
    exec.Command("hoopoo", "notification", "send", "Window gaps enabled").Run()
  }

  // Restart niri
  cmd := exec.Command("hoopoo", "restart", "niri")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
