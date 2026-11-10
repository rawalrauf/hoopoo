// hoopoo:summary=Reload the niri compositor configuration

package main

import (
  "fmt"
  "os"
  "os/exec"
)

func main() {
  cmd := exec.Command("niri", "msg", "action", "load-config-file")
  err := cmd.Run()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reloading niri: %v\n", err)
    os.Exit(1)
  }
}
