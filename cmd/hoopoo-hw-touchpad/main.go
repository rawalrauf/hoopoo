// hoopoo:summary=Check if a touchpad or trackpad device is present

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.HasTouchpad() {
    os.Exit(0)
  }
  os.Exit(1)
}
