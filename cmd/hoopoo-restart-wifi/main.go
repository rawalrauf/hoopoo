// hoopoo:summary=Unblock and restart the Wi-Fi service.

package main

import (
  "fmt"
  "os/exec"
)

func main() {
  fmt.Println("Unblocking wifi...\n")
  
  // Unblock wifi
  exec.Command("rfkill", "unblock", "wifi").Run()
  
  // List wifi devices
  cmd := exec.Command("rfkill", "list", "wifi")
  cmd.Stdout = nil
  cmd.Stderr = nil
  cmd.Run()
}
