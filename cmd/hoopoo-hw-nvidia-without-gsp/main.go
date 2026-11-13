// hoopoo:summary=Detect whether the computer has an NVIDIA GPU without GSP firmware (Maxwell/Pascal/Volta).

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsNvidiaWithoutGSP() {
    os.Exit(0)
  }
  os.Exit(1)
}
