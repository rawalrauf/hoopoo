// hoopoo:summary=Restart Waybar
// hoopoo:examples=hoopoo restart waybar

package main

import (
  "os/exec"
  "syscall"

  "hoopoo/pkg/system/process"
)

func main() {
  // Kill waybar with SIGKILL
  process.KillForce("waybar")

  // Relaunch waybar
  cmd := exec.Command("waybar")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
