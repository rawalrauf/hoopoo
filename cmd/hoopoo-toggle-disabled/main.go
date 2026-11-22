// hoopoo:summary=Check if a hoopoo toggle flag is currently disabled
// hoopoo:group=toggle
// hoopoo:name=disabled
// hoopoo:args=<flag-name>

package main

import (
  "fmt"
  "os"

  "hoopoo/pkg/toggle"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo toggle disabled <flag-name>")
    os.Exit(1)
  }

  if toggle.Disabled(os.Args[1]) {
    os.Exit(0)
  }
  os.Exit(1)
}