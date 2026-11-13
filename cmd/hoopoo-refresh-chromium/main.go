// hoopoo:summary=Refresh the ~/.config/chromium-flags.conf file from the Hoopoo defaults.

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

func main() {
  home, err := os.UserHomeDir()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
    os.Exit(1)
  }

  configFile := filepath.Join(home, ".config", "chromium-flags.conf")
  installGoogleAccounts := false

  // Check if google accounts were installed
  if fileContent, err := os.ReadFile(configFile); err == nil {
    content := string(fileContent)
    if strings.Contains(content, "--oauth2-client-id") && strings.Contains(content, "--oauth2-client-secret") {
      installGoogleAccounts = true
    }
  }

  // Refresh the Chromium configuration
  if err := exec.Command("hoopoo", "refresh", "config", "chromium-flags.conf").Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error refreshing chromium config: %v\n", err)
    os.Exit(1)
  }

  // Install/refresh the native messaging host for the yt-dlp extension
  if err := exec.Command("hoopoo", "install", "chromium", "ytdlp").Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error installing chromium ytdlp: %v\n", err)
    os.Exit(1)
  }

  // Re-install Google accounts if previously configured
  if installGoogleAccounts {
    if err := exec.Command("hoopoo", "install", "chromium", "google", "account").Run(); err != nil {
      fmt.Fprintf(os.Stderr, "Error installing chromium google account: %v\n", err)
      os.Exit(1)
    }
  }
}
