// hoopoo:summary=Reset Waybar config to Hoopoo defaults
// hoopoo:examples=hoopoo refresh waybar

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Refresh config files
  configs := []string{"waybar/config.jsonc", "waybar/style.css"}
  
  for _, config := range configs {
    cmd := exec.Command("hoopoo", "refresh", "config", config)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Run()
  }

  // Restart waybar
  cmd := exec.Command("hoopoo", "restart", "waybar")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
