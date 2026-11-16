// hoopoo:summary=Overwrite the user configs for swayosd (controls on-screen feedback for changing volume/songs etc) with the Hoopoo defaults and restart the service.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Refresh config files
  configs := []string{"swayosd/config.toml", "swayosd/style.css"}
  
  for _, config := range configs {
    cmd := exec.Command("hoopoo", "refresh", "config", config)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Run()
  }

  // Restart swayosd
  cmd := exec.Command("hoopoo", "restart", "swayosd")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
