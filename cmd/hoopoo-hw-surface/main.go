// hoopoo:summary=Detect whether the computer is a Microsoft Surface device.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsSurface() {
    os.Exit(0)
  }
  os.Exit(1)
}
