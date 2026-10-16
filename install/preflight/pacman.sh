if [[ -n ${HOOPOO_ONLINE_INSTALL:-} ]]; then
  # Install build tools
  hoopoo-pkg-add base-devel

  # Configure pacman
  sudo cp -f ~/.local/share/hoopoo/default/pacman/pacman-${HOOPOO_MIRROR:-stable}.conf /etc/pacman.conf
  sudo cp -f ~/.local/share/hoopoo/default/pacman/mirrorlist /etc/pacman.d/mirrorlist

  # Import Hoopoo GPG key and install keyring
  sudo pacman-key --recv-keys 070030E5D914817BA797F66AAF1C8F036A5E73FD --keyserver keys.openpgp.org || \
    sudo pacman-key --recv-keys 070030E5D914817BA797F66AAF1C8F036A5E73FD --keyserver keyserver.ubuntu.com
  sudo pacman-key --lsign-key 070030E5D914817BA797F66AAF1C8F036A5E73FD
  hoopoo-pkg-add hoopoo-keyring
  sudo pacman -Sy

  # Refresh all repos
  sudo pacman -Syyuu --noconfirm
fi
