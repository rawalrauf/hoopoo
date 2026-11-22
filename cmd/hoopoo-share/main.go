package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
)

// hoopoo:summary=Share clipboard, files, or folders with LocalSend
// hoopoo:group=share
// hoopoo:name=
// hoopoo:args=<clipboard|file|folder> [path...]
// hoopoo:examples=hoopoo share clipboard | hoopoo share file ~/Downloads/example.txt

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo share [clipboard|file|folder]")
    os.Exit(1)
  }
  
  mode := os.Args[1]
  var files []string
  
  switch mode {
  case "clipboard":
    // Get clipboard content and save to temp file
    cmd := exec.Command("wl-paste")
    content, err := cmd.Output()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to read clipboard: %v\n", err)
      os.Exit(1)
    }
    
    // Create temp file
    tmpFile, err := os.CreateTemp("", "clipboard-*.txt")
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to create temp file: %v\n", err)
      os.Exit(1)
    }
    
    if _, err := tmpFile.Write(content); err != nil {
      fmt.Fprintf(os.Stderr, "Failed to write temp file: %v\n", err)
      os.Exit(1)
    }
    tmpFile.Close()
    
    files = []string{tmpFile.Name()}
    
  case "file", "folder":
    if len(os.Args) > 2 {
      // Files provided as arguments
      files = os.Args[2:]
    } else {
      // Interactive file/folder selection using portal
      cmd := exec.Command("hoopoo", "file", "select")
      
      if mode == "folder" {
        cmd.Args = append(cmd.Args, "--title", "Select folder to share", "--directory")
      } else {
        cmd.Args = append(cmd.Args, "--title", "Select files to share", "--multiple")
      }
      
      output, err := cmd.Output()
      if err != nil {
        // User cancelled or error - exit silently
        os.Exit(0)
      }
      
      selected := strings.TrimSpace(string(output))
      if selected == "" {
        os.Exit(0)
      }
      
      // Split by newlines for multiple files
      files = strings.Split(selected, "\n")
    }
    
  default:
    fmt.Fprintf(os.Stderr, "Invalid mode: %s\nUse: clipboard, file, or folder\n", mode)
    os.Exit(1)
  }
  
  if len(files) == 0 {
    os.Exit(0)
  }
  
  // Run LocalSend via systemd-run (detached from terminal)
  args := []string{"--user", "--quiet", "--collect", "localsend", "--headless", "send"}
  args = append(args, files...)
  
  cmd := exec.Command("systemd-run", args...)
  if err := cmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to start LocalSend: %v\n", err)
    os.Exit(1)
  }
}
