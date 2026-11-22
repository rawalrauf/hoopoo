// Audio device detection for PipeWire/PulseAudio

package hardware

import (
  "os/exec"
  "strings"
)

// Audio device detection

// GetDefaultSink returns the default audio output device
func GetDefaultSink() (string, error) {
  cmd := exec.Command("pactl", "get-default-sink")
  output, err := cmd.Output()
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(output)), nil
}

// GetDefaultSource returns the default audio input device (microphone)
func GetDefaultSource() (string, error) {
  cmd := exec.Command("pactl", "get-default-source")
  output, err := cmd.Output()
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(output)), nil
}

// BuildDeviceString returns the device string for gpu-screen-recorder
// e.g., "default_output" or "default_output|default_input"
func BuildDeviceString(desktop, microphone bool) string {
  var devices []string
  
  if desktop {
    devices = append(devices, "default_output")
  }
  
  if microphone {
    devices = append(devices, "default_input")
  }
  
  return strings.Join(devices, "|")
}
