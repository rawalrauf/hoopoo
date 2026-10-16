ICON_DIR="$HOME/.local/share/icons/hicolor/512x512@2/apps"

hoopoo-tui-install "Disk Usage" "bash -c 'dust -r; read -n 1 -s'" float "$ICON_DIR/Disk-Usage.png"
hoopoo-tui-install "Docker" "lazydocker" tile "$ICON_DIR/Docker.png"
