// hoopoo:summary=Match Dell XPS systems with the Synaptics haptic touchpad.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.HasDellXPSHapticTouchpad() {
    os.Exit(0)
  }
  os.Exit(1)
}
