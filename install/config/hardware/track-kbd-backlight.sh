# Firmware on some laptops (ThinkPads, Dell Latitudes) owns the keyboard
# backlight and zeroes it when the lid closes, before any userspace process
# runs. This poller tracks backlight changes (firmware Fn keys, keybinds,
# sleep) into a persistent state file so hoopoo-brightness-keyboard restore
# can recover after wake. On other hardware the Hoopoo keybinds already
# save state, so the poller is not needed and is not installed.

if hoopoo-hw-lenovo-thinkpad || hoopoo-hw-dell-latitude-7490; then
  mkdir -p ~/.config/systemd/user
  cp "$HOOPOO_PATH/default/systemd/user/hoopoo-brightness-keyboard-track.service" ~/.config/systemd/user/
  systemctl --user daemon-reload
  systemctl --user enable --now hoopoo-brightness-keyboard-track.service
fi
