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
  
  // Get current font and move it to the top
  current := font.Current()
  if current != "" {
    // Remove current font from list if it exists
    for i, f := range fonts {
      if f == current {
        fonts = append(fonts[:i], fonts[i+1:]...)
        break
      }
    }
    // Print current font first
    fmt.Println(current)
  }
  
  // Print remaining fonts
  for _, f := range fonts {
    fmt.Println(f)
  }
}
