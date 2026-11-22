package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  
  "hoopoo/pkg/media/transcode"
  "hoopoo/pkg/menu"
)

// hoopoo:group=transcode
// hoopoo:name=
// hoopoo:summary=Transcode pictures and videos for sharing
// hoopoo:args=[input] [format] [resolution]
// hoopoo:examples=hoopoo transcode|hoopoo transcode ~/Videos/demo.mov mp4 1080p|hoopoo transcode ~/Pictures/wallpaper.heic jpg medium

func main() {
  var input, format, resolution string
  var positional []string
  
  // Parse arguments
  i := 1
  for i < len(os.Args) {
    arg := os.Args[i]
    
    switch arg {
    case "--help", "-h":
      printUsage()
      os.Exit(0)
    default:
      positional = append(positional, arg)
      i++
    }
  }
  
  // Get positional args
  if len(positional) > 0 {
    input = positional[0]
  }
  if len(positional) > 1 {
    format = positional[1]
  }
  if len(positional) > 2 {
    resolution = positional[2]
  }
  
  // Interactive file selection if no input
  if input == "" {
    cmd := exec.Command("hoopoo", "file", "select",
      "--title", "Transcode picture or video",
      "--extensions", "jpg jpeg png webp gif heic avif mp4 mov m4v mkv webm avi",
    )
    output, err := cmd.Output()
    if err != nil {
      // User cancelled
      os.Exit(0)
    }
    input = strings.TrimSpace(string(output))
  }
  
  if input == "" {
    os.Exit(0)
  }
  
  // Check file exists
  if _, err := os.Stat(input); os.IsNotExist(err) {
    fmt.Fprintf(os.Stderr, "File not found: %s\n", input)
    os.Exit(1)
  }
  
  // Detect media type
  mediaType, err := transcode.GetMediaType(input)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v\n", err)
    os.Exit(1)
  }
  
  // Interactive format selection
  if format == "" {
    var formats []menu.MenuItem
    if mediaType == transcode.Picture {
      formats = []menu.MenuItem{
        {Label: "jpg"},
        {Label: "png"},
      }
    } else {
      formats = []menu.MenuItem{
        {Label: "mp4"},
        {Label: "gif"},
      }
    }
    
    selected, err := menu.Show("Select format", formats)
    if err != nil {
      os.Exit(1)
    }
    format = selected.Label
  }
  
  // Interactive resolution selection
  if resolution == "" {
    var resolutions []menu.MenuItem
    if mediaType == transcode.Picture {
      resolutions = []menu.MenuItem{
        {Label: "high"},
        {Label: "medium"},
        {Label: "low"},
      }
    } else {
      resolutions = []menu.MenuItem{
        {Label: "4k"},
        {Label: "1080p"},
        {Label: "720p"},
      }
    }
    
    selected, err := menu.Show("Select resolution", resolutions)
    if err != nil {
      os.Exit(1)
    }
    resolution = selected.Label
  }
  
  // Generate output path
  output := transcode.GenerateOutputPath(input, format, resolution)
  
  // Transcode
  if mediaType == transcode.Video {
    // Show progress notification for videos
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰧸",
      "Transcoding video…",
      fmt.Sprintf("%s to %s (%s)", filepath.Base(input), format, resolution),
    ).Run()
    
    if err := transcode.TranscodeVideo(input, format, resolution, output); err != nil {
      exec.Command("hoopoo", "notification", "send",
        "--glyph=󰀦",
        "Transcode failed",
        err.Error(),
        "-u", "critical",
      ).Run()
      os.Exit(1)
    }
  } else {
    if err := transcode.TranscodePicture(input, format, resolution, output); err != nil {
      exec.Command("hoopoo", "notification", "send",
        "--glyph=󰀦",
        "Transcode failed",
        err.Error(),
        "-u", "critical",
      ).Run()
      os.Exit(1)
    }
  }
  
  // Copy to clipboard as file URI
  uri := "file://" + output
  wlCopyCmd := exec.Command("wl-copy", "--type", "text/uri-list")
  wlCopyCmd.Stdin = strings.NewReader(uri + "\n")
  wlCopyCmd.Run()
  
  // Show completion notification
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰧸",
    fmt.Sprintf("Transcoded to %s %s", resolution, format),
    "Saved and copied to clipboard.",
  ).Run()
}

func printUsage() {
  fmt.Println(`Usage:
  hoopoo transcode [input] [format] [resolution]

With no input, opens a native file picker to select a picture or video.
Then pick the output format and resolution.

Formats:
  Pictures: jpg, png
  Videos: mp4, gif

Resolutions:
  Pictures: high, medium, low
  Videos: 4k, 1080p, 720p

Examples:
  hoopoo transcode
  hoopoo transcode ~/Videos/demo.mov mp4 1080p
  hoopoo transcode ~/Pictures/wallpaper.heic jpg medium`)
}
