# xair-cli

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Use](#use)
- [Shell Completion](#shell-completion)
- [License](#license)

### Installation

*xair-cli*

```console
go install github.com/onyx-and-iris/xair-cli/cmd/xair-cli@latest
```

*x32-cli*
```console
go install github.com/onyx-and-iris/xair-cli/cmd/x32-cli@latest
```

### Configuration

#### Flags

- --host/-H: Host of the mixer.
- --port/-P: Port of the mixer.
- --timeout/-T: Timeout for OSC operations.
- --loglevel/-L: The application's logging verbosity.

Pass `--host` and any other configuration as flags on the root commmand:

```console
xair-cli --host mixer.local --timeout 50ms --help
```

#### Environment Variables

Or you may load them from your environment:

Example xair .envrc:

```bash
#!/usr/bin/env bash

export XAIR_CLI_HOST=mixer.local
export XAIR_CLI_PORT=10024
export XAIR_CLI_TIMEOUT=100ms
export XAIR_CLI_LOGLEVEL=warn
```

Example x32 .envrc:

```bash
#!/usr/bin/env bash

export X32_CLI_HOST=x32.local
export X32_CLI_PORT=10023
export X32_CLI_TIMEOUT=100ms
export X32_CLI_LOGLEVEL=warn
```

---

### Use

For each command/subcommand in the tree there exists a `--help` flag, use it to print usage information.

- [xair-cli](./tools/xair-help.md)
- [x32-cli](./tools/x32-help.md)

### Examples

*Fade out main LR all the way to -∞ over a 5s duration*

```console
xair-cli main fadeout
```

*enable phantom power and set the gain to 28.0dB over a 10s duration for headamp (strip) 09*
```console
xair-cli headamp 9 phantom on

xair-cli headamp 9 gain --duration 10s 18.0
```

*set strip 09 send level for bus 5 to -18.0dB*
```console
xair-cli strip 9 send 5 -- -18.0
```

*enable eq for strip 01*
```console
xair-cli strip 1 eq on true
```

*rename bus 01 to 'vocal mix'*
```console
xair-cli bus 1 name 'vocal mix'
```

*set bus 03 eq band 03 (LoMid) gain*
```console
xair-cli bus 3 eq 3 gain -- -3.5
```

*Send a raw OSC message to the mixer*
```console
xair-cli raw /xinfo

xair-cli raw /ch/01/config/name 'rode podmic'
xair-cli raw /ch/01/config/name
```

*Save current mixer state to a snapshot*
```console
xair-cli snapshot 20 save 'twitch live'
```

---

## Shell Completion

-   completion:

    *optional*
    -   args: Shell
        - If no shell is passed completion will attempt to detect the current login shell.

*xair-cli*

```console
xair-cli completion

xair-cli completion bash
```

*x32-cli*

```console
x32-cli completion

x32-cli completion bash
```

Currently supported shells: *bash* *zsh* *fish*.

---

### License

`xair-cli` is distributed under the terms of the [MIT](https://spdx.org/licenses/MIT.html) license.


### Notes

For an alternative, python implementation consider checking out [dc-xair-cli](https://pypi.org/project/dc-xair-cli/). It supports some operations like batch commands and network discovery which this CLI doesn't (and I have no plans to implement them).
