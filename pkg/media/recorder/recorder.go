// Screen recorder wrapper for gpu-screen-recorder

package recorder

import (
  "fmt"
  "os"
  "path/filepath"
  "time"
  
  "hoopoo/pkg/hardware"
  "hoopoo/pkg/niri"
  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
)

// Options for screen recording
type Options struct {
  OutputDir        string
  DesktopAudio     bool
  MicrophoneAudio  bool
  Webcam           bool
  WebcamDevice     string
  Resolution       string
  Framerate        int
}

// IsRecording checks if a recording is currently active
func IsRecording() bool {
  return process.IsRunningByPattern("^gpu-screen-recorder")
}

// Start begins a new screen recording
func Start(opts Options) (string, error) {
  // Generate filename
  timestamp := time.Now().Format("2006-01-02_15-04-05")
  filename := filepath.Join(opts.OutputDir, fmt.Sprintf("screenrecording-%s.mp4", timestamp))
  
  // Build audio device string
  audioDevices := hardware.BuildDeviceString(opts.DesktopAudio, opts.MicrophoneAudio)
  
  // Get resolution
  resolution := opts.Resolution
  if resolution == "" {
    var err error
    resolution, err = getDefaultResolution()
    if err != nil {
      return "", err
    }
  }
  
  // Get focused monitor
  output, err := niri.GetFocusedOutputInfo()
  if err != nil {
    return "", err
  }
  
  // Build command
  args := []string{
    "-w", output.Name,
    "-k", "auto",
    "-s", resolution,
    "-f", fmt.Sprintf("%d", opts.Framerate),
    "-fm", "cfr",
    "-fallback-cpu-encoding", "yes",
    "-o", filename,
  }
  
  // Add audio if configured
  if audioDevices != "" {
    args = append(args, "-a", audioDevices, "-ac", "aac")
  }
  
  if err := process.Start("gpu-screen-recorder", args...); err != nil {
    return "", err
  }
  
  // Wait for file to be created
  timeout := time.After(5 * time.Second)
  ticker := time.NewTicker(200 * time.Millisecond)
  defer ticker.Stop()
  
  for {
    select {
    case <-timeout:
      return "", fmt.Errorf("timeout waiting for recording to start")
    case <-ticker.C:
      if _, err := os.Stat(filename); err == nil {
        return filename, nil
      }
    }
  }
}

// Stop terminates the current recording
func Stop() error {
  if !IsRecording() {
    return fmt.Errorf("no recording in progress")
  }
  
  // Send SIGINT for clean shutdown
  if err := process.KillByPattern("^gpu-screen-recorder"); err != nil {
    return err
  }
  
  // Wait for process to exit (max 5 seconds)
  for i := 0; i < 50; i++ {
    if !IsRecording() {
      return nil
    }
    time.Sleep(100 * time.Millisecond)
  }
  
  // Force kill if still running
  if IsRecording() {
    process.KillForceByPattern("^gpu-screen-recorder")
    return fmt.Errorf("recording was force-killed, video may be corrupted")
  }
  
  return nil
}

// getDefaultResolution returns appropriate resolution based on display
func getDefaultResolution() (string, error) {
  output, err := niri.GetFocusedOutputInfo()
  if err != nil {
    return "0x0", err
  }
  
  width := output.Logical.Width
  height := output.Logical.Height
  
  // Cap at 4K for performance
  if width > 3840 || height > 2160 {
    return "3840x2160", nil
  }
  
  return "0x0", nil
}

// GetRecordingFile reads the recording filename from state file
func GetRecordingFile() (string, error) {
  stateFile := filepath.Join(paths.Runtime(), "hoopoo-screenrecord-filename")
  data, err := os.ReadFile(stateFile)
  if err != nil {
    return "", err
  }
  return string(data), nil
}

// SaveRecordingFile saves the recording filename to state file
func SaveRecordingFile(filename string) error {
  stateFile := filepath.Join(paths.Runtime(), "hoopoo-screenrecord-filename")
  return os.WriteFile(stateFile, []byte(filename), 0644)
}

// DeleteRecordingFile removes the recording state file
func DeleteRecordingFile() error {
  stateFile := filepath.Join(paths.Runtime(), "hoopoo-screenrecord-filename")
  return os.Remove(stateFile)
}
