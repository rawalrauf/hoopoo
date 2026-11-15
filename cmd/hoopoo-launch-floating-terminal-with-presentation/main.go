// hoopoo:summary=Launch a floating terminal with the Hoopoo presentation wrapper
// hoopoo:args=<command>

package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "syscall"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch floating terminal with presentation <command>")
    os.Exit(1)
  }

  // Join all args into a single command string
  cmdStr := strings.Join(os.Args[1:], " ")

  // Build the bash wrapper command
  // Shows logo, runs command, shows done (unless exit code 130 = Ctrl+C)
  bashCmd := fmt.Sprintf("hoopoo show logo; %s; if (( $? != 130 )); then hoopoo show done; fi", cmdStr)

  // Launch terminal with presentation wrapper
  termArgs := []string{
    "--app-id=org.hoopoo.terminal",
    "--title=Hoopoo",
    "-e",
    "bash",
    "-c",
    bashCmd,
  }

  cmd := exec.Command("xdg-terminal-exec", termArgs...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
