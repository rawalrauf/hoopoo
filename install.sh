#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -eEo pipefail

# Define Hoopoo locations
export HOOPOO_PATH="$HOME/.local/share/hoopoo"
export HOOPOO_INSTALL="$HOOPOO_PATH/install"
export HOOPOO_INSTALL_LOG_FILE="/var/log/hoopoo-install.log"
export PATH="$HOOPOO_PATH/bin:$PATH"

# Install
source "$HOOPOO_INSTALL/helpers/all.sh"
source "$HOOPOO_INSTALL/preflight/all.sh"
source "$HOOPOO_INSTALL/packaging/all.sh"
source "$HOOPOO_INSTALL/config/all.sh"
source "$HOOPOO_INSTALL/login/all.sh"
source "$HOOPOO_INSTALL/post-install/all.sh"
