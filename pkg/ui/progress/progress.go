// Reusable progress widgets for GTK4
//
// This package provides various progress display widgets that can be
// embedded in overlays or other UI components.

package progress

import (
	"fmt"
)

// Bar represents a progress bar widget
type Bar struct {
	// Progress state
	Percent int    // 0-100
	Title   string // Optional title above bar
	
	// Styling
	ShowPercentage bool
	Animated       bool
	
	// GTK widget (will be actual GTK4 widget)
	widget interface{}
}

// Config holds progress bar configuration
type Config struct {
	Title          string
	ShowPercentage bool
	Animated       bool
	Width          int
	Height         int
}

// NewBar creates a new progress bar widget
func NewBar(config Config) (*Bar, error) {
	// TODO: Create GTK4 progress bar widget
	// TODO: Apply Hoopoo theme styling
	// TODO: Set up animations if enabled
	
	b := &Bar{
		Title:          config.Title,
		ShowPercentage: config.ShowPercentage,
		Animated:       config.Animated,
		Percent:        0,
	}
	
	return b, nil
}

// SetProgress updates the progress bar value
func (b *Bar) SetProgress(percent int) error {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	
	b.Percent = percent
	
	// TODO: Update GTK widget
	// TODO: Trigger animation if enabled
	
	return nil
}

// SetTitle updates the progress bar title
func (b *Bar) SetTitle(title string) error {
	b.Title = title
	
	// TODO: Update GTK label
	
	return nil
}

// GetWidget returns the underlying GTK4 widget
func (b *Bar) GetWidget() interface{} {
	// TODO: Return actual GTK4 widget
	return b.widget
}

// Spinner represents an indeterminate spinner widget
type Spinner struct {
	Running bool
	widget  interface{}
}

// NewSpinner creates a new spinner widget
func NewSpinner() (*Spinner, error) {
	// TODO: Create GTK4 spinner
	// TODO: Apply theme
	
	return &Spinner{Running: false}, nil
}

// Start starts the spinner animation
func (s *Spinner) Start() error {
	if s.Running {
		return nil
	}
	
	s.Running = true
	// TODO: Start GTK spinner
	
	return nil
}

// Stop stops the spinner animation
func (s *Spinner) Stop() error {
	if !s.Running {
		return nil
	}
	
	s.Running = false
	// TODO: Stop GTK spinner
	
	return nil
}

// GetWidget returns the underlying GTK4 widget
func (s *Spinner) GetWidget() interface{} {
	return s.widget
}

// Status represents a status display with icon and text
type Status struct {
	Icon    string
	Message string
	widget  interface{}
}

// NewStatus creates a new status widget
func NewStatus(icon, message string) (*Status, error) {
	// TODO: Create GTK4 box with icon and label
	// TODO: Apply theme
	
	return &Status{
		Icon:    icon,
		Message: message,
	}, nil
}

// Update updates the status display
func (s *Status) Update(icon, message string) error {
	s.Icon = icon
	s.Message = message
	
	// TODO: Update GTK widgets
	
	return nil
}

// GetWidget returns the underlying GTK4 widget
func (s *Status) GetWidget() interface{} {
	return s.widget
}
