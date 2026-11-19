// hoopoo:summary=Reload Helix configuration

package main

import (
  "os/exec"

  "hoopoo/pkg/system/process"
)

func main() {
  // Check if helix is running
  if process.IsRunning("helix") {
    // Send SIGUSR1 to reload config
    exec.Command("pkill", "-USR1", "helix").Run()
  }
}
