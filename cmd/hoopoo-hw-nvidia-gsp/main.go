// hoopoo:summary=Detect whether the computer has an NVIDIA GPU with GSP firmware (Turing or newer).

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.IsNvidiaGSP() {
    os.Exit(0)
  }
  os.Exit(1)
}
