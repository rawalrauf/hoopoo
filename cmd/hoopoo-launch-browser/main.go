// hoopoo:summary=Launch the default browser as determined by xdg-settings.
// hoopoo:args=[url]

package main

import (
  "os"
  "os/exec"
  "syscall"

  "hoopoo/pkg/browser"
)

func main() {
  // Get default browser
  defaultBrowser, err := browser.GetDefault()
  if err != nil {
    os.Exit(1)
  }

  // Get browser executable
  browserExec, err := browser.GetExecutable(defaultBrowser)
  if err != nil {
    os.Exit(1)
  }

  // Detect private flag
  privateFlag := browser.DetectPrivateFlag(browserExec)

  // Replace --private in args with browser-specific flag
  args := make([]string, 0, len(os.Args)-1)
  for _, arg := range os.Args[1:] {
    if arg == "--private" {
      args = append(args, privateFlag)
    } else {
      args = append(args, arg)
    }
  }

  // Launch browser
  cmd := exec.Command(browserExec, args...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
