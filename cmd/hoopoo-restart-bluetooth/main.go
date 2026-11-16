// hoopoo:summary=Unblock and restart the bluetooth service.

package main

import (
  "fmt"
  "os/exec"
)

func main() {
  fmt.Println("Unblocking bluetooth...\n")
  
  // Unblock bluetooth
  exec.Command("rfkill", "unblock", "bluetooth").Run()
  
  // List bluetooth devices
  cmd := exec.Command("rfkill", "list", "bluetooth")
  cmd.Stdout = nil
  cmd.Stderr = nil
  cmd.Run()
}
