// hoopoo:summary=Take a screenshot
// hoopoo:group=capture
// hoopoo:args=[smart|region|window|fullscreen] [--editor=<name>]
// hoopoo:examples=hoopoo screenshot | hoopoo capture screenshot region
// hoopoo:aliases=hoopoo screenshot

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "time"
)

func main() {
  outputDir := filepath.Join(os.Getenv("HOME"), "Pictures", "Screenshots")
  os.MkdirAll(outputDir, 0755)

  editor := os.Getenv("HOOPOO_SCREENSHOT_EDITOR")
  if editor == "" {
    editor = "satty"
  }

  mode := "smart"
  if len(os.Args) > 1 {
    mode = os.Args[1]
  }

  // Parse --editor flag
  for _, arg := range os.Args {
    if strings.HasPrefix(arg, "--editor=") {
      editor = strings.TrimPrefix(arg, "--editor=")
    }
  }

  // Create marker file to detect new screenshots
  markerPath := filepath.Join(outputDir, fmt.Sprintf(".marker.%d", time.Now().UnixNano()))
  markerFile, _ := os.Create(markerPath)
  if markerFile != nil {
    markerFile.Close()
  }
  defer os.Remove(markerPath)

  // Take screenshot using niri
  var cmd *exec.Cmd
  switch mode {
  case "fullscreen":
    cmd = exec.Command("niri", "msg", "action", "screenshot-screen")
  case "window":
    cmd = exec.Command("niri", "msg", "action", "screenshot-window")
  default: // smart, region, or anything else
    cmd = exec.Command("niri", "msg", "action", "screenshot")
  }

  if err := cmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error taking screenshot: %v\n", err)
    os.Exit(1)
  }

  // Wait for screenshot file to appear
  screenshotPath := waitForScreenshot(outputDir, markerPath)
  if screenshotPath == "" {
    os.Exit(0)
  }

  // Send notification with action to edit
  cmd = exec.Command("hoopoo", "notification", "send",
    "Screenshot saved",
    "Edit with Super + Alt + , (or click this)",
    "-t", "10000",
    "-i", screenshotPath,
    "-A", "default=edit",
  )
  output, err := cmd.Output()
  if err != nil {
    os.Exit(0)
  }

  action := strings.TrimSpace(string(output))
  if action == "default" {
    openEditor(screenshotPath, editor)
  }
}

func waitForScreenshot(dir, markerPath string) string {
  markerInfo, err := os.Stat(markerPath)
  if err != nil {
    return ""
  }

  timeout := 5 * time.Second
  start := time.Now()

  for time.Since(start) < timeout {
    files, _ := filepath.Glob(filepath.Join(dir, "*.png"))
    for _, file := range files {
      info, err := os.Stat(file)
      if err != nil {
        continue
      }
      if info.ModTime().After(markerInfo.ModTime()) {
        return file
      }
    }
    time.Sleep(100 * time.Millisecond)
  }

  return ""
}

func openEditor(filepath, editor string) {
  if editor == "satty" {
    cmd := exec.Command("satty",
      "--filename", filepath,
      "--output-filename", filepath,
      "--actions-on-enter", "save-to-clipboard",
      "--save-after-copy",
      "--copy-command", "wl-copy",
    )
    cmd.Run()
  } else {
    exec.Command(editor, filepath).Run()
  }
}
