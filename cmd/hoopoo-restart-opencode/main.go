// hoopoo:summary=Reload opencode configuration (used by the Hoopoo theme switching).

package main

import (
  "os/exec"

  "hoopoo/pkg/system/process"
)

func main() {
  // Check if opencode is running
  if process.IsRunning("opencode") {
    // Send SIGUSR2 to reload config
    exec.Command("killall", "-SIGUSR2", "opencode").Run()
  }
}
