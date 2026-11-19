// Font-aware walker dimension configurations

package walker

import (
  "strings"
  
  "hoopoo/pkg/font"
)

// Dimensions holds walker window sizing configuration
type Dimensions struct {
  Width      string
  ItemHeight int
  MinHeight  int
  MaxHeight  int
}

// Font-specific dimension presets
var fontDimensions = map[string]Dimensions{
  "0xproto nerd font": {
    Width:      "295",
    ItemHeight: 58,  // 348px ÷ 6 items = 58px per item
    MinHeight:  150,
    MaxHeight:  630,
  },
  "caskaydiamono nerd font": {
    Width:      "295",
    ItemHeight: 49,  // 294px ÷ 6 items = 49px per item
    MinHeight:  150,
    MaxHeight:  630,
  },
  "jetbrainsmono nerd font": {
    Width:      "295",
    ItemHeight: 52,  // 312px ÷ 6 items = 52px per item
    MinHeight:  150,
    MaxHeight:  630,
  },
  "ubuntumono nerd font": {
    Width:      "295",
    ItemHeight: 46,  // 324px ÷ 7 items ≈ 46px per item (shows 7 items)
    MinHeight:  150,
    MaxHeight:  630,
  },
  "monospace": {
    Width:      "295",
    ItemHeight: 58,
    MinHeight:  150,
    MaxHeight:  630,
  },
  "default": {
    Width:      "295",
    ItemHeight: 58,
    MinHeight:  150,
    MaxHeight:  630,
  },
}

// Get returns dimensions for current font
func Get() Dimensions {
  fontName := strings.ToLower(font.Current())
  if dims, ok := fontDimensions[fontName]; ok {
    return dims
  }
  return fontDimensions["default"]
}

// GetLauncherDefaults returns launcher dimensions based on current font
func GetLauncherDefaults() Dimensions {
  dims := Get() // Get font-specific dimensions
  
  // Calculate launcher height based on font
  // Most fonts: 6 items, UbuntuMono: 7 items (better rendering)
  itemCount := 6
  fontName := strings.ToLower(font.Current())
  if fontName == "ubuntumono nerd font" {
    itemCount = 7
  }
  
  launcherHeight := itemCount * dims.ItemHeight
  
  return Dimensions{
    Width:     "644",
    MinHeight: launcherHeight,
    MaxHeight: launcherHeight,
  }
}
