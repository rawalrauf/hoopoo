// Font detection and management

package font

import (
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "strings"
)

// Current returns the currently configured monospace font
func Current() string {
  // Read from waybar config (primary source)
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return ""
  }
  
  waybarStyle := filepath.Join(homeDir, ".config", "waybar", "style.css")
  content, err := os.ReadFile(waybarStyle)
  if err != nil {
    return ""
  }
  
  // Extract font-family: "FontName" or font-family: FontName
  re := regexp.MustCompile(`font-family:\s*["']?([^;"']+)["']?`)
  matches := re.FindStringSubmatch(string(content))
  if len(matches) > 1 {
    return strings.TrimSpace(matches[1])
  }
  
  return ""
}

// List returns all available monospace fonts
func List() ([]string, error) {
  cmd := exec.Command("fc-list", ":spacing=100", "-f", "%{family[0]}\n")
  output, err := cmd.Output()
  if err != nil {
    return nil, err
  }
  
  lines := strings.Split(string(output), "\n")
  var fonts []string
  seen := make(map[string]bool)
  
  for _, line := range lines {
    line = strings.TrimSpace(line)
    if line == "" || seen[line] {
      continue
    }
    
    // Skip emoji/special fonts
    lower := strings.ToLower(line)
    if strings.Contains(lower, "emoji") ||
       strings.Contains(lower, "signwriting") ||
       strings.Contains(lower, "hoopoo") {
      continue
    }
    
    fonts = append(fonts, line)
    seen[line] = true
  }
  
  return fonts, nil
}

// Type returns font category for height calculation
// Returns: "monospace", "nerd", or "default"
func Type() string {
  current := Current()
  lower := strings.ToLower(current)
  
  if strings.Contains(lower, "nerd") {
    return "nerd"
  }
  
  // Add more font type detection here
  return "monospace"
}
