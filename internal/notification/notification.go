// Notification utilities with hoopoo-style glyph formatting

package notification

import (
  "fmt"
  "os/exec"
  "regexp"
  "strings"
)

type Options struct {
  Glyph       string
  Headline    string
  Description string
  Urgency     string
  Icon        string
  Timeout     int
  SoundFile   string
}

func Send(opts Options) error {
  args := []string{}
  
  if opts.Urgency != "" {
    args = append(args, "-u", opts.Urgency)
  }
  
  if opts.Icon != "" {
    args = append(args, "-i", opts.Icon)
  }
  
  if opts.Timeout > 0 {
    args = append(args, "-t", fmt.Sprintf("%d", opts.Timeout))
  }
  
  var summary string
  if opts.Glyph != "" {
    summary = fmt.Sprintf("%s    %s", opts.Glyph, opts.Headline)
  } else {
    summary = opts.Headline
  }
  args = append(args, summary)
  
  if opts.Description != "" {
    lines := strings.Split(opts.Description, "\n")
    for i, line := range lines {
      lines[i] = "      " + line
    }
    indentedDescription := strings.Join(lines, "\n")
    args = append(args, indentedDescription)
  }
  
  if opts.SoundFile != "" {
    exec.Command("systemd-run", "--user", "--no-block", "paplay", opts.SoundFile).Start()
  }
  
  return exec.Command("notify-send", args...).Run()
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
