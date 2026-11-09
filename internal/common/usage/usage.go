// Handles --help flag and validates argument count

package usage

import (
  "fmt"
  "os"
)

func Message(message string, minArgs ...int) {
  actualArgs := len(os.Args) - 1
  
  requiredArgs := 0
  if len(minArgs) > 0 {
    requiredArgs = minArgs[0]
  }
  
  for _, arg := range os.Args[1:] {
    if arg == "--help" || arg == "-h" {
      fmt.Fprintln(os.Stderr, message)
      os.Exit(0)
    }
  }
  
  if actualArgs < requiredArgs {
    fmt.Fprintln(os.Stderr, message)
    os.Exit(1)
  }
}
