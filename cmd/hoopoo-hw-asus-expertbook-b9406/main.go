// hoopoo:summary=Detect ASUS ExpertBook B9406 series laptops on Intel Panther Lake.

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.MatchHardware("B9406") && hw.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
