#!/bin/bash

# File to load the environmental variables from
envFile=launch.env

# Ensure launch.env exists, create default if not
if [ ! -f $envFile ]; then
  cat >$envFile <<EOF
PAUTH_CDN="https://cdn.paymentauth.me:8823"
PAUTH_PLACEMENT_TOKEN="CHANGE-ME"
PAUTH_PROJECT_ID="CHANGE-ME"
EOF
  echo "Missing $envFile created. Edit the file and re-run."
  exit
fi

# Apply environmental variables
echo "Using $envFile template"
set -a
source $envFile
set +a

# Check if set and not default values
[ -z "$PAUTH_CDN" ] && echo "WARNING: PAUTH_CDN variable not set!"
[[ -z "$PAUTH_PLACEMENT_TOKEN" ]] || [ "$PAUTH_PLACEMENT_TOKEN" == "CHANGE-ME" ] && echo "WARNING: PAUTH_PLACEMENT_TOKEN variable not set!"
[[ -z "$PAUTH_PROJECT_ID" ]] || [ "$PAUTH_PROJECT_ID" == "CHANGE-ME" ] && echo "WARNING: PAUTH_PROJECT_ID variable not set!"

# Launch example
php -S 127.0.0.1:9180 -t ./src
