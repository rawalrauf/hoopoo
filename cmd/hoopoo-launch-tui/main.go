// hoopoo:summary=Launch a TUI command in default terminal
// hoopoo:args=[--app-id=ID] [--title=TITLE] <command> [args...]

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "syscall"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch tui [--app-id=ID] [--title=TITLE] <command> [args...]")
    os.Exit(1)
  }

  args := os.Args[1:]
  appID := ""
  title := ""

  // Parse flags
  for len(args) > 0 && strings.HasPrefix(args[0], "--") {
    if strings.HasPrefix(args[0], "--app-id=") {
      appID = strings.TrimPrefix(args[0], "--app-id=")
      args = args[1:]
    } else if strings.HasPrefix(args[0], "--title=") {
      title = strings.TrimPrefix(args[0], "--title=")
      args = args[1:]
    } else {
      break
    }
  }

  if len(args) == 0 {
    fmt.Fprintln(os.Stderr, "Error: command is required")
    os.Exit(1)
  }

  // Generate defaults from command name
  cmdName := filepath.Base(args[0])
  if appID == "" {
    appID = fmt.Sprintf("org.hoopoo.%s", cmdName)
  }
  if title == "" {
    title = cmdName
  }

  // Build xdg-terminal-exec command
  termArgs := []string{
    fmt.Sprintf("--app-id=%s", appID),
    fmt.Sprintf("--title=%s", title),
    "-e",
  }
  termArgs = append(termArgs, args...)

  // Launch with setsid (native Go way)
  cmd := exec.Command("xdg-terminal-exec", termArgs...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
