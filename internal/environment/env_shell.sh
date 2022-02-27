#!/usr/bin/env bash
# __DESCRIPTION__

GSDEV_BIN="_GSDEV_BIN_"
ENV_SHELL_FILE="_ENV_SHELL_FILE_"

function dev() {
    output=$(mktmp _TOOLNAME_.XXX)
    $GSDEV_BIN dev go --output $output $@
    source $output
    rm $output
}