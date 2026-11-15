// hoopoo:summary=Launch the Hoopoo audio controls TUI (provided by wiremix).

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Call hoopoo launch or focus tui with wiremix
  cmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    "tui",
    "--title=Volume Control",
    "wiremix",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
