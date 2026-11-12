// Progress bar formatting for notifications

package progress

import (
  "fmt"
  "strings"
)

func Format(percent int) string {
  if percent < 0 {
    percent = 0
  }
  if percent > 100 {
    percent = 100
  }
  
  totalBlocks := 20
  filledBlocks := (percent * totalBlocks) / 100
  emptyBlocks := totalBlocks - filledBlocks
  
  bar := strings.Repeat("■", filledBlocks) + strings.Repeat("□", emptyBlocks)
  return fmt.Sprintf("%s %d%%", bar, percent)
}
