# Install all base packages
mapfile -t packages < <(grep -v '^#' "$HOOPOO_INSTALL/hoopoo-base.packages" | grep -v '^$')
hoopoo-pkg-add "${packages[@]}"
