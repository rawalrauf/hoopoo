// Hook execution utilities

package hooks

import (
  "os/exec"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func Run(hookName string, args ...string) error {
  hookPath := filepath.Join(paths.Bin(), "hoopoo-hook")
  hookArgs := append([]string{hookName}, args...)
  return exec.Command(hookPath, hookArgs...).Run()
}
