// Notification utilities with hoopoo-style glyph formatting

package notify

import (
  "fmt"
  "os/exec"
  "strings"
)

func SendWithGlyph(glyph, headline, description, urgency, icon string, timeout int) error {
  args := []string{}
  
  if urgency != "" {
    args = append(args, "-u", urgency)
  }
  
  if icon != "" {
    args = append(args, "-i", icon)
  }
  
  if timeout > 0 {
    args = append(args, "-t", fmt.Sprintf("%d", timeout))
  }
  
  var summary string
  if glyph != "" {
    summary = fmt.Sprintf("%s    %s", glyph, headline)
  } else {
    summary = headline
  }
  args = append(args, summary)
  
  if description != "" {
    lines := strings.Split(description, "\n")
    for i, line := range lines {
      lines[i] = "      " + line
    }
    indentedDescription := strings.Join(lines, "\n")
    args = append(args, indentedDescription)
  }
  
  return exec.Command("notify-send", args...).Run()
}

func SendWithGlyphAndSound(glyph, headline, description, urgency, icon string, timeout int, soundFile string) error {
  exec.Command("systemd-run", "--user", "--no-block", "paplay", soundFile).Start()
  return SendWithGlyph(glyph, headline, description, urgency, icon, timeout)
}

func Simple(title, message string) error {
  if message != "" {
    return exec.Command("notify-send", title, message).Run()
  }
  return exec.Command("notify-send", title).Run()
}
