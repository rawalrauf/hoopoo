// hoopoo:summary=Detect ASUS Zenbook UX5406AA series laptops on Intel Panther Lake.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.MatchHardware("ux5406aa") && hw.IsIntelPantherLake() {
    os.Exit(0)
  }
  os.Exit(1)
}
