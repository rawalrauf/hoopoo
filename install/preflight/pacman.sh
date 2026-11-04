if [[ -n ${HOOPOO_ONLINE_INSTALL:-} ]]; then
  # Install build tools
  hoopoo-pkg-add base-devel

  # Configure pacman
  sudo cp -f ~/.local/share/hoopoo/default/pacman/pacman-${HOOPOO_MIRROR:-stable}.conf /etc/pacman.conf
  sudo cp -f ~/.local/share/hoopoo/default/pacman/mirrorlist /etc/pacman.d/mirrorlist

  # Import Hoopoo GPG key and install keyring
  sudo pacman-key --recv-keys D3931CBE26A484582489EEB7394AFBE6E502DB8F --keyserver keyserver.ubuntu.com || \
    sudo pacman-key --recv-keys D3931CBE26A484582489EEB7394AFBE6E502DB8F --keyserver keys.openpgp.org
  sudo pacman-key --lsign-key D3931CBE26A484582489EEB7394AFBE6E502DB8F
  hoopoo-pkg-add hoopoo-keyring
  sudo pacman -Sy

  # Refresh all repos
  sudo pacman -Syyuu --noconfirm
fi
