#!/bin/bash
# shellcheck disable=SC2068
__godev_binds=""

function _dev() {
  if [[ ! -f "{GO_OUTPUT}" ]]; then
    echo "No output file found: {GO_OUTPUT}"
    exit 1
  fi
  source "{GO_OUTPUT}"
  rm "{GO_OUTPUT}"
  stty sane
}

dev() {
  {GO_DEV} go --output="{GO_OUTPUT}" $@ && _dev  
}

devc() {
  {GO_DEV} console --output="{GO_OUTPUT}" && _dev
}

devdbg() {
  {GO_DEV} --debug go --output="{GO_OUTPUT}" $@ && _dev
}

__go_dev_key_binding() {
  # bind -r "\C-y"
  bind -r "\C-e"
  # bind -r '"^[d"'
  # bind '"^[d":"echo 'alt-d'\n"'
  bind -x '"\C-e":devc'
  # bind -x '"\eOS":"fortune | cowsay"'
  # bind '"^[e":"devc\n"'
  # bind -x '"\e\d": devc'

  echo "go-dev: Ctrl-e is bound to console"
}

__go_dev_key_binding

# [ "$(go-dev config-get enable-alt-d)" = true ] && __go_dev_key_binding
# __history_control_r() {
#         READLINE_LINE=$(HISHTORY_TERM_INTEGRATION=1 hishtory tquery "$READLINE_LINE" | tr -d '\n')
#         READLINE_POINT=0x7FFFFFFF
# }

# __hishtory_bind_control_r() {
#   bind -x '"\C-r": __history_control_r'
# }

# [ "$(hishtory config-get enable-control-r)" = true ] && __hishtory_bind_control_r


echo "go-dev is ready to use (dev, devdbg)"