// GTK theme integration with Hoopoo's theme system
//
// This package loads Hoopoo theme colors from the generated gtk.css
// and provides styling utilities for GTK4 widgets.

package theme

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetGTKCSSPath returns the path to the generated GTK CSS with theme colors
func GetGTKCSSPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	cssPath := filepath.Join(home, ".config", "hoopoo", "current", "theme", "gtk.css")
	if _, err := os.Stat(cssPath); err != nil {
		return "", fmt.Errorf("gtk.css not found (theme not set?): %w", err)
	}
	
	return cssPath, nil
}

// GenerateWidgetCSS generates GTK4 CSS for Hoopoo widgets
// This imports the theme colors and defines all widget styles
func GenerateWidgetCSS() string {
	return `
/* Import theme colors */
@import url('file:///home/${USER}/.config/hoopoo/current/theme/gtk.css');

/* Universal overlay window styling */
.hoopoo-overlay {
  background-color: @background;
  color: @foreground;
  border: 2px solid @accent;
  border-radius: 12px;
  padding: 16px;
}

.hoopoo-overlay-collapsed {
  opacity: 0.8;
  padding: 8px;
}

.hoopoo-overlay-header {
  font-weight: bold;
  font-size: 14px;
  color: @accent;
  margin-bottom: 12px;
}

/* Progress widgets */
.hoopoo-progress-bar {
  min-height: 8px;
  border-radius: 4px;
}

.hoopoo-progress-bar progress {
  background-color: @color8;
}

.hoopoo-progress-bar trough {
  background-color: @color0;
  border-radius: 4px;
}

.hoopoo-progress-percentage {
  font-size: 12px;
  color: @foreground;
  margin-top: 4px;
}

.hoopoo-progress-title {
  font-size: 13px;
  color: @foreground;
  margin-bottom: 8px;
}

/* Status indicators */
.hoopoo-status {
  padding: 8px;
  border-radius: 8px;
  background-color: alpha(@accent, 0.1);
}

.hoopoo-status-icon {
  color: @accent;
  font-size: 16px;
  margin-right: 8px;
}

.hoopoo-status-message {
  color: @foreground;
  font-size: 12px;
}

/* Spinner animation */
.hoopoo-spinner {
  color: @accent;
}

/* Buttons */
.hoopoo-button {
  background-color: @accent;
  color: @background;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-weight: 500;
}

.hoopoo-button:hover {
  background-color: lighter(@accent);
}

.hoopoo-button-secondary {
  background-color: @color8;
  color: @foreground;
}

/* Collapse/expand button */
.hoopoo-collapse-button {
  background-color: transparent;
  border: 1px solid @accent;
  color: @accent;
  border-radius: 4px;
  padding: 4px 8px;
}

.hoopoo-collapse-button:hover {
  background-color: alpha(@accent, 0.1);
}
`
}

// ApplyToProvider applies Hoopoo theme to a GTK CSS provider
func ApplyToProvider(provider interface{}) error {
	// TODO: Accept GTK4 CssProvider
	// Load both the color definitions and widget styles
	return fmt.Errorf("not implemented yet")
}

// WatchThemeChanges watches for theme changes and reloads
func WatchThemeChanges(callback func()) error {
	// TODO: Use fsnotify to watch ~/.config/hoopoo/current/theme/gtk.css
	// Call callback when file changes
	return fmt.Errorf("not implemented yet")
}
