if [[ -n ${OMARCHY_ONLINE_INSTALL:-} ]]; then
  # Install build tools
  omarchy-pkg-add base-devel

  # Configure pacman
  sudo cp -f ~/.local/share/omarchy/default/pacman/pacman-${OMARCHY_MIRROR:-stable}.conf /etc/pacman.conf
  sudo cp -f ~/.local/share/omarchy/default/pacman/mirrorlist-${OMARCHY_MIRROR:-stable} /etc/pacman.d/mirrorlist

  # TODO: Once hoopoo-keyring is created, import GPG key and install keyring here:
  #   sudo pacman-key --recv-keys <HOOPOO_KEY_ID> --keyserver keys.openpgp.org
  #   sudo pacman-key --lsign-key <HOOPOO_KEY_ID>
  #   omarchy-pkg-add hoopoo-keyring
  sudo pacman -Sy

  # Refresh all repos
  sudo pacman -Syyuu --noconfirm
fi
