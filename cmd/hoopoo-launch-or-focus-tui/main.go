// hoopoo:summary=Launch or focus a TUI window matching base, float, or tile mode
// hoopoo:args=[--app-id=ID] [--title=TITLE] <command> [args...]

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "strings"
  "syscall"

  "hoopoo/pkg/niri"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch or focus tui [--app-id=ID] [--title=TITLE] <command> [args...]")
    os.Exit(1)
  }

  args := os.Args[1:]
  explicitID := ""
  explicitTitle := ""

  // Parse flags
  for len(args) > 0 && strings.HasPrefix(args[0], "--") {
    if strings.HasPrefix(args[0], "--app-id=") {
      explicitID = strings.TrimPrefix(args[0], "--app-id=")
      args = args[1:]
    } else if strings.HasPrefix(args[0], "--title=") {
      explicitTitle = strings.TrimPrefix(args[0], "--title=")
      args = args[1:]
    } else {
      break
    }
  }

  if len(args) == 0 {
    fmt.Fprintln(os.Stderr, "Error: command is required")
    os.Exit(1)
  }

  // Extract command name from first argument
  cmdName := filepath.Base(args[0])

  // Query Niri for running window matching org.hoopoo.<CMD_NAME>, org.hoopoo.float.<CMD_NAME>, or org.hoopoo.tile.<CMD_NAME>
  var runningID string
  windows, err := niri.GetWindows()
  if err == nil {
    pattern := regexp.MustCompile(fmt.Sprintf(`^org\.hoopoo\.(float\.|tile\.)?%s$`, regexp.QuoteMeta(cmdName)))
    for _, win := range windows {
      if pattern.MatchString(win.AppID) {
        runningID = win.AppID
        break
      }
    }
  }

  // Determine target app ID
  var targetAppID string
  if runningID != "" {
    targetAppID = runningID
  } else if explicitID != "" {
    targetAppID = explicitID
  } else {
    targetAppID = fmt.Sprintf("org.hoopoo.%s", cmdName)
  }

  // Determine title
  title := explicitTitle
  if title == "" {
    title = cmdName
  }

  // Build launch command arguments (shell-quoted since they get run via sh -c)
  launchArgs := []string{
    "hoopoo",
    "launch",
    "tui",
    fmt.Sprintf("--app-id=%s", shellQuote(targetAppID)),
    fmt.Sprintf("--title=%s", shellQuote(title)),
  }
  for _, arg := range args {
    launchArgs = append(launchArgs, shellQuote(arg))
  }
  launchCommand := strings.Join(launchArgs, " ")

  // Call hoopoo launch or focus
  finalCmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    targetAppID,
    launchCommand,
  )
  
  finalCmd.Stdin = os.Stdin
  finalCmd.Stdout = os.Stdout
  finalCmd.Stderr = os.Stderr
  finalCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  finalCmd.Run()
}

// shellQuote wraps a string in single quotes for safe use in sh -c
func shellQuote(s string) string {
  return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

