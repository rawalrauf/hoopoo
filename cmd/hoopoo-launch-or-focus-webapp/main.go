// hoopoo:summary=Construct app-id and launch or focus a web app in Niri
// hoopoo:args=<url> [flags...]

package main

import (
  "encoding/json"
  "fmt"
  "os"
  "os/exec"
  "regexp"
  "strings"
  "syscall"

  "hoopoo/pkg/browser"
)

type Window struct {
  ID    int    `json:"id"`
  AppID string `json:"app_id"`
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch or focus webapp <url> [flags...]")
    os.Exit(1)
  }

  url := os.Args[1]

  // Extract clean domain
  cleanDomain := browser.ExtractDomain(url)

  // Check if window for this domain is already running
  cmd := exec.Command("niri", "msg", "-j", "windows")
  output, err := cmd.Output()

  var runningAppID string
  if err == nil {
    var windows []Window
    if json.Unmarshal(output, &windows) == nil {
      // Match pattern: *-<domain>__-Default
      pattern := regexp.MustCompile(fmt.Sprintf(`^.*-%s__-Default$`, regexp.QuoteMeta(cleanDomain)))
      for _, win := range windows {
        if pattern.MatchString(win.AppID) {
          runningAppID = win.AppID
          break
        }
      }
    }
  }

  var appID string
  if runningAppID != "" {
    appID = runningAppID
  } else {
    // Determine browser prefix
    defaultBrowser, _ := browser.GetDefault()
    prefix := browser.GetBrowserPrefix(defaultBrowser)
    appID = browser.BuildWebAppID(prefix, cleanDomain)
  }

  // Build launch command
  launchArgs := []string{"hoopoo", "launch", "webapp"}
  launchArgs = append(launchArgs, os.Args[1:]...)
  launchCommand := strings.Join(launchArgs, " ")

  // Call hoopoo launch or focus
  finalCmd := exec.Command(
    "hoopoo",
    "launch",
    "or",
    "focus",
    appID,
    launchCommand,
  )

  finalCmd.Stdin = os.Stdin
  finalCmd.Stdout = os.Stdout
  finalCmd.Stderr = os.Stderr
  finalCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

  finalCmd.Run()
}
