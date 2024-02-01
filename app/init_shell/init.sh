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

_last_calendar_show=0
show_calendar() {
  local _elapsed
  _elapsed=$(( $SECONDS - $_last_calendar_show ))
  if [ $_elapsed -gt 300 ]; then
    GS_DEV calendar list
    _last_calendar_show=$SECONDS
  fi
}

trap show_calendar DEBUG

echo "GS_TOOL is ready to use (dev, devdbg)"
