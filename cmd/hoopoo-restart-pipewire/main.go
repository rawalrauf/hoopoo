// hoopoo:summary=Restart the PipeWire audio service to fix audio issues or apply new configuration.

package main

import (
  "fmt"

  "hoopoo/pkg/system/systemd"
)

func main() {
  fmt.Println("Restarting PipeWire audio services...\n")
  
  // Restart wireplumber and pipewire
  systemd.Restart("wireplumber.service")
  systemd.Restart("pipewire.service")
  
  // Check if pipewire-pulse is active, then restart it
  if systemd.IsActive("pipewire-pulse.service") {
    systemd.Restart("pipewire-pulse.service")
  }
}
