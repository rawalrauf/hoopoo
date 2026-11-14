// Universal overlay window system for persistent UI elements
//
// This package provides a GTK4-based overlay system that can show
// persistent, collapsible UI elements over the desktop. Designed to
// be content-agnostic - any widget can be displayed in the overlay.

package overlay

import (
	"fmt"
)

// Overlay represents a persistent, collapsible overlay window
type Overlay struct {
	// Position settings
	X, Y   int
	Width  int
	Height int
	
	// State
	Collapsed bool
	Visible   bool
	
	// Callbacks
	OnClose    func()
	OnCollapse func()
	OnExpand   func()
}

// Config holds overlay configuration
type Config struct {
	Title     string
	X, Y      int    // Position on screen
	Width     int    // Window width
	Height    int    // Window height
	Resizable bool   // Can user resize?
}

// New creates a new overlay window
func New(config Config) (*Overlay, error) {
	// TODO: Initialize GTK4 window
	// TODO: Set up layer-shell for proper positioning
	// TODO: Apply theme
	
	o := &Overlay{
		X:      config.X,
		Y:      config.Y,
		Width:  config.Width,
		Height: config.Height,
	}
	
	return o, nil
}

// Show displays the overlay
func (o *Overlay) Show() error {
	if o.Visible {
		return nil
	}
	
	// TODO: Show GTK window
	o.Visible = true
	return nil
}

// Hide hides the overlay
func (o *Overlay) Hide() error {
	if !o.Visible {
		return nil
	}
	
	// TODO: Hide GTK window
	o.Visible = false
	return nil
}

// Collapse minimizes the overlay to a small indicator
func (o *Overlay) Collapse() error {
	if o.Collapsed {
		return nil
	}
	
	// TODO: Animate collapse
	o.Collapsed = true
	if o.OnCollapse != nil {
		o.OnCollapse()
	}
	return nil
}

// Expand restores the overlay to full size
func (o *Overlay) Expand() error {
	if !o.Collapsed {
		return nil
	}
	
	// TODO: Animate expand
	o.Collapsed = false
	if o.OnExpand != nil {
		o.OnExpand()
	}
	return nil
}

// SetContent sets the widget content of the overlay
func (o *Overlay) SetContent(content interface{}) error {
	// TODO: Accept any GTK4 widget
	// This is where download progress, status, etc. will be shown
	return fmt.Errorf("not implemented yet")
}

// Close closes and destroys the overlay
func (o *Overlay) Close() error {
	if o.OnClose != nil {
		o.OnClose()
	}
	
	// TODO: Destroy GTK window
	o.Visible = false
	return nil
}
