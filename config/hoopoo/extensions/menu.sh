# Overwrite parts of the hoopoo-menu with user-specific submenus.
# See $HOOPOO_PATH/bin/hoopoo-menu for functions that can be overwritten.
#
# WARNING: Overwritten functions will obviously not be updated when Hoopoo changes.
#
# Example of minimal system menu:
#
# show_system_menu() {
#   case $(menu "System" "  Lock\n󰐥  Shutdown") in
#   *Lock*) hoopoo-system-lock ;;
#   *Shutdown*) hoopoo-system-shutdown ;;
#   *) back_to show_main_menu ;;
#   esac
# }
#
# Example of overriding just the about menu action: (Using zsh instead of bash (default))
#
# show_about() {
#   exec hoopoo-launch-or-focus-tui "zsh -c 'fastfetch; read -k 1'"
# }
