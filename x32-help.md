```console
Usage: x32-cli <command> [flags]

A CLI to control Behringer X32 mixers.

Flags:
  -h, --help                  Show context-sensitive help.
  -H, --host="mixer.local"    The host of the X32 device ($X32_CLI_HOST).
  -P, --port=10023            The port of the X32 device ($X32_CLI_PORT).
  -T, --timeout=100ms         Timeout for OSC operations ($X32_CLI_TIMEOUT).
  -L, --loglevel="warn"       Log level for the CLI ($X32_CLI_LOGLEVEL).
  -v, --version               Print x32-cli version information and quit

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

MainMono
  mainmono mute              Get or set the mute state of the Main Mono output.
  mainmono fader             Get or set the fader level of the Main Mono output.
  mainmono fadein            Fade in the Main Mono output over a specified
                             duration.
  mainmono fadeout           Fade out the Main Mono output over a specified
                             duration.
  mainmono eq on             Get or set the EQ on/off state of the Main Mono
                             output.
  mainmono eq <band> gain    Get or set the gain of the specified EQ band.
  mainmono eq <band> freq    Get or set the frequency of the specified EQ band.
  mainmono eq <band> q       Get or set the Q factor of the specified EQ band.
  mainmono eq <band> type    Get or set the type of the specified EQ band.
  mainmono comp on           Get or set the compressor on/off state of the Main
                             Mono output.
  mainmono comp mode         Get or set the compressor mode of the Main Mono
                             output.
  mainmono comp threshold    Get or set the compressor threshold of the Main
                             Mono output.
  mainmono comp ratio        Get or set the compressor ratio of the Main Mono
                             output.
  mainmono comp mix          Get or set the compressor mix level of the Main
                             Mono output.
  mainmono comp makeup       Get or set the compressor makeup gain of the Main
                             Mono output.
  mainmono comp attack       Get or set the compressor attack time of the Main
                             Mono output.
  mainmono comp hold         Get or set the compressor hold time of the Main
                             Mono output.
  mainmono comp release      Get or set the compressor release time of the Main
                             Mono output.

Matrix
  matrix <index> mute              Get or set the mute state of the Matrix
                                   output.
  matrix <index> fader             Get or set the fader level of the Matrix
                                   output.
  matrix <index> fadein            Fade in the Matrix output over a specified
                                   duration.
  matrix <index> fadeout           Fade out the Matrix output over a specified
                                   duration.
  matrix <index> eq on             Get or set the EQ on/off state of the Matrix
                                   output.
  matrix <index> eq <band> gain    Get or set the gain of the specified EQ band.
  matrix <index> eq <band> freq    Get or set the frequency of the specified EQ
                                   band.
  matrix <index> eq <band> q       Get or set the Q factor of the specified EQ
                                   band.
  matrix <index> eq <band> type    Get or set the type of the specified EQ band.
  matrix <index> comp on           Get or set the compressor on/off state of the
                                   Matrix output.
  matrix <index> comp mode         Get or set the compressor mode of the Matrix
                                   output.
  matrix <index> comp threshold    Get or set the compressor threshold of the
                                   Matrix output.
  matrix <index> comp ratio        Get or set the compressor ratio of the Matrix
                                   output.
  matrix <index> comp mix          Get or set the compressor mix level of the
                                   Matrix output.
  matrix <index> comp makeup       Get or set the compressor makeup gain of the
                                   Matrix output.
  matrix <index> comp attack       Get or set the compressor attack time of the
                                   Matrix output.
  matrix <index> comp hold         Get or set the compressor hold time of the
                                   Matrix output.
  matrix <index> comp release      Get or set the compressor release time of the
                                   Matrix output.

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

DCA
  dca <index> mute    Get or set the mute status of the DCA group.
  dca <index> name    Get or set the name of the DCA group.

Run "x32-cli <command> --help" for more information on a command.
```
