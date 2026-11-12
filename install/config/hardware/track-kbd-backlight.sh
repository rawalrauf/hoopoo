# Firmware on some laptops (ThinkPads, Dell Latitudes) owns the keyboard
# backlight and zeroes it when the lid closes, before any userspace process
# runs. This tracker monitors backlight changes (firmware Fn keys, keybinds,
# sleep) into a persistent state file so hoopoo-brightness-keyboard restore
# can recover after wake. On other hardware the Hoopoo keybinds already
# save state, so the tracker is not needed and is not installed.

# Exit if no keyboard backlight exists
brightness_file=""
for class in leds backlight; do
  for candidate in /sys/class/$class/*kbd_backlight*; do
    if [[ -e $candidate/brightness ]]; then
      brightness_file="$candidate/brightness"
      break 2
    fi
  done
done
[[ -z $brightness_file ]] && exit 0

# Auto-detect which tracker to install based on hardware capabilities
service=""

# 1. Try to find hw_changed file (event-driven, best performance)
for hw_changed in \
  /sys/devices/platform/*/leds/*kbd_backlight*/brightness_hw_changed \
  /sys/devices/platform/*/backlight/*kbd_backlight*/brightness_hw_changed \
  /sys/class/leds/*kbd_backlight*/brightness_hw_changed \
  /sys/class/backlight/*kbd_backlight*/brightness_hw_changed; do
  if [[ -f $hw_changed ]]; then
    service="kbd-track-inotify.service"
    echo "Installing event-driven keyboard backlight tracker (inotify)"
    break
  fi
done

# 2. If no hw_changed found, use polling (safe default for all laptops)
if [[ -z $service ]]; then
  service="kbd-track-polling.service"
  echo "Installing polling keyboard backlight tracker"
fi

# Install selected service
mkdir -p ~/.config/systemd/user
cp "$HOOPOO_PATH/default/systemd/backlight/$service" ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now "$service"
