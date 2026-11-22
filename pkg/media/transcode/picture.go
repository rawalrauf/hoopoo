// Picture transcoding using ImageMagick

package transcode

import (
  "fmt"
  "os/exec"
)

// PictureFormat represents supported picture formats
type PictureFormat string

const (
  JPG PictureFormat = "jpg"
  PNG PictureFormat = "png"
)

// PictureResolution represents supported picture resolutions
type PictureResolution string

const (
  PictureHigh   PictureResolution = "high"
  PictureMedium PictureResolution = "medium"
  PictureLow    PictureResolution = "low"
)

// TranscodePicture converts a picture to specified format and resolution
func TranscodePicture(input, format, resolution, output string) error {
  var resize string
  
  switch PictureResolution(resolution) {
  case PictureHigh:
    resize = "3160x>"
  case PictureMedium:
    resize = "2160x>"
  case PictureLow:
    resize = "1080x>"
  default:
    return fmt.Errorf("invalid picture resolution: %s", resolution)
  }
  
  var cmd *exec.Cmd
  
  switch PictureFormat(format) {
  case JPG:
    cmd = exec.Command("magick", input,
      "-resize", resize,
      "-quality", "85",
      "-strip",
      output,
    )
  case PNG:
    cmd = exec.Command("magick", input,
      "-resize", resize,
      "-strip",
      "-define", "png:compression-filter=5",
      "-define", "png:compression-level=9",
      "-define", "png:compression-strategy=1",
      "-define", "png:exclude-chunk=all",
      output,
    )
  default:
    return fmt.Errorf("invalid picture format: %s", format)
  }
  
  return cmd.Run()
}
