package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// rawCmd represents the raw command
var rawCmd = &cobra.Command{
	Short: "Send a raw OSC message to the mixer",
	Long: `Send a raw OSC message to the mixer.
You need to provide the OSC address and any parameters as arguments.

Optionally provide a timeout duration to wait for a response from the mixer. Default is 200ms.`,
	Use: "raw",
	Example: `  xair-cli raw /xinfo
  xair-cli raw /ch/01/mix/fader 0.75
  xair-cli raw --timeout 500ms /bus/02/mix/on 1`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("no client found in context")
		}

		command := args[0]
		params := make([]any, len(args[1:]))
		for i, arg := range args[1:] {
			params[i] = arg
		}
		if err := client.SendMessage(command, params...); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}

		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return fmt.Errorf("error getting timeout flag: %v", err)
		}

		msg, err := client.ReceiveMessage(timeout)
		if err != nil {
			return fmt.Errorf("error receiving message: %v", err)
		}

		if msg != nil {
			fmt.Printf("Received response: %v\n", msg)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)

	rawCmd.Flags().DurationP("timeout", "t", 200*time.Millisecond, "Timeout duration for receiving a response")
}
