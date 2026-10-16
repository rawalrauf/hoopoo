HOOPOO_MIGRATIONS_STATE_PATH=~/.local/state/hoopoo/migrations
mkdir -p $HOOPOO_MIGRATIONS_STATE_PATH

for file in ~/.local/share/hoopoo/migrations/*.sh; do
  touch "$HOOPOO_MIGRATIONS_STATE_PATH/$(basename "$file")"
done
