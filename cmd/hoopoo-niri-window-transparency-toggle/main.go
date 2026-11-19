// hoopoo:summary=Toggle window transparency on/off in niri

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
    "--invert",
    "--enabled-notification", "  Window transparency enabled",
    "--disabled-notification", "  Window transparency disabled",
    "transparency-disabled",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
