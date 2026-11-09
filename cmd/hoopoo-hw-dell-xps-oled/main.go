// hoopoo:summary=Match Dell XPS systems with LG OLED panel on Intel Panther Lake (Xe3) GPU.

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.IsDellXPSOLED() {
    os.Exit(0)
  }
  os.Exit(1)
}
