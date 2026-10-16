# Set default XCompose that is triggered with CapsLock
tee ~/.XCompose >/dev/null <<EOF
# ~/.XCompose — Custom key composition sequences
#
# XCompose lets you type special characters by pressing a compose key
# followed by a short sequence of keys. Your compose key is Right Alt
# (set in ~/.config/niri/input.kdl via options "compose:ralt").
#
# How it works:
#   Right Alt + m + s  →  😄
#   Right Alt + c + o  →  © (from the system compose table)
#   Right Alt + - + -  →  — (em dash)
#
# This file includes two sources:
#   1. The system compose table (%L) — hundreds of special characters
#      built into Linux (accented letters, math symbols, currency, etc.)
#   2. Hoopoo's default sequences — emoji shortcuts via Right Alt + m + letter
#
# To add your own sequences, add lines below in this format:
#   <Multi_key> <key1> <key2> : "result"   # optional comment
#
# Examples:
#   <Multi_key> <space> <n> : "Your Name"    # Right Alt + Space + n
#   <Multi_key> <space> <e> : "you@email.com" # Right Alt + Space + e
#   <Multi_key> <m> <z> : "🤪"               # Right Alt + m + z
#
# After making changes, run: hoopoo-restart-xcompose

# Include hoopoo emoji shortcuts (Right Alt + m + letter)
include "%H/.local/share/hoopoo/default/xcompose"

# Personal shortcuts — customize these
<Multi_key> <space> <n> : "$HOOPOO_USER_NAME"
<Multi_key> <space> <e> : "$HOOPOO_USER_EMAIL"
EOF
