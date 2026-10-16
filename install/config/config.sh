# Copy over Hoopoo configs
mkdir -p ~/.config
cp -R ~/.local/share/hoopoo/config/* ~/.config/

# Use default bashrc from Hoopoo
cp ~/.local/share/hoopoo/default/bashrc ~/.bashrc
