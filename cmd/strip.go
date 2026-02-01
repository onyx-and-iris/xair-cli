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

	stripCmd.AddCommand(stripEqCmd)
	stripEqCmd.AddCommand(stripEqOnCmd)

	stripCmd.AddCommand(stripCompCmd)
	stripCompCmd.AddCommand(stripCompOnCmd)
}
