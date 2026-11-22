// Core menu engine with walker integration

package menu

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "hoopoo/pkg/walker"
  "hoopoo/pkg/system/process"
)

// Show displays a menu and returns the selected item (or error if cancelled)
func Show(title string, items []MenuItem) (*MenuItem, error) {
  return ShowWithLimit(title, items, 0) // 0 = no limit (use actual item count)
}

// ShowWithLimit displays a menu with a maximum visible item limit
func ShowWithLimit(title string, items []MenuItem, maxVisible int) (*MenuItem, error) {
  if len(items) == 0 {
    return nil, fmt.Errorf("no menu items provided")
  }

  // Build menu text
  var options []string
  for _, item := range items {
    display := item.Label
    if item.Icon != "" {
      display = item.Icon + "  " + item.Label
    }
    options = append(options, display)
  }

  // Determine visible item count for walker height calculation
  visibleItems := len(items)
  if maxVisible > 0 && visibleItems > maxVisible {
    visibleItems = maxVisible
  }
  
  choice, err := walker.ShowDmenu(title, strings.Join(options, "\n"), visibleItems)
  if err != nil || choice == "" {
    return nil, fmt.Errorf("cancelled")
  }

  // Find selected item
  for i, item := range items {
    if strings.Contains(choice, item.Label) {
      return &items[i], nil
    }
  }

  return nil, fmt.Errorf("no match found")
}

// ToggleExistingMenu closes walker if already open
func ToggleExistingMenu() {
  if process.IsRunningByPattern("walker.*--dmenu") {
    exec.Command("walker", "--close").Run()
    os.Exit(0)
  }
}

// Execute command action
func (a *CommandAction) Execute() error {
  cmd := exec.Command(a.Command, a.Args...)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin
  return cmd.Run()
}

// Execute submenu action
func (a *SubmenuAction) Execute() error {
  // Return special marker that caller should show submenu
  return nil
}

// Execute function action
func (a *FuncAction) Execute() error {
  return a.Fn()
}
