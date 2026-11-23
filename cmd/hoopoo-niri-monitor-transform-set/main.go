// hoopoo:summary=Set niri output transform (rotation)
// hoopoo:group=niri
// hoopoo:name=monitor transform set
// hoopoo:args=--output <name> <transform>

package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "hoopoo/pkg/niri"
  "hoopoo/pkg/niri/config"
)

func main() {
  var output string
  var transform string
  
  // Parse arguments
  for i := 1; i < len(os.Args); i++ {
    arg := os.Args[i]
    
    if strings.HasPrefix(arg, "--output=") {
      output = strings.TrimPrefix(arg, "--output=")
    } else if arg == "--output" && i+1 < len(os.Args) {
      output = os.Args[i+1]
      i++
    } else if transform == "" {
      transform = arg
    }
  }
  
  if output == "" || transform == "" {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor transform set --output <name> <transform>")
    os.Exit(1)
  }
  
  // Set the transform via niri msg
  if err := niri.SetOutputTransform(output, transform); err != nil {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰑴  Failed to set transform").Run()
    os.Exit(1)
  }
  
  // Update config file
  if err := config.SetProperty("monitors.kdl", "output", output, "transform", transform); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to update config: %v\n", err)
  }
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", fmt.Sprintf("󰑴  %s transform set to %s", output, transform)).Run()
}
