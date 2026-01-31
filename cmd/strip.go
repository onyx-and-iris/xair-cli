/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var stripCmd = &cobra.Command{
	Use:   "strip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("strip called")
	},
}

var stripMuteCmd = &cobra.Command{
	Use:   "mute",
	Short: "Get or set the mute status of a strip",
	Long: `Get or set the mute status of a specific strip.

If no argument is provided, the current mute status is retrieved.
If "true" or "1" is provided as an argument, the strip is muted.
If "false" or "0" is provided, the strip is unmuted.

For example:
  # Get the current mute status of strip 1
  xair-cli strip mute 1
  
  # Mute strip 1
  xair-cli strip mute 1 true
  # Unmute strip 1
  xair-cli strip mute 1 false
`,
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
			resp, err := client.StripMute(stripIndex)
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

		err := client.SetStripMute(stripIndex, muted)
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

func init() {
	rootCmd.AddCommand(stripCmd)

	stripCmd.AddCommand(stripMuteCmd)
}
