#!/usr/bin/env bash

# Install pre-commit

color_red="\033[31m"
color_yellow="\033[33m"
color_green="\033[32m"
color_blue="\033[34m"
color_magenta="\033[35m"
color_cyan="\033[36m"
color_white="\033[37m"
color_reset="\033[0m"

if [ "$EUID" -ne 0 ]; then
  as_root=false
else
  as_root=true
fi
  
assert_command() {
    cmd="$1"
    install_cmd="$2"
    echo -e -n "${color_green}Checking${color_reset} if $cmd is installed... "
    if command -v $cmd &> /dev/null
    then
        echo -e "${color_green}✔${color_reset}"
        return 0
    fi

    if [ ! -z "$3" ]
    then
        echo -e "${color_red}✘${color_reset} : $3"
        exit 1
    fi
    
    echo -e -n "${color_yellow}✘${color_reset} : Trying to install... "
    output=$($install_cmd 2>&1)
    if [ $? -ne 0 ]; then
        echo -e "${color_red}✘${color_reset}"
        echo "$output"
        exit 1
    fi
    echo -e "${color_green}✔${color_reset}"
    return 0
}

assert_sudo() {
    if [ "$EUID" -eq 0 ]; then
        echo -e "${color_red}✘${color_reset} : Please do not run as root"
        exit 1
    fi    
}


assert_pre_commit() {
    assert_command "pre-commit" "Install pre-commit using 'sudo apt install -y pre-commit'" 1
}

assert_golangci_lint() {
    assert_command "golangci-lint" "Cannot install golangci-lint. Access https://golangci-lint.run/install/ to install it manually." "1"
}

assert_commitlint() {
    assert_command "commitlint" "go install github.com/conventionalcommit/commitlint@latest"
}

assert_govulncheck() {
    assert_command "govulncheck" "go install golang.org/x/vuln/cmd/govulncheck@latest"
}

install_pre_commit() {
    pre-commit autoupdate
    if pre-commit install -t commit-msg -t pre-commit; then
        echo -e "${color_green}✔${color_reset} : Pre-commit hooks installed"
    else
        echo -e "${color_red}✘${color_reset} : Pre-commit hooks installation skipped (may need manual setup)"
        exit 1
    fi
}

echo -e "${color_yellow}GS-DEV: Setting up your project for development...${color_reset}"

assert_sudo
assert_pre_commit
assert_golangci_lint
assert_commitlint
assert_govulncheck
install_pre_commit




# pre-commit install
# pre-commit install-hooks

# echo "pre-commit hooks installed"