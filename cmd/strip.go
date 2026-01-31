package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide a strip number")
			return
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			resp, err := client.Strip.Mute(stripIndex)
			if err != nil {
				cmd.PrintErrln("Error getting strip mute status:", err)
				return
			}
			cmd.Printf("Strip %d mute: %v\n", stripIndex, resp)
			return
		}

		var muted bool
		switch args[1] {
		case "true", "1":
			muted = true
		case "false", "0":
			muted = false
		default:
			cmd.PrintErrln("Invalid mute status. Use true/false or 1/0")
			return
		}

		err := client.Strip.SetMute(stripIndex, muted)
		if err != nil {
			cmd.PrintErrln("Error setting strip mute status:", err)
			return
		}
		if muted {
			cmd.Printf("Strip %d muted successfully\n", stripIndex)
		} else {
			cmd.Printf("Strip %d unmuted successfully\n", stripIndex)
		}
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
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide a strip number")
			return
		}

		stripIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			level, err := client.Strip.Fader(stripIndex)
			if err != nil {
				cmd.PrintErrln("Error getting strip fader level:", err)
				return
			}
			cmd.Printf("Strip %d fader level: %.2f\n", stripIndex, level)
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide a fader level in dB")
			return
		}

		level := mustConvToFloat64(args[1])

		err := client.Strip.SetFader(stripIndex, level)
		if err != nil {
			cmd.PrintErrln("Error setting strip fader level:", err)
			return
		}
		cmd.Printf("Strip %d fader set to %.2f dB\n", stripIndex, level)
	},
}

// stripFadeOutCmd represents the strip fade out command.
var stripFadeOutCmd = &cobra.Command{
	Short: "Fade out the strip over a specified duration",
	Long:  "Fade out the strip over a specified duration in seconds.",
	Use:   "fadeout [strip number] --duration [seconds] [target level in dB]",
	Example: `  # Fade out strip 1 over 5 seconds
  xair-cli strip fadeout 1 --duration 5.0 -- -90.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide strip number")
			return
		}

		stripIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		target := -90.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Strip.Fader(stripIndex)
		if err != nil {
			cmd.PrintErrln("Error getting current strip fader level:", err)
			return
		}

		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Strip is already at or below target level")
			return
		}

		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err := client.Strip.SetFader(stripIndex, currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting strip fader level:", err)
				return
			}
			time.Sleep(stepDelay)
		}

		cmd.Printf("Strip %d faded out to %.2f dB over %.2f seconds\n", stripIndex, target, duration)
	},
}

// stripFadeInCmd represents the strip fade in command.
var stripFadeInCmd = &cobra.Command{
	Short: "Fade in the strip over a specified duration",
	Long:  "Fade in the strip over a specified duration in seconds.",
	Use:   "fadein [strip number] --duration [seconds] [target level in dB]",
	Example: `  # Fade in strip 1 over 5 seconds
  xair-cli strip fadein 1 --duration 5.0 0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide strip number")
			return
		}

		stripIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		target := 0.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Strip.Fader(stripIndex)
		if err != nil {
			cmd.PrintErrln("Error getting current strip fader level:", err)
			return
		}

		totalSteps := float64(target - currentFader)
		if totalSteps <= 0 {
			cmd.Println("Strip is already at or above target level")
			return
		}

		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader < target {
			currentFader += 1.0
			err := client.Strip.SetFader(stripIndex, currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting strip fader level:", err)
				return
			}
			time.Sleep(stepDelay)
		}

		cmd.Printf("Strip %d faded in to %.2f dB over %.2f seconds\n", stripIndex, target, duration)
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
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide strip number and bus number")
			return
		}

		stripIndex, busIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			currentLevel, err := client.Strip.SendLevel(stripIndex, busIndex)
			if err != nil {
				cmd.PrintErrln("Error getting strip send level:", err)
				return
			}
			cmd.Printf("Strip %d send level to bus %d: %.2f dB\n", stripIndex, busIndex, currentLevel)
			return
		}

		if len(args) < 3 {
			cmd.PrintErrln("Please provide a send level in dB")
			return
		}

		level := mustConvToFloat64(args[2])

		err := client.Strip.SetSendLevel(stripIndex, busIndex, level)
		if err != nil {
			cmd.PrintErrln("Error setting strip send level:", err)
			return
		}
		cmd.Printf("Strip %d send level to bus %d set to %.2f dB\n", stripIndex, busIndex, level)
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)

	stripCmd.AddCommand(stripMuteCmd)

	stripCmd.AddCommand(stripFaderCmd)
	stripCmd.AddCommand(stripFadeOutCmd)
	stripFadeOutCmd.Flags().Float64P("duration", "d", 5.0, "Duration of the fade out in seconds")
	stripCmd.AddCommand(stripFadeInCmd)
	stripFadeInCmd.Flags().Float64P("duration", "d", 5.0, "Duration of the fade in in seconds")

	stripCmd.AddCommand(stripSendCmd)
}
