# Set identification from install inputs
if [[ -n ${HOOPOO_USER_NAME//[[:space:]]/} ]]; then
  git config --global user.name "$HOOPOO_USER_NAME"
fi

if [[ -n ${HOOPOO_USER_EMAIL//[[:space:]]/} ]]; then
  git config --global user.email "$HOOPOO_USER_EMAIL"
fi
