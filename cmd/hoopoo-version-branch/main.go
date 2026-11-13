// hoopoo:summary=Print the current Hoopoo git branch

package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "hoopoo/pkg/system/paths"
)

func main() {
  cmd := exec.Command("git", "-C", paths.Hoopoo(), "rev-parse", "--abbrev-ref", "HEAD")
  output, err := cmd.Output()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error getting git branch: %v\n", err)
    os.Exit(1)
  }
  fmt.Println(strings.TrimSpace(string(output)))
}
