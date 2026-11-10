// hoopoo:summary=Print when system packages were last upgraded

package main

import (
  "fmt"
  "os"
  "os/exec"
  "regexp"
  "strings"
)

func main() {
  data, err := os.ReadFile("/var/log/pacman.log")
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading pacman log: %v\n", err)
    os.Exit(1)
  }
  
  lines := strings.Split(string(data), "\n")
  var lastUpgrade string
  
  for _, line := range lines {
    if strings.Contains(line, "upgraded") {
      lastUpgrade = line
    }
  }
  
  if lastUpgrade == "" {
    fmt.Println("No upgrade records found")
    os.Exit(1)
  }
  
  re := regexp.MustCompile(`\[([^\]]+)\]`)
  matches := re.FindStringSubmatch(lastUpgrade)
  
  if len(matches) < 2 {
    fmt.Fprintf(os.Stderr, "Could not parse timestamp from log\n")
    os.Exit(1)
  }
  
  timestamp := matches[1]
  cmd := exec.Command("date", "-d", timestamp, "+%A, %B %d %Y at %H:%M")
  output, err := cmd.Output()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error formatting date: %v\n", err)
    os.Exit(1)
  }
  
  fmt.Print(string(output))
}
