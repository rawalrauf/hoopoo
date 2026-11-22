// Media transcoding for pictures and videos

package transcode

import (
  "fmt"
  "os/exec"
  "path/filepath"
  "strings"
)

// MediaType represents the type of media file
type MediaType string

const (
  Picture MediaType = "picture"
  Video   MediaType = "video"
)

// GetMediaType detects if a file is a picture or video
func GetMediaType(filePath string) (MediaType, error) {
  cmd := exec.Command("file", "-b", "--mime-type", filePath)
  output, err := cmd.Output()
  if err != nil {
    return "", err
  }
  
  mimeType := strings.TrimSpace(string(output))
  
  if strings.HasPrefix(mimeType, "image/") {
    return Picture, nil
  }
  
  if strings.HasPrefix(mimeType, "video/") {
    return Video, nil
  }
  
  return "", fmt.Errorf("unsupported file type: %s", mimeType)
}

// GenerateOutputPath creates output path with format and resolution suffix
func GenerateOutputPath(input, format, resolution string) string {
  dir := filepath.Dir(input)
  base := filepath.Base(input)
  ext := filepath.Ext(base)
  stem := strings.TrimSuffix(base, ext)
  
  return filepath.Join(dir, fmt.Sprintf("%s-%s-%s.%s", stem, resolution, format, format))
}
