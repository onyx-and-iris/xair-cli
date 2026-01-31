package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

// mainCmd represents the main command.
var mainCmd = &cobra.Command{
	Short: "Commands to control the main output",
	Long:  `Commands to control the main output of the XAir mixer, including fader level and mute status.`,
	Use:   "main",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// mainMuteCmd represents the main mute command.
var mainMuteCmd = &cobra.Command{
	Short: "Get or set the main LR mute status",
	Long: `Get or set the main L/R mute status.

If no argument is provided, the current mute status is retrieved.
If "true" or "1" is provided as an argument, the main output is muted.
If "false" or "0" is provided, the main output is unmuted.`,
	Use: "mute [true|false]",
	Example: `  # Get the current main LR mute status
  xair-cli main mute

  # Mute the main output
  xair-cli main mute true

  # Unmute the main output
  xair-cli main mute false`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) == 0 {
			resp, err := client.Main.Mute()
			if err != nil {
				cmd.PrintErrln("Error getting main LR mute status:", err)
				return
			}
			cmd.Printf("Main LR mute: %v\n", resp)
			return
		}

		var muted bool
		switch args[0] {
		case "true", "1":
			muted = true
		case "false", "0":
			muted = false
		default:
			cmd.PrintErrln("Invalid mute status. Use true/false or 1/0")
			return
		}

		err := client.Main.SetMute(muted)
		if err != nil {
			cmd.PrintErrln("Error setting main LR mute status:", err)
			return
		}
		cmd.Println("Main LR mute status set successfully")
	},
}

// mainFaderCmd represents the main fader command.
var mainFaderCmd = &cobra.Command{
	Short: "Set or get the main LR fader level",
	Long: `Set or get the main L/R fader level in dB.

If no argument is provided, the current fader level is retrieved.
If a dB value is provided as an argument, the fader level is set to that value.`,
	Use: "fader [level in dB]",
	Example: `  # Get the current main LR fader level
  xair-cli main fader

  # Set the main LR fader level to -10.0 dB
  xair-cli main fader -- -10.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) == 0 {
			resp, err := client.Main.Fader()
			if err != nil {
				cmd.PrintErrln("Error getting main LR fader:", err)
				return
			}
			cmd.Printf("Main LR fader: %.1f dB\n", resp)
			return
		}

		err := client.Main.SetFader(mustConvToFloat64(args[0]))
		if err != nil {
			cmd.PrintErrln("Error setting main LR fader:", err)
			return
		}
		cmd.Println("Main LR fader set successfully")
	},
}

// mainFadeOutCmd represents the main fadeout command.
var mainFadeOutCmd = &cobra.Command{
	Short: "Fade out the main output",
	Long: `Fade out the main output over a specified duration.

This command will fade out the main output to the specified dB level.
`,
	Use: "fadeout --duration [seconds] [target_db]",
	Example: `  # Fade out main output over 5 seconds
  xair-cli main fadeout --duration 5 -- -90.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		// Default target for fadeout
		target := -90.0
		if len(args) > 0 {
			target = mustConvToFloat64(args[0])
		}

		currentFader, err := client.Main.Fader()
		if err != nil {
			cmd.PrintErrln("Error getting current main LR fader:", err)
			return
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Main output is already faded out")
			return
		}

		// Calculate delay per step to achieve exact duration
		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err = client.Main.SetFader(currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting main LR fader:", err)
				return
			}
			time.Sleep(stepDelay)
		}
		cmd.Println("Main output faded out successfully")
	},
}

// mainFadeInCmd represents the main fadein command.
var mainFadeInCmd = &cobra.Command{
	Short: "Fade in the main output",
	Long: `Fade in the main output over a specified duration.

This command will fade in the main output to the specified dB level.
`,
	Use: "fadein --duration [seconds] [target_db]",
	Example: `  # Fade in main output over 5 seconds
  xair-cli main fadein --duration 5 -- 0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			cmd.PrintErrln("Error getting duration flag:", err)
			return
		}

		target := 0.0
		if len(args) > 0 {
			target = mustConvToFloat64(args[0])
		}

		currentFader, err := client.Main.Fader()
		if err != nil {
			cmd.PrintErrln("Error getting current main LR fader:", err)
			return
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(target - currentFader)
		if totalSteps <= 0 {
			cmd.Println("Main output is already at or above target level")
			return
		}

		// Calculate delay per step to achieve exact duration
		stepDelay := time.Duration(duration*1000/totalSteps) * time.Millisecond

		for currentFader < target {
			currentFader += 1.0
			err = client.Main.SetFader(currentFader)
			if err != nil {
				cmd.PrintErrln("Error setting main LR fader:", err)
				return
			}
			time.Sleep(stepDelay)
		}
		cmd.Println("Main output faded in successfully")
	},
}

func init() {
	rootCmd.AddCommand(mainCmd)

	mainCmd.AddCommand(mainMuteCmd)

	mainCmd.AddCommand(mainFaderCmd)
	mainCmd.AddCommand(mainFadeOutCmd)
	mainFadeOutCmd.Flags().Float64P("duration", "d", 5, "Duration for fade out in seconds")
	mainCmd.AddCommand(mainFadeInCmd)
	mainFadeInCmd.Flags().Float64P("duration", "d", 5, "Duration for fade in in seconds")
}
