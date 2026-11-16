// hoopoo:summary=Reload opencode configuration (used by the Hoopoo theme switching).

package main

import (
  "os/exec"
)

func main() {
  // Check if opencode is running
  if exec.Command("pgrep", "-x", "opencode").Run() == nil {
    // Send SIGUSR2 to reload config
    exec.Command("killall", "-SIGUSR2", "opencode").Run()
  }
}
