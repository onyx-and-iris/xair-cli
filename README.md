# xair-cli

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

### Use

```console
Usage: xair-cli <command> [flags]

A CLI to control Behringer X-Air mixers.

Flags:
  -h, --help                  Show context-sensitive help.
  -H, --host="mixer.local"    The host of the X-Air device ($XAIR_CLI_HOST).
  -P, --port=10024            The port of the X-Air device ($XAIR_CLI_PORT).
  -T, --timeout=100ms         Timeout for OSC operations ($XAIR_CLI_TIMEOUT).
  -L, --loglevel="warn"       Log level for the CLI ($XAIR_CLI_LOGLEVEL).
  -v, --version               Print xair-cli version information and quit

Commands:
  completion (c)    Generate shell completion scripts.

Raw
  raw    Send raw OSC messages to the mixer.

Main
  main mute              Get or set the mute state of the Main L/R output.
  main fader             Get or set the fader level of the Main L/R output.
  main fadein            Fade in the Main L/R output over a specified duration.
  main fadeout           Fade out the Main L/R output over a specified duration.
  main eq on             Get or set the EQ on/off state of the Main L/R output.
  main eq <band> gain    Get or set the gain of the specified EQ band.
  main eq <band> freq    Get or set the frequency of the specified EQ band.
  main eq <band> q       Get or set the Q factor of the specified EQ band.
  main eq <band> type    Get or set the type of the specified EQ band.
  main comp on           Get or set the compressor on/off state of the Main L/R
                         output.
  main comp mode         Get or set the compressor mode of the Main L/R output.
  main comp threshold    Get or set the compressor threshold of the Main L/R
                         output.
  main comp ratio        Get or set the compressor ratio of the Main L/R output.
  main comp mix          Get or set the compressor mix level of the Main L/R
                         output.
  main comp makeup       Get or set the compressor makeup gain of the Main L/R
                         output.
  main comp attack       Get or set the compressor attack time of the Main L/R
                         output.
  main comp hold         Get or set the compressor hold time of the Main L/R
                         output.
  main comp release      Get or set the compressor release time of the Main L/R
                         output.

Strip
  strip <index> mute              Get or set the mute state of the strip.
  strip <index> fader             Get or set the fader level of the strip.
  strip <index> fadein            Fade in the strip over a specified duration.
  strip <index> fadeout           Fade out the strip over a specified duration.
  strip <index> send              Get or set the send level for a specific bus.
  strip <index> name              Get or set the name of the strip.
  strip <index> gate on           Get or set the gate on/off state of the strip.
  strip <index> gate mode         Get or set the gate mode of the strip.
  strip <index> gate threshold    Get or set the gate threshold of the strip.
  strip <index> gate range        Get or set the gate range of the strip.
  strip <index> gate attack       Get or set the gate attack time of the strip.
  strip <index> gate hold         Get or set the gate hold time of the strip.
  strip <index> gate release      Get or set the gate release time of the strip.
  strip <index> eq on             Get or set the EQ on/off state of the strip.
  strip <index> eq <band> gain    Get or set the gain of the EQ band.
  strip <index> eq <band> freq    Get or set the frequency of the EQ band.
  strip <index> eq <band> q       Get or set the Q factor of the EQ band.
  strip <index> eq <band> type    Get or set the type of the EQ band.
  strip <index> comp on           Get or set the compressor on/off state of the
                                  strip.
  strip <index> comp mode         Get or set the compressor mode of the strip.
  strip <index> comp threshold    Get or set the compressor threshold of the
                                  strip.
  strip <index> comp ratio        Get or set the compressor ratio of the strip.
  strip <index> comp mix          Get or set the compressor mix of the strip.
  strip <index> comp makeup       Get or set the compressor makeup gain of the
                                  strip.
  strip <index> comp attack       Get or set the compressor attack time of the
                                  strip.
  strip <index> comp hold         Get or set the compressor hold time of the
                                  strip.
  strip <index> comp release      Get or set the compressor release time of the
                                  strip.

Bus
  bus <index> mute              Get or set the mute state of the bus.
  bus <index> fader             Get or set the fader level of the bus.
  bus <index> fadein            Fade in the bus over a specified duration.
  bus <index> fadeout           Fade out the bus over a specified duration.
  bus <index> name              Get or set the name of the bus.
  bus <index> eq on             Get or set the EQ on/off state of the bus.
  bus <index> eq mode           Get or set the EQ mode of the bus (peq, geq or
                                teq).
  bus <index> eq <band> gain    Get or set the gain of the EQ band.
  bus <index> eq <band> freq    Get or set the frequency of the EQ band.
  bus <index> eq <band> q       Get or set the Q factor of the EQ band.
  bus <index> eq <band> type    Get or set the type of the EQ band (lcut, lshv,
                                peq, veq, hshv, hcut).
  bus <index> comp on           Get or set the compressor on/off state of the
                                bus.
  bus <index> comp mode         Get or set the compressor mode of the bus (comp,
                                exp).
  bus <index> comp threshold    Get or set the compressor threshold of the bus
                                (in dB).
  bus <index> comp ratio        Get or set the compressor ratio of the bus.
  bus <index> comp mix          Get or set the compressor mix level of the bus
                                (in %).
  bus <index> comp makeup       Get or set the compressor makeup gain of the bus
                                (in dB).
  bus <index> comp attack       Get or set the compressor attack time of the bus
                                (in ms).
  bus <index> comp hold         Get or set the compressor hold time of the bus
                                (in ms).
  bus <index> comp release      Get or set the compressor release time of the
                                bus (in ms).

Headamp
  headamp <index> gain       Get or set the gain of the headamp.
  headamp <index> phantom    Get or set the phantom power state of the headamp.

Snapshot
  snapshot list              List all snapshots.
  snapshot <index> name      Get or set the name of a snapshot.
  snapshot <index> save      Save the current mixer state to a snapshot.
  snapshot <index> load      Load a mixer state from a snapshot.
  snapshot <index> delete    Delete a snapshot.

Run "xair-cli <command> --help" for more information on a command.
```

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


### License

`xair-cli` is distributed under the terms of the [MIT](https://spdx.org/licenses/MIT.html) license.


### Notes

For an alternative, python implementation consider checking out [dc-xair-cli](https://pypi.org/project/dc-xair-cli/). It supports some operations like batch commands and network discovery which this CLI doesn't (and I have no plans to implement them).
