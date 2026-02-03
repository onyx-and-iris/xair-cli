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

### Examples

*Fade out main LR all the way to -∞ over a 5s duration*

```console
xair-cli main fadeout
```

*enable phantom power and set the gain to 28.0dB over a 10s duration for strip 09*
```console
xair-cli headamp phantom 9 on

xair-cli headamp gain 9 28.0 --duration 10s
```

*set strip 09 send level for bus 5 to -18.0dB*
```console
xair-cli strip send 9 5 -- -18.0
```

*enable eq for strip 01*
```console
xair-cli strip eq on 1 true
```

*rename bus 01 to 'vocal mix'*
```console
xair-cli bus name 1 'vocal mix'
```

For every command/subcommand there exists a `--help` flag which you can use to get usage information.


### Notes

This CLI is useful if just want to run some commands on the terminal using a single binary, no further downloads. However, there exists an alternative you should check out that has wider support of the XAir OSC protocol including support for operations like batch commands and network discovery (which I have no plans to implement on this CLI). Check out [dc-xair-cli](https://pypi.org/project/dc-xair-cli/) on pypi.


### License

`xair-cli` is distributed under the terms of the [MIT](https://spdx.org/licenses/MIT.html) license.
