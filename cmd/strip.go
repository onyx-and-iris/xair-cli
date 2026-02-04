package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command.
var stripCmd = &cobra.Command{
	Short: "Commands to control individual strips",
	Long:  `Commands to control individual strips of the XAir mixer, including fader level and mute status.`,
	Use:   "strip",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// stripMuteCmd represents the strip mute command.
var stripMuteCmd = &cobra.Command{
	Short: "Get or set the mute status of a strip",
	Long: `Get or set the mute status of a specific strip.

If no argument is provided, the current mute status is retrieved.
If "true" or "1" is provided as an argument, the strip is muted.
If "false" or "0" is provided, the strip is unmuted.`,
	Use: "mute [strip number] [true|false]",
	Example: `  # Get the current mute status of strip 1
  xair-cli strip mute 1
  
  # Mute strip 1
  xair-cli strip mute 1 true
  # Unmute strip 1
  xair-cli strip mute 1 false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			resp, err := client.Strip.Mute(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip mute status: %w", err)
			}
			cmd.Printf("Strip %d mute: %v\n", stripIndex, resp)
			return nil
		}

		var muted bool
		switch args[1] {
		case "true", "1":
			muted = true
		case "false", "0":
			muted = false
		default:
			return fmt.Errorf("Invalid mute status. Use true/false or 1/0")
		}

		err := client.Strip.SetMute(stripIndex, muted)
		if err != nil {
			return fmt.Errorf("Error setting strip mute status: %w", err)
		}

		if muted {
			cmd.Printf("Strip %d muted successfully\n", stripIndex)
		} else {
			cmd.Printf("Strip %d unmuted successfully\n", stripIndex)
		}
		return nil
	},
}

// stripFaderCmd represents the strip fader command.
var stripFaderCmd = &cobra.Command{
	Short: "Get or set the fader level of a strip",
	Long: `Get or set the fader level of a specific strip.

If no level argument is provided, the current fader level is retrieved.
If a level argument (in dB) is provided, the strip fader is set to that level.`,
	Use: "fader [strip number] [level in dB]",
	Example: `  # Get the current fader level of strip 1
  xair-cli strip fader 1
  
  # Set the fader level of strip 1 to -10.0 dB
  xair-cli strip fader 1 -10.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			level, err := client.Strip.Fader(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip fader level: %w", err)
			}
			cmd.Printf("Strip %d fader level: %.2f\n", stripIndex, level)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a fader level in dB")
		}

		level := mustConvToFloat64(args[1])

		err := client.Strip.SetFader(stripIndex, level)
		if err != nil {
			return fmt.Errorf("Error setting strip fader level: %w", err)
		}

		cmd.Printf("Strip %d fader set to %.2f dB\n", stripIndex, level)
		return nil
	},
}

// stripFadeOutCmd represents the strip fade out command.
var stripFadeOutCmd = &cobra.Command{
	Short: "Fade out the strip over a specified duration",
	Long:  "Fade out the strip over a specified duration in seconds.",
	Use:   "fadeout [strip number] --duration [seconds] [target level in dB]",
	Example: `  # Fade out strip 1 over 5 seconds
  xair-cli strip fadeout 1 --duration 5s -- -90.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide strip number")
		}

		stripIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetDuration("duration")
		if err != nil {
			return fmt.Errorf("Error getting duration flag: %w", err)
		}

		target := -90.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Strip.Fader(stripIndex)
		if err != nil {
			return fmt.Errorf("Error getting current strip fader level: %w", err)
		}

		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Strip is already at or below target level")
			return nil
		}

		stepDelay := time.Duration(duration.Seconds()*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err := client.Strip.SetFader(stripIndex, currentFader)
			if err != nil {
				return fmt.Errorf("Error setting strip fader level: %w", err)
			}
			time.Sleep(stepDelay)
		}

		cmd.Printf("Strip %d faded out to %.2f dB over %.2f seconds\n", stripIndex, target, duration.Seconds())
		return nil
	},
}

// stripFadeInCmd represents the strip fade in command.
var stripFadeInCmd = &cobra.Command{
	Short: "Fade in the strip over a specified duration",
	Long:  "Fade in the strip over a specified duration in seconds.",
	Use:   "fadein [strip number] --duration [seconds] [target level in dB]",
	Example: `  # Fade in strip 1 over 5 seconds
  xair-cli strip fadein 1 --duration 5s 0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide strip number")
		}

		stripIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetDuration("duration")
		if err != nil {
			return fmt.Errorf("Error getting duration flag: %w", err)
		}

		target := 0.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Strip.Fader(stripIndex)
		if err != nil {
			return fmt.Errorf("Error getting current strip fader level: %w", err)
		}

		totalSteps := float64(target - currentFader)
		if totalSteps <= 0 {
			cmd.Println("Strip is already at or above target level")
			return nil
		}

		stepDelay := time.Duration(duration.Seconds()*1000/totalSteps) * time.Millisecond

		for currentFader < target {
			currentFader += 1.0
			err := client.Strip.SetFader(stripIndex, currentFader)
			if err != nil {
				return fmt.Errorf("Error setting strip fader level: %w", err)
			}
			time.Sleep(stepDelay)
		}

		cmd.Printf("Strip %d faded in to %.2f dB over %.2f seconds\n", stripIndex, target, duration.Seconds())
		return nil
	},
}

// stripSendCmd represents the strip send command.
var stripSendCmd = &cobra.Command{
	Short: "Get or set the send levels for individual strips",
	Long:  "Get or set the send level from a specific strip to a specific bus.",
	Use:   "send [strip number] [bus number] [level in dB]",
	Example: `  # Get the send level of strip 1 to bus 1
  xair-cli strip send 1 1
  
  # Set the send level of strip 1 to bus 1 to -5.0 dB
  xair-cli strip send 1 1 -- -5.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide strip number and bus number")
		}

		stripIndex, busIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			currentLevel, err := client.Strip.SendLevel(stripIndex, busIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip send level: %w", err)
			}
			cmd.Printf("Strip %d send level to bus %d: %.2f dB\n", stripIndex, busIndex, currentLevel)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide a send level in dB")
		}

		level := mustConvToFloat64(args[2])

		err := client.Strip.SetSendLevel(stripIndex, busIndex, level)
		if err != nil {
			return fmt.Errorf("Error setting strip send level: %w", err)
		}
		cmd.Printf("Strip %d send level to bus %d set to %.2f dB\n", stripIndex, busIndex, level)
		return nil
	},
}

// stripNameCmd represents the strip name command.
var stripNameCmd = &cobra.Command{
	Short: "Get or set the name of a strip",
	Long: `Get or set the name of a specific strip.

If no name argument is provided, the current strip name is retrieved.
If a name argument is provided, the strip name is set to that value.`,
	Use: "name [strip number] [name]",
	Example: `  # Get the current name of strip 1
  xair-cli strip name 1
  
  # Set the name of strip 1 to "Guitar"
  xair-cli strip name 1 "Guitar"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			name, err := client.Strip.Name(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip name: %w", err)
			}
			cmd.Printf("Strip %d name: %s\n", stripIndex, name)
			return nil
		}

		name := args[1]

		err := client.Strip.SetName(stripIndex, name)
		if err != nil {
			return fmt.Errorf("Error setting strip name: %w", err)
		}
		cmd.Printf("Strip %d name set to: %s\n", stripIndex, name)
		return nil
	},
}

// stripGateCmd represents the strip Gate command.
var stripGateCmd = &cobra.Command{
	Short: "Commands to control the Gate of individual strips.",
	Long:  `Commands to control the Gate of individual strips, including turning the Gate on or off.`,
	Use:   "gate",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// stripGateOnCmd represents the strip Gate on command.
var stripGateOnCmd = &cobra.Command{
	Short: "Get or set the Gate on/off status of a strip",
	Long: `Get or set the Gate on/off status of a specific strip.

If no status argument is provided, the current Gate status is retrieved.
If "true" or "1" is provided as an argument, the Gate is turned on.
If "false" or "0" is provided, the Gate is turned off.`,
	Use: "on [strip number] [true|false]",
	Example: `  # Get the current Gate status of strip 1
  xair-cli strip gate on 1
  
  # Turn on Gate for strip 1
  xair-cli strip gate on 1 true
  # Turn off Gate for strip 1
  xair-cli strip gate on 1 false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			on, err := client.Strip.Gate.On(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate on status: %w", err)
			}
			cmd.Printf("Strip %d Gate on: %v\n", stripIndex, on)
			return nil
		}

		var on bool
		switch args[1] {
		case "true", "1":
			on = true
		case "false", "0":
			on = false
		default:
			return fmt.Errorf("Invalid Gate status. Use true/false or 1/0")
		}

		err := client.Strip.Gate.SetOn(stripIndex, on)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate on status: %w", err)
		}

		if on {
			cmd.Printf("Strip %d Gate turned on successfully\n", stripIndex)
		} else {
			cmd.Printf("Strip %d Gate turned off successfully\n", stripIndex)
		}
		return nil
	},
}

// stripGateModeCmd represents the strip Gate Mode command.
var stripGateModeCmd = &cobra.Command{
	Short: "Get or set the Gate mode for a strip",
	Long:  "Get or set the Gate mode for a specific strip.",
	Use:   "mode [strip number] [mode]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}
		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentMode, err := client.Strip.Gate.Mode(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate mode: %w", err)
			}
			cmd.Printf("Strip %d Gate mode: %s\n", stripIndex, currentMode)
			return nil
		}
		if len(args) < 2 {
			return fmt.Errorf("Please provide a mode")
		}

		mode := args[1]
		possibleModes := []string{"exp2", "exp3", "exp4", "gate", "duck"}
		if !contains(possibleModes, mode) {
			return fmt.Errorf("Invalid mode value. Valid values are: %v", possibleModes)
		}

		err := client.Strip.Gate.SetMode(stripIndex, mode)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate mode: %w", err)
		}

		cmd.Printf("Strip %d Gate mode set to %s\n", stripIndex, mode)
		return nil
	},
}

// stripGateThresholdCmd represents the strip Gate Threshold command.
var stripGateThresholdCmd = &cobra.Command{
	Short: "Get or set the Gate threshold for a strip",
	Long:  "Get or set the Gate threshold for a specific strip.",
	Use:   "threshold [strip number] [threshold in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentThreshold, err := client.Strip.Gate.Threshold(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate threshold: %w", err)
			}
			cmd.Printf("Strip %d Gate threshold: %.2f dB\n", stripIndex, currentThreshold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a threshold in dB")
		}

		threshold := mustConvToFloat64(args[1])
		err := client.Strip.Gate.SetThreshold(stripIndex, threshold)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate threshold: %w", err)
		}

		cmd.Printf("Strip %d Gate threshold set to %.2f dB\n", stripIndex, threshold)
		return nil
	},
}

// stripGateRangeCmd represents the strip Gate Range command.
var stripGateRangeCmd = &cobra.Command{
	Short: "Get or set the Gate range for a strip",
	Long:  "Get or set the Gate range for a specific strip.",
	Use:   "range [strip number] [range in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentRange, err := client.Strip.Gate.Range(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate range: %w", err)
			}
			cmd.Printf("Strip %d Gate range: %.2f dB\n", stripIndex, currentRange)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a range in dB")
		}

		rangeDb := mustConvToFloat64(args[1])
		err := client.Strip.Gate.SetRange(stripIndex, rangeDb)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate range: %w", err)
		}

		cmd.Printf("Strip %d Gate range set to %.2f dB\n", stripIndex, rangeDb)
		return nil
	},
}

// stripGateAttackCmd represents the strip Gate Attack command.
var stripGateAttackCmd = &cobra.Command{
	Short: "Get or set the Gate attack time for a strip",
	Long:  "Get or set the Gate attack time for a specific strip.",
	Use:   "attack [strip number] [attack time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentAttack, err := client.Strip.Gate.Attack(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate attack time: %w", err)
			}
			cmd.Printf("Strip %d Gate attack time: %.2f ms\n", stripIndex, currentAttack)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide an attack time in ms")
		}

		attack := mustConvToFloat64(args[1])
		err := client.Strip.Gate.SetAttack(stripIndex, attack)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate attack time: %w", err)
		}

		cmd.Printf("Strip %d Gate attack time set to %.2f ms\n", stripIndex, attack)
		return nil
	},
}

// stripGateHoldCmd represents the strip Gate Hold command.
var stripGateHoldCmd = &cobra.Command{
	Short: "Get or set the Gate hold time for a strip",
	Long:  "Get or set the Gate hold time for a specific strip.",
	Use:   "hold [strip number] [hold time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentHold, err := client.Strip.Gate.Hold(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate hold time: %w", err)
			}
			cmd.Printf("Strip %d Gate hold time: %.2f ms\n", stripIndex, currentHold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a hold time in ms")
		}

		hold := mustConvToFloat64(args[1])
		err := client.Strip.Gate.SetHold(stripIndex, hold)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate hold time: %w", err)
		}

		cmd.Printf("Strip %d Gate hold time set to %.2f ms\n", stripIndex, hold)
		return nil
	},
}

// stripGateReleaseCmd represents the strip Gate Release command.
var stripGateReleaseCmd = &cobra.Command{
	Short: "Get or set the Gate release time for a strip",
	Long:  "Get or set the Gate release time for a specific strip.",
	Use:   "release [strip number] [release time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])
		if len(args) == 1 {
			currentRelease, err := client.Strip.Gate.Release(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Gate release time: %w", err)
			}
			cmd.Printf("Strip %d Gate release time: %.2f ms\n", stripIndex, currentRelease)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a release time in ms")
		}

		release := mustConvToFloat64(args[1])
		err := client.Strip.Gate.SetRelease(stripIndex, release)
		if err != nil {
			return fmt.Errorf("Error setting strip Gate release time: %w", err)
		}

		cmd.Printf("Strip %d Gate release time set to %.2f ms\n", stripIndex, release)
		return nil
	},
}

// stripEqCmd represents the strip EQ command.
var stripEqCmd = &cobra.Command{
	Short: "Commands to control the EQ of individual strips.",
	Long:  `Commands to control the EQ of individual strips, including turning the EQ on or off.`,
	Use:   "eq",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// stripEqOnCmd represents the strip EQ on command.
var stripEqOnCmd = &cobra.Command{
	Short: "Get or set the EQ on/off status of a strip",
	Long: `Get or set the EQ on/off status of a specific strip.

If no status argument is provided, the current EQ status is retrieved.
If "true" or "1" is provided as an argument, the EQ is turned on.
If "false" or "0" is provided, the EQ is turned off.`,
	Use: "on [strip number] [true|false]",
	Example: `  # Get the current EQ status of strip 1
  xair-cli strip eq on 1
  
  # Turn on EQ for strip 1
  xair-cli strip eq on 1 true
  # Turn off EQ for strip 1
  xair-cli strip eq on 1 false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			on, err := client.Strip.Eq.On(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip EQ on status: %w", err)
			}
			cmd.Printf("Strip %d EQ on: %v\n", stripIndex, on)
			return nil
		}

		var on bool
		switch args[1] {
		case "true", "1":
			on = true
		case "false", "0":
			on = false
		default:
			return fmt.Errorf("Invalid EQ status. Use true/false or 1/0")
		}

		err := client.Strip.Eq.SetOn(stripIndex, on)
		if err != nil {
			return fmt.Errorf("Error setting strip EQ on status: %w", err)
		}

		if on {
			cmd.Printf("Strip %d EQ turned on successfully\n", stripIndex)
		} else {
			cmd.Printf("Strip %d EQ turned off successfully\n", stripIndex)
		}
		return nil
	},
}

// stripEqGainCmd represents the strip EQ Gain command.
var stripEqGainCmd = &cobra.Command{
	Short: "Get or set the EQ band gain for a strip",
	Long:  "Get or set the EQ band gain for a specific strip and band.",
	Use:   "gain [strip number] [band number] [gain in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide strip number and band number")
		}

		stripIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			currentGain, err := client.Strip.Eq.Gain(stripIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip EQ band gain: %w", err)
			}
			cmd.Printf("Strip %d EQ band %d gain: %.2f dB\n", stripIndex, bandIndex, currentGain)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide a gain in dB")
		}

		gain := mustConvToFloat64(args[2])

		err := client.Strip.Eq.SetGain(stripIndex, bandIndex, gain)
		if err != nil {
			return fmt.Errorf("Error setting strip EQ band gain: %w", err)
		}

		cmd.Printf("Strip %d EQ band %d gain set to %.2f dB\n", stripIndex, bandIndex, gain)
		return nil
	},
}

// stripEqFreqCmd represents the strip EQ Frequency command.
var stripEqFreqCmd = &cobra.Command{
	Short: "Get or set the EQ band frequency for a strip",
	Long:  "Get or set the EQ band frequency for a specific strip and band.",
	Use:   "freq [strip number] [band number] [frequency in Hz]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide strip number and band number")
		}

		stripIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			currentFreq, err := client.Strip.Eq.Frequency(stripIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip EQ band frequency: %w", err)
			}
			cmd.Printf("Strip %d EQ band %d frequency: %.2f Hz\n", stripIndex, bandIndex, currentFreq)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide a frequency in Hz")
		}

		freq := mustConvToFloat64(args[2])

		err := client.Strip.Eq.SetFrequency(stripIndex, bandIndex, freq)
		if err != nil {
			return fmt.Errorf("Error setting strip EQ band frequency: %w", err)
		}

		cmd.Printf("Strip %d EQ band %d frequency set to %.2f Hz\n", stripIndex, bandIndex, freq)
		return nil
	},
}

// stripEqQCmd represents the strip EQ Q command.
var stripEqQCmd = &cobra.Command{
	Short: "Get or set the EQ band Q factor for a strip",
	Long:  "Get or set the EQ band Q factor for a specific strip and band.",
	Use:   "q [strip number] [band number] [Q factor]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide strip number and band number")
		}

		stripIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			currentQ, err := client.Strip.Eq.Q(stripIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip EQ band Q factor: %w", err)
			}
			cmd.Printf("Strip %d EQ band %d Q factor: %.2f\n", stripIndex, bandIndex, currentQ)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide a Q factor")
		}

		q := mustConvToFloat64(args[2])

		err := client.Strip.Eq.SetQ(stripIndex, bandIndex, q)
		if err != nil {
			return fmt.Errorf("Error setting strip EQ band Q factor: %w", err)
		}

		cmd.Printf("Strip %d EQ band %d Q factor set to %.2f\n", stripIndex, bandIndex, q)
		return nil
	},
}

// stripEqTypeCmd represents the strip EQ Type command.
var stripEqTypeCmd = &cobra.Command{
	Short: "Get or set the EQ band type for a strip",
	Long:  "Get or set the EQ band type for a specific strip and band.",
	Use:   "type [strip number] [band number] [type]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide strip number and band number")
		}

		stripIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		eqTypeNames := []string{"lcut", "lshv", "peq", "veq", "hshv", "hcut"}

		if len(args) == 2 {
			currentType, err := client.Strip.Eq.Type(stripIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip EQ band type: %w", err)
			}
			cmd.Printf("Strip %d EQ band %d type: %s\n", stripIndex, bandIndex, eqTypeNames[currentType])
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide a type")
		}

		eqType := indexOf(eqTypeNames, args[2])
		if eqType == -1 {
			return fmt.Errorf("Invalid EQ band type. Valid types are: %v", eqTypeNames)
		}

		err := client.Strip.Eq.SetType(stripIndex, bandIndex, eqType)
		if err != nil {
			return fmt.Errorf("Error setting strip EQ band type: %w", err)
		}

		cmd.Printf("Strip %d EQ band %d type set to %s\n", stripIndex, bandIndex, eqTypeNames[eqType])
		return nil
	},
}

// stripCompCmd represents the strip Compressor command.
var stripCompCmd = &cobra.Command{
	Short: "Commands to control the Compressor of individual strips.",
	Long:  `Commands to control the Compressor of individual strips, including turning the Compressor on or off.`,
	Use:   "comp",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// stripCompOnCmd represents the strip Compressor on command.
var stripCompOnCmd = &cobra.Command{
	Short: "Get or set the Compressor on/off status of a strip",
	Long: `Get or set the Compressor on/off status of a specific strip.

If no status argument is provided, the current Compressor status is retrieved.
If "true" or "1" is provided as an argument, the Compressor is turned on.
If "false" or "0" is provided, the Compressor is turned off.`,
	Use: "on [strip number] [true|false]",
	Example: `  # Get the current Compressor status of strip 1
  xair-cli strip comp on 1
  
  # Turn on Compressor for strip 1
  xair-cli strip comp on 1 true
  # Turn off Compressor for strip 1
  xair-cli strip comp on 1 false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			on, err := client.Strip.Comp.On(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor on status: %w", err)
			}
			cmd.Printf("Strip %d Compressor on: %v\n", stripIndex, on)
			return nil
		}

		var on bool
		switch args[1] {
		case "true", "1":
			on = true
		case "false", "0":
			on = false
		default:
			return fmt.Errorf("Invalid Compressor status. Use true/false or 1/0")
		}

		err := client.Strip.Comp.SetOn(stripIndex, on)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor on status: %w", err)
		}

		if on {
			cmd.Printf("Strip %d Compressor turned on successfully\n", stripIndex)
		} else {
			cmd.Printf("Strip %d Compressor turned off successfully\n", stripIndex)
		}
		return nil
	},
}

// stripCompModeCmd represents the strip Compressor Mode command.
var stripCompModeCmd = &cobra.Command{
	Short: "Get or set the Compressor mode for a strip",
	Long:  "Get or set the Compressor mode for a specific strip.",
	Use:   "mode [strip number] [mode]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentMode, err := client.Strip.Comp.Mode(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor mode: %w", err)
			}
			cmd.Printf("Strip %d Compressor mode: %s\n", stripIndex, currentMode)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a mode")
		}

		mode := args[1]
		if !contains([]string{"comp", "exp"}, mode) {
			return fmt.Errorf("Invalid mode value. Valid values are: comp, exp")
		}

		err := client.Strip.Comp.SetMode(stripIndex, mode)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor mode: %w", err)
		}

		cmd.Printf("Strip %d Compressor mode set to %s\n", stripIndex, mode)
		return nil
	},
}

// stripCompThresholdCmd represents the strip Compressor Threshold command.
var stripCompThresholdCmd = &cobra.Command{
	Short: "Get or set the Compressor threshold for a strip",
	Long:  "Get or set the Compressor threshold for a specific strip.",
	Use:   "threshold [strip number] [threshold in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentThreshold, err := client.Strip.Comp.Threshold(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor threshold: %w", err)
			}
			cmd.Printf("Strip %d Compressor threshold: %.2f dB\n", stripIndex, currentThreshold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a threshold in dB")
		}

		threshold := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetThreshold(stripIndex, threshold)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor threshold: %w", err)
		}

		cmd.Printf("Strip %d Compressor threshold set to %.2f dB\n", stripIndex, threshold)
		return nil
	},
}

// stripCompRatioCmd represents the strip Compressor Ratio command.
var stripCompRatioCmd = &cobra.Command{
	Short: "Get or set the Compressor ratio for a strip",
	Long:  "Get or set the Compressor ratio for a specific strip.",
	Use:   "ratio [strip number] [ratio]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentRatio, err := client.Strip.Comp.Ratio(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor ratio: %w", err)
			}
			cmd.Printf("Strip %d Compressor ratio: %.2f\n", stripIndex, currentRatio)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a ratio")
		}

		ratio := mustConvToFloat64(args[1])
		possibleValues := []float64{1.1, 1.3, 1.5, 2.0, 2.5, 3.0, 4.0, 5.0, 7.0, 10, 20, 100}
		if !contains(possibleValues, ratio) {
			return fmt.Errorf("Invalid ratio value. Valid values are: %v", possibleValues)
		}

		err := client.Strip.Comp.SetRatio(stripIndex, ratio)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor ratio: %w", err)
		}

		cmd.Printf("Strip %d Compressor ratio set to %.2f\n", stripIndex, ratio)
		return nil
	},
}

// stripCompMixCmd represents the strip Compressor Mix command.
var stripCompMixCmd = &cobra.Command{
	Short: "Get or set the Compressor mix for a strip",
	Long:  "Get or set the Compressor mix for a specific strip.",
	Use:   "mix [strip number] [mix percentage]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentMix, err := client.Strip.Comp.Mix(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor mix: %w", err)
			}
			cmd.Printf("Strip %d Compressor mix: %.2f%%\n", stripIndex, currentMix)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a mix percentage")
		}

		mix := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetMix(stripIndex, mix)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor mix: %w", err)
		}

		cmd.Printf("Strip %d Compressor mix set to %.2f%%\n", stripIndex, mix)
		return nil
	},
}

// stripCompMakeUpCmd represents the strip Compressor Make-Up Gain command.
var stripCompMakeUpCmd = &cobra.Command{
	Short: "Get or set the Compressor make-up gain for a strip",
	Long:  "Get or set the Compressor make-up gain for a specific strip.",
	Use:   "makeup [strip number] [make-up gain in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentMakeUp, err := client.Strip.Comp.MakeUp(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor make-up gain: %w", err)
			}
			cmd.Printf("Strip %d Compressor make-up gain: %.2f dB\n", stripIndex, currentMakeUp)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a make-up gain in dB")
		}

		makeUp := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetMakeUp(stripIndex, makeUp)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor make-up gain: %w", err)
		}

		cmd.Printf("Strip %d Compressor make-up gain set to %.2f dB\n", stripIndex, makeUp)
		return nil
	},
}

// stripCompAttackCmd represents the strip Compressor Attack command.
var stripCompAttackCmd = &cobra.Command{
	Short: "Get or set the Compressor attack time for a strip",
	Long:  "Get or set the Compressor attack time for a specific strip.",
	Use:   "attack [strip number] [attack time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentAttack, err := client.Strip.Comp.Attack(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor attack time: %w", err)
			}
			cmd.Printf("Strip %d Compressor attack time: %.2f ms\n", stripIndex, currentAttack)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide an attack time in ms")
		}

		attack := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetAttack(stripIndex, attack)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor attack time: %w", err)
		}

		cmd.Printf("Strip %d Compressor attack time set to %.2f ms\n", stripIndex, attack)
		return nil
	},
}

// stripCompHoldCmd represents the strip Compressor Hold command.
var stripCompHoldCmd = &cobra.Command{
	Short: "Get or set the Compressor hold time for a strip",
	Long:  "Get or set the Compressor hold time for a specific strip.",
	Use:   "hold [strip number] [hold time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentHold, err := client.Strip.Comp.Hold(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor hold time: %w", err)
			}
			cmd.Printf("Strip %d Compressor hold time: %.2f ms\n", stripIndex, currentHold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a hold time in ms")
		}

		hold := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetHold(stripIndex, hold)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor hold time: %w", err)
		}

		cmd.Printf("Strip %d Compressor hold time set to %.2f ms\n", stripIndex, hold)
		return nil
	},
}

// stripCompReleaseCmd represents the strip Compressor Release command.
var stripCompReleaseCmd = &cobra.Command{
	Short: "Get or set the Compressor release time for a strip",
	Long:  "Get or set the Compressor release time for a specific strip.",
	Use:   "release [strip number] [release time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide a strip number")
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			currentRelease, err := client.Strip.Comp.Release(stripIndex)
			if err != nil {
				return fmt.Errorf("Error getting strip Compressor release time: %w", err)
			}
			cmd.Printf("Strip %d Compressor release time: %.2f ms\n", stripIndex, currentRelease)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide a release time in ms")
		}

		release := mustConvToFloat64(args[1])

		err := client.Strip.Comp.SetRelease(stripIndex, release)
		if err != nil {
			return fmt.Errorf("Error setting strip Compressor release time: %w", err)
		}

		cmd.Printf("Strip %d Compressor release time set to %.2f ms\n", stripIndex, release)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)

	stripCmd.AddCommand(stripMuteCmd)
	stripCmd.AddCommand(stripFaderCmd)
	stripCmd.AddCommand(stripFadeOutCmd)
	stripFadeOutCmd.Flags().DurationP("duration", "d", 5*time.Second, "Duration of the fade out in seconds")
	stripCmd.AddCommand(stripFadeInCmd)
	stripFadeInCmd.Flags().DurationP("duration", "d", 5*time.Second, "Duration of the fade in in seconds")
	stripCmd.AddCommand(stripSendCmd)
	stripCmd.AddCommand(stripNameCmd)

	stripCmd.AddCommand(stripGateCmd)
	stripGateCmd.AddCommand(stripGateOnCmd)
	stripGateCmd.AddCommand(stripGateModeCmd)
	stripGateCmd.AddCommand(stripGateThresholdCmd)
	stripGateCmd.AddCommand(stripGateRangeCmd)
	stripGateCmd.AddCommand(stripGateAttackCmd)
	stripGateCmd.AddCommand(stripGateHoldCmd)
	stripGateCmd.AddCommand(stripGateReleaseCmd)

	stripCmd.AddCommand(stripEqCmd)
	stripEqCmd.AddCommand(stripEqOnCmd)
	stripEqCmd.AddCommand(stripEqGainCmd)
	stripEqCmd.AddCommand(stripEqFreqCmd)
	stripEqCmd.AddCommand(stripEqQCmd)
	stripEqCmd.AddCommand(stripEqTypeCmd)

	stripCmd.AddCommand(stripCompCmd)
	stripCompCmd.AddCommand(stripCompOnCmd)
	stripCompCmd.AddCommand(stripCompModeCmd)
	stripCompCmd.AddCommand(stripCompThresholdCmd)
	stripCompCmd.AddCommand(stripCompRatioCmd)
	stripCompCmd.AddCommand(stripCompMixCmd)
	stripCompCmd.AddCommand(stripCompMakeUpCmd)
	stripCompCmd.AddCommand(stripCompAttackCmd)
	stripCompCmd.AddCommand(stripCompHoldCmd)
	stripCompCmd.AddCommand(stripCompReleaseCmd)
}
