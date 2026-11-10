// hoopoo:summary=Detect whether the computer is a Lenovo ThinkPad.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.IsLenovoThinkPad() {
    os.Exit(0)
  }
  os.Exit(1)
}
