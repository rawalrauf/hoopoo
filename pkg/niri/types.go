// Niri types

package niri

// OutputMode represents a display mode
type OutputMode struct {
  Width        int  `json:"width"`
  Height       int  `json:"height"`
  RefreshRate  int  `json:"refresh_rate"`
  IsPreferred  bool `json:"is_preferred"`
}

// LogicalOutput represents the logical output dimensions
type LogicalOutput struct {
  X         int     `json:"x"`
  Y         int     `json:"y"`
  Width     int     `json:"width"`
  Height    int     `json:"height"`
  Scale     float64 `json:"scale"`
  Transform string  `json:"transform"`
}

// Output represents a niri output (display)
type Output struct {
  Name         string        `json:"name"`
  Make         string        `json:"make"`
  Model        string        `json:"model"`
  Modes        []OutputMode  `json:"modes"`
  CurrentMode  int           `json:"current_mode"`
  Logical      LogicalOutput `json:"logical"`
}

// Window represents a niri window
type Window struct {
  ID    int    `json:"id"`
  PID   int    `json:"pid"`
  AppID string `json:"app_id"`
  Title string `json:"title"`
}
