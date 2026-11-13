// hoopoo:summary=Print the installed Hoopoo version

package main

import (
  "fmt"
  "os"

  "hoopoo/pkg/system/paths"
)

func main() {
  versionFile := paths.Hoopoo() + "/version"
  data, err := os.ReadFile(versionFile)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading version file: %v\n", err)
    os.Exit(1)
  }
  fmt.Print(string(data))
}
