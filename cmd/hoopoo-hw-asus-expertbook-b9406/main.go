// hoopoo:summary=Detect ASUS ExpertBook B9406 series laptops on Intel Panther Lake.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.MatchHardware("B9406") && hardware.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
