// hoopoo:summary=Restart the SwayOSD server

package main

import (
  "os/exec"
)

func main() {
  // Daemon reload
  exec.Command("systemctl", "--user", "daemon-reload").Run()
  
  // Stop service (ignore errors)
  exec.Command("systemctl", "--user", "stop", "swayosd-server.service").Run()
  
  // Kill any remaining processes
  exec.Command("pkill", "-x", "swayosd-server").Run()
  
  // Reset failed state
  exec.Command("systemctl", "--user", "reset-failed", "swayosd-server.service").Run()
  
  // Enable and start service
  exec.Command("systemctl", "--user", "enable", "--now", "swayosd-server.service").Run()
}
