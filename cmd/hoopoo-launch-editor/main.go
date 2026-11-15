// hoopoo:summary=Launch the default editor as determined by $EDITOR (or nvim if missing).
// hoopoo:args=<path>

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  editor := os.Getenv("EDITOR")
  
  // Default to nvim if EDITOR not set
  if editor == "" || exec.Command("command", "-v", editor).Run() != nil {
    editor = "nvim"
  }

  // Check if editor is terminal-based
  terminalEditors := []string{"nvim", "vim", "nano", "micro", "hx", "helix", "fresh"}
  isTerminal := false
  for _, te := range terminalEditors {
    if editor == te {
      isTerminal = true
      break
    }
  }

  if isTerminal {
    // Launch in terminal
    args := []string{"hoopoo", "launch", "tui", editor}
    args = append(args, os.Args[1:]...)
    
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Run()
  } else {
    // Launch GUI editor directly
    args := []string{editor}
    args = append(args, os.Args[1:]...)
    
    cmd := exec.Command(args[0], args[1:]...)
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Start()
  }
}
