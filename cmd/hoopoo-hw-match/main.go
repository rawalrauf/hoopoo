// hoopoo:summary=Match against the computer's DMI product name or product family (case-insensitive).
// hoopoo:args=<pattern>

package main

import (
  "fmt"
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-hw-match <pattern>")
    os.Exit(1)
  }

  if hardware.MatchHardware(os.Args[1]) {
    os.Exit(0)
  }
  os.Exit(1)
}
