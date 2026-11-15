// hoopoo:summary=Launch the system monitor TUI (provided by btop).

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Call hoopoo launch or focus tui with btop
  cmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    "tui",
    "--title=System Monitor",
    "btop",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
