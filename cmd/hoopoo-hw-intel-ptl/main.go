// hoopoo:summary=Detect whether the computer has an Intel Panther Lake GPU.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
