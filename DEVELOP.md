# Development notes

* https://github.com/manifoldco/promptui
* https://github.com/Songmu/prompter
* https://github.com/go-survey/survey
* https://dev.to/divrhino/building-an-interactive-cli-app-with-go-cobra-promptui-346n

Github actions good ideas!

* https://github.com/daytonaio/daytona/tree/main/.github/workflows
* https://github.com/gofr-dev/gofr/blob/development/.github/workflows/go.yml

## Protobuf details

### Install protobuf-compiler

```bash
sudo nala install protobuf-compiler
```

### Install Code generator

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Github spec kit

```bash
pipx install git+https://github.com/github/spec-kit.git@v0.8.1
```
