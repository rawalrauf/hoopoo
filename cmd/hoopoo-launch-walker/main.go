// hoopoo:summary=Launch Walker and ensure its Elephant data provider is running

package main

import (
  "os"

  "hoopoo/pkg/walker"
)

func main() {
  walker.Launch(os.Args[1:])
}
