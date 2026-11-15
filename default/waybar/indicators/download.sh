#!/bin/bash

# Multiple downloads - show count if more than 1
PROGRESS_DIR="/tmp/hoopoo-downloads"

if [[ ! -d $PROGRESS_DIR ]]; then
  echo '{"text": ""}'
  exit 0
fi

count=0
single_percent=""

for file in "$PROGRESS_DIR"/*; do
  if [[ -f $file ]]; then
    percent=$(cat "$file" 2>/dev/null)
    if [[ $percent =~ ^[0-9]+$ ]]; then
      ((count++))
      single_percent="$percent"
    fi
  fi
done

if [[ $count -eq 0 ]]; then
  echo '{"text": ""}'
elif [[ $count -eq 1 ]]; then
  echo "{\"text\": \"󰇚 ${single_percent}%\", \"tooltip\": \"Downloading video\", \"class\": \"active\"}"
else
  echo "{\"text\": \"󰇚 ${count}\", \"tooltip\": \"${count} downloads in progress\", \"class\": \"active\"}"
fi
