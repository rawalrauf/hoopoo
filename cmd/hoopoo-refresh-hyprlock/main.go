// hoopoo:summary=Overwrite the user config for hyprlock with the Hoopoo default.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  cmd := exec.Command("hoopoo", "refresh", "config", "hypr/hyprlock.conf")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
