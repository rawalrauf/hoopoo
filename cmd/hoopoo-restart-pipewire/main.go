// hoopoo:summary=Restart the PipeWire audio service to fix audio issues or apply new configuration.

package main

import (
  "fmt"
  "os/exec"
)

func main() {
  fmt.Println("Restarting PipeWire audio services...\n")
  
  // Restart wireplumber and pipewire
  exec.Command("systemctl", "--user", "restart", "wireplumber.service", "pipewire.service").Run()
  
  // Check if pipewire-pulse is active, then restart it
  checkCmd := exec.Command("systemctl", "--user", "--quiet", "is-active", "pipewire-pulse.service")
  if checkCmd.Run() == nil {
    exec.Command("systemctl", "--user", "restart", "pipewire-pulse.service").Run()
  }
}
