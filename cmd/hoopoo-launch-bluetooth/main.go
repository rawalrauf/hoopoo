// hoopoo:summary=Launch the Hoopoo bluetooth controls TUI (provided by bluetui).

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Unblock bluetooth radio
  exec.Command("rfkill", "unblock", "bluetooth").Run()
  
  // Call hoopoo launch or focus tui with bluetui
  cmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    "tui",
    "--title=Bluetooth",
    "bluetui",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
