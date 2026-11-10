// hoopoo:summary=Detect whether the computer is a Framework Laptop 16.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsFramework16() {
    os.Exit(0)
  }
  os.Exit(1)
}
