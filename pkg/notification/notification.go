// Notification utilities with hoopoo-style glyph formatting

package notification

import (
  "fmt"
  "os/exec"
  "regexp"
  "strings"
  "sync"
  "time"

  "hoopoo/pkg/notification/progress"
)

var (
  lastUpdateTime = make(map[string]time.Time)
  updateMutex    sync.Mutex
)

type Options struct {
  Glyph       string
  Headline    string
  Description string
  Urgency     string
  Icon        string
  Timeout     int
  SoundFile   string
  Progress    int    // 0-100, generates progress bar
  ReplaceID   string // Replace existing notification
}

func Send(opts Options) error {
  _, err := sendWithID(opts)
  return err
}

func SendWithID(opts Options) (string, error) {
  return sendWithID(opts)
}

func sendWithID(opts Options) (string, error) {
  args := []string{}
  
  args = append(args, "--print-id")
  
  if opts.ReplaceID != "" {
    // Throttle progress updates to 250ms (4 updates/sec)
    if opts.Progress > 0 && opts.Progress <= 100 {
      updateMutex.Lock()
      lastTime, exists := lastUpdateTime[opts.ReplaceID]
      now := time.Now()
      
      if exists && now.Sub(lastTime) < 250*time.Millisecond {
        // Too soon, skip this update
        updateMutex.Unlock()
        return opts.ReplaceID, nil
      }
      
      lastUpdateTime[opts.ReplaceID] = now
      updateMutex.Unlock()
    }
    
    // Check if notification exists before trying to replace
    if !notificationExists(opts.ReplaceID) {
      return "", fmt.Errorf("notification %s not found", opts.ReplaceID)
    }
    args = append(args, "--replace-id", opts.ReplaceID)
  }
  
  if opts.Urgency != "" {
    args = append(args, "-u", opts.Urgency)
  }
  
  if opts.Icon != "" {
    args = append(args, "-i", opts.Icon)
  }
  
  // Auto-set timeout for progress bars
  timeout := opts.Timeout
  if opts.Progress > 0 && opts.Progress <= 100 {
    if opts.Progress == 100 && timeout == 0 {
      timeout = 3000  // 3 seconds for completion
    } else if timeout == 0 {
      timeout = 300000  // 5 minutes for in-progress
    }
  }
  
  if timeout > 0 {
    args = append(args, "-t", fmt.Sprintf("%d", timeout))
  }
  
  var summary string
  if opts.Glyph != "" {
    summary = fmt.Sprintf("%s    %s", opts.Glyph, opts.Headline)
  } else {
    summary = opts.Headline
  }
  args = append(args, summary)
  
  description := opts.Description
  
  // If Progress is set and > 0, use progress bar as description (ignore original description)
  if opts.Progress > 0 && opts.Progress <= 100 {
    description = progress.Format(opts.Progress)
  }
  
  if description != "" {
    var formattedDescription string
    if opts.Glyph != "" {
      lines := strings.Split(description, "\n")
      for i, line := range lines {
        lines[i] = "      " + line
      }
      formattedDescription = strings.Join(lines, "\n")
    } else {
      formattedDescription = description
    }
    args = append(args, formattedDescription)
  }
  
  if opts.SoundFile != "" {
    exec.Command("systemd-run", "--user", "--no-block", "paplay", opts.SoundFile).Start()
  }
  
  cmd := exec.Command("notify-send", args...)
  output, err := cmd.CombinedOutput()
  if err != nil {
    return "", fmt.Errorf("%w: %s", err, string(output))
  }
  
  return strings.TrimSpace(string(output)), nil
}

func notificationExists(id string) bool {
  cmd := exec.Command("makoctl", "list")
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  
  // Check if notification ID exists in makoctl list output
  lines := strings.Split(string(output), "\n")
  re := regexp.MustCompile(`^Notification ([0-9]+):.*`)
  
  for _, line := range lines {
    matches := re.FindStringSubmatch(line)
    if len(matches) >= 2 && matches[1] == id {
      return true
    }
  }
  
  return false
}

func Dismiss(summary string) error {
  cmd := exec.Command("makoctl", "list")
  output, err := cmd.Output()
  if err != nil {
    return fmt.Errorf("error listing notifications: %w", err)
  }

  lines := strings.Split(string(output), "\n")
  re := regexp.MustCompile(`^Notification ([0-9]+):.*`)

  for _, line := range lines {
    if strings.Contains(line, summary) {
      matches := re.FindStringSubmatch(line)
      if len(matches) >= 2 {
        notificationID := matches[1]
        return exec.Command("makoctl", "dismiss", "-n", notificationID).Run()
      }
    }
  }
  return nil
}

func ReloadMako() error {
  return exec.Command("makoctl", "reload").Run()
}
