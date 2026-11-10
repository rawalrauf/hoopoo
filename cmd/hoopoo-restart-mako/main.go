// hoopoo:summary=Reload mako configuration (used by the Hoopoo theme switching).

package main

import (
  "fmt"
  "os"
  "os/exec"
)

func main() {
  cmd := exec.Command("makoctl", "reload")
  err := cmd.Run()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reloading mako: %v\n", err)
    os.Exit(1)
  }
}
