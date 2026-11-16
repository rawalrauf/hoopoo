// hoopoo:summary=Reset niri user configs to Hoopoo defaults

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Refresh all niri config files
  configs := []string{
    "niri/config.kdl",
    "niri/autostart.kdl",
    "niri/bindings.kdl",
    "niri/looknfeel.kdl",
    "niri/monitors.kdl",
    "niri/input.kdl",
    "niri/flags.kdl",
    "niri/window-rules.kdl",
  }
  
  for _, config := range configs {
    cmd := exec.Command("hoopoo", "refresh", "config", config)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Run()
  }

  // Restart niri (ignore errors)
  cmd := exec.Command("hoopoo", "restart", "niri")
  cmd.Stdout = nil
  cmd.Stderr = nil
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
