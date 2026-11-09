// hoopoo:summary=Check if a touchscreen or touch device is present

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.HasTouchscreen() {
    os.Exit(0)
  }
  os.Exit(1)
}
