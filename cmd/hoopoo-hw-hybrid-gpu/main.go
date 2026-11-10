// hoopoo:summary=Detect whether the system has an active hybrid GPU configuration

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsHybridGPU() {
    os.Exit(0)
  }
  os.Exit(1)
}
