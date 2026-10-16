# Show installation environment variables
gum log --level info "Installation Environment:"

env | grep -E "^(HOOPOO_CHROOT_INSTALL|HOOPOO_ONLINE_INSTALL|HOOPOO_USER_NAME|HOOPOO_USER_EMAIL|USER|HOME|HOOPOO_REPO|HOOPOO_REF|HOOPOO_PATH)=" | sort | while IFS= read -r var; do
  gum log --level info "  $var"
done
