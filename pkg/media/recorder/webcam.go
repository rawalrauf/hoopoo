// Webcam overlay management using ffplay

package recorder

import (
  "fmt"
  "time"
  
  "hoopoo/pkg/hardware"
  "hoopoo/pkg/niri"
  "hoopoo/pkg/system/process"
)

// StartWebcamOverlay starts ffplay webcam overlay window
func StartWebcamOverlay(devicePath string) error {
  // Stop any existing overlay
  StopWebcamOverlay()
  
  // Auto-detect device if not specified
  if devicePath == "" {
    var err error
    devicePath, err = hardware.GetFirstWebcam()
    if err != nil {
      return fmt.Errorf("no webcam devices found")
    }
  }
  
  // Get screen dimensions
  output, err := niri.GetFocusedOutputInfo()
  if err != nil {
    return err
  }
  
  screenHeight := output.Logical.Height
  
  // Detect best video size
  videoSize := "640x480"
  preferredResolutions := []string{"640x480", "640x360", "1280x720", "1920x1080"}
  supportedResolutions, _ := hardware.GetWebcamResolutions(devicePath)
  
  for _, pref := range preferredResolutions {
    for _, supported := range supportedResolutions {
      if pref == supported {
        videoSize = pref
        break
      }
    }
  }
  
  // Calculate window size based on camera's actual aspect ratio
  // Parse detected resolution to get aspect ratio
  var resWidth, resHeight int
  fmt.Sscanf(videoSize, "%dx%d", &resWidth, &resHeight)
  
  // Window height: 30% of screen height
  camHeight := screenHeight * 3 / 10
  // Window width: maintain camera's aspect ratio (no letterboxing)
  camWidth := camHeight * resWidth / resHeight
  
  // Check for MJPEG support
  args := []string{
    "-f", "v4l2",
    "-video_size", videoSize,
  }
  
  if hardware.WebcamSupportsFormat(devicePath, "MJPG") {
    args = append(args, "-input_format", "mjpeg")
  }
  
  args = append(args,
    "-framerate", "30",
    devicePath,
    "-vf", "hflip",
    "-x", fmt.Sprintf("%d", camWidth),
    "-y", fmt.Sprintf("%d", camHeight),
    "-window_title", "WebcamOverlay",
    "-noborder",
    "-fflags", "nobuffer",
    "-flags", "low_delay",
    "-probesize", "32",
    "-analyzeduration", "0",
    "-loglevel", "quiet",
  )
  
  if err := process.Start("ffplay", args...); err != nil {
    return err
  }
  
  // Give it time to start
  time.Sleep(1 * time.Second)
  
  return nil
}

// StopWebcamOverlay kills the webcam overlay process
func StopWebcamOverlay() {
  process.KillByName("WebcamOverlay")
}

// IsWebcamOverlayRunning checks if webcam overlay is active
func IsWebcamOverlayRunning() bool {
  return process.IsRunningByName("WebcamOverlay")
}
