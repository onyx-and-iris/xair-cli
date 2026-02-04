package cmd

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and mute status (true/false)")
		}

		busNum := mustConvToInt(args[0])
		var muted bool
		switch args[1] {
		case "true", "1":
			muted = true
		case "false", "0":
			muted = false
		default:
			return fmt.Errorf("Invalid mute status. Use true/false or 1/0")
		}

		err := client.Bus.SetMute(busNum, muted)
		if err != nil {
			return fmt.Errorf("Error setting bus mute status: %w", err)
		}

		cmd.Printf("Bus %d mute set to %v\n", busNum, muted)
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			level, err := client.Bus.Fader(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus fader level: %w", err)
			}
			cmd.Printf("Bus %d fader level: %.1f dB\n", busIndex, level)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and fader level (in dB)")
		}

		level := mustConvToFloat64(args[1])

		err := client.Bus.SetFader(busIndex, level)
		if err != nil {
			return fmt.Errorf("Error setting bus fader level: %w", err)
		}

		cmd.Printf("Bus %d fader set to %.2f dB\n", busIndex, level)
		return nil
	},
}

// busFadeOutCmd represents the bus fade out command.
var busFadeOutCmd = &cobra.Command{
	Short: "Fade out the bus fader over a specified duration",
	Long:  "Fade out the bus fader to minimum level over a specified duration in seconds.",
	Use:   "fadeout [bus number] --duration [seconds] [target level in dB]",
	Example: `  # Fade out bus 1 over 5 seconds
  xair-cli bus fadeout 1 --duration 5s -- -90.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetDuration("duration")
		if err != nil {
			return fmt.Errorf("Error getting duration flag: %w", err)
		}

		target := -90.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Bus.Fader(busIndex)
		if err != nil {
			return fmt.Errorf("Error getting current bus fader level: %w", err)
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(currentFader - target)
		if totalSteps <= 0 {
			cmd.Println("Bus is already at or below target level")
			return nil
		}

		stepDelay := time.Duration(duration.Seconds()*1000/totalSteps) * time.Millisecond

		for currentFader > target {
			currentFader -= 1.0
			err := client.Bus.SetFader(busIndex, currentFader)
			if err != nil {
				return fmt.Errorf("Error setting bus fader level: %w", err)
			}
			time.Sleep(stepDelay)
		}

		cmd.Println("Bus fade out completed")
		return nil
	},
}

// BusFadeInCmd represents the bus fade in command.
var busFadeInCmd = &cobra.Command{
	Short: "Fade in the bus fader over a specified duration",
	Long:  "Fade in the bus fader to maximum level over a specified duration in seconds.",
	Use:   "fadein [bus number] --duration [seconds] [target level in dB]",
	Example: `  # Fade in bus 1 over 5 seconds
  xair-cli bus fadein 1 --duration 5s -- 0.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		duration, err := cmd.Flags().GetDuration("duration")
		if err != nil {
			return fmt.Errorf("Error getting duration flag: %w", err)
		}

		target := 0.0
		if len(args) > 1 {
			target = mustConvToFloat64(args[1])
		}

		currentFader, err := client.Bus.Fader(busIndex)
		if err != nil {
			return fmt.Errorf("Error getting current bus fader level: %w", err)
		}

		// Calculate total steps needed to reach target dB
		totalSteps := float64(target - currentFader)
		if totalSteps <= 0 {
			cmd.Println("Bus is already at or above target level")
			return nil
		}

		stepDelay := time.Duration(duration.Seconds()*1000/totalSteps) * time.Millisecond

		for currentFader < target {
			currentFader += 1.0
			err := client.Bus.SetFader(busIndex, currentFader)
			if err != nil {
				return fmt.Errorf("Error setting bus fader level: %w", err)
			}
			time.Sleep(stepDelay)
		}

		cmd.Println("Bus fade in completed")
		return nil
	},
}

// busNameCmd represents the bus name command.
var busNameCmd = &cobra.Command{
	Short: "Get or set the bus name",
	Long:  `Get or set the name of a specific bus.`,
	Use:   "name [bus number] [new name]",
	Example: `  # Get the name of bus 1
  xair-cli bus name 1

  # Set the name of bus 1 to "Vocals"
  xair-cli bus name 1 Vocals`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			name, err := client.Bus.Name(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus name: %w", err)
			}
			cmd.Printf("Bus %d name: %s\n", busIndex, name)
			return nil
		}

		newName := args[1]
		err := client.Bus.SetName(busIndex, newName)
		if err != nil {
			return fmt.Errorf("Error setting bus name: %w", err)
		}

		cmd.Printf("Bus %d name set to: %s\n", busIndex, newName)
		return nil
	},
}

// busEqCmd represents the bus EQ command.
var busEqCmd = &cobra.Command{
	Short: "Commands to control bus EQ settings",
	Long:  `Commands to control the EQ of individual buses, including turning the EQ on or off.`,
	Use:   "eq",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// busEqOnCmd represents the bus EQ on/off command.
var busEqOnCmd = &cobra.Command{
	Short: "Get or set the bus EQ on/off status",
	Long:  `Get or set the EQ on/off status of a specific bus.`,
	Use:   "on [bus number] [true|false]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and EQ on status (true/false)")
		}

		busNum := mustConvToInt(args[0])
		var eqOn bool
		switch args[1] {
		case "true", "1":
			eqOn = true
		case "false", "0":
			eqOn = false
		default:
			return fmt.Errorf("Invalid EQ on status. Use true/false or 1/0")
		}

		err := client.Bus.Eq.SetOn(busNum, eqOn)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ on status: %w", err)
		}

		cmd.Printf("Bus %d EQ on set to %v\n", busNum, eqOn)
		return nil
	},
}

// busEqModeCmd represents the bus EQ mode command.
var busEqModeCmd = &cobra.Command{
	Short: "Get or set the bus EQ mode",
	Long:  `Get or set the EQ mode of a specific bus.`,
	Use:   "mode [bus number] [mode]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		modeNames := []string{"peq", "geq", "teq"}

		if len(args) == 1 {
			mode, err := client.Bus.Eq.Mode(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus EQ mode: %w", err)
			}
			cmd.Printf("Bus %d EQ mode: %s\n", busIndex, modeNames[mode])
			return nil
		}

		mode := indexOf(modeNames, args[1])
		if mode == -1 {
			return fmt.Errorf("Invalid EQ mode. Valid modes are: %v", modeNames)
		}

		err := client.Bus.Eq.SetMode(busIndex, mode)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ mode: %w", err)
		}

		cmd.Printf("Bus %d EQ mode set to %s\n", busIndex, modeNames[mode])
		return nil
	},
}

// busEqGainCmd represents the bus EQ gain command.
var busEqGainCmd = &cobra.Command{
	Short: "Get or set the bus EQ gain for a specific band",
	Long: `Get or set the EQ gain (in dB) for a specific band of a bus.
	
	Gain values range from -15.0 dB to +15.0 dB.`,
	Use: "gain [bus number] [band number] [gain in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and band number")
		}

		busIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			gain, err := client.Bus.Eq.Gain(busIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus EQ gain: %w", err)
			}
			cmd.Printf("Bus %d EQ band %d gain: %.1f dB\n", busIndex, bandIndex, gain)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide bus number, band number, and gain (in dB)")
		}

		gain := mustConvToFloat64(args[2])

		err := client.Bus.Eq.SetGain(busIndex, bandIndex, gain)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ gain: %w", err)
		}

		cmd.Printf("Bus %d EQ band %d gain set to %.1f dB\n", busIndex, bandIndex, gain)
		return nil
	},
}

// busEqFreqCmd represents the bus EQ frequency command.
var busEqFreqCmd = &cobra.Command{
	Short: "Get or set the bus EQ frequency for a specific band",
	Long:  `Get or set the EQ frequency (in Hz) for a specific band of a bus.`,
	Use:   "freq [bus number] [band number] [frequency in Hz]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and band number")
		}

		busIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			freq, err := client.Bus.Eq.Frequency(busIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus EQ frequency: %w", err)
			}
			cmd.Printf("Bus %d EQ band %d frequency: %.1f Hz\n", busIndex, bandIndex, freq)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide bus number, band number, and frequency (in Hz)")
		}

		freq := mustConvToFloat64(args[2])

		err := client.Bus.Eq.SetFrequency(busIndex, bandIndex, freq)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ frequency: %w", err)
		}

		cmd.Printf("Bus %d EQ band %d frequency set to %.1f Hz\n", busIndex, bandIndex, freq)
		return nil
	},
}

// busEqQCmd represents the bus EQ Q command.
var busEqQCmd = &cobra.Command{
	Short: "Get or set the bus EQ Q factor for a specific band",
	Long:  `Get or set the EQ Q factor for a specific band of a bus.`,
	Use:   "q [bus number] [band number] [Q factor]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and band number")
		}

		busIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		if len(args) == 2 {
			qFactor, err := client.Bus.Eq.Q(busIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus EQ Q factor: %w", err)
			}
			cmd.Printf("Bus %d EQ band %d Q factor: %.2f\n", busIndex, bandIndex, qFactor)
			return nil
		}

		if len(args) < 3 {
			return fmt.Errorf("Please provide bus number, band number, and Q factor")
		}

		qFactor := mustConvToFloat64(args[2])

		err := client.Bus.Eq.SetQ(busIndex, bandIndex, qFactor)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ Q factor: %w", err)
		}

		cmd.Printf("Bus %d EQ band %d Q factor set to %.2f\n", busIndex, bandIndex, qFactor)
		return nil
	},
}

// busEqTypeCmd represents the bus EQ type command.
var busEqTypeCmd = &cobra.Command{
	Short: "Get or set the bus EQ band type",
	Long:  `Get or set the EQ band type for a specific band of a bus.`,
	Use:   "type [bus number] [band number] [type]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and band number")
		}

		busIndex, bandIndex := func() (int, int) {
			return mustConvToInt(args[0]), mustConvToInt(args[1])
		}()

		eqTypeNames := []string{"lcut", "lshv", "peq", "veq", "hshv", "hcut"}

		if len(args) == 2 {
			currentType, err := client.Bus.Eq.Type(busIndex, bandIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus EQ band type: %w", err)
			}
			cmd.Printf("Bus %d EQ band %d type: %s\n", busIndex, bandIndex, eqTypeNames[currentType])
			return nil
		}

		eqType := indexOf(eqTypeNames, args[2])
		if eqType == -1 {
			return fmt.Errorf("Invalid EQ band type. Valid types are: %v", eqTypeNames)
		}

		err := client.Bus.Eq.SetType(busIndex, bandIndex, eqType)
		if err != nil {
			return fmt.Errorf("Error setting bus EQ band type: %w", err)
		}

		cmd.Printf("Bus %d EQ band %d type set to %s\n", busIndex, bandIndex, eqTypeNames[eqType])
		return nil
	},
}

// busCompCmd represents the bus Compressor command.
var busCompCmd = &cobra.Command{
	Short: "Commands to control bus Compressor settings",
	Long:  `Commands to control the Compressor of individual buses, including turning the Compressor on or off.`,
	Use:   "comp",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

// busCompOnCmd represents the bus Compressor on/off command.
var busCompOnCmd = &cobra.Command{
	Short: "Get or set the bus Compressor on/off status",
	Long:  `Get or set the Compressor on/off status of a specific bus.`,
	Use:   "on [bus number] [true|false]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and Compressor on status (true/false)")
		}

		busNum := mustConvToInt(args[0])
		var compOn bool
		switch args[1] {
		case "true", "1":
			compOn = true
		case "false", "0":
			compOn = false
		default:
			return fmt.Errorf("Invalid Compressor on status. Use true/false or 1/0")
		}

		err := client.Bus.Comp.SetOn(busNum, compOn)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor on status: %w", err)
		}

		cmd.Printf("Bus %d Compressor on set to %v\n", busNum, compOn)
		return nil
	},
}

// busCompThresholdCmd represents the bus Compressor threshold command.
var busCompThresholdCmd = &cobra.Command{
	Short: "Get or set the bus Compressor threshold",
	Long:  `Get or set the Compressor threshold (in dB) for a specific bus.`,
	Use:   "threshold [bus number] [threshold in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			threshold, err := client.Bus.Comp.Threshold(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor threshold: %w", err)
			}
			cmd.Printf("Bus %d Compressor threshold: %.1f dB\n", busIndex, threshold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and threshold (in dB)")
		}

		threshold := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetThreshold(busIndex, threshold)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor threshold: %w", err)
		}

		cmd.Printf("Bus %d Compressor threshold set to %.1f dB\n", busIndex, threshold)
		return nil
	},
}

// busCompRatioCmd represents the bus Compressor ratio command.
var busCompRatioCmd = &cobra.Command{
	Short: "Get or set the bus Compressor ratio",
	Long:  `Get or set the Compressor ratio for a specific bus.`,
	Use:   "ratio [bus number] [ratio]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			ratio, err := client.Bus.Comp.Ratio(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor ratio: %w", err)
			}
			cmd.Printf("Bus %d Compressor ratio: %.2f\n", busIndex, ratio)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and ratio")
		}

		ratio := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetRatio(busIndex, ratio)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor ratio: %w", err)
		}

		cmd.Printf("Bus %d Compressor ratio set to %.2f\n", busIndex, ratio)
		return nil
	},
}

// busMixCmd represents the bus Compressor mix command.
var busCompMixCmd = &cobra.Command{
	Short: "Get or set the bus Compressor mix",
	Long:  `Get or set the Compressor mix (0-100%) for a specific bus.`,
	Use:   "mix [bus number] [mix percentage]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			mix, err := client.Bus.Comp.Mix(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor mix: %w", err)
			}
			cmd.Printf("Bus %d Compressor mix: %.1f%%\n", busIndex, mix)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and mix percentage")
		}

		mix := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetMix(busIndex, mix)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor mix: %w", err)
		}

		cmd.Printf("Bus %d Compressor mix set to %.1f%%\n", busIndex, mix)
		return nil
	},
}

// busMakeUpCmd represents the bus Compressor make-up gain command.
var busCompMakeUpCmd = &cobra.Command{
	Short: "Get or set the bus Compressor make-up gain",
	Long:  `Get or set the Compressor make-up gain (in dB) for a specific bus.`,
	Use:   "makeup [bus number] [make-up gain in dB]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			makeUp, err := client.Bus.Comp.MakeUp(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor make-up gain: %w", err)
			}
			cmd.Printf("Bus %d Compressor make-up gain: %.1f dB\n", busIndex, makeUp)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and make-up gain (in dB)")
		}

		makeUp := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetMakeUp(busIndex, makeUp)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor make-up gain: %w", err)
		}

		cmd.Printf("Bus %d Compressor make-up gain set to %.1f dB\n", busIndex, makeUp)
		return nil
	},
}

// busAttackCmd represents the bus Compressor attack time command.
var busCompAttackCmd = &cobra.Command{
	Short: "Get or set the bus Compressor attack time",
	Long:  `Get or set the Compressor attack time (in milliseconds) for a specific bus.`,
	Use:   "attack [bus number] [attack time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			attack, err := client.Bus.Comp.Attack(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor attack time: %w", err)
			}
			cmd.Printf("Bus %d Compressor attack time: %.1f ms\n", busIndex, attack)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and attack time (in ms)")
		}

		attack := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetAttack(busIndex, attack)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor attack time: %w", err)
		}

		cmd.Printf("Bus %d Compressor attack time set to %.1f ms\n", busIndex, attack)
		return nil
	},
}

// busHoldCmd represents the bus Compressor hold time command.
var busCompHoldCmd = &cobra.Command{
	Short: "Get or set the bus Compressor hold time",
	Long:  `Get or set the Compressor hold time (in milliseconds) for a specific bus.`,
	Use:   "hold [bus number] [hold time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			hold, err := client.Bus.Comp.Hold(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor hold time: %w", err)
			}
			cmd.Printf("Bus %d Compressor hold time: %.2f ms\n", busIndex, hold)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and hold time (in ms)")
		}

		hold := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetHold(busIndex, hold)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor hold time: %w", err)
		}

		cmd.Printf("Bus %d Compressor hold time set to %.2f ms\n", busIndex, hold)
		return nil
	},
}

// busReleaseCmd represents the bus Compressor release time command.
var busCompReleaseCmd = &cobra.Command{
	Short: "Get or set the bus Compressor release time",
	Long:  `Get or set the Compressor release time (in milliseconds) for a specific bus.`,
	Use:   "release [bus number] [release time in ms]",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ClientFromContext(cmd.Context())
		if client == nil {
			return fmt.Errorf("OSC client not found in context")
		}

		if len(args) < 1 {
			return fmt.Errorf("Please provide bus number")
		}

		busIndex := mustConvToInt(args[0])

		if len(args) == 1 {
			release, err := client.Bus.Comp.Release(busIndex)
			if err != nil {
				return fmt.Errorf("Error getting bus Compressor release time: %w", err)
			}
			cmd.Printf("Bus %d Compressor release time: %.1f ms\n", busIndex, release)
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Please provide bus number and release time (in ms)")
		}

		release := mustConvToFloat64(args[1])

		err := client.Bus.Comp.SetRelease(busIndex, release)
		if err != nil {
			return fmt.Errorf("Error setting bus Compressor release time: %w", err)
		}

		cmd.Printf("Bus %d Compressor release time set to %.1f ms\n", busIndex, release)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(busCmd)

	busCmd.AddCommand(busMuteCmd)
	busCmd.AddCommand(busFaderCmd)
	busCmd.AddCommand(busFadeOutCmd)
	busFadeOutCmd.Flags().DurationP("duration", "d", 5*time.Second, "Duration for fade out in seconds")
	busCmd.AddCommand(busFadeInCmd)
	busFadeInCmd.Flags().DurationP("duration", "d", 5*time.Second, "Duration for fade in in seconds")
	busCmd.AddCommand(busNameCmd)

	busCmd.AddCommand(busEqCmd)
	busEqCmd.AddCommand(busEqOnCmd)
	busEqCmd.AddCommand(busEqModeCmd)
	busEqCmd.AddCommand(busEqGainCmd)
	busEqCmd.AddCommand(busEqFreqCmd)
	busEqCmd.AddCommand(busEqQCmd)
	busEqCmd.AddCommand(busEqTypeCmd)

	busCmd.AddCommand(busCompCmd)
	busCompCmd.AddCommand(busCompOnCmd)
	busCompCmd.AddCommand(busCompThresholdCmd)
	busCompCmd.AddCommand(busCompRatioCmd)
	busCompCmd.AddCommand(busCompMixCmd)
	busCompCmd.AddCommand(busCompMakeUpCmd)
	busCompCmd.AddCommand(busCompAttackCmd)
	busCompCmd.AddCommand(busCompHoldCmd)
	busCompCmd.AddCommand(busCompReleaseCmd)
}
