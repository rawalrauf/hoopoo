// hoopoo:summary=Reload btop configuration (used by the Hoopoo theme switching).

package main

import (
  "os/exec"
)

func main() {
  // Send SIGUSR2 to btop to reload config
  exec.Command("pkill", "-SIGUSR2", "btop").Run()
}
