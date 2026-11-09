// hoopoo:summary=Returns true when an external monitor is physically connected.

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.HasExternalMonitors() {
    os.Exit(0)
  }
  os.Exit(1)
}
