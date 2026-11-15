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
    "--progress",
    "--newline",
    "--format", "bestvideo[protocol^=https]+bestaudio[protocol^=https]/best[protocol^=https]/bestvideo+bestaudio/best",
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
        if len(videoTitle) > 45 {
          videoTitle = videoTitle[:45] + "…"
        }
        titleUpdated = true
      }
    }

    if strings.HasPrefix(line, "HOOPOO_FILE\t") {
      filePath = strings.TrimPrefix(line, "HOOPOO_FILE\t")
    }

    if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
      if percent, err := strconv.Atoi(matches[1]); err == nil {
        // Only update notification if we've made progress
        if percent > lastPercent {
          if opts.OnProgress != nil {
            displayTitle := videoTitle
            if displayTitle == "" {
              displayTitle = opts.URL
            }
            opts.OnProgress(percent, displayTitle)
          }
          lastPercent = percent
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
    return nil, fmt.Errorf("%s", errMsg)
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

  if title == "" || title[0] == '-' {
    return "", fmt.Errorf("invalid title")
  }

  return title, nil
}

func TitleFromFile(filePath string) string {
  name := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
  name = strings.Map(func(r rune) rune {
    if r < 32 { return -1 }
    return r
  }, name)
  
  if name == "" || name[0] == '-' {
    return "Video"
  }
  return name
}
