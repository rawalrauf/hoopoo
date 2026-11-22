// hoopoo:summary=Reload Helix configuration

package main

import (
  "hoopoo/pkg/system/process"
)

func main() {
  // Check if helix is running
  if process.IsRunning("helix") {
    // Send SIGUSR1 to reload config
    process.SendSignal("USR1", "helix")
  }
}
