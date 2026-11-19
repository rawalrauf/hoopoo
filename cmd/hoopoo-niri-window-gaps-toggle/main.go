// hoopoo:summary=Toggle window gaps on/off in niri

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  cmd := exec.Command(
    "hoopoo",
    "niri",
    "toggle",
    "--enabled-notification", "  Window gaps enabled",
    "--disabled-notification", "  Window gaps disabled",
    "gaps-enabled",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
