// hoopoo:summary=Restart Waybar
// hoopoo:examples=hoopoo restart waybar

package main

import (
  "os/exec"
  "syscall"
)

func main() {
  // Kill waybar with SIGKILL
  exec.Command("pkill", "-9", "-x", "waybar").Run()

  // Relaunch waybar
  cmd := exec.Command("waybar")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
