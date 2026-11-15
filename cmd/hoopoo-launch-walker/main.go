// hoopoo:summary=Launch Walker and ensure its Elephant data provider is running

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Check and start both services in parallel
  go func() {
    if exec.Command("pgrep", "-x", "elephant").Run() != nil {
      cmd := exec.Command("elephant")
      cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
      cmd.Start()
    }
  }()

  go func() {
    if exec.Command("pgrep", "-f", "walker --gapplication-service").Run() != nil {
      cmd := exec.Command("walker", "--gapplication-service")
      cmd.Env = append(os.Environ(), "GSK_RENDERER=cairo")
      cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
      cmd.Start()
    }
  }()

  // Launch walker immediately
  args := []string{"--width", "644", "--maxheight", "345", "--minheight", "345"}
  args = append(args, os.Args[1:]...)
  
  cmd := exec.Command("walker", args...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()
}
