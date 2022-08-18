#!/usr/bin/env bash
echo "Running splice router"

API_VERSION="--api-version=2"

if [[ "$DLV_API_VERSION1" == "true" ]]; then
  API_VERSION="--api-version=1"
fi

if ! [[ "$DEBUG_FLAG" == "false" ]]; then
  if [[ "$WAIT_FOR_ATTACH_FLAG" == "true" ]]; then
    echo "Starting splice router in Wait mode. Service will start executing once debugger is connected."
  fi
  pkill dlv
  #shellcheck disable=2046
  dlv exec$([[ "$WAIT_FOR_ATTACH_FLAG" == "true" ]] && echo "" || echo " --continue") \
    --listen=:40000 --log --headless $API_VERSION --accept-multiclient "/tmp/localdev router"
else
  "/tmp/localdev router"
fi
