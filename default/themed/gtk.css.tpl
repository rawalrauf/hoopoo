/* Hoopoo GTK4 Color Definitions
 * Generated from current theme colors
 */

@define-color accent {{ accent }};
@define-color background {{ background }};
@define-color foreground {{ foreground }};
@define-color color0 {{ color0 }};
@define-color color1 {{ color1 }};
@define-color color2 {{ color2 }};
@define-color color3 {{ color3 }};
@define-color color4 {{ color4 }};
@define-color color5 {{ color5 }};
@define-color color6 {{ color6 }};
@define-color color7 {{ color7 }};
@define-color color8 {{ color8 }};
@define-color color9 {{ color9 }};
@define-color color10 {{ color10 }};
@define-color color11 {{ color11 }};
@define-color color12 {{ color12 }};
@define-color color13 {{ color13 }};
@define-color color14 {{ color14 }};
@define-color color15 {{ color15 }};

/* Apply colors to GTK/Adwaita widgets */
@define-color accent_bg_color @accent;
@define-color accent_fg_color @background;
@define-color accent_color @accent;
@define-color destructive_bg_color @color1;
@define-color destructive_fg_color @background;
@define-color success_bg_color @color2;
@define-color success_fg_color @background;
@define-color warning_bg_color @color3;
@define-color warning_fg_color @background;
@define-color error_bg_color @color1;
@define-color error_fg_color @background;
@define-color window_bg_color @background;
@define-color window_fg_color @foreground;
@define-color view_bg_color @background;
@define-color view_fg_color @foreground;
@define-color headerbar_bg_color @background;
@define-color headerbar_fg_color @foreground;
@define-color headerbar_border_color @color8;
@define-color headerbar_backdrop_color @background;
@define-color headerbar_shade_color alpha(@color0, 0.5);
@define-color card_bg_color alpha(@foreground, 0.05);
@define-color card_fg_color @foreground;
@define-color card_shade_color alpha(@color0, 0.5);
@define-color dialog_bg_color @background;
@define-color dialog_fg_color @foreground;
@define-color popover_bg_color @background;
@define-color popover_fg_color @foreground;
@define-color shade_color alpha(@color0, 0.5);
@define-color scrollbar_outline_color alpha(@color0, 0.5);

/* Add this to your Adwaita widget variable mapping section */
@define-color sidebar_bg_color @background;
@define-color sidebar_fg_color @foreground;
@define-color sidebar_backdrop_color @background;
/* Modern GTK4 / Libadwaita Sidebar Styling */
.nautilus-window navigation-sidebar,
.nautilus-window sidebar,
.nautilus-window split-view > sidebar,
sidebar list,
sidebar row {
  background-color: @sidebar_bg_color;
  color: @sidebar_fg_color;
}

/* Hover and selection states */
sidebar row:selected,
navigation-sidebar row:selected {
  background-color: @accent;
  color: @background;
}

sidebar row:hover,
navigation-sidebar row:hover {
  background-color: alpha(@accent, 0.15);
}
