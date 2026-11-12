// Video download functionality using yt-dlp

package ytdlp

import (
  "bufio"
  "bytes"
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
  IconPath    string
  OnProgress  func(percent int, title string)
  OnComplete  func(filePath, title string)
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
    "--quiet",
    "--no-warnings",
    "--no-exec",
    "--no-exec-before-download",
    "--progress",
    "--newline",
    "--paths", opts.DownloadDir,
    "-o", "%(title)s.%(ext)s",
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
  lastPercent := -1
  progressRegex := regexp.MustCompile(`\[download\]\s+(\d+)(?:\.\d+)?%`)

  for scanner.Scan() {
    line := scanner.Text()

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
      if percent, err := strconv.Atoi(matches[1]); err == nil && percent != lastPercent {
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

  cmd.Wait()

  if opts.OnComplete != nil {
    opts.OnComplete(filePath, videoTitle)
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
