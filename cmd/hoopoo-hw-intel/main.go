// hoopoo:summary=Detect whether the computer has an Intel CPU.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsIntelCPU() {
    os.Exit(0)
  }
  os.Exit(1)
}
