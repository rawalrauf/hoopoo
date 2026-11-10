// hoopoo:summary=Detect whether the computer has an Intel Panther Lake GPU.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
