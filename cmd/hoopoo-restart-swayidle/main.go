// hoopoo:summary=Restart the swayidle service (used for idle detection and auto-lock).

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "syscall"
)

func main() {
  configPath := filepath.Join(os.Getenv("HOME"), ".config/swayidle/swayidle.conf")
  
  // Call hoopoo restart app with swayidle
  cmd := exec.Command(
    "hoopoo",
    "restart",
    "app",
    "swayidle",
    "-w",
    "-C",
    configPath,
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
