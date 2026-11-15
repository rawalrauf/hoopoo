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

type restoreData struct {
  Percent int    `json:"percent"`
  Title   string `json:"title"`
  URL     string `json:"url"`
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

  if len(os.Args) > 1 && os.Args[1] == "--restore-notifications" {
    restoreNotifications()
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

  // Track if waybar indicator is active
  waybarIndicatorActive := false
  progressDir := "/tmp/hoopoo-downloads"
  progressFile := fmt.Sprintf("%s/%d", progressDir, os.Getpid())
  restoreFile := fmt.Sprintf("%s/%d.restore", progressDir, os.Getpid())
  signalFile := fmt.Sprintf("%s/%d.signal", progressDir, os.Getpid())

  opts := ytdlp.DownloadOptions{
    URL:         url,
    DownloadDir: downloadDir,
    OnProgress: func(percent int, title string) {
      // Check for signal to restore notification
      if waybarIndicatorActive {
        if data, err := os.ReadFile(signalFile); err == nil {
          var signal map[string]string
          if json.Unmarshal(data, &signal) == nil {
            notifID = signal["notifID"]
            waybarIndicatorActive = false
            os.Remove(signalFile)
          }
        }
      }
      
      if notifID != "" {
        updateCmd := exec.Command("hoopoo", "notification", "send",
          "--replace-id="+notifID,
          "-i", iconPath,
          "--progress="+strconv.Itoa(percent),
          title,
        )
        err := updateCmd.Run()
        if err != nil {
          // Notification was dismissed - switch to waybar indicator
          notifID = ""
          waybarIndicatorActive = true
          os.MkdirAll(progressDir, 0755)
          
          // Write progress and restore data
          os.WriteFile(progressFile, []byte(strconv.Itoa(percent)), 0644)
          restoreJSON, _ := json.Marshal(restoreData{
            Percent: percent,
            Title:   title,
            URL:     url,
          })
          os.WriteFile(restoreFile, restoreJSON, 0644)
          
          exec.Command("pkill", "-RTMIN+11", "waybar").Run()
        }
      } else if waybarIndicatorActive {
        // Update waybar indicator and restore data
        os.WriteFile(progressFile, []byte(strconv.Itoa(percent)), 0644)
        restoreJSON, _ := json.Marshal(restoreData{
          Percent: percent,
          Title:   title,
          URL:     url,
        })
        os.WriteFile(restoreFile, restoreJSON, 0644)
      }
    },
  }

  result, err := ytdlp.Download(opts)
  
  // Clear waybar indicator and signal files
  if waybarIndicatorActive {
    os.Remove(progressFile)
    os.Remove(restoreFile)
    os.Remove(signalFile)
    exec.Command("pkill", "-RTMIN+11", "waybar").Run()
  } else {
    // Clean up any leftover signal file
    os.Remove(signalFile)
  }
  
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
  description := err.Error()
  
  // Check for common error patterns in yt-dlp output
  errStr := strings.ToLower(description)
  switch {
  case strings.Contains(errStr, "video unavailable"), strings.Contains(errStr, "private video"):
    headline = "Video not available"
    description = "This video cannot be downloaded"
  case strings.Contains(errStr, "no space left"), strings.Contains(errStr, "disk quota"):
    headline = "No disk space"
    description = "Free up space and try again"
  case strings.Contains(errStr, "network"), strings.Contains(errStr, "timed out"), strings.Contains(errStr, "connection"):
    headline = "Network error"
    description = "Check your internet connection"
  case strings.Contains(errStr, "no video file downloaded"):
    description = "Could not download video"
  case strings.Contains(errStr, "error starting yt-dlp"):
    headline = "yt-dlp not found"
    description = "Run: hoopoo install chromium-ytdlp"
  }
  
  exec.Command("hoopoo", "notification", "send",
    "-u", "critical",
    "-t", "10000",
    headline,
    description,
  ).Run()
}

func restoreNotifications() {
  progressDir := "/tmp/hoopoo-downloads"
  
  entries, err := os.ReadDir(progressDir)
  if err != nil {
    return
  }
  
  iconPath := filepath.Join(paths.Hoopoo(), "default", "chromium", "extensions", "yt-dlp", "icon.png")
  
  // Find all restore files
  var restoreFiles []string
  for _, entry := range entries {
    if strings.HasSuffix(entry.Name(), ".restore") {
      restoreFiles = append(restoreFiles, entry.Name())
    }
  }
  
  if len(restoreFiles) == 0 {
    return
  }
  
  // Restore all downloads back to notifications
  for _, fileName := range restoreFiles {
    restorePath := filepath.Join(progressDir, fileName)
    data, err := os.ReadFile(restorePath)
    if err != nil {
      continue
    }
    
    var restore restoreData
    if err := json.Unmarshal(data, &restore); err != nil {
      continue
    }
    
    // Create notification with current progress
    notifCmd := exec.Command("hoopoo", "notification", "send",
      "--print-id",
      "-i", iconPath,
      "-t", "300000",
      "--progress="+strconv.Itoa(restore.Percent),
      restore.Title,
    )
    output, err := notifCmd.Output()
    if err != nil {
      continue
    }
    
    notifID := strings.TrimSpace(string(output))
    
    // Extract PID from filename (remove .restore extension)
    pidStr := strings.TrimSuffix(fileName, ".restore")
    
    // Create signal file to tell the download process to switch back
    signalFile := filepath.Join(progressDir, pidStr+".signal")
    signalData, _ := json.Marshal(map[string]string{
      "notifID": notifID,
    })
    os.WriteFile(signalFile, signalData, 0644)
    
    // Remove restore file and progress file (waybar indicator will disappear)
    os.Remove(restorePath)
    os.Remove(filepath.Join(progressDir, pidStr))
  }
  
  // Refresh waybar
  exec.Command("pkill", "-RTMIN+11", "waybar").Run()
}
