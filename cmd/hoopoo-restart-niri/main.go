// hoopoo:summary=Reload the niri compositor configuration

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/niri"
)

func main() {
  err := niri.Reload()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reloading niri: %v\n", err)
    os.Exit(1)
  }
}
