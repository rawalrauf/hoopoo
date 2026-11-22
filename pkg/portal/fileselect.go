// File selection via XDG FileChooser portal

package portal

import (
  "errors"
  "fmt"
  "net/url"
  "strings"
  "time"
  
  "github.com/godbus/dbus/v5"
)

var (
  // ErrNoFilesSelected indicates user cancelled or selected nothing
  ErrNoFilesSelected = errors.New("no files selected")
  // ErrSelectionCancelled indicates user explicitly cancelled
  ErrSelectionCancelled = errors.New("file selection cancelled")
)

// FileSelectOptions configures file selection dialog
type FileSelectOptions struct {
  Title      string   // Dialog title
  Multiple   bool     // Allow multiple file selection
  Directory  bool     // Select directories instead of files
  Extensions []string // File extensions to filter (e.g., "jpg", "png")
}

// SelectFiles opens a file selection dialog and returns selected paths
func SelectFiles(opts FileSelectOptions) ([]string, error) {
  conn, err := NewConnection()
  if err != nil {
    return nil, err
  }
  defer conn.Close()
  
  if opts.Title == "" {
    opts.Title = "Select file"
  }
  
  // Generate token and predicted path
  token := conn.generateToken("hoopoo")
  predictedPath := conn.predictPath(token)
  
  // Subscribe to response signal
  responseChan := make(chan *dbus.Signal, 1)
  conn.conn.Signal(responseChan)
  
  if err := conn.subscribePath(predictedPath); err != nil {
    return nil, fmt.Errorf("failed to subscribe: %w", err)
  }
  defer conn.unsubscribePath(predictedPath)
  
  // Build options
  options := map[string]dbus.Variant{
    "handle_token": dbus.MakeVariant(token),
    "multiple":     dbus.MakeVariant(opts.Multiple),
  }
  
  if opts.Directory {
    options["directory"] = dbus.MakeVariant(true)
  }
  
  // Add file filters (only for files, not directories)
  if len(opts.Extensions) > 0 && !opts.Directory {
    // Build pattern list: array of (uint32, string) tuples
    type filterPattern struct {
      Type    uint32
      Pattern string
    }
    
    var patterns []filterPattern
    for _, ext := range opts.Extensions {
      ext = strings.TrimPrefix(ext, ".")
      // Add both lowercase and uppercase patterns
      patterns = append(patterns, filterPattern{0, "*." + strings.ToLower(ext)})
      patterns = append(patterns, filterPattern{0, "*." + strings.ToUpper(ext)})
    }
    
    // Build filter: (string, array of patterns)
    type fileFilter struct {
      Label    string
      Patterns []filterPattern
    }
    
    label := strings.Join(opts.Extensions, " ")
    filter := fileFilter{Label: label, Patterns: patterns}
    
    options["filters"] = dbus.MakeVariant([]fileFilter{filter})
    options["current_filter"] = dbus.MakeVariant(filter)
  }
  
  // Call OpenFile portal method
  obj := conn.conn.Object(portalBus, portalPath)
  call := obj.Call("org.freedesktop.portal.FileChooser.OpenFile", 0, "", opts.Title, options)
  if call.Err != nil {
    return nil, fmt.Errorf("failed to open file chooser: %w", call.Err)
  }
  
  var handle string
  if err := call.Store(&handle); err != nil {
    return nil, fmt.Errorf("failed to get handle: %w", err)
  }
  
  // Subscribe to actual handle if different from predicted
  if handle != predictedPath {
    if err := conn.subscribePath(handle); err != nil {
      return nil, fmt.Errorf("failed to subscribe to handle: %w", err)
    }
    defer conn.unsubscribePath(handle)
  }
  
  // Wait for response or timeout
  timeout := time.After(answerTimeout)
  
  select {
  case sig := <-responseChan:
    return parseFileResponse(sig)
    
  case <-timeout:
    return nil, fmt.Errorf("timeout waiting for file selection")
  }
}

// parseFileResponse extracts file paths from portal response signal
func parseFileResponse(sig *dbus.Signal) ([]string, error) {
  if len(sig.Body) < 2 {
    return nil, ErrNoFilesSelected
  }
  
  code, ok := sig.Body[0].(uint32)
  if !ok || code != 0 {
    return nil, ErrSelectionCancelled
  }
  
  results, ok := sig.Body[1].(map[string]dbus.Variant)
  if !ok {
    return nil, fmt.Errorf("invalid response format")
  }
  
  urisVariant, ok := results["uris"]
  if !ok {
    return nil, ErrNoFilesSelected
  }
  
  uris, ok := urisVariant.Value().([]string)
  if !ok || len(uris) == 0 {
    return nil, ErrNoFilesSelected
  }
  
  // Convert file:// URIs to local paths
  var paths []string
  for _, uri := range uris {
    u, err := url.Parse(uri)
    if err != nil {
      continue
    }
    paths = append(paths, u.Path)
  }
  
  if len(paths) == 0 {
    return nil, ErrNoFilesSelected
  }
  
  return paths, nil
}
