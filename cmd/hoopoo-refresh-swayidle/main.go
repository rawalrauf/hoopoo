// hoopoo:summary=Overwrite the user config for swayidle with the Hoopoo default and restart the service.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Refresh config
  cmd := exec.Command("hoopoo", "refresh", "config", "swayidle/swayidle.conf")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()

  // Restart swayidle
  cmd = exec.Command("hoopoo", "restart", "swayidle")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
