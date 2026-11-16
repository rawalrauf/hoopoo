// hoopoo:summary=Reload Helix configuration

package main

import (
  "os/exec"
)

func main() {
  // Check if helix is running
  if exec.Command("pgrep", "-x", "helix").Run() == nil {
    // Send SIGUSR1 to reload config
    exec.Command("pkill", "-USR1", "helix").Run()
  }
}
