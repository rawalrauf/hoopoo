// hoopoo:summary=Install the native messaging host for the yt-dlp Chromium extension

package main

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "hoopoo/pkg/system/paths"
)

func main() {
  home, err := os.UserHomeDir()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
    os.Exit(1)
  }

  hostName := "com.hoopoo.ytdlp"
  hostPath := filepath.Join(paths.Hoopoo(), "bin", "hoopoo-chromium-ytdlp-host")
  templatePath := filepath.Join(paths.Hoopoo(), "default", "chromium", "native-messaging-hosts", hostName+".json")

  browserDirs := []string{
    filepath.Join(home, ".config", "chromium"),
    filepath.Join(home, ".config", "google-chrome"),
    filepath.Join(home, ".config", "google-chrome-beta"),
    filepath.Join(home, ".config", "google-chrome-unstable"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Browser"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Browser-Beta"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Browser-Nightly"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Origin"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Origin-Beta"),
    filepath.Join(home, ".config", "BraveSoftware", "Brave-Origin-Nightly"),
    filepath.Join(home, ".config", "microsoft-edge"),
    filepath.Join(home, ".config", "microsoft-edge-dev"),
  }

  templateContent, err := os.ReadFile(templatePath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading template: %v\n", err)
    os.Exit(1)
  }

  manifest := strings.ReplaceAll(string(templateContent), "__HOST_PATH__", hostPath)

  successCount := 0
  for _, dir := range browserDirs {
    targetDir := filepath.Join(dir, "NativeMessagingHosts")
    if err := os.MkdirAll(targetDir, 0755); err != nil {
      fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", targetDir, err)
      continue
    }

    targetFile := filepath.Join(targetDir, hostName+".json")
    if err := os.WriteFile(targetFile, []byte(manifest), 0644); err != nil {
      fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", targetFile, err)
      continue
    }
    successCount++
  }

  if successCount == 0 {
    fmt.Fprintln(os.Stderr, "Failed to install native messaging host to any browser")
    os.Exit(1)
  }
}
