// hoopoo:summary=Native messaging host: download the URL sent by the yt-dlp Chromium extension
// hoopoo:hidden=true

package main

import (
  "encoding/binary"
  "encoding/json"
  "fmt"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "strings"
  "syscall"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/ytdlp"
  "hoopoo/pkg/media"
)

type nativeMessage struct {
  URL string `json:"url"`
}

func main() {
  if len(os.Args) > 1 && os.Args[1] == "--download" {
    if len(os.Args) < 3 {
      fmt.Fprintln(os.Stderr, "Error: --download requires URL argument")
      os.Exit(1)
    }
    downloadURL(os.Args[2])
    return
  }

  handleNativeMessaging()
}

func handleNativeMessaging() {
  var length uint32
  if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
    return
  }

  if length == 0 {
    return
  }

  payload := make([]byte, length)
  if _, err := io.ReadFull(os.Stdin, payload); err != nil {
    return
  }

  var msg nativeMessage
  if err := json.Unmarshal(payload, &msg); err != nil {
    return
  }

  if msg.URL == "" || !strings.HasPrefix(msg.URL, "http://") && !strings.HasPrefix(msg.URL, "https://") {
    return
  }

  response := []byte("{}")
  responseLen := uint32(len(response))
  binary.Write(os.Stdout, binary.LittleEndian, responseLen)
  os.Stdout.Write(response)

  cmd := exec.Command(os.Args[0], "--download", msg.URL)
  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Setsid: true,
  }
  cmd.Start()
}

func downloadURL(url string) {
  downloadDir := os.Getenv("HOOPOO_YTDLP_DIR")
  if downloadDir == "" {
    home, _ := os.UserHomeDir()
    downloadDir = filepath.Join(home, "Videos")
  }

  iconPath := filepath.Join(paths.Hoopoo(), "default", "chromium", "extensions", "yt-dlp", "icon.png")

  notifCmd := exec.Command("hoopoo", "notification", "send",
    "--print-id",
    "-i", iconPath,
    "-t", "300000",
    "Processing video",
    url,
  )
  output, err := notifCmd.Output()
  notifID := strings.TrimSpace(string(output))
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error sending initial notification: %v\n", err)
    notifID = ""
  }

  opts := ytdlp.DownloadOptions{
    URL:         url,
    DownloadDir: downloadDir,
    OnProgress: func(percent int, title string) {
      if notifID != "" {
        updateCmd := exec.Command("hoopoo", "notification", "send",
          "--replace-id="+notifID,
          "-i", iconPath,
          "--progress="+strconv.Itoa(percent),
          title,
        )
        if err := updateCmd.Run(); err != nil {
          notifID = ""
        }
      }
    },
  }

  result, err := ytdlp.Download(opts)
  
  // Generate thumbnail while download just completed (before dismissing progress)
  var thumbnailPath string
  if err == nil && result != nil && result.FilePath != "" {
    thumbnailPath, _ = media.GenerateVideoThumbnail(result.FilePath)
  }
  
  // Dismiss progress notification
  if notifID != "" {
    exec.Command("hoopoo", "notification", "dismiss", notifID).Run()
  }
  
  // Handle errors
  if err != nil {
    // Always clean up all partial downloads (safer and simpler)
    cleanupPartialDownloads(downloadDir)
    showErrorNotification(url, err)
    return
  }
  
  // Show completion with pre-generated thumbnail
  if result != nil && result.FilePath != "" {
    showCompletionNotification(result.FilePath, result.Title, thumbnailPath)
  } else {
    showErrorNotification(url, fmt.Errorf("no video file downloaded"))
  }
}

func cleanupPartialDownloads(downloadDir string) {
  // Clean up all partial download artifacts in the download directory
  entries, err := os.ReadDir(downloadDir)
  if err != nil {
    return
  }
  
  for _, entry := range entries {
    if entry.IsDir() {
      continue
    }
    
    name := entry.Name()
    path := filepath.Join(downloadDir, name)
    
    // Remove ANY temporary files created by yt-dlp (universal patterns)
    shouldRemove := false
    
    // Standard yt-dlp temporary files
    if strings.HasSuffix(name, ".part") ||
       strings.HasSuffix(name, ".ytdl") ||
       strings.HasSuffix(name, ".temp") ||
       strings.Contains(name, ".part-Frag") {
      shouldRemove = true
    }
    
    // Fragment files (format code in filename like .f616., .f251., etc.)
    if strings.Contains(name, ".f") && (strings.Contains(name, ".mp4") || strings.Contains(name, ".webm")) {
      // Check if it has a format code pattern: .fXXX.
      parts := strings.Split(name, ".")
      for i, part := range parts {
        if len(part) > 1 && part[0] == 'f' && i+1 < len(parts) {
          // Check if rest is numeric
          if _, err := strconv.Atoi(part[1:]); err == nil {
            shouldRemove = true
            break
          }
        }
      }
    }
    
    if shouldRemove {
      os.Remove(path)
    }
  }
}

func showCompletionNotification(filePath, title, thumbnailPath string) {
  if title == "" {
    title = ytdlp.TitleFromFile(filePath)
  }

  if len(title) > 45 {
    title = title[:45] + "…"
  }

  // Clean up thumbnail after notification
  if thumbnailPath != "" {
    defer os.Remove(thumbnailPath)
  }

  // Build notification command
  args := []string{"notification", "send", "-A", "default=Open"}
  
  // Only add icon if thumbnail was generated
  if thumbnailPath != "" {
    args = append(args, "-i", thumbnailPath)
  }
  
  args = append(args, "Download complete", title)
  
  cmd := exec.Command("hoopoo", args...)
  output, _ := cmd.Output()
  action := strings.TrimSpace(string(output))

  if action == "default" {
    exec.Command("mpv", filePath).Start()
  }
}

func showErrorNotification(url string, err error) {
  headline := "Download failed"
  description := "Unknown error occurred"
  
  // Provide specific error messages based on error type
  errStr := err.Error()
  if strings.Contains(errStr, "video not available") || strings.Contains(errStr, "Video unavailable") || strings.Contains(errStr, "Private video") {
    headline = "Video not available"
    description = "This video cannot be downloaded"
  } else if strings.Contains(errStr, "file system error") {
    headline = "File system error"
    description = "Check disk space and permissions"
  } else if strings.Contains(errStr, "network error") || strings.Contains(errStr, "timed out") || strings.Contains(errStr, "Unable to download") {
    headline = "Network error"
    description = "Check your internet connection"
  } else if strings.Contains(errStr, "no video file downloaded") {
    headline = "Download failed"
    description = "Could not download video from this URL"
  } else if strings.Contains(errStr, "error starting yt-dlp") {
    headline = "yt-dlp not found"
    description = "Run: hoopoo install chromium-ytdlp"
  } else {
    description = url
  }
  
  exec.Command("hoopoo", "notification", "send",
    "-u", "critical",
    "-t", "10000",
    headline,
    description,
  ).Run()
}
