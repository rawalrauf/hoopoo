// Font detection and management

package font

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
)

// StateFile returns the path to the font state file
func StateFile() string {
  return filepath.Join(paths.State(), "font")
}

// Current returns the currently configured monospace font
func Current() string {
  content, err := os.ReadFile(StateFile())
  if err != nil {
    // Fallback: try to detect from waybar if state file doesn't exist
    return detectFromWaybar()
  }
  return strings.TrimSpace(string(content))
}

// detectFromWaybar reads font from waybar config (fallback only)
func detectFromWaybar() string {
  waybarStyle := filepath.Join(paths.UserConfig(), "waybar", "style.css")
  content, err := os.ReadFile(waybarStyle)
  if err != nil {
    return ""
  }
  
  // Simple extraction: font-family: 'Font Name';
  lines := strings.Split(string(content), "\n")
  for _, line := range lines {
    if strings.Contains(line, "font-family:") {
      // Extract between quotes or after colon
      line = strings.TrimSpace(line)
      if idx := strings.Index(line, "'"); idx != -1 {
        end := strings.Index(line[idx+1:], "'")
        if end != -1 {
          return line[idx+1 : idx+1+end]
        }
      }
    }
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

// Set changes the system monospace font and updates all config files
func Set(fontName string) error {
  if fontName == "" {
    return fmt.Errorf("font name cannot be empty")
  }

  // Verify font exists
  cmd := exec.Command("fc-list")
  output, err := cmd.Output()
  if err != nil {
    return fmt.Errorf("failed to list fonts: %w", err)
  }
  
  if !strings.Contains(strings.ToLower(string(output)), strings.ToLower(fontName)) {
    return fmt.Errorf("font '%s' not found", fontName)
  }

  // Write to state file (source of truth)
  stateFile := StateFile()
  if err := os.MkdirAll(filepath.Dir(stateFile), 0755); err != nil {
    return fmt.Errorf("failed to create state directory: %w", err)
  }
  if err := os.WriteFile(stateFile, []byte(fontName+"\n"), 0644); err != nil {
    return fmt.Errorf("failed to write state file: %w", err)
  }

  // Update config files
  updateAlacritty(fontName)
  updateKitty(fontName)
  updateGhostty(fontName)
  updateFoot(fontName)
  updateHyprlock(fontName)
  updateWaybar(fontName)
  updateSwayosd(fontName)
  updateFontconfig(fontName)

  return nil
}

// updateAlacritty updates alacritty config
func updateAlacritty(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "alacritty", "alacritty.toml")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.Contains(line, `family = "`) {
        lines[i] = `family = "` + fontName + `"`
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateKitty updates kitty config and reloads
func updateKitty(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "kitty", "kitty.conf")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.HasPrefix(strings.TrimSpace(line), "font_family ") {
        lines[i] = "font_family " + fontName
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
    process.SendSignal("USR1", "kitty")
  }
}

// updateGhostty updates ghostty config
func updateGhostty(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "ghostty", "config")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.Contains(line, `font-family = "`) {
        lines[i] = `font-family = "` + fontName + `"`
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateFoot updates foot config
func updateFoot(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "foot", "foot.ini")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    currentSize := "11" // default
    
    // Extract current size
    for _, line := range lines {
      if strings.HasPrefix(strings.TrimSpace(line), "font=") {
        if idx := strings.Index(line, ":size="); idx != -1 {
          sizeStr := line[idx+6:]
          if end := strings.IndexAny(sizeStr, " \t\n"); end != -1 {
            currentSize = sizeStr[:end]
          } else {
            currentSize = sizeStr
          }
        }
      }
    }
    
    // Update font line
    for i, line := range lines {
      if strings.HasPrefix(strings.TrimSpace(line), "font=") {
        lines[i] = "font=" + fontName + ":size=" + currentSize
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateHyprlock updates hyprlock config
func updateHyprlock(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "hypr", "hyprlock.conf")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.Contains(line, "font_family =") {
        lines[i] = "    font_family = " + fontName
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateWaybar updates waybar style.css
func updateWaybar(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "waybar", "style.css")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.Contains(line, "font-family:") {
        lines[i] = "  font-family: '" + fontName + "';"
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateSwayosd updates swayosd style.css
func updateSwayosd(fontName string) {
  configPath := filepath.Join(paths.UserConfig(), "swayosd", "style.css")
  if content, err := os.ReadFile(configPath); err == nil {
    lines := strings.Split(string(content), "\n")
    for i, line := range lines {
      if strings.Contains(line, "font-family:") {
        lines[i] = "  font-family: '" + fontName + "';"
      }
    }
    os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
  }
}

// updateFontconfig updates fontconfig XML
func updateFontconfig(fontName string) {
  exec.Command("xmlstarlet", "ed", "-L",
    "-u", "//match[@target=\"pattern\"][test/string=\"monospace\"]/edit[@name=\"family\"]/string",
    "-v", fontName,
    filepath.Join(paths.UserConfig(), "fontconfig", "fonts.conf")).Run()
}
