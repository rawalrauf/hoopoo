// Video processing operations using ffmpeg

package video

import (
  "os"
  "os/exec"
)

// TrimFirstFrame removes the first 0.1 seconds from a video file
func TrimFirstFrame(filePath string) error {
  if filePath == "" {
    return nil
  }
  
  tmpFile := filePath + ".trimmed.mp4"
  
  cmd := exec.Command("ffmpeg",
    "-y",
    "-ss", "0.1",
    "-i", filePath,
    "-c", "copy",
    tmpFile,
    "-loglevel", "quiet",
  )
  
  if err := cmd.Run(); err != nil {
    os.Remove(tmpFile)
    return err
  }
  
  // Replace original with trimmed version
  if err := os.Rename(tmpFile, filePath); err != nil {
    os.Remove(tmpFile)
    return err
  }
  
  return nil
}
