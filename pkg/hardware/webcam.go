// Webcam device detection and management using v4l2

package hardware

import (
  "bufio"
  "bytes"
  "os/exec"
  "strings"
)

// Webcam device detection

// WebcamDevice represents a webcam device
type WebcamDevice struct {
  Path string
  Name string
}

// ListWebcams returns all available webcam devices
func ListWebcams() ([]WebcamDevice, error) {
  cmd := exec.Command("v4l2-ctl", "--list-devices")
  output, err := cmd.Output()
  if err != nil {
    return nil, err
  }
  
  var devices []WebcamDevice
  scanner := bufio.NewScanner(bytes.NewReader(output))
  
  var currentName string
  for scanner.Scan() {
    line := scanner.Text()
    
    // Device name line (not indented)
    if len(line) > 0 && line[0] != '\t' && line[0] != ' ' {
      currentName = strings.TrimSpace(line)
      continue
    }
    
    // Device path line (indented with tab or spaces)
    trimmed := strings.TrimSpace(line)
    if strings.HasPrefix(trimmed, "/dev/video") {
      devices = append(devices, WebcamDevice{
        Path: trimmed,
        Name: currentName,
      })
      break // Take first device path for each name
    }
  }
  
  return devices, nil
}

// GetFirstWebcam returns the first available webcam device path
func GetFirstWebcam() (string, error) {
  devices, err := ListWebcams()
  if err != nil || len(devices) == 0 {
    return "", err
  }
  return devices[0].Path, nil
}

// GetWebcamResolutions returns available resolutions for a device
func GetWebcamResolutions(devicePath string) ([]string, error) {
  cmd := exec.Command("v4l2-ctl", "--list-formats-ext", "-d", devicePath)
  output, err := cmd.Output()
  if err != nil {
    return nil, err
  }
  
  var resolutions []string
  scanner := bufio.NewScanner(bytes.NewReader(output))
  
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if strings.Contains(line, "Size:") {
      // Extract resolution like "640x480"
      parts := strings.Fields(line)
      for _, part := range parts {
        if strings.Contains(part, "x") && len(strings.Split(part, "x")) == 2 {
          resolutions = append(resolutions, part)
          break
        }
      }
    }
  }
  
  return resolutions, nil
}

// WebcamSupportsFormat checks if device supports a specific format (e.g., "MJPG")
func WebcamSupportsFormat(devicePath, format string) bool {
  cmd := exec.Command("v4l2-ctl", "--list-formats", "-d", devicePath)
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  return strings.Contains(string(output), format)
}
