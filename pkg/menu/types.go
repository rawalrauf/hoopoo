// Menu system types and interfaces

package menu

// MenuItem represents a single menu option
type MenuItem struct {
  Icon   string
  Label  string
  Action Action
}

// Action is executed when a menu item is selected
type Action interface {
  Execute() error
}

// CommandAction runs a shell command
type CommandAction struct {
  Command string
  Args    []string
}

// SubmenuAction shows another menu
type SubmenuAction struct {
  MenuFunc func() []MenuItem
}

// FuncAction runs a custom Go function
type FuncAction struct {
  Fn func() error
}
