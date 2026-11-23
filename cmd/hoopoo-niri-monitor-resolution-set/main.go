// hoopoo:summary=Set niri output resolution (mode)
// hoopoo:group=niri
// hoopoo:name=monitor resolution set
// hoopoo:args=--output <name> <resolution>

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
  var resolution string
  
  // Parse arguments
  for i := 1; i < len(os.Args); i++ {
    arg := os.Args[i]
    
    if strings.HasPrefix(arg, "--output=") {
      output = strings.TrimPrefix(arg, "--output=")
    } else if arg == "--output" && i+1 < len(os.Args) {
      output = os.Args[i+1]
      i++
    } else if resolution == "" {
      resolution = arg
    }
  }
  
  if output == "" || resolution == "" {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor resolution set --output <name> <resolution>")
    os.Exit(1)
  }
  
  // Set the resolution via niri msg
  if err := niri.SetOutputMode(output, resolution); err != nil {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰹑  Failed to set resolution").Run()
    os.Exit(1)
  }
  
  // Update config file
  if err := config.SetProperty("monitors.kdl", "output", output, "mode", resolution); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to update config: %v\n", err)
  }
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", fmt.Sprintf("󰹑  %s resolution set to %s", output, resolution)).Run()
}
