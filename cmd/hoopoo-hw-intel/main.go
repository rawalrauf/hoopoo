// hoopoo:summary=Detect whether the computer has an Intel CPU.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsIntelCPU() {
    os.Exit(0)
  }
  os.Exit(1)
}
