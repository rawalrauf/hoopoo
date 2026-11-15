// hoopoo:summary=Launch a URL as a web app in the default supported browser
// hoopoo:args=<url>

package main

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"

  "hoopoo/pkg/browser"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch webapp <url>")
    os.Exit(1)
  }

  // Get default browser
  defaultBrowser, err := browser.GetDefault()
  if err != nil {
    os.Exit(1)
  }

  // If not Chromium-based, find one that is
  browserDesktopFile := defaultBrowser
  if !browser.IsChromiumBased(defaultBrowser) {
    chromiumBrowser, err := browser.FindChromiumBased()
    if err != nil {
      os.Exit(1)
    }
    browserDesktopFile = chromiumBrowser
  }

  // Get browser executable
  browserBin, err := browser.GetExecutable(browserDesktopFile)
  if err != nil {
    os.Exit(1)
  }

  // Launch webapp
  args := []string{fmt.Sprintf("--app=%s", os.Args[1])}
  args = append(args, os.Args[2:]...)

  cmd := exec.Command(browserBin, args...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
