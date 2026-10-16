# Place in each assistant's global skills directory so the Hoopoo skill is available on first install
mkdir -p ~/.agents/skills ~/.claude/skills ~/.codex/skills ~/.pi/agent/skills
ln -sfn "$HOOPOO_PATH/default/hoopoo-skill" ~/.agents/skills/hoopoo
ln -sfn "$HOOPOO_PATH/default/hoopoo-skill" ~/.claude/skills/hoopoo
ln -sfn "$HOOPOO_PATH/default/hoopoo-skill" ~/.codex/skills/hoopoo
ln -sfn "$HOOPOO_PATH/default/hoopoo-skill" ~/.pi/agent/skills/hoopoo
