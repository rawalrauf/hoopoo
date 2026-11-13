// hoopoo:summary=Check if a touchscreen or touch device is present

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.HasTouchscreen() {
    os.Exit(0)
  }
  os.Exit(1)
}
