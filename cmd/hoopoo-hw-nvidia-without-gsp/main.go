// hoopoo:summary=Detect whether the computer has an NVIDIA GPU without GSP firmware (Maxwell/Pascal/Volta).

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsNvidiaWithoutGSP() {
    os.Exit(0)
  }
  os.Exit(1)
}
