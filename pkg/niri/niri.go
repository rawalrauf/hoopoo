// Niri compositor utilities

package niri

import (
  "bufio"
  "encoding/json"
  "fmt"
  "os/exec"
  "regexp"
  "strconv"
  "strings"
)

func Reload() error {
  return exec.Command("niri", "msg", "action", "load-config-file").Run()
}

// FocusWorkspace focuses a workspace by number
func FocusWorkspace(num int) error {
  return exec.Command("niri", "msg", "action", "focus-workspace", strconv.Itoa(num)).Run()
}

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

// OutputMode represents a display mode
type OutputMode struct {
  Width        int  `json:"width"`
  Height       int  `json:"height"`
  RefreshRate  int  `json:"refresh_rate"`
  IsPreferred  bool `json:"is_preferred"`
}

// LogicalOutput represents the logical output dimensions
type LogicalOutput struct {
  X         int     `json:"x"`
  Y         int     `json:"y"`
  Width     int     `json:"width"`
  Height    int     `json:"height"`
  Scale     float64 `json:"scale"`
  Transform string  `json:"transform"`
}

// Output represents a niri output (display)
type Output struct {
  Name         string        `json:"name"`
  Make         string        `json:"make"`
  Model        string        `json:"model"`
  Modes        []OutputMode  `json:"modes"`
  CurrentMode  int           `json:"current_mode"`
  Logical      LogicalOutput `json:"logical"`
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

// Window represents a niri window
type Window struct {
  ID    int    `json:"id"`
  PID   int    `json:"pid"`
  AppID string `json:"app_id"`
  Title string `json:"title"`
}

// GetWindows returns all open windows
func GetWindows() ([]Window, error) {
  out, err := exec.Command("niri", "msg", "-j", "windows").Output()
  if err != nil {
    return nil, err
  }

  var windows []Window
  if err := json.Unmarshal(out, &windows); err != nil {
    return nil, err
  }

  return windows, nil
}

// GetFocusedWindow returns the currently focused window
func GetFocusedWindow() (*Window, error) {
  out, err := exec.Command("niri", "msg", "--json", "focused-window").Output()
  if err != nil {
    return nil, err
  }

  var window Window
  if err := json.Unmarshal(out, &window); err != nil {
    return nil, err
  }

  return &window, nil
}

// CloseWindow closes a window by ID
func CloseWindow(id int) error {
  return exec.Command("niri", "msg", "action", "close-window", "--id", strconv.Itoa(id)).Run()
}

// FocusWindow focuses a window by ID
func FocusWindow(id int) error {
  return exec.Command("niri", "msg", "action", "focus-window", "--id", strconv.Itoa(id)).Run()
}

// SetOutputScale sets the scale for an output
func SetOutputScale(output string, scale float64) error {
  return exec.Command("niri", "msg", "output", output, "scale", fmt.Sprintf("%g", scale)).Run()
}

// EnableOutput enables an output
func EnableOutput(output string) error {
  return exec.Command("niri", "msg", "output", output, "on").Run()
}

// DisableOutput disables an output
func DisableOutput(output string) error {
  return exec.Command("niri", "msg", "output", output, "off").Run()
}

// EventStream returns a scanner for the niri event stream
func EventStream() (*bufio.Scanner, *exec.Cmd, error) {
  cmd := exec.Command("niri", "msg", "--json", "event-stream")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return nil, nil, err
  }

  if err := cmd.Start(); err != nil {
    return nil, nil, err
  }

  scanner := bufio.NewScanner(stdout)
  return scanner, cmd, nil
}
