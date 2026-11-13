// hoopoo:summary=Detect whether the computer is a Microsoft Surface device.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsSurface() {
    os.Exit(0)
  }
  os.Exit(1)
}
