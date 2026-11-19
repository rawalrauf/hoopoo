// hoopoo:summary=Display niri keybindings defined in your configuration using walker for an interactive search menu
// hoopoo:group=menu
// hoopoo:name=keybindings
// hoopoo:args=[--print|-p]

package main

import (
  "bufio"
  "crypto/sha256"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "sort"
  "strings"

  "hoopoo/pkg/system/paths"
)

type binding struct {
  display string
  desc    string
  action  string
  prio    int
}

var bindingFiles = []string{
  "default/niri/bindings/utilities.kdl",
  "default/niri/bindings/applications.kdl",
  "default/niri/bindings/tiling.kdl",
  "default/niri/bindings/media.kdl",
}

func main() {
  printOnly := len(os.Args) > 1 && (os.Args[1] == "--print" || os.Args[1] == "-p")

  records := getBindingRecords()
  
  if printOnly {
    for _, r := range records {
      fmt.Println(r.display)
    }
    return
  }

  // Show in walker
  var input strings.Builder
  for _, r := range records {
    input.WriteString(r.display + "\n")
  }

  cmd := exec.Command("hoopoo", "launch", "walker", "--dmenu", "-p", "Keybindings", "--width", "950")
  cmd.Stdin = strings.NewReader(input.String())
  output, err := cmd.Output()
  if err != nil {
    return
  }

  selection := strings.TrimSpace(string(output))
  if selection == "" {
    return
  }

  // Find and execute action
  for _, r := range records {
    if r.display == selection {
      dispatchBinding(r.action)
      break
    }
  }
}

func getBindingRecords() []binding {
  cacheDir := filepath.Join(paths.UserConfig(), "..", ".cache", "hoopoo")
  cacheKey := computeCacheKey()
  cacheFile := filepath.Join(cacheDir, "keybindings-"+cacheKey+".records")

  // Try cache first
  if data, err := os.ReadFile(cacheFile); err == nil {
    return parseRecords(string(data))
  }

  // Generate fresh
  records := generateBindings()
  
  // Save to cache
  os.MkdirAll(cacheDir, 0755)
  var buf strings.Builder
  for _, r := range records {
    buf.WriteString(fmt.Sprintf("%s\t%s\t%s\n", r.display, r.desc, r.action))
  }
  os.WriteFile(cacheFile, []byte(buf.String()), 0644)
  
  // Clean old cache files
  cleanOldCache(cacheDir, "keybindings-"+cacheKey+".records")

  return records
}

func computeCacheKey() string {
  h := sha256.New()
  h.Write([]byte("v1\n"))
  
  hoopooPath := paths.Hoopoo()
  for _, relPath := range bindingFiles {
    absPath := filepath.Join(hoopooPath, relPath)
    if info, err := os.Stat(absPath); err == nil {
      h.Write([]byte(fmt.Sprintf("%d %s\n", info.ModTime().Unix(), absPath)))
    }
  }
  
  // Also check user config
  userConfig := filepath.Join(paths.UserConfig(), "niri", "bindings.kdl")
  if info, err := os.Stat(userConfig); err == nil {
    h.Write([]byte(fmt.Sprintf("%d %s\n", info.ModTime().Unix(), userConfig)))
  }
  
  return fmt.Sprintf("%x", h.Sum(nil))
}

func generateBindings() []binding {
  var bindings []binding
  
  hoopooPath := paths.Hoopoo()
  
  // Parse all binding files
  for _, relPath := range bindingFiles {
    absPath := filepath.Join(hoopooPath, relPath)
    bindings = append(bindings, parseNiriBindings(absPath)...)
  }
  
  // Parse user config
  userConfig := filepath.Join(paths.UserConfig(), "niri", "bindings.kdl")
  bindings = append(bindings, parseNiriBindings(userConfig)...)
  
  // Add static bindings
  bindings = append(bindings, binding{
    display: fmt.Sprintf("%-37s → %s", "Alt + Shift + L", "Copy URL from Web App"),
    desc:    "Copy URL from Web App",
    action:  "",
    prio:    50,
  })
  bindings = append(bindings, binding{
    display: fmt.Sprintf("%-37s → %s", "Alt + Shift + D", "Download Video from Web App"),
    desc:    "Download Video from Web App",
    action:  "",
    prio:    50,
  })
  
  // Remove duplicates
  seen := make(map[string]bool)
  unique := []binding{}
  for _, b := range bindings {
    if !seen[b.display] {
      seen[b.display] = true
      unique = append(unique, b)
    }
  }
  
  // Prioritize and sort
  for i := range unique {
    unique[i].prio = getPriority(unique[i].desc)
  }
  
  sort.Slice(unique, func(i, j int) bool {
    if unique[i].prio != unique[j].prio {
      return unique[i].prio < unique[j].prio
    }
    return unique[i].display < unique[j].display
  })
  
  return unique
}

func parseNiriBindings(path string) []binding {
  file, err := os.Open(path)
  if err != nil {
    return nil
  }
  defer file.Close()
  
  var bindings []binding
  var skipNext bool
  
  // Regex to match: Mod+Shift+X { ... } hotkey-overlay-title="Description"
  re := regexp.MustCompile(`^\s*([A-Za-z0-9+]+)\s+.*hotkey-overlay-title="([^"]+)"`)
  actionRe := regexp.MustCompile(`\{([^}]+)\}`)
  
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    trimmed := strings.TrimSpace(line)
    
    // Skip line comments
    if strings.HasPrefix(trimmed, "//") {
      continue
    }
    
    // Check for /- slashdash (comments out next node)
    if strings.HasPrefix(trimmed, "/-") {
      skipNext = true
      continue
    }
    
    // If we're in skip mode and hit a node, skip it
    if skipNext {
      // Skip this line and any continuation until we hit closing brace at same level
      // For simplicity, if line has opening brace, we're skipping a block
      if strings.Contains(line, "{") {
        // Skip until matching }
        depth := strings.Count(line, "{") - strings.Count(line, "}")
        for depth > 0 && scanner.Scan() {
          nextLine := scanner.Text()
          depth += strings.Count(nextLine, "{") - strings.Count(nextLine, "}")
        }
      }
      skipNext = false
      continue
    }
    
    // Strip inline comments
    if idx := strings.Index(line, "//"); idx != -1 {
      line = line[:idx]
    }
    
    matches := re.FindStringSubmatch(line)
    if matches != nil {
      key := matches[1]
      desc := matches[2]
      
      // Format key: Mod+Shift+X → Mod + Shift + X
      key = strings.ReplaceAll(key, "+", " + ")
      key = strings.Join(strings.Fields(key), " ")
      
      // Extract action
      action := ""
      actionMatches := actionRe.FindStringSubmatch(line)
      if len(actionMatches) > 1 {
        action = strings.TrimSpace(strings.ReplaceAll(actionMatches[1], ";", ""))
      }
      
      display := fmt.Sprintf("%-37s → %s", key, desc)
      bindings = append(bindings, binding{
        display: display,
        desc:    desc,
        action:  action,
      })
    }
  }
  
  return bindings
}

func getPriority(desc string) int {
  lower := strings.ToLower(desc)
  
  if strings.Contains(lower, "terminal") { return 0 }
  if strings.Contains(lower, "application") { return 1 }
  if strings.Contains(lower, "hoopoo menu") { return 2 }
  if strings.Contains(lower, "apps menu") { return 3 }
  if strings.Contains(lower, "lock") { return 4 }
  if strings.Contains(lower, "emoji") { return 5 }
  if strings.Contains(lower, "clipboard") { return 6 }
  if strings.Contains(lower, "screenshot") { return 7 }
  if strings.Contains(lower, "close") { return 8 }
  if strings.Contains(lower, "overview") { return 9 }
  if strings.Contains(lower, "fullscreen") { return 10 }
  if strings.Contains(lower, "maximize") { return 11 }
  if strings.Contains(lower, "floating") { return 12 }
  if strings.Contains(lower, "opacity") { return 13 }
  if strings.Contains(lower, "focus") && strings.Contains(lower, "column") { return 20 }
  if strings.Contains(lower, "focus") && strings.Contains(lower, "window") { return 21 }
  if strings.Contains(lower, "focus") && strings.Contains(lower, "monitor") { return 22 }
  if strings.Contains(lower, "focus") && strings.Contains(lower, "workspace") { return 23 }
  if strings.Contains(lower, "move") && strings.Contains(lower, "column") { return 30 }
  if strings.Contains(lower, "move") && strings.Contains(lower, "window") { return 31 }
  if strings.Contains(lower, "move") && strings.Contains(lower, "workspace") { return 32 }
  if strings.Contains(lower, "width") { return 40 }
  if strings.Contains(lower, "height") { return 41 }
  if strings.Contains(lower, "volume") { return 60 }
  if strings.Contains(lower, "brightness") { return 61 }
  if strings.Contains(lower, "play") || strings.Contains(lower, "pause") || strings.Contains(lower, "media") { return 62 }
  if strings.Contains(lower, "quit") { return 90 }
  if strings.Contains(lower, "power") { return 91 }
  
  return 50
}

func parseRecords(data string) []binding {
  var records []binding
  lines := strings.Split(data, "\n")
  for _, line := range lines {
    if line == "" {
      continue
    }
    parts := strings.Split(line, "\t")
    if len(parts) >= 3 {
      records = append(records, binding{
        display: parts[0],
        desc:    parts[1],
        action:  parts[2],
      })
    }
  }
  return records
}

func cleanOldCache(cacheDir, keepFile string) {
  entries, err := os.ReadDir(cacheDir)
  if err != nil {
    return
  }
  
  for _, entry := range entries {
    if strings.HasPrefix(entry.Name(), "keybindings-") &&
       strings.HasSuffix(entry.Name(), ".records") &&
       entry.Name() != filepath.Base(keepFile) {
      os.Remove(filepath.Join(cacheDir, entry.Name()))
    }
  }
}

func dispatchBinding(action string) {
  if action == "" {
    return
  }
  
  // Built-in niri action (no spawn)
  if !strings.Contains(action, "spawn") {
    args := append([]string{"msg", "action"}, strings.Fields(action)...)
    exec.Command("niri", args...).Run()
    return
  }
  
  // spawn-sh or spawn command
  cmd := strings.TrimSpace(action)
  cmd = strings.TrimPrefix(cmd, "spawn-sh ")
  cmd = strings.TrimPrefix(cmd, "spawn ")
  cmd = strings.Trim(cmd, `"`)
  
  if cmd != "" {
    exec.Command("niri", "msg", "action", "spawn-sh", "--", cmd).Run()
  }
}
