// Notification utilities with hoopoo-style glyph formatting

package notify

import (
  "fmt"
  "os/exec"
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

func Simple(title, message string) error {
  if message != "" {
    return exec.Command("notify-send", title, message).Run()
  }
  return exec.Command("notify-send", title).Run()
}
