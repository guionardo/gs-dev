__godev_binds=""

function _dev() {
  if [[ ! -f "{GS_OUTPUT}" ]]; then
    echo "No output file found: {GS_OUTPUT}"
    exit 1
  fi
  source "{GS_OUTPUT}"
  rm "{GS_OUTPUT}"
  stty sane
}

dev() {
  {GS_DEV} dev --output="{GS_OUTPUT}" $@ && _dev  
}

devdbg() {
  {GS_DEV} --debug dev --output="{GS_OUTPUT}" $@ && _dev
}

__GS_DEV_key_binding() {
  bind -r "\C-e"
  bind -x '"\C-e":devc'
  echo "gs-dev: Ctrl-e is bound to console"
}

__GS_DEV_key_binding

echo "gs-dev is ready to use (dev, devdbg)"