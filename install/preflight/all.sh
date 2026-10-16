source $HOOPOO_INSTALL/preflight/guard.sh
source $HOOPOO_INSTALL/preflight/begin.sh
run_logged $HOOPOO_INSTALL/preflight/show-env.sh
run_logged $HOOPOO_INSTALL/preflight/pacman.sh
run_logged $HOOPOO_INSTALL/preflight/migrations.sh
run_logged $HOOPOO_INSTALL/preflight/first-run-mode.sh
run_logged $HOOPOO_INSTALL/preflight/disable-mkinitcpio.sh
