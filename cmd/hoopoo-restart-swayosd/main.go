// hoopoo:summary=Restart the SwayOSD server

package main

import (
  "os/exec"

  "hoopoo/pkg/system/systemd"
)

func main() {
  // Daemon reload
  systemd.DaemonReload()
  
  // Stop service (ignore errors)
  systemd.Stop("swayosd-server.service")
  
  // Kill any remaining processes
  exec.Command("pkill", "-x", "swayosd-server").Run()
  
  // Reset failed state
  systemd.ResetFailed("swayosd-server.service")
  
  // Enable and start service
  systemd.EnableNow("swayosd-server.service")
}
