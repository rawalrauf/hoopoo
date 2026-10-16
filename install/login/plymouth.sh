if [[ $(plymouth-set-default-theme) != "hoopoo" ]]; then
  sudo cp -r "$HOME/.local/share/hoopoo/default/plymouth" /usr/share/plymouth/themes/hoopoo/
  sudo plymouth-set-default-theme hoopoo
fi
