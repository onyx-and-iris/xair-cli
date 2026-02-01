package cmd

import (
	"github.com/spf13/cobra"
)

// headampCmd represents the headamp command
var headampCmd = &cobra.Command{
	Short: "Commands to control headamp gain and phantom power",
	Long: `Commands to control the headamp gain and phantom power settings of the XAir mixer.

You can get or set the gain level for individual headamps, as well as enable or disable phantom power.`,
	Use: "headamp",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// headampGainCmd represents the headamp gain command
var headampGainCmd = &cobra.Command{
	Use:   "gain",
	Short: "Get or set headamp gain level",
	Long: `Get or set the gain level for a specified headamp index.

Examples:
  # Get gain level for headamp index 1
  xairctl headamp gain 1
  # Set gain level for headamp index 1 to 3.5 dB
  xairctl headamp gain 1 3.5`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide a headamp index")
			return
		}

		index := mustConvToInt(args[0])

		if len(args) == 1 {
			gain, err := client.HeadAmp.Gain(index)
			if err != nil {
				cmd.PrintErrln("Error getting headamp gain level:", err)
				return
			}
			cmd.Printf("Headamp %d Gain: %.2f dB\n", index, gain)
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide a gain level in dB")
			return
		}

		level := mustConvToFloat64(args[1])

		err := client.HeadAmp.SetGain(index, level)
		if err != nil {
			cmd.PrintErrln("Error setting headamp gain level:", err)
			return
		}
		cmd.Printf("Headamp %d Gain set to %.2f dB\n", index, level)
	},
}

// headampPhantomPowerCmd represents the headamp phantom power command
var headampPhantomPowerCmd = &cobra.Command{
	Use:   "phantom",
	Short: "Get or set headamp phantom power status",
	Long: `Get or set the phantom power status for a specified headamp index.
Examples:
  # Get phantom power status for headamp index 1
  xairctl headamp phantom 1
  # Enable phantom power for headamp index 1
  xairctl headamp phantom 1 on
  # Disable phantom power for headamp index 1
  xairctl headamp phantom 1 off`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			cmd.PrintErrln("OSC client not found in context")
			return
		}

		if len(args) < 1 {
			cmd.PrintErrln("Please provide a headamp index")
			return
		}

		index := mustConvToInt(args[0])

		if len(args) == 1 {
			enabled, err := client.HeadAmp.PhantomPower(index)
			if err != nil {
				cmd.PrintErrln("Error getting headamp phantom power status:", err)
				return
			}
			status := "disabled"
			if enabled {
				status = "enabled"
			}
			cmd.Printf("Headamp %d Phantom Power is %s\n", index, status)
			return
		}

		if len(args) < 2 {
			cmd.PrintErrln("Please provide phantom power status: on or off")
			return
		}

		var enable bool
		switch args[1] {
		case "on", "enable":
			enable = true
		case "off", "disable":
			enable = false
		default:
			cmd.PrintErrln("Invalid phantom power status. Use 'on' or 'off'")
			return
		}

		err := client.HeadAmp.SetPhantomPower(index, enable)
		if err != nil {
			cmd.PrintErrln("Error setting headamp phantom power status:", err)
			return
		}
		status := "disabled"
		if enable {
			status = "enabled"
		}
		cmd.Printf("Headamp %d Phantom Power %s successfully\n", index, status)
	},
}

func init() {
	rootCmd.AddCommand(headampCmd)

	headampCmd.AddCommand(headampGainCmd)
	headampCmd.AddCommand(headampPhantomPowerCmd)
}
