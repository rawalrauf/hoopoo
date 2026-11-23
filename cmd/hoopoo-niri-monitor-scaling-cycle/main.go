// hoopoo:summary=Set or cycle niri output scaling through 1x, 1.25x, 1.5x, 2x
// hoopoo:group=niri
// hoopoo:name=monitor scaling cycle
// hoopoo:args=[--output <name>] [--set <scale>|--reverse|--current]

package main

import (
  "fmt"
  "math"
  "os"
  "os/exec"
  "strconv"
  "strings"

  "hoopoo/pkg/niri"
  "hoopoo/pkg/niri/config"
)

var scales = []float64{1, 1.25, 1.5, 2}

func getOutput(explicitOutput string) (string, error) {
  // Use explicit output if provided
  if explicitOutput != "" {
    return explicitOutput, nil
  }
  
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
  var explicitOutput string
  var setScale string
  var reverse bool
  var current bool
  
  // Parse arguments
  for i := 1; i < len(os.Args); i++ {
    arg := os.Args[i]
    
    if strings.HasPrefix(arg, "--output=") {
      explicitOutput = strings.TrimPrefix(arg, "--output=")
    } else if arg == "--output" && i+1 < len(os.Args) {
      explicitOutput = os.Args[i+1]
      i++
    } else if strings.HasPrefix(arg, "--set=") {
      setScale = strings.TrimPrefix(arg, "--set=")
    } else if arg == "--set" && i+1 < len(os.Args) {
      setScale = os.Args[i+1]
      i++
    } else if arg == "--reverse" {
      reverse = true
    } else if arg == "--current" {
      current = true
    }
  }
  
  output, err := getOutput(explicitOutput)
  if err != nil {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  No output found").Run()
    os.Exit(1)
  }

  currentScale, err := niri.GetCurrentScale(output)
  if err != nil {
    currentScale = 1.0
  }

  // Handle --current flag
  if current {
    fmt.Println(currentScale)
    return
  }

  var newScale float64

  // Handle --set flag
  if setScale != "" {
    newScale, err = strconv.ParseFloat(setScale, 64)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Invalid scale value")
      os.Exit(1)
    }
  } else {
    // Cycle through scales
    currentIdx := findClosestScaleIndex(currentScale)
    
    var newIdx int
    if reverse {
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
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", fmt.Sprintf("󰍹  %s scaling set to %gx", output, newScale)).Run()
}
