// hoopoo:summary=Reset the sudo lockout/faillock for the current user.

package main

import (
  "fmt"
  "os"
  "os/exec"
)

func main() {
  user := os.Getenv("USER")
  if user == "" {
    fmt.Fprintf(os.Stderr, "USER environment variable not set\n")
    os.Exit(1)
  }
  
  cmd := exec.Command("su", "-c", fmt.Sprintf("faillock --reset --user %s", user))
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  
  err := cmd.Run()
  if err != nil {
    os.Exit(1)
  }
}
