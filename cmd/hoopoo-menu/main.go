// hoopoo:summary=Launch the Hoopoo Menu or jump straight to a submenu
// hoopoo:args=[submenu]

package main

import (
  "fmt"
  "os"
  "strings"

  "hoopoo/pkg/menu"
)

func main() {
  // Toggle existing menu if already open
  menu.ToggleExistingMenu()

  // Direct submenu mode: hoopoo menu trigger
  if len(os.Args) > 1 {
    submenu := os.Args[1]
    menuFunc := getSubmenu(submenu)
    if menuFunc == nil {
      fmt.Fprintf(os.Stderr, "Unknown submenu: %s\n", submenu)
      os.Exit(1)
    }
    
    // Get menu items
    items := menuFunc()
    
    // If single item, execute directly without showing menu
    if len(items) == 1 {
      if err := items[0].Action.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Action error: %v\n", err)
        os.Exit(1)
      }
      os.Exit(0)
    }
    
    // Show submenu with appropriate limit and exit on cancel
    showMenuLoopWithLimit(submenu, menuFunc, true, getMenuLimit(submenu))
    return
  }

  // Show main menu with navigation loop
  showMenuLoop("Go", MainMenu, false)
}

// showMenuLoop displays menu and handles navigation
func showMenuLoop(title string, menuFunc func() []menu.MenuItem, exitOnCancel bool) {
  showMenuLoopWithLimit(title, menuFunc, exitOnCancel, 0)
}

// showMenuLoopWithLimit displays menu with visible item limit and handles navigation
func showMenuLoopWithLimit(title string, menuFunc func() []menu.MenuItem, exitOnCancel bool, maxVisible int) {
  for {
    items := menuFunc()

    var selected *menu.MenuItem
    var err error

    if maxVisible > 0 {
      selected, err = menu.ShowWithLimit(title, items, maxVisible)
    } else {
      selected, err = menu.Show(title, items)
    }

    if err != nil {
      // User cancelled
      if exitOnCancel {
        os.Exit(0)
      }
      return // Go back to parent
    }

    // Execute the selected action
    if submenu, ok := selected.Action.(*menu.SubmenuAction); ok {
      // Navigate to submenu with appropriate limit
      showMenuLoopWithLimit(selected.Label, submenu.MenuFunc, false, getMenuLimit(selected.Label))
    } else {
      // Execute command/function and EXIT (close walker)
      if err := selected.Action.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Action error: %v\n", err)
      }
      // Exit after executing action (fixes walker staying open)
      os.Exit(0)
    }
  }
}

// getMenuLimit returns the visible item limit for specific menus (0 = no limit)
func getMenuLimit(name string) int {
  if strings.EqualFold(name, "font") {
    return 6 // Font menu shows max 6 items for scaled displays
  }
  return 0
}

// getSubmenu returns the menu function for a given name (simple switch statement)
func getSubmenu(name string) func() []menu.MenuItem {
  switch strings.ToLower(name) {
  // Main menu items (direct access via keybinds)
  case "apps":
    // Extract Apps item from MainMenu
    return func() []menu.MenuItem {
      items := MainMenu()
      for _, item := range items {
        if item.Label == "Apps" {
          return []menu.MenuItem{item}
        }
      }
      return nil
    }
  case "learn":
    return LearnMenu
  case "trigger":
    return TriggerMenu
  case "appearance":
    return AppearanceMenu
  case "setup":
    return SetupMenu
  case "install":
    return InstallMenu
  case "remove":
    return RemoveMenu
  case "update":
    return UpdateMenu
  case "about":
    return AboutMenu
  case "system":
    return SystemMenu

  // Trigger submenu shortcuts
  case "reminder":
    return ReminderMenu
  case "capture":
    return CaptureMenu
  case "share":
    return ShareMenu
  case "toggle":
    return TogglesMenu
  case "hardware":
    return HardwareMenu

  // Capture submenu shortcuts
  case "screenshot":
    return ScreenshotMenu
  case "screenrecord":
    return ScreenrecordMenu

  // Appearance submenu shortcuts
  case "theme":
    // Extract Theme item from AppearanceMenu
    return func() []menu.MenuItem {
      items := AppearanceMenu()
      for _, item := range items {
        if item.Label == "Theme" {
          return []menu.MenuItem{item}
        }
      }
      return nil
    }
  case "background":
    // Extract Background item from AppearanceMenu
    return func() []menu.MenuItem {
      items := AppearanceMenu()
      for _, item := range items {
        if item.Label == "Background" {
          return []menu.MenuItem{item}
        }
      }
      return nil
    }
  case "font":
    return FontMenu
  case "windows":
    return WindowsMenu

  // Setup submenu shortcuts
  case "power":
    return PowerProfileMenu

  default:
    return nil
  }
}
