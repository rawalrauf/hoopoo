// hoopoo:summary=Launch the Hoopoo wifi controls (provided by the Impala TUI).

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Unblock wifi radio
  exec.Command("rfkill", "unblock", "wifi").Run()
  
  // Call hoopoo launch or focus tui with impala
  cmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    "tui",
    "--title=Wi-Fi",
    "impala",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
