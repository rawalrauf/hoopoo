// hoopoo:summary=Check if a niri toggle flag is currently disabled.
// hoopoo:args=<flag-name>

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/common/toggle"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-niri-toggle-disabled <flag-name>")
    os.Exit(1)
  }

  if toggle.NiriDisabled(os.Args[1]) {
    os.Exit(0)
  }
  os.Exit(1)
}
