#!/bin/bash

_run_gs_dev() {
  doesnt_use_output=(GS_DOESNT_USE_OUTPUT)

  # debug=(( "$1" == "DEBUG" ? "--debug" : ""))
  if [[ "$1" == "DEBUG" ]]; then
    debug="--debug"
  else
    debug=""
  fi

  shift 1

  if [[ ${doesnt_use_output[@]} =~ $1 ]]; then
    GS_DEV $debug $@
  else
    GS_DEV $debug $@ && _gsdev_treat_output
  fi
}


_gsdev_treat_output() {
  if [[ ! -f "GS_OUTPUT" ]]; then
    echo "No output file found: GS_OUTPUT"
    return
  fi
  source "GS_OUTPUT"
  rm "GS_OUTPUT"
  stty sane
}

dev() {
  _run_gs_dev NODEBUG $@
}

devdbg() {
  _run_gs_dev DEBUG $@
}

dev_1() {
  GS_DEV --output="GS_OUTPUT" $@ && _dev
}

devdbg_1() {
  GS_DEV --debug --output="GS_OUTPUT" $@ && _dev
}

echo "GS_TOOL is ready to use (dev, devdbg)"
