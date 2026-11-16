// hoopoo:summary=Overwrite the user tmux config with the Hoopoo default and reload tmux.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Refresh config
  cmd := exec.Command("hoopoo", "refresh", "config", "tmux/tmux.conf")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()

  // Restart tmux
  cmd = exec.Command("hoopoo", "restart", "tmux")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
