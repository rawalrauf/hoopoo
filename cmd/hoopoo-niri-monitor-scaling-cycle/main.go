// hoopoo:summary=Set or cycle niri output scaling through 1x, 1.25x, 1.5x, 2x
// hoopoo:group=niri
// hoopoo:name=monitor scaling cycle
// hoopoo:args=[--set <scale>|--reverse|--current]

package main

import (
  "fmt"
  "math"
  "os"
  "os/exec"
  "strconv"

  "hoopoo/pkg/niri"
  "hoopoo/pkg/niri/config"
)

var scales = []float64{1, 1.25, 1.5, 2}

func getOutput() (string, error) {
  // Try to get eDP output first
  internal, err := niri.GetInternalDisplay()
  if err == nil {
    return internal, nil
  }

  // Fall back to focused output
  return niri.GetFocusedOutput()
}

func findClosestScaleIndex(current float64) int {
  bestIdx := 0
  bestDiff := math.Abs(current - scales[0])
  
  for i, scale := range scales {
    diff := math.Abs(current - scale)
    if diff < bestDiff {
      bestDiff = diff
      bestIdx = i
    }
  }
  
  return bestIdx
}

func updateConfigFile(output string, newScale float64) error {
  return config.SetProperty("monitors.kdl", "output", output, "scale", fmt.Sprintf("%g", newScale))
}

func main() {
  output, err := getOutput()
  if err != nil {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  No output found").Run()
    os.Exit(1)
  }

  currentScale, err := niri.GetCurrentScale(output)
  if err != nil {
    currentScale = 1.0
  }

  // Handle --current flag
  if len(os.Args) > 1 && os.Args[1] == "--current" {
    fmt.Println(currentScale)
    return
  }

  var newScale float64

  // Handle --set flag
  if len(os.Args) > 2 && os.Args[1] == "--set" {
    newScale, err = strconv.ParseFloat(os.Args[2], 64)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Invalid scale value")
      os.Exit(1)
    }
  } else {
    // Cycle through scales
    currentIdx := findClosestScaleIndex(currentScale)
    
    var newIdx int
    if len(os.Args) > 1 && os.Args[1] == "--reverse" {
      newIdx = (currentIdx - 1 + len(scales)) % len(scales)
    } else {
      newIdx = (currentIdx + 1) % len(scales)
    }
    
    newScale = scales[newIdx]
  }

  // Apply the scale
  niri.SetOutputScale(output, newScale)
  
  // Update config file
  if err := updateConfigFile(output, newScale); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to update config: %v\n", err)
    os.Exit(1)
  }
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", fmt.Sprintf("󰍹  Display scaling set to %gx", newScale)).Run()
}
