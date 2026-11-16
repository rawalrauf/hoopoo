// hoopoo:summary=Restart an application by killing it and relaunching.
// hoopoo:args=<application-name> [application-args...]

package main

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo restart app <application-name> [application-args...]")
    os.Exit(1)
  }

  appName := os.Args[1]

  // Kill existing instances
  exec.Command("pkill", "-x", appName).Run()

  // Relaunch application
  cmd := exec.Command(os.Args[1], os.Args[2:]...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
