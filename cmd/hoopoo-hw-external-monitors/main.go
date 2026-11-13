// hoopoo:summary=Returns true when an external monitor is physically connected.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.HasExternalMonitors() {
    os.Exit(0)
  }
  os.Exit(1)
}
