// hoopoo:summary=Copies the named config from ~/.local/share/hoopoo/config/X/Y/Z -> ~/.config/X/Y/Z.
// hoopoo:args=<config-path>

package main

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "time"

  "hoopoo/pkg/system/paths"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, `  Usage: %s [config_file]

  Must provide a file path from the .config directory to be refreshed.
  To copy ~/.local/share/hoopoo/config/hypr/hyprlock.conf to ~/.config/hypr/hyprlock.conf

    %s hypr/hyprlock.conf
`, os.Args[0], os.Args[0])
    os.Exit(1)
  }

  configFile := os.Args[1]
  
  userConfigFile := filepath.Join(paths.UserConfig(), configFile)
  defaultConfigFile := filepath.Join(paths.Hoopoo(), "config", configFile)
  backupConfigFile := fmt.Sprintf("%s.bak.%d", userConfigFile, time.Now().Unix())

  // Create parent directory
  os.MkdirAll(filepath.Dir(userConfigFile), 0755)

  // Check if user config exists
  if _, err := os.Stat(userConfigFile); err == nil {
    // Create preliminary backup
    userContent, _ := os.ReadFile(userConfigFile)
    os.WriteFile(backupConfigFile, userContent, 0644)

    // Replace config with new default
    defaultContent, _ := os.ReadFile(defaultConfigFile)
    os.WriteFile(userConfigFile, defaultContent, 0644)

    // Compare and delete/inform accordingly
    newContent, _ := os.ReadFile(userConfigFile)
    if bytes.Equal(newContent, userContent) {
      os.Remove(backupConfigFile)
    } else {
      fmt.Printf("\033[31mReplaced %s with new Hoopoo default.\nSaved backup as %s.\n\n\033[32mChanges:\033[0m\n", userConfigFile, backupConfigFile)
      
      // Show diff
      cmd := exec.Command("diff", userConfigFile, backupConfigFile)
      cmd.Stdout = os.Stdout
      cmd.Stderr = os.Stderr
      cmd.Run()
    }
  } else {
    // Config file did not exist already
    defaultContent, _ := os.ReadFile(defaultConfigFile)
    os.WriteFile(userConfigFile, defaultContent, 0644)
  }
}
