// hoopoo:summary=List available monospace fonts
// hoopoo:group=font
// hoopoo:name=list
// hoopoo:examples=hoopoo font list | hoopoo font set "CaskaydiaMono Nerd Font"

package main

import (
  "fmt"
  "os"
  "sort"

  "hoopoo/pkg/font"
)

func main() {
  fonts, err := font.List()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error listing fonts: %v\n", err)
    os.Exit(1)
  }

  sort.Strings(fonts)
  for _, f := range fonts {
    fmt.Println(f)
  }
}
