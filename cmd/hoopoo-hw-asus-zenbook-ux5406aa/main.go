// hoopoo:summary=Detect ASUS Zenbook UX5406AA series laptops on Intel Panther Lake.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.MatchHardware("ux5406aa") && hardware.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
