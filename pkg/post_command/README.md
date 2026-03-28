# Post command feature

This allow the program to execute a command in the shell after the program exits.

## Components

### Program wrapper

A shell function that will call the program, wait until it exits, and get the post command from a temporary file and execute it.
