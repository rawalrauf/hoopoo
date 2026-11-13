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
  
  totalBlocks := 200
  filledBlocks := (percent * totalBlocks) / 100
  emptyBlocks := totalBlocks - filledBlocks
  
  filled := strings.Repeat("▎", filledBlocks)
  empty := strings.Repeat("▎", emptyBlocks)
  
  // Use theme color for filled, medium gray for empty (works with light and dark themes)
  bar := fmt.Sprintf("<span size='small' letter_spacing='-9500'>%s<span foreground='#888888'>%s</span></span>", filled, empty)
  return fmt.Sprintf("%s %d%%", bar, percent)
}
