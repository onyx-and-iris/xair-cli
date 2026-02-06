package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

// BusCmdGroup defines the commands related to controlling the buses of the X-Air device.
type BusCmdGroup struct {
	Index struct {
		Index   int           `arg:"" help:"The index of the bus. (1-based indexing)"`
		Mute    BusMuteCmd    `       help:"Get or set the mute state of the bus." cmd:""`
		Fader   BusFaderCmd   `     help:"Get or set the fader level of the bus." cmd:""`
		Fadein  BusFadeinCmd  `      help:"Fade in the bus over a specified duration." cmd:""`
		Fadeout BusFadeoutCmd `     help:"Fade out the bus over a specified duration." cmd:""`
		Name    BusNameCmd    `       help:"Get or set the name of the bus." cmd:""`

		Eq   BusEqCmdGroup   `       help:"Commands related to the bus EQ." cmd:"eq"`
		Comp BusCompCmdGroup `     help:"Commands related to the bus compressor." cmd:"comp"`
	} `arg:"" help:"Control a specific bus by index."`
}

// BusMuteCmd defines the command for getting or setting the mute state of a bus.
type BusMuteCmd struct {
	State *string `arg:"" help:"The mute state to set (true or false). If not provided, the current mute state will be returned." optional:"" enum:"true,false"`
}

// Run executes the BusMuteCmd command, either retrieving the current mute state or setting it based on the provided argument.
func (cmd *BusMuteCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.Bus.Mute(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d mute state: %t\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.SetMute(bus.Index.Index, *cmd.State == "true"); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d mute state set to: %s\n", bus.Index.Index, *cmd.State)
	return nil
}

// BusFaderCmd defines the command for getting or setting the fader level of a bus.
type BusFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set (in dB). If not provided, the current fader level will be returned." optional:""`
}

// Run executes the BusFaderCmd command, either retrieving the current fader level or setting it based on the provided argument.
func (cmd *BusFaderCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Bus.Fader(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d fader level: %.2f dB\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.SetFader(bus.Index.Index, *cmd.Level); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d fader level set to: %.2f dB\n", bus.Index.Index, *cmd.Level)
	return nil
}

// BusFadeinCmd defines the command for fading in a bus over a specified duration to a target fader level.
type BusFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in effect." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."     default:"0.0" arg:""`
}

// Run executes the BusFadeinCmd command, gradually increasing the fader level of the bus from its current level to the target level over the specified duration.
func (cmd *BusFadeinCmd) Run(ctx *context, bus *BusCmdGroup) error {
	currentLevel, err := ctx.Client.Bus.Fader(bus.Index.Index)
	if err != nil {
		return fmt.Errorf("failed to get current fader level: %w", err)
	}

	if currentLevel >= cmd.Target {
		return fmt.Errorf(
			"current fader level (%.2f dB) is already at or above the target level (%.2f dB)",
			currentLevel,
			cmd.Target,
		)
	}

	totalSteps := float64(cmd.Target - currentLevel)
	stepDuration := time.Duration(cmd.Duration.Seconds()*1000/totalSteps) * time.Millisecond
	for currentLevel < cmd.Target {
		currentLevel += totalSteps / float64(cmd.Duration.Seconds()*1000/stepDuration.Seconds())
		if currentLevel > cmd.Target {
			currentLevel = cmd.Target
		}

		if err := ctx.Client.Bus.SetFader(bus.Index.Index, currentLevel); err != nil {
			return fmt.Errorf("failed to set fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}

	fmt.Fprintf(ctx.Out, "Bus %d fade-in complete. Final level: %.2f dB\n", bus.Index.Index, cmd.Target)
	return nil
}

// BusFadeoutCmd defines the command for fading out a bus over a specified duration to a target fader level.
type BusFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out effect." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."      default:"-90.0" arg:""`
}

// Run executes the BusFadeoutCmd command, gradually decreasing the fader level of the bus from its current level to the target level over the specified duration.
func (cmd *BusFadeoutCmd) Run(ctx *context, bus *BusCmdGroup) error {
	currentLevel, err := ctx.Client.Bus.Fader(bus.Index.Index)
	if err != nil {
		return fmt.Errorf("failed to get current fader level: %w", err)
	}

	if currentLevel <= cmd.Target {
		return fmt.Errorf(
			"current fader level (%.2f dB) is already at or below the target level (%.2f dB)",
			currentLevel,
			cmd.Target,
		)
	}

	totalSteps := float64(currentLevel - cmd.Target)
	stepDuration := time.Duration(cmd.Duration.Seconds()*1000/totalSteps) * time.Millisecond
	for currentLevel > cmd.Target {
		currentLevel -= totalSteps / float64(cmd.Duration.Seconds()*1000/stepDuration.Seconds())
		if currentLevel < cmd.Target {
			currentLevel = cmd.Target
		}

		if err := ctx.Client.Bus.SetFader(bus.Index.Index, currentLevel); err != nil {
			return fmt.Errorf("failed to set fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}

	fmt.Fprintf(ctx.Out, "Bus %d fade-out complete. Final level: %.2f dB\n", bus.Index.Index, cmd.Target)
	return nil
}

// BusNameCmd defines the command for getting or setting the name of a bus.
type BusNameCmd struct {
	Name *string `arg:"" help:"The name to set for the bus. If not provided, the current name will be returned." optional:""`
}

// Run executes the BusNameCmd command, either retrieving the current name of the bus or setting it based on the provided argument.
func (cmd *BusNameCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Name == nil {
		resp, err := ctx.Client.Bus.Name(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d name: %s\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.SetName(bus.Index.Index, *cmd.Name); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d name set to: %s\n", bus.Index.Index, *cmd.Name)
	return nil
}

// BusEqCmdGroup defines the commands related to controlling the EQ of a bus.
type BusEqCmdGroup struct {
	On   BusEqOnCmd   `help:"Get or set the EQ on/off state of the bus."              cmd:"on"`
	Mode BusEqModeCmd `help:"Get or set the EQ mode of the bus (peq, geq or teq)."    cmd:"mode"`
	Band struct {
		Band int              `arg:"" help:"The EQ band number."`
		Gain BusEqBandGainCmd `help:"Get or set the gain of the EQ band." cmd:"gain"`
		Freq BusEqBandFreqCmd `help:"Get or set the frequency of the EQ band." cmd:"freq"`
		Q    BusEqBandQCmd    `help:"Get or set the Q factor of the EQ band." cmd:"q"`
		Type BusEqBandTypeCmd `help:"Get or set the type of the EQ band (lcut, lshv, peq, veq, hshv, hcut)." cmd:"type"`
	} `help:"Commands for controlling a specific EQ band of the bus."            arg:""`
}

// Validate checks that the provided EQ band number is within the valid range (1-6).
func (cmd *BusEqCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 6 {
		return fmt.Errorf("EQ band number must be between 1 and 6")
	}
	return nil
}

// BusCompCmdGroup defines the commands related to controlling the compressor of a bus.
type BusEqOnCmd struct {
	State *string `arg:"" help:"The EQ on/off state to set (true or false). If not provided, the current EQ state will be returned." optional:"" enum:"true,false"`
}

// Run executes the BusEqOnCmd command, either retrieving the current EQ on/off state of the bus or setting it based on the provided argument.
func (cmd *BusEqOnCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.Bus.Eq.On(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ on state: %t\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetOn(bus.Index.Index, *cmd.State == "true"); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ on state set to: %s\n", bus.Index.Index, *cmd.State)
	return nil
}

// BusEqModeCmd defines the command for getting or setting the EQ mode of a bus.
type BusEqModeCmd struct {
	Mode *string `arg:"" help:"The EQ mode to set (peq, geq or teq). If not provided, the current EQ mode will be returned." optional:"" enum:"peq,geq,teq"`
}

// Run executes the BusEqModeCmd command, either retrieving the current EQ mode of the bus or setting it based on the provided argument.
func (cmd *BusEqModeCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Bus.Eq.Mode(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ mode: %s\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetMode(bus.Index.Index, *cmd.Mode); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ mode set to: %s\n", bus.Index.Index, *cmd.Mode)
	return nil
}

// BusEqBandGainCmd defines the command for getting or setting the gain of a specific EQ band of a bus.
type BusEqBandGainCmd struct {
	Gain *float64 `arg:"" help:"The gain to set for the EQ band (in dB). If not provided, the current gain will be returned." optional:""`
}

// Run executes the BusEqBandGainCmd command, either retrieving the current gain of the specified EQ band of the bus or setting it based on the provided argument.
func (cmd *BusEqBandGainCmd) Run(ctx *context, bus *BusCmdGroup, busEq *BusEqCmdGroup) error {
	if cmd.Gain == nil {
		resp, err := ctx.Client.Bus.Eq.Gain(bus.Index.Index, busEq.Band.Band)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ band %d gain: %.2f dB\n", bus.Index.Index, busEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetGain(bus.Index.Index, busEq.Band.Band, *cmd.Gain); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ band %d gain set to: %.2f dB\n", bus.Index.Index, busEq.Band.Band, *cmd.Gain)
	return nil
}

// BusEqBandFreqCmd defines the command for getting or setting the frequency of a specific EQ band of a bus.
type BusEqBandFreqCmd struct {
	Freq *float64 `arg:"" help:"The frequency to set for the EQ band (in Hz). If not provided, the current frequency will be returned." optional:""`
}

// Run executes the BusEqBandFreqCmd command, either retrieving the current frequency of the specified EQ band of the bus or setting it based on the provided argument.
func (cmd *BusEqBandFreqCmd) Run(ctx *context, bus *BusCmdGroup, busEq *BusEqCmdGroup) error {
	if cmd.Freq == nil {
		resp, err := ctx.Client.Bus.Eq.Frequency(bus.Index.Index, busEq.Band.Band)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ band %d frequency: %.2f Hz\n", bus.Index.Index, busEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetFrequency(bus.Index.Index, busEq.Band.Band, *cmd.Freq); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ band %d frequency set to: %.2f Hz\n", bus.Index.Index, busEq.Band.Band, *cmd.Freq)
	return nil
}

// BusEqBandQCmd defines the command for getting or setting the Q factor of a specific EQ band of a bus.
type BusEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the EQ band. If not provided, the current Q factor will be returned." optional:""`
}

// Run executes the BusEqBandQCmd command, either retrieving the current Q factor of the specified EQ band of the bus or setting it based on the provided argument.
func (cmd *BusEqBandQCmd) Run(ctx *context, bus *BusCmdGroup, busEq *BusEqCmdGroup) error {
	if cmd.Q == nil {
		resp, err := ctx.Client.Bus.Eq.Q(bus.Index.Index, busEq.Band.Band)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ band %d Q factor: %.2f\n", bus.Index.Index, busEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetQ(bus.Index.Index, busEq.Band.Band, *cmd.Q); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ band %d Q factor set to: %.2f\n", bus.Index.Index, busEq.Band.Band, *cmd.Q)
	return nil
}

// BusEqBandTypeCmd defines the command for getting or setting the type of a specific EQ band of a bus.
type BusEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the EQ band (lcut, lshv, peq, veq, hshv, hcut). If not provided, the current type will be returned." optional:"" enum:"lcut,lshv,peq,veq,hshv,hcut"`
}

// Run executes the BusEqBandTypeCmd command, either retrieving the current type of the specified EQ band of the bus or setting it based on the provided argument.
func (cmd *BusEqBandTypeCmd) Run(ctx *context, bus *BusCmdGroup, busEq *BusEqCmdGroup) error {
	if cmd.Type == nil {
		resp, err := ctx.Client.Bus.Eq.Type(bus.Index.Index, busEq.Band.Band)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d EQ band %d type: %s\n", bus.Index.Index, busEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Bus.Eq.SetType(bus.Index.Index, busEq.Band.Band, *cmd.Type); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d EQ band %d type set to: %s\n", bus.Index.Index, busEq.Band.Band, *cmd.Type)
	return nil
}

// BusCompCmdGroup defines the commands related to controlling the compressor of a bus.
type BusCompCmdGroup struct {
	On        BusCompOnCmd        `help:"Get or set the compressor on/off state of the bus."         cmd:"on"`
	Mode      BusCompModeCmd      `help:"Get or set the compressor mode of the bus (comp, exp)."     cmd:"mode"`
	Threshold BusCompThresholdCmd `help:"Get or set the compressor threshold of the bus (in dB)."    cmd:"threshold"`
	Ratio     BusCompRatioCmd     `help:"Get or set the compressor ratio of the bus."                cmd:"ratio"`
	Mix       BusCompMixCmd       `help:"Get or set the compressor mix level of the bus (in %)."     cmd:"mix"`
	Makeup    BusCompMakeupCmd    `help:"Get or set the compressor makeup gain of the bus (in dB)."  cmd:"makeup"`
	Attack    BusCompAttackCmd    `help:"Get or set the compressor attack time of the bus (in ms)."  cmd:"attack"`
	Hold      BusCompHoldCmd      `help:"Get or set the compressor hold time of the bus (in ms)."    cmd:"hold"`
	Release   BusCompReleaseCmd   `help:"Get or set the compressor release time of the bus (in ms)." cmd:"release"`
}

// BusCompOnCmd defines the command for getting or setting the compressor on/off state of a bus.
type BusCompOnCmd struct {
	State *string `arg:"" help:"The compressor on/off state to set (true or false). If not provided, the current compressor state will be returned." optional:"" enum:"true,false"`
}

// Run executes the BusCompOnCmd command, either retrieving the current compressor on/off state of the bus or setting it based on the provided argument.
func (cmd *BusCompOnCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.Bus.Comp.On(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor on state: %t\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetOn(bus.Index.Index, *cmd.State == "true"); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor on state set to: %s\n", bus.Index.Index, *cmd.State)
	return nil
}

// BusCompModeCmd defines the command for getting or setting the compressor mode of a bus.
type BusCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set (comp, exp). If not provided, the current compressor mode will be returned." optional:"" enum:"comp,exp"`
}

// Run executes the BusCompModeCmd command, either retrieving the current compressor mode of the bus or setting it based on the provided argument.
func (cmd *BusCompModeCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Bus.Comp.Mode(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor mode: %s\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetMode(bus.Index.Index, *cmd.Mode); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor mode set to: %s\n", bus.Index.Index, *cmd.Mode)
	return nil
}

// BusCompThresholdCmd defines the command for getting or setting the compressor threshold of a bus.
type BusCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set (in dB). If not provided, the current compressor threshold will be returned." optional:""`
}

// Run executes the BusCompThresholdCmd command, either retrieving the current compressor threshold of the bus or setting it based on the provided argument.
func (cmd *BusCompThresholdCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.Bus.Comp.Threshold(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor threshold: %.2f dB\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetThreshold(bus.Index.Index, *cmd.Threshold); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor threshold set to: %.2f dB\n", bus.Index.Index, *cmd.Threshold)
	return nil
}

// BusCompRatioCmd defines the command for getting or setting the compressor ratio of a bus.
type BusCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set. If not provided, the current compressor ratio will be returned." optional:""`
}

// Run executes the BusCompRatioCmd command, either retrieving the current compressor ratio of the bus or setting it based on the provided argument.
func (cmd *BusCompRatioCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Ratio == nil {
		resp, err := ctx.Client.Bus.Comp.Ratio(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor ratio: %.2f\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetRatio(bus.Index.Index, *cmd.Ratio); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor ratio set to: %.2f\n", bus.Index.Index, *cmd.Ratio)
	return nil
}

// BusCompMixCmd defines the command for getting or setting the compressor mix level of a bus.
type BusCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix level to set (in %). If not provided, the current compressor mix level will be returned." optional:""`
}

// Run executes the BusCompMixCmd command, either retrieving the current compressor mix level of the bus or setting it based on the provided argument.
func (cmd *BusCompMixCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Mix == nil {
		resp, err := ctx.Client.Bus.Comp.Mix(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor mix level: %.2f%%\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetMix(bus.Index.Index, *cmd.Mix); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor mix level set to: %.2f%%\n", bus.Index.Index, *cmd.Mix)
	return nil
}

// BusCompMakeupCmd defines the command for getting or setting the compressor makeup gain of a bus.
type BusCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set (in dB). If not provided, the current compressor makeup gain will be returned." optional:""`
}

// Run executes the BusCompMakeupCmd command, either retrieving the current compressor makeup gain of the bus or setting it based on the provided argument.
func (cmd *BusCompMakeupCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Makeup == nil {
		resp, err := ctx.Client.Bus.Comp.Makeup(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor makeup gain: %.2f dB\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetMakeup(bus.Index.Index, *cmd.Makeup); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor makeup gain set to: %.2f dB\n", bus.Index.Index, *cmd.Makeup)
	return nil
}

// BusCompAttackCmd defines the command for getting or setting the compressor attack time of a bus.
type BusCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set (in ms). If not provided, the current compressor attack time will be returned." optional:""`
}

// Run executes the BusCompAttackCmd command, either retrieving the current compressor attack time of the bus or setting it based on the provided argument.
func (cmd *BusCompAttackCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.Bus.Comp.Attack(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor attack time: %.2f ms\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetAttack(bus.Index.Index, *cmd.Attack); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor attack time set to: %.2f ms\n", bus.Index.Index, *cmd.Attack)
	return nil
}

// BusCompHoldCmd defines the command for getting or setting the compressor hold time of a bus.
type BusCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set (in ms). If not provided, the current compressor hold time will be returned." optional:""`
}

// Run executes the BusCompHoldCmd command, either retrieving the current compressor hold time of the bus or setting it based on the provided argument.
func (cmd *BusCompHoldCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.Bus.Comp.Hold(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor hold time: %.2f ms\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetHold(bus.Index.Index, *cmd.Hold); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor hold time set to: %.2f ms\n", bus.Index.Index, *cmd.Hold)
	return nil
}

// BusCompReleaseCmd defines the command for getting or setting the compressor release time of a bus.
type BusCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set (in ms). If not provided, the current compressor release time will be returned." optional:""`
}

// Run executes the BusCompReleaseCmd command, either retrieving the current compressor release time of the bus or setting it based on the provided argument.
func (cmd *BusCompReleaseCmd) Run(ctx *context, bus *BusCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.Bus.Comp.Release(bus.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintf(ctx.Out, "Bus %d compressor release time: %.2f ms\n", bus.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Bus.Comp.SetRelease(bus.Index.Index, *cmd.Release); err != nil {
		return err
	}
	fmt.Fprintf(ctx.Out, "Bus %d compressor release time set to: %.2f ms\n", bus.Index.Index, *cmd.Release)
	return nil
}
