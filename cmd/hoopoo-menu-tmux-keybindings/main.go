// hoopoo:summary=Display annotated Tmux keybindings using an interactive search menu
// hoopoo:group=menu
// hoopoo:name=tmux-keybindings
// hoopoo:args=[--print|-p] [--config <path>]

package main

import (
  "bufio"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "hoopoo/pkg/system/paths"
)

func main() {
  printOnly := false
  configFile := os.Getenv("TMUX_CONF")
  if configFile == "" {
    configFile = filepath.Join(paths.UserConfig(), "tmux", "tmux.conf")
  }

  // Parse args
  for i := 1; i < len(os.Args); i++ {
    switch os.Args[i] {
    case "--print", "-p":
      printOnly = true
    case "--config":
      if i+1 < len(os.Args) {
        configFile = os.Args[i+1]
        i++
      }
    default:
      configFile = os.Args[i]
    }
  }

  // Check if config exists
  if _, err := os.Stat(configFile); err != nil {
    defaultConfig := filepath.Join(paths.Hoopoo(), "config", "tmux", "tmux.conf")
    if _, err := os.Stat(defaultConfig); err == nil {
      configFile = defaultConfig
    } else {
      fmt.Fprintf(os.Stderr, "Tmux config not found: %s\n", configFile)
      os.Exit(1)
    }
  }

  lines, err := extractKeybindings(configFile)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error extracting keybindings: %v\n", err)
    os.Exit(1)
  }

  if printOnly {
    for _, line := range lines {
      fmt.Println(line)
    }
    return
  }

  // Show in walker
  input := strings.Join(lines, "\n")
  cmd := exec.Command("hoopoo", "launch", "walker", "--dmenu", "-p", "Tmux keybindings", "--width", "950")
  cmd.Stdin = strings.NewReader(input)
  cmd.Stdout = nil // Discard output
  cmd.Run()
}

func extractKeybindings(configFile string) ([]string, error) {
  // Create temporary socket
  socketDir, err := os.MkdirTemp("", "tmux-socket-*")
  if err != nil {
    return nil, err
  }
  defer os.RemoveAll(socketDir)

  socket := filepath.Join(socketDir, "tmux")

  // Start tmux server with config
  cmd := exec.Command("tmux",
    "-S", socket,
    "-f", "/dev/null",
    "start-server", ";",
    "unbind-key", "-a", "-T", "prefix", ";",
    "unbind-key", "-a", "-T", "root", ";",
    "unbind-key", "-a", "-T", "copy-mode", ";",
    "unbind-key", "-a", "-T", "copy-mode-vi", ";",
    "source-file", configFile, ";",
    "show-options", "-g", "prefix", ";",
    "show-options", "-g", "prefix2", ";",
    "list-keys", "-aN", "-T", "prefix", "-P", "prefix", ";",
    "list-keys", "-aN", "-T", "root", "-P", "root", ";",
    "list-keys", "-aN", "-T", "copy-mode-vi", "-P", "copy-mode-vi",
  )

  output, err := cmd.Output()
  if err != nil {
    return nil, err
  }

  // Kill server
  exec.Command("tmux", "-S", socket, "kill-server").Run()

  return formatKeybindings(string(output)), nil
}

func formatKeybindings(output string) []string {
  var lines []string
  var prefix, prefix2 string
  headerPrinted := false

  scanner := bufio.NewScanner(strings.NewReader(output))
  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    if len(fields) < 2 {
      continue
    }

    // Parse prefix options
    if fields[0] == "prefix" && len(fields) == 2 {
      prefix = fields[1]
      continue
    }
    if fields[0] == "prefix2" && len(fields) == 2 {
      prefix2 = fields[1]
      continue
    }

    // Print header once
    if !headerPrinted {
      prefixDesc := formatKey(prefix)
      if prefix2 != "" && prefix2 != "None" {
        prefixDesc = prefixDesc + " / " + formatKey(prefix2)
      }
      lines = append(lines, fmt.Sprintf("%-32s → %s", "PREFIX", prefixDesc))
      headerPrinted = true
    }

    // Parse keybinding
    table := fields[0]
    key := fields[1]
    note := strings.Join(fields[2:], " ")

    var combo string
    if table == "root" {
      combo = formatKey(key)
    } else {
      combo = formatTable(table) + " + " + formatKey(key)
    }

    lines = append(lines, fmt.Sprintf("%-32s → %s", combo, note))
  }

  return lines
}

func formatKey(key string) string {
  // Remove backslashes
  key = strings.ReplaceAll(key, "\\", "")

  // Split on dashes
  parts := strings.Split(key, "-")
  formatted := make([]string, 0, len(parts))

  for i, part := range parts {
    isModifier := i < len(parts)-1
    formatted = append(formatted, formatKeyPart(part, isModifier, len(parts)))
  }

  return strings.Join(formatted, " + ")
}

func formatKeyPart(part string, isModifier bool, partCount int) string {
  // Handle modifiers
  if isModifier {
    switch part {
    case "C":
      return "CTRL"
    case "M":
      return "ALT"
    case "S":
      return "SHIFT"
    }
  }

  // Handle special keys
  switch part {
  case "Space":
    return "SPACE"
  case "BSpace":
    return "BACKSPACE"
  case "BTab":
    return "SHIFT + TAB"
  case "PPage":
    return "PAGE UP"
  case "NPage":
    return "PAGE DOWN"
  case "DC":
    return "DELETE"
  case "IC":
    return "INSERT"
  }

  // Single lowercase letter without modifiers
  if partCount == 1 && len(part) == 1 && part >= "a" && part <= "z" {
    return part
  }

  return strings.ToUpper(part)
}

func formatTable(table string) string {
  switch table {
  case "copy-mode-vi", "copy-mode":
    return "COPY MODE"
  case "prefix":
    return "PREFIX"
  default:
    return strings.ToUpper(strings.ReplaceAll(table, "-", " "))
  }
}
