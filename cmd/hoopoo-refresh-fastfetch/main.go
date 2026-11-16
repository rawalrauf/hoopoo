// hoopoo:summary=Overwrite the user config for fastfetch with the Hoopoo default.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  cmd := exec.Command("hoopoo", "refresh", "config", "fastfetch/config.jsonc")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
