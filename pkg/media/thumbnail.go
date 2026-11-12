// Video thumbnail generation for media files

package media

import (
  "bytes"
  "os"
  "os/exec"
)

func GenerateVideoThumbnail(filePath string) (string, error) {
  tmpFile, err := os.CreateTemp("", "preview-*.jpg")
  if err != nil {
    return "", err
  }
  preview := tmpFile.Name()
  tmpFile.Close()

  cmd := exec.Command("ffmpeg",
    "-y",
    "-i", filePath,
    "-ss", "00:00:00.1",
    "-vframes", "1",
    "-vf", "crop='min(iw,ih)':'min(iw,ih)',scale=256:256",
    "-q:v", "2",
    preview,
    "-loglevel", "quiet",
  )

  var stderr bytes.Buffer
  cmd.Stderr = &stderr

  if err := cmd.Run(); err != nil {
    os.Remove(preview)
    return "", err
  }

  return preview, nil
}
