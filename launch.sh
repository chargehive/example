#!/bin/bash

### You must set these parameters
export PAUTH_CDN="https://cdn.paymentauth.me:8823"
export PAUTH_PLACEMENT_TOKEN="CHANGE-ME"
export PAUTH_PROJECT_ID="CHANGE-ME2"

php -S 127.0.0.1:9180 -t ./src
