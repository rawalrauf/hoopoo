// hoopoo:summary=Reload btop configuration (used by the Hoopoo theme switching).

package main

import (
  "hoopoo/pkg/system/process"
)

func main() {
  // Send SIGUSR2 to btop to reload config
  process.SendSignal("SIGUSR2", "btop")
}
