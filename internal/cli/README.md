# `internal/cli`

Helpers to build `cobra.Command` values from typed command structs.

## Overview

`GenerateCobraCommand` maps exported struct fields into Cobra flags, positional
args, and nested subcommands, then executes `Run(ctx, output)` after parsing.

## Command contract

Command structs must implement:

```go
type CommandStruct interface {
    Run(ctx context.Context, output io.Writer) error
}
```

`GenerateCobraCommand` expects a pointer to a struct (`*MyCommand`). It panics
for unsupported input kinds.

## Supported struct tags

- `flag:"name[,short]"`: defines a flag name and optional shorthand.
- `description:"text"`: flag help description.
- `default:"value"`: default flag value (string format).
- `persistent:""`: marks the flag as persistent.
- `validate:"..."`: enables validation with `validator/v10`.
- `subcommand:"name,usage"`: creates a nested command from a struct field.
- `args:"..."`: binds positional args into a `[]string` field.

## Supported flag field types

- `string`
- `int`
- `bool`
- `[]string` (comma-separated default values)
- `time.Duration` (Go duration format, e.g. `30s`, `1m`)
- `time.Time` (all time.Time templates)

Any unsupported type or invalid default value causes panic while building the
command.

## Example

```go
type HelloCommand struct {
    Name    string        `flag:"name,n" description:"User name" validate:"required"`
    Timeout time.Duration `flag:"timeout" default:"30s"`
    Extra   []string      `flag:"extra" default:"a,b"`
    Args    []string      `args:"*"`
}

func (c *HelloCommand) Run(ctx context.Context, out io.Writer) error {
    _, _ = fmt.Fprintf(out, "name=%s timeout=%s args=%v extra=%v\n", c.Name, c.Timeout, c.Args, c.Extra)
    return nil
}

func build() *cobra.Command {
    cmd := &HelloCommand{}
    return GenerateCobraCommand(cmd, "hello", "Example command")
}
```
