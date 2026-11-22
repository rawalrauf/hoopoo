// hoopoo:summary=Clear a hoopoo toggle flag (remove state file)
// hoopoo:group=toggle
// hoopoo:name=clear
// hoopoo:args=<flag-name>

package main

import (
  "fmt"
  "os"

  "hoopoo/pkg/toggle"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo toggle clear <flag-name>")
    os.Exit(1)
  }

  if err := toggle.Clear(os.Args[1]); err != nil {
    fmt.Fprintf(os.Stderr, "Error clearing toggle: %v\n", err)
    os.Exit(1)
  }
}