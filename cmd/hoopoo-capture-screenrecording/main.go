package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  
  "hoopoo/pkg/media/recorder"
  "hoopoo/pkg/media/video"
  "hoopoo/pkg/notification"
  "hoopoo/pkg/system/process"
)

// hoopoo:summary=Start or stop screen recording
// hoopoo:group=capture
// hoopoo:args=[--with-desktop-audio] [--with-microphone-audio] [--with-webcam] [--webcam-device=<device>] [--resolution=<size>] [--stop-recording]
// hoopoo:examples=hoopoo screenrecord | hoopoo capture screenrecording --with-desktop-audio
// hoopoo:aliases=hoopoo screenrecord

func main() {
  // Get output directory
  outputDir := os.Getenv("HOOPOO_SCREENRECORD_DIR")
  if outputDir == "" {
    outputDir = os.Getenv("XDG_VIDEOS_DIR")
  }
  if outputDir == "" {
    outputDir = filepath.Join(os.Getenv("HOME"), "Videos")
  }
  
  // Check directory exists
  if _, err := os.Stat(outputDir); os.IsNotExist(err) {
    notification.Send(notification.Options{
      Glyph:       "󰀦",
      Headline:    "Screen recording directory does not exist: " + outputDir,
      Urgency:     "critical",
      Timeout:     3000,
    })
    os.Exit(1)
  }
  
  // Parse arguments
  opts := recorder.Options{
    OutputDir: outputDir,
    Framerate: 60,
  }
  
  stopRecording := false
  
  for _, arg := range os.Args[1:] {
    switch {
    case arg == "--with-desktop-audio":
      opts.DesktopAudio = true
    case arg == "--with-microphone-audio":
      opts.MicrophoneAudio = true
    case arg == "--with-webcam":
      opts.Webcam = true
    case strings.HasPrefix(arg, "--webcam-device="):
      opts.WebcamDevice = strings.TrimPrefix(arg, "--webcam-device=")
    case strings.HasPrefix(arg, "--resolution="):
      opts.Resolution = strings.TrimPrefix(arg, "--resolution=")
    case arg == "--stop-recording":
      stopRecording = true
    }
  }
  
  // Check if recording is active
  if recorder.IsRecording() {
    stopAndNotify()
    return
  }
  
  // If --stop-recording but nothing is recording, exit
  if stopRecording {
    os.Exit(1)
  }
  
  // Start webcam overlay if requested
  if opts.Webcam {
    if err := recorder.StartWebcamOverlay(opts.WebcamDevice); err != nil {
      notification.Send(notification.Options{
        Glyph:       "󰀦",
        Headline:    "Failed to start webcam overlay",
        Description: err.Error(),
        Urgency:     "critical",
        Timeout:     3000,
      })
      os.Exit(1)
    }
  }
  
  // Play start sound
  process.Start("paplay", "/usr/share/sounds/freedesktop/stereo/dialog-warning.oga")
  
  // Start recording
  filename, err := recorder.Start(opts)
  if err != nil {
    recorder.StopWebcamOverlay()
    notification.Send(notification.Options{
      Glyph:       "󰀦",
      Headline:    "Failed to start recording",
      Description: err.Error(),
      Urgency:     "critical",
      Timeout:     3000,
    })
    os.Exit(1)
  }
  
  // Save recording filename
  recorder.SaveRecordingFile(filename)
  
  // Toggle waybar indicator
  process.SendSignal("RTMIN+8", "waybar")
}

func stopAndNotify() {
  err := recorder.Stop()
  
  // Toggle waybar indicator
  process.SendSignal("RTMIN+8", "waybar")
  
  // Stop webcam overlay
  recorder.StopWebcamOverlay()
  
  if err != nil {
    notification.Send(notification.Options{
      Glyph:       "󰀦",
      Headline:    "Screen recording error",
      Description: "Recording was force-killed. Video may be corrupted.",
      Urgency:     "critical",
      Timeout:     5000,
    })
    recorder.DeleteRecordingFile()
    return
  }
  
  // Get recording filename
  filename, err := recorder.GetRecordingFile()
  if err != nil {
    return
  }
  
  // Trim first frame
  video.TrimFirstFrame(filename)
  
  // Generate thumbnail (cropped square for notification)
  thumbnailPath, _ := video.GenerateVideoThumbnail(filename)
  if thumbnailPath != "" {
    defer os.Remove(thumbnailPath)
  }
  
  // Send notification with action
  icon := filename
  if thumbnailPath != "" {
    icon = thumbnailPath
  }
  
  cmd := exec.Command("hoopoo", "notification", "send",
    "Screen recording saved",
    "Open with Mod + Alt + , (or click this)",
    "-t", "10000",
    "-i", icon,
    "-A", "default=open",
  )
  output, err := cmd.Output()
  if err == nil && strings.TrimSpace(string(output)) == "default" {
    process.Start("mpv", filename)
  }
  
  // Clean up state file
  recorder.DeleteRecordingFile()
}
