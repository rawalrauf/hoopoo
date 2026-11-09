// hoopoo:summary=Match Dell XPS systems with the Synaptics haptic touchpad.

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.HasDellXPSHapticTouchpad() {
    os.Exit(0)
  }
  os.Exit(1)
}
