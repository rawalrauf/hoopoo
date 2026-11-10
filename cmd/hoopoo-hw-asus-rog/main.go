// hoopoo:summary=Detect whether the computer is an Asus ROG machine.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsAsusROG() {
    os.Exit(0)
  }
  os.Exit(1)
}
