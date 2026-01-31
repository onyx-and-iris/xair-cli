/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"fmt"

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

func init() {
	rootCmd.AddCommand(busCmd)

	busCmd.AddCommand(busMuteCmd)
}
