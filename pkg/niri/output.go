// Niri output/display management

package niri

import (
  "encoding/json"
  "fmt"
  "os/exec"
  "regexp"
  "strconv"
  "strings"
)

// GetInternalDisplay returns the internal display name (eDP)
func GetInternalDisplay() (string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return "", err
  }

  re := regexp.MustCompile(`\(([^)]+)\)`)
  matches := re.FindAllStringSubmatch(string(out), -1)
  
  for _, match := range matches {
    if len(match) > 1 && strings.Contains(match[1], "eDP") {
      return match[1], nil
    }
  }
  
  return "", fmt.Errorf("no internal display found")
}

// GetExternalDisplay returns the first external display name (non-eDP)
func GetExternalDisplay() (string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return "", err
  }

  re := regexp.MustCompile(`\(([^)]+)\)`)
  matches := re.FindAllStringSubmatch(string(out), -1)
  
  for _, match := range matches {
    if len(match) > 1 && !strings.Contains(match[1], "eDP") && 
       !strings.Contains(match[1], "preferred") && 
       !strings.Contains(match[1], "current") {
      return match[1], nil
    }
  }
  
  return "", fmt.Errorf("no external display found")
}

// GetAllOutputs returns all connected output names
func GetAllOutputs() ([]string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return nil, err
  }

  var outputs []string
  lines := strings.Split(string(out), "\n")
  
  for _, line := range lines {
    // Only match output names on lines starting with "Output"
    if strings.HasPrefix(line, "Output") {
      re := regexp.MustCompile(`\(([^)]+)\)`)
      match := re.FindStringSubmatch(line)
      if len(match) > 1 {
        outputs = append(outputs, match[1])
      }
    }
  }
  
  if len(outputs) == 0 {
    return nil, fmt.Errorf("no outputs found")
  }
  
  return outputs, nil
}

// GetFocusedOutput returns the name of the currently focused output
func GetFocusedOutput() (string, error) {
  out, err := exec.Command("niri", "msg", "--json", "focused-output").Output()
  if err != nil {
    return "", err
  }

  var data map[string]interface{}
  if json.Unmarshal(out, &data) != nil {
    return "", fmt.Errorf("failed to parse output")
  }

  if name, ok := data["name"].(string); ok {
    return name, nil
  }

  return "", fmt.Errorf("no focused output")
}

// GetFocusedOutputInfo returns full info about the currently focused output
func GetFocusedOutputInfo() (*Output, error) {
  out, err := exec.Command("niri", "msg", "--json", "focused-output").Output()
  if err != nil {
    return nil, err
  }

  var output Output
  if err := json.Unmarshal(out, &output); err != nil {
    return nil, err
  }

  return &output, nil
}

// GetCurrentScale returns the current scale for the given output
func GetCurrentScale(output string) (float64, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return 1.0, err
  }

  lines := strings.Split(string(out), "\n")
  inTarget := false
  
  for _, line := range lines {
    if strings.Contains(line, "("+output+")") {
      inTarget = true
      continue
    }
    if inTarget && strings.Contains(line, "Scale:") {
      parts := strings.Fields(line)
      for i, part := range parts {
        if part == "Scale:" && i+1 < len(parts) {
          scale, err := strconv.ParseFloat(parts[i+1], 64)
          if err != nil {
            return 1.0, err
          }
          return scale, nil
        }
      }
    }
    if inTarget && strings.HasPrefix(line, "Output") {
      break
    }
  }

  return 1.0, fmt.Errorf("scale not found for output %s", output)
}

// GetCurrentResolution returns the current resolution for the given output (e.g., "1920x1080@60")
func GetCurrentResolution(output string) (string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return "", err
  }

  lines := strings.Split(string(out), "\n")
  inTarget := false
  
  for _, line := range lines {
    if strings.Contains(line, "("+output+")") {
      inTarget = true
      continue
    }
    if inTarget && strings.Contains(line, "Current mode:") {
      // Parse "Current mode: 1920x1080 @ 60.008 Hz (preferred)"
      parts := strings.Fields(line)
      for i, part := range parts {
        if part == "mode:" && i+1 < len(parts) {
          // Extract resolution@refresh (e.g., "1920x1080")
          resolution := parts[i+1]
          
          // If refresh rate exists in next field, append it
          if i+3 < len(parts) && parts[i+2] == "@" {
            // Reconstruct as "1920x1080@60.008"
            modeStr := resolution + "@" + parts[i+3]
            
            // Round refresh rate to whole number
            return roundRefreshRate(modeStr), nil
          }
          
          return resolution, nil
        }
      }
    }
    if inTarget && strings.HasPrefix(line, "Output") {
      break
    }
  }

  return "", fmt.Errorf("resolution not found for output %s", output)
}

// GetCurrentTransform returns the current transform for the given output (e.g., "normal", "90", "180", "270")
func GetCurrentTransform(output string) (string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return "", err
  }

  lines := strings.Split(string(out), "\n")
  inTarget := false
  
  for _, line := range lines {
    if strings.Contains(line, "("+output+")") {
      inTarget = true
      continue
    }
    if inTarget && strings.Contains(line, "Transform:") {
      // Parse "Transform: 90° counter-clockwise" or "Transform: normal"
      parts := strings.Fields(line)
      for i, part := range parts {
        if part == "Transform:" && i+1 < len(parts) {
          transform := parts[i+1]
          
          // If it's a rotation like "90°", strip the degree symbol
          transform = strings.TrimSuffix(transform, "°")
          
          return transform, nil
        }
      }
    }
    if inTarget && strings.HasPrefix(line, "Output") {
      break
    }
  }

  return "normal", fmt.Errorf("transform not found for output %s", output)
}

// GetAvailableModes returns all available modes for the given output with original refresh rates
// Returns a map: display string -> original mode string
// Example: "1920x1080@60" -> "1920x1080@60.008"
func GetAvailableModes(output string) (map[string]string, error) {
  out, err := exec.Command("niri", "msg", "outputs").Output()
  if err != nil {
    return nil, err
  }

  lines := strings.Split(string(out), "\n")
  inTarget := false
  inModes := false
  modes := make(map[string]string) // displayString -> originalString
  
  for _, line := range lines {
    // Check if we're at the target output
    if strings.Contains(line, "("+output+")") {
      inTarget = true
      continue
    }
    
    // Check if we've moved to another output
    if inTarget && strings.HasPrefix(line, "Output") {
      break
    }
    
    // Check if we've reached the "Available modes:" section
    if inTarget && strings.Contains(line, "Available modes:") {
      inModes = true
      continue
    }
    
    // Parse mode lines (they're indented and contain resolution@refresh)
    if inModes {
      trimmed := strings.TrimSpace(line)
      if trimmed == "" {
        continue
      }
      
      // Mode lines look like: "1920x1080@60.008 (current, preferred)"
      // Extract the mode string (resolution@refresh) before any annotations
      parts := strings.Fields(trimmed)
      if len(parts) > 0 {
        originalMode := parts[0]
        
        // Create display version with rounded refresh rate
        displayMode := roundRefreshRate(originalMode)
        
        // Store mapping (deduplicated by displayMode key)
        modes[displayMode] = originalMode
      }
    }
  }

  if len(modes) == 0 {
    return nil, fmt.Errorf("no modes found for output %s", output)
  }

  return modes, nil
}

// GetAvailableModesDisplay returns just the display strings (for menu)
func GetAvailableModesDisplay(output string) ([]string, error) {
  modesMap, err := GetAvailableModes(output)
  if err != nil {
    return nil, err
  }
  
  var modes []string
  for displayMode := range modesMap {
    modes = append(modes, displayMode)
  }
  
  return modes, nil
}

// GetOriginalMode returns the original mode string for a display string
func GetOriginalMode(output, displayMode string) (string, error) {
  modesMap, err := GetAvailableModes(output)
  if err != nil {
    return "", err
  }
  
  if originalMode, ok := modesMap[displayMode]; ok {
    return originalMode, nil
  }
  
  return "", fmt.Errorf("mode %s not found for output %s", displayMode, output)
}

// roundRefreshRate rounds the refresh rate to a whole number
// Example: "1920x1080@60.008" → "1920x1080@60"
func roundRefreshRate(mode string) string {
  parts := strings.Split(mode, "@")
  if len(parts) != 2 {
    return mode
  }
  
  resolution := parts[0]
  refreshStr := parts[1]
  
  // Parse and round the refresh rate
  refresh, err := strconv.ParseFloat(refreshStr, 64)
  if err != nil {
    return mode
  }
  
  roundedRefresh := int(refresh + 0.5) // Round to nearest integer
  
  return fmt.Sprintf("%s@%d", resolution, roundedRefresh)
}

// SetOutputScale sets the scale for an output
func SetOutputScale(output string, scale float64) error {
  return exec.Command("niri", "msg", "output", output, "scale", fmt.Sprintf("%g", scale)).Run()
}

// SetOutputMode sets the display mode (resolution) for the given output
func SetOutputMode(output, mode string) error {
  return exec.Command("niri", "msg", "output", output, "mode", mode).Run()
}

// SetOutputTransform sets the transform (rotation) for the given output
func SetOutputTransform(output, transform string) error {
  return exec.Command("niri", "msg", "output", output, "transform", transform).Run()
}

// EnableOutput enables an output
func EnableOutput(output string) error {
  return exec.Command("niri", "msg", "output", output, "on").Run()
}

// DisableOutput disables an output
func DisableOutput(output string) error {
  return exec.Command("niri", "msg", "output", output, "off").Run()
}
