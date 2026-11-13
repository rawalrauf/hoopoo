// Video download functionality using yt-dlp

package ytdlp

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "strconv"
  "strings"
)

type DownloadOptions struct {
  URL         string
  DownloadDir string
  OnProgress  func(percent int, title string)
}

type DownloadResult struct {
  FilePath string
  Title    string
}

func Download(opts DownloadOptions) (*DownloadResult, error) {
  if err := os.MkdirAll(opts.DownloadDir, 0755); err != nil {
    return nil, fmt.Errorf("error creating download directory: %w", err)
  }

  cmd := exec.Command("yt-dlp",
    "--no-playlist",
    "--no-simulate",
    "--no-exec",
    "--no-exec-before-download",
    "--progress",
    "--newline",
    "--paths", opts.DownloadDir,
    "-o", "%(title)s [%(id)s].%(ext)s",
    "--print", "before_dl:HOOPOO_TITLE\t%(title)j",
    "--print", "after_move:HOOPOO_FILE\t%(filepath)s",
    "--", opts.URL,
  )

  cmd.Env = append(os.Environ(), "PYTHONUNBUFFERED=1")

  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return nil, fmt.Errorf("error creating pipe: %w", err)
  }

  stderr, err := cmd.StderrPipe()
  if err != nil {
    return nil, fmt.Errorf("error creating stderr pipe: %w", err)
  }

  if err := cmd.Start(); err != nil {
    return nil, fmt.Errorf("error starting yt-dlp: %w", err)
  }

  combined := io.MultiReader(stdout, stderr)
  scanner := bufio.NewScanner(combined)

  var videoTitle string
  var filePath string
  var titleUpdated bool
  var stderrOutput strings.Builder
  lastPercent := 0
  maxPercent := 0
  progressRegex := regexp.MustCompile(`\[download\]\s+(\d+)(?:\.\d+)?%`)

  for scanner.Scan() {
    line := scanner.Text()

    // Capture stderr for error detection
    if strings.Contains(line, "ERROR:") || strings.Contains(line, "WARNING:") {
      stderrOutput.WriteString(line)
      stderrOutput.WriteString("\n")
    }

    if strings.HasPrefix(line, "HOOPOO_TITLE\t") {
      title := strings.TrimPrefix(line, "HOOPOO_TITLE\t")
      if decoded, err := DecodeTitle(title); err == nil && !titleUpdated {
        videoTitle = decoded
        if len(videoTitle) > 50 {
          videoTitle = videoTitle[:50] + "…"
        }
        titleUpdated = true
      }
    }

    if strings.HasPrefix(line, "HOOPOO_FILE\t") {
      path := strings.TrimPrefix(line, "HOOPOO_FILE\t")
      if resolved, err := resolveDownloadFile(path, opts.DownloadDir); err == nil {
        filePath = resolved
      }
    }

    if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
      if percent, err := strconv.Atoi(matches[1]); err == nil {
        // Detect fragment reset (progress dropped significantly, like 100% -> 0%)
        if percent < 10 && maxPercent > 80 {
          // New fragment starting, reset tracking
          maxPercent = 0
          lastPercent = 0
        }
        
        // Track max progress seen (for multi-fragment downloads)
        if percent > maxPercent {
          maxPercent = percent
        }
        
        // Only update notification if we've made progress
        if maxPercent > lastPercent {
          if opts.OnProgress != nil {
            displayTitle := videoTitle
            if displayTitle == "" {
              displayTitle = opts.URL
            }
            opts.OnProgress(maxPercent, displayTitle)
          }
          lastPercent = maxPercent
        }
      }
    }
  }

  // Check for scanner errors
  if err := scanner.Err(); err != nil {
    return nil, fmt.Errorf("error reading yt-dlp output: %w", err)
  }

  if err := cmd.Wait(); err != nil {
    errMsg := stderrOutput.String()
    if errMsg == "" {
      errMsg = err.Error()
    }
    
    // Normalize for case-insensitive matching
    errMsgLower := strings.ToLower(errMsg)
    
    // Classify errors
    if strings.Contains(errMsgLower, "video unavailable") || 
       strings.Contains(errMsgLower, "private video") || 
       strings.Contains(errMsgLower, "not available") {
      return nil, fmt.Errorf("video not available: %s", errMsg)
    } else if strings.Contains(errMsgLower, "no such file") || 
       strings.Contains(errMsgLower, "permission denied") ||
       strings.Contains(errMsgLower, "cannot open") ||
       strings.Contains(errMsgLower, "disk quota") ||
       strings.Contains(errMsgLower, "no space left") {
      return nil, fmt.Errorf("file system error: %s", errMsg)
    } else if strings.Contains(errMsgLower, "network") || 
       strings.Contains(errMsgLower, "timed out") || 
       strings.Contains(errMsgLower, "unable to download") ||
       strings.Contains(errMsgLower, "connection") ||
       strings.Contains(errMsgLower, "failed to") ||
       strings.Contains(errMsgLower, "http error") {
      return nil, fmt.Errorf("network error: %s", errMsg)
    } else {
      return nil, fmt.Errorf("download failed: %s", errMsg)
    }
  }

  // If no file was downloaded, return error
  if filePath == "" {
    return nil, fmt.Errorf("no video file downloaded")
  }

  return &DownloadResult{
    FilePath: filePath,
    Title:    videoTitle,
  }, nil
}

func DecodeTitle(jsonStr string) (string, error) {
  var title string
  if err := json.Unmarshal([]byte(jsonStr), &title); err != nil {
    return "", err
  }

  title = strings.TrimFunc(title, func(r rune) bool {
    return r < 32
  })

  if title == "" || strings.HasPrefix(title, "-") {
    return "", fmt.Errorf("invalid title")
  }

  return title, nil
}

func resolveDownloadFile(candidate, downloadDir string) (string, error) {
  if candidate == "" {
    return "", fmt.Errorf("empty path")
  }

  if strings.ContainsAny(candidate, "\n\r\t") {
    return "", fmt.Errorf("path contains control characters")
  }

  info, err := os.Stat(candidate)
  if err != nil || !info.Mode().IsRegular() {
    return "", fmt.Errorf("not a regular file")
  }

  fileReal, err := filepath.EvalSymlinks(candidate)
  if err != nil {
    return "", err
  }

  dirReal, err := filepath.EvalSymlinks(downloadDir)
  if err != nil {
    return "", err
  }

  if !strings.HasPrefix(fileReal, dirReal+string(filepath.Separator)) {
    return "", fmt.Errorf("file outside download directory")
  }

  return fileReal, nil
}

func TitleFromFile(filePath string) string {
  name := filepath.Base(filePath)
  ext := filepath.Ext(name)
  name = strings.TrimSuffix(name, ext)
  name = strings.Map(func(r rune) rune {
    if r == '\n' || r == '\r' || r == '\t' {
      return -1
    }
    return r
  }, name)

  if name == "" || strings.HasPrefix(name, "-") {
    return "Video"
  }
  return name
}
