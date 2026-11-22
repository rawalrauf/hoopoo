// Video transcoding using FFmpeg

package transcode

import (
  "fmt"
  "os/exec"
)

// VideoFormat represents supported video formats
type VideoFormat string

const (
  MP4 VideoFormat = "mp4"
  GIF VideoFormat = "gif"
)

// VideoResolution represents supported video resolutions
type VideoResolution string

const (
  Video4K   VideoResolution = "4k"
  Video1080 VideoResolution = "1080p"
  Video720  VideoResolution = "720p"
)

// TranscodeVideo converts a video to specified format and resolution
func TranscodeVideo(input, format, resolution, output string) error {
  var scale string
  
  switch VideoResolution(resolution) {
  case Video4K:
    scale = "scale=-2:2160"
  case Video1080:
    scale = "scale=-2:1080"
  case Video720:
    scale = "scale=-2:720"
  default:
    return fmt.Errorf("invalid video resolution: %s", resolution)
  }
  
  var cmd *exec.Cmd
  
  switch VideoFormat(format) {
  case MP4:
    if resolution == "4k" {
      // Use H.265 for 4K
      cmd = exec.Command("ffmpeg",
        "-y",
        "-i", input,
        "-vf", scale,
        "-c:v", "libx265",
        "-preset", "slow",
        "-crf", "24",
        "-c:a", "aac",
        "-b:a", "192k",
        output,
      )
    } else {
      // Use H.264 for lower resolutions
      cmd = exec.Command("ffmpeg",
        "-y",
        "-i", input,
        "-vf", scale,
        "-c:v", "libx264",
        "-preset", "fast",
        "-crf", "23",
        "-c:a", "copy",
        output,
      )
    }
  case GIF:
    // Generate optimized GIF with palette
    cmd = exec.Command("ffmpeg",
      "-y",
      "-i", input,
      "-vf", fmt.Sprintf("fps=10,%s:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse", scale),
      output,
    )
  default:
    return fmt.Errorf("invalid video format: %s", format)
  }
  
  return cmd.Run()
}
