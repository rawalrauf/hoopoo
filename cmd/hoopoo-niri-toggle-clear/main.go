// hoopoo:summary=Clear a niri toggle flag (remove state file)
// hoopoo:group=niri
// hoopoo:name=toggle clear
// hoopoo:args=<flag-name>

package main

import (
  "fmt"
  "os"

  "hoopoo/pkg/niri"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri toggle clear <flag-name>")
    os.Exit(1)
  }

  if err := niri.ClearState(os.Args[1]); err != nil {
    fmt.Fprintf(os.Stderr, "Error clearing toggle: %v\n", err)
    os.Exit(1)
  }
}