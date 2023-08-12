#!/bin/bash

_dev() {
  if [[ ! -f "GS_OUTPUT" ]]; then
    echo "No output file found: GS_OUTPUT"
    exit 1
  fi
  source "GS_OUTPUT"
  rm "GS_OUTPUT"
  stty sane
}

dev() {
  GS_DEV --output="GS_OUTPUT" $@ && _dev
}

devdbg() {
  GS_DEV --debug --output="GS_OUTPUT" $@ && _dev
}

echo "GS_TOOL is ready to use (dev, devdbg)"
