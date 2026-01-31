package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// busCmd represents the bus command.
var busCmd = &cobra.Command{
	Short: "Commands to control individual buses",
	Long:  `Commands to control individual buses of the XAir mixer, including mute status.`,
	Use:   "bus",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// busMuteCmd represents the bus mute command.
var busMuteCmd = &cobra.Command{
	Short: "Get or set the bus mute status",
	Long:  `Get or set the mute status of a specific bus.`,
	Use:   "mute [bus number] [true|false]",
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide bus number and mute status (true/false)")
			return
		}

		busNum := mustConvToInt(args[0])
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

		err := client.Bus.SetMute(busNum, muted)
		if err != nil {
			cmd.PrintErrln("Error setting bus mute status:", err)
			return
		}
		cmd.Printf("Bus %d mute set to %v\n", busNum, muted)
	},
}

// busFaderCmd represents the bus fader command.
var busFaderCmd = &cobra.Command{
	Short: "Get or set the bus fader level",
	Long: `Get or set the fader level of a specific bus.
If no level argument is provided, the current fader level is retrieved.
If a level argument (in dB) is provided, the bus fader is set to that level.`,
	Use: "fader [bus number] [level in dB]",
	Example: `	# Get the current fader level of bus 1
	xair-cli bus fader 1
	
	# Set the fader level of bus 1 to -10.0 dB
	xair-cli bus fader 1 -10.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			level, err := client.Bus.Fader(busIndex)
			if err != nil {
				cmd.PrintErrln("Error getting bus fader level:", err)
				return
			}
			cmd.Printf("Bus %d fader level: %.1f dB\n", busIndex, level)
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide bus number and fader level (in dB)")
			return
		}

		level := mustConvToFloat64(args[1])

		err := client.Bus.SetFader(busIndex, level)
		if err != nil {
			cmd.PrintErrln("Error setting bus fader level:", err)
			return
		}
		cmd.Printf("Bus %d fader set to %.2f dB\n", busIndex, level)
	},
}

// busFadeOutCmd represents the bus fade out command.
var busFadeOutCmd = &cobra.Command{
	Short: "Fade out the bus fader over a specified duration",
	Long:  "Fade out the bus fader to minimum level over a specified duration in seconds.",
	Use:   "fadeout [bus number] --duration [seconds] [target level in dB]",
	Example: `  # Fade out bus 1 over 5 seconds
  xair-cli bus fadeout 1 --duration 5 -- -90.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide bus number")
			return
		}

		busIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		target := -90.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Bus.Fader(busIndex)
		if err != nil {
			cmd.PrintErrln("Error getting current bus fader level:", err)
			return
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Bus is already at or below target level")
			return
		}

		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err := client.Bus.SetFader(busIndex, currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting bus fader level:", err)
				return
			}
			time.Sleep(stepDelay)
		}

		cmd.Println("Bus fade out completed")
	},
}

// BusFadeInCmd represents the bus fade in command.
var busFadeInCmd = &cobra.Command{
	Short: "Fade in the bus fader over a specified duration",
	Long:  "Fade in the bus fader to maximum level over a specified duration in seconds.",
	Use:   "fadein [bus number] --duration [seconds] [target level in dB]",
	Example: `  # Fade in bus 1 over 5 seconds
  xair-cli bus fadein 1 --duration 5 -- 0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide bus number")
			return
		}

		busIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		target := 0.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Bus.Fader(busIndex)
		if err != nil {
			cmd.PrintErrln("Error getting current bus fader level:", err)
			return
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(target - currentFader)
		if totalSteps <= 0 {
			cmd.Println("Bus is already at or above target level")
			return
		}

		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader < target {
			currentFader += 1.0
			err := client.Bus.SetFader(busIndex, currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting bus fader level:", err)
				return
			}
			time.Sleep(stepDelay)
		}

		cmd.Println("Bus fade in completed")
	},
}

func init() {
	rootCmd.AddCommand(busCmd)

	busCmd.AddCommand(busMuteCmd)

	busCmd.AddCommand(busFaderCmd)
	busCmd.AddCommand(busFadeOutCmd)
	busFadeOutCmd.Flags().Float64P("duration", "d", 5.0, "Duration for fade out in seconds")
	busCmd.AddCommand(busFadeInCmd)
	busFadeInCmd.Flags().Float64P("duration", "d", 5.0, "Duration for fade in in seconds")
}
