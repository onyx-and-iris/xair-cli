# xair-cli

### Installation

```console
go install github.com/onyx-and-iris/xair-cli@latest
```

### Use

```console
xair-cli is a command-line tool that allows users to send OSC messages
to Behringer X Air mixers for remote control and configuration. It supports
various commands to manage mixer settings directly from the terminal.

Usage:
  xair-cli [flags]
  xair-cli [command]

Available Commands:
  bus         Commands to control individual buses
  completion  Generate the autocompletion script for the specified shell
  headamp     Commands to control headamp gain and phantom power
  help        Help about any command
  main        Commands to control the main output
  strip       Commands to control individual strips

Flags:
  -h, --help              help for xair-cli
  -H, --host string       host address of the X Air mixer (default "mixer.local")
  -k, --kind string       Kind of mixer (xair, x32) (default "xair")
  -l, --loglevel string   Log level (debug, info, warn, error, fatal, panic) (default "warn")
  -p, --port int          Port number of the X Air mixer (default 10024)
  -v, --version           version for xair-cli

Use "xair-cli [command] --help" for more information about a command.
```
