/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// busCmd represents the bus command
var busCmd = &cobra.Command{
	Use:   "bus",
	Short: "Commands to control individual buses",
	Long:  `Commands to control individual buses of the XAir mixer, including mute status.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bus called")
	},
}

// busMuteCmd represents the bus mute command
var busMuteCmd = &cobra.Command{
	Use:   "mute",
	Short: "Get or set the bus mute status",
	Long:  `Get or set the mute status of a specific bus.`,
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

		err := client.SetBusMute(busNum, muted)
		if err != nil {
			cmd.PrintErrln("Error setting bus mute status:", err)
			return
		}
		cmd.Printf("Bus %d mute set to %v\n", busNum, muted)
	},
}

// busFaderCmd represents the bus fader command
var busFaderCmd = &cobra.Command{
	Use:   "fader",
	Short: "Get or set the bus fader level",
	Long:  `Get or set the fader level of a specific bus.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide bus number and fader level (in dB)")
			return
		}

		busNum := mustConvToInt(args[0])
		level := mustConvToFloat64(args[1])

		err := client.SetBusFader(busNum, level)
		if err != nil {
			cmd.PrintErrln("Error setting bus fader level:", err)
			return
		}
		cmd.Printf("Bus %d fader set to %.2f dB\n", busNum, level)
	},
}

// busFadeOutCmd represents the bus fade out command
var busFadeOutCmd = &cobra.Command{
	Use:   "fadeout",
	Short: "Fade out the bus fader over a specified duration",
	Long: `Fade out the bus fader to minimum level over a specified duration in seconds.

For example:
  # Fade out bus 1 over 5 seconds
  xair-cli bus fadeout 1 --duration 5
`,
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

		currentFader, err := client.BusFader(busIndex)
		if err != nil {
			cmd.PrintErrln("Error getting current bus fader level:", err)
			return
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Bus is already faded out")
			return
		}

		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err := client.SetBusFader(busIndex, currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting bus fader level:", err)
				return
			}
			time.Sleep(stepDelay)
		}

		cmd.Println("Bus fade out completed")
	},
}

// BusFadeInCmd represents the bus fade in command
var busFadeInCmd = &cobra.Command{
	Use:   "fadein",
	Short: "Fade in the bus fader over a specified duration",
	Long: `Fade in the bus fader to maximum level over a specified duration in seconds.

For example:
  # Fade in bus 1 over 5 seconds
  xair-cli bus fadein 1 --duration 5
`,
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

		currentFader, err := client.BusFader(busIndex)
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
			err := client.SetBusFader(busIndex, currentFader)
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
	busFadeOutCmd.Flags().Float64P("duration", "d", 5, "Duration for fade out in seconds")
	busCmd.AddCommand(busFadeInCmd)
	busFadeInCmd.Flags().Float64P("duration", "d", 5, "Duration for fade in in seconds")
}
