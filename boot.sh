#!/bin/bash

set -e

# Set install mode to online since boot.sh is used for curl installations
export HOOPOO_ONLINE_INSTALL=true

ansi_art='
████      ████  ██████████████  ██████████████  █████████████  ██████████████  ██████████████
████      ████  ████      ████  ████      ████  ████     ████  ████      ████  ████      ████
██████████████  ████      ████  ████      ████  ████     ████  ████      ████  ████      ████
                ████      ████  ████      ████  █████████████  ████      ████  ████      ████
██████████████  ████      ████  ████      ████  ████           ████      ████  ████      ████
████      ████  ████      ████  ████      ████  ████           ████      ████  ████      ████
████      ████  ██████████████  ██████████████  ████           ██████████████  ██████████████
'

clear
echo -e "\n$ansi_art\n"

# Use custom branch if instructed, otherwise default to master
HOOPOO_REF="${HOOPOO_REF:-master}"

# Set mirror based on branch
if [[ $HOOPOO_REF == "dev" ]]; then
  export HOOPOO_MIRROR=edge
else
  export HOOPOO_MIRROR=stable
fi

printf 'Server = https://geo.mirror.pkgbuild.com/$repo/os/$arch\nServer = https://mirror.rackspace.com/archlinux/$repo/os/$arch\nServer = https://mirrors.kernel.org/archlinux/$repo/os/$arch\n' | sudo tee /etc/pacman.d/mirrorlist >/dev/null

sudo pacman -Syu --noconfirm --needed git

# Use custom repo if specified, otherwise default to rawalrauf/hoopoo
HOOPOO_REPO="${HOOPOO_REPO:-rawalrauf/hoopoo}"

echo -e "\nCloning Hoopoo from: https://github.com/${HOOPOO_REPO}.git"
echo -e "\e[32mUsing branch: $HOOPOO_REF\e[0m"
rm -rf ~/.local/share/hoopoo/
git clone --branch "$HOOPOO_REF" "https://github.com/${HOOPOO_REPO}.git" ~/.local/share/hoopoo >/dev/null

echo -e "\nInstallation starting..."
source ~/.local/share/hoopoo/install.sh
