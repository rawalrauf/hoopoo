// hoopoo:summary=Check if a touchpad or trackpad device is present

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.HasTouchpad() {
    os.Exit(0)
  }
  os.Exit(1)
}
