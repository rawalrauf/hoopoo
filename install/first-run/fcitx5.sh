systemctl --user enable --now hoopoo-fcitx5.service
cp $HOOPOO_PATH/config/fcitx5/conf/xcb.conf ~/.config/fcitx5/conf/xcb.conf 2>/dev/null || true
cp $HOOPOO_PATH/config/fcitx5/conf/clipboard.conf ~/.config/fcitx5/conf/clipboard.conf 2>/dev/null || true
