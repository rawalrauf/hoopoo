// hoopoo:summary=Show current monospace font
// hoopoo:group=font
// hoopoo:name=current
// hoopoo:examples=hoopoo font current

package main

import (
  "fmt"

  "hoopoo/pkg/font"
)

func main() {
  current := font.Current()
  if current == "" {
    fmt.Println("No font configured")
    return
  }
  fmt.Println(current)
}
