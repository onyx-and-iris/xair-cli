package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

// StripCmdGroup defines the command group for controlling the strips of the mixer, including commands for getting and setting various parameters such as mute state, fader level, send levels, and EQ settings.
type StripCmdGroup struct {
	Index struct {
		Index   int             `arg:"" help:"The index of the strip. (1-based indexing)"`
		Mute    StripMuteCmd    `       help:"Get or set the mute state of the strip." cmd:""`
		Fader   StripFaderCmd   `     help:"Get or set the fader level of the strip." cmd:""`
		Fadein  StripFadeinCmd  `      help:"Fade in the strip over a specified duration." cmd:""`
		Fadeout StripFadeoutCmd `     help:"Fade out the strip over a specified duration." cmd:""`
		Send    StripSendCmd    `      help:"Get or set the send level for a specific bus." cmd:""`
		Name    StripNameCmd    `      help:"Get or set the name of the strip." cmd:""`

		Gate StripGateCmdGroup `     help:"Commands related to the strip gate." cmd:"gate"`
		Eq   StripEqCmdGroup   `       help:"Commands related to the strip EQ." cmd:"eq"`
		Comp StripCompCmdGroup `      help:"Commands related to the strip compressor." cmd:"comp"`
	} `arg:"" help:"Control a specific strip by index."`
}

// StripMuteCmd defines the command for getting or setting the mute state of a strip.
type StripMuteCmd struct {
	State *string `arg:"" help:"The mute state to set (true or false). If not provided, the current mute state will be returned." optional:"" enum:"true,false"`
}

// Run executes the StripMuteCmd command, either retrieving the current mute state of the strip or setting it based on the provided argument.
func (cmd *StripMuteCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.Strip.Mute(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get mute state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d mute state: %t\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.SetMute(strip.Index.Index, *cmd.State == "true"); err != nil {
		return fmt.Errorf("failed to set mute state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d mute state set to: %s\n", strip.Index.Index, *cmd.State)
	return nil
}

// StripFaderCmd defines the command for getting or setting the fader level of a strip.
type StripFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set (in dB)." optional:""`
}

// Run executes the StripFaderCmd command, either retrieving the current fader level of the strip or setting it based on the provided argument.
func (cmd *StripFaderCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Strip.Fader(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get fader level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d fader level: %.2f dB\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.SetFader(strip.Index.Index, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set fader level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d fader level set to: %.2f dB\n", strip.Index.Index, *cmd.Level)
	return nil
}

// StripFadeinCmd defines the command for fading in a strip over a specified duration, gradually increasing the fader level from its current value to a target value.
type StripFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in (in seconds)." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."           default:"0.0" arg:""`
}

// Run executes the StripFadeinCmd command, gradually increasing the fader level of the strip from its current value to the specified target value over the specified duration.
func (cmd *StripFadeinCmd) Run(ctx *context, strip *StripCmdGroup) error {
	currentLevel, err := ctx.Client.Strip.Fader(strip.Index.Index)
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
		currentLevel++
		if err := ctx.Client.Strip.SetFader(strip.Index.Index, currentLevel); err != nil {
			return fmt.Errorf("failed to set fader level during fade-in: %w", err)
		}
		time.Sleep(stepDuration)
	}

	fmt.Fprintf(ctx.Out, "Strip %d fade-in complete. Final level: %.2f dB\n", strip.Index.Index, cmd.Target)
	return nil
}

// StripFadeoutCmd defines the command for fading out a strip over a specified duration, gradually decreasing the fader level from its current value to a target value.
type StripFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out (in seconds)." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."            default:"-90.0" arg:""`
}

// Run executes the StripFadeoutCmd command, gradually decreasing the fader level of the strip from its current value to the specified target value over the specified duration.
func (cmd *StripFadeoutCmd) Run(ctx *context, strip *StripCmdGroup) error {
	{
		currentLevel, err := ctx.Client.Strip.Fader(strip.Index.Index)
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
			currentLevel--
			if err := ctx.Client.Strip.SetFader(strip.Index.Index, currentLevel); err != nil {
				return fmt.Errorf("failed to set fader level during fade-out: %w", err)
			}
			time.Sleep(stepDuration)
		}

		fmt.Fprintf(ctx.Out, "Strip %d fade-out complete. Final level: %.2f dB\n", strip.Index.Index, cmd.Target)
		return nil
	}
}

// StripSendCmd defines the command for getting or setting the send level for a specific bus on a strip, allowing users to control the level of the signal being sent from the strip to a particular bus.
type StripSendCmd struct {
	BusNum int      `arg:"" help:"The bus number to get or set the send level for."`
	Level  *float64 `arg:"" help:"The send level to set (in dB)."                   optional:""`
}

// Run executes the StripSendCmd command, either retrieving the current send level for the specified bus on the strip or setting it based on the provided argument.
func (cmd *StripSendCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Strip.SendLevel(strip.Index.Index, cmd.BusNum)
		if err != nil {
			return fmt.Errorf("failed to get send level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d send level for bus %d: %.2f dB\n", strip.Index.Index, cmd.BusNum, resp)
		return nil
	}

	if err := ctx.Client.Strip.SetSendLevel(strip.Index.Index, cmd.BusNum, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set send level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d send level for bus %d set to: %.2f dB\n", strip.Index.Index, cmd.BusNum, *cmd.Level)
	return nil
}

// StripNameCmd defines the command for getting or setting the name of a strip, allowing users to assign custom names to strips for easier identification and organization.
type StripNameCmd struct {
	Name *string `arg:"" help:"The name to set for the strip." optional:""`
}

// Run executes the StripNameCmd command, either retrieving the current name of the strip or setting it based on the provided argument.
func (cmd *StripNameCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Name == nil {
		resp, err := ctx.Client.Strip.Name(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get strip name: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d name: %s\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.SetName(strip.Index.Index, *cmd.Name); err != nil {
		return fmt.Errorf("failed to set strip name: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d name set to: %s\n", strip.Index.Index, *cmd.Name)
	return nil
}

// StripGateCmdGroup defines the command group for controlling the gate settings of a strip, including commands for getting and setting the gate on/off state, mode, threshold, range, attack time, hold time, and release time.
type StripGateCmdGroup struct {
	On        StripGateOnCmd        `help:"Get or set the gate on/off state of the strip." cmd:""`
	Mode      StripGateModeCmd      `help:"Get or set the gate mode of the strip."         cmd:""`
	Threshold StripGateThresholdCmd `help:"Get or set the gate threshold of the strip."    cmd:""`
	Range     StripGateRangeCmd     `help:"Get or set the gate range of the strip."        cmd:""`
	Attack    StripGateAttackCmd    `help:"Get or set the gate attack time of the strip."  cmd:""`
	Hold      StripGateHoldCmd      `help:"Get or set the gate hold time of the strip."    cmd:""`
	Release   StripGateReleaseCmd   `help:"Get or set the gate release time of the strip." cmd:""`
}

// StripGateOnCmd defines the command for getting or setting the gate on/off state of a strip, allowing users to enable or disable the gate effect on the strip.
type StripGateOnCmd struct {
	Enable *string `arg:"" help:"Whether to enable or disable the gate." optional:"" enum:"true,false"`
}

// Run executes the StripGateOnCmd command, either retrieving the current gate on/off state of the strip or setting it based on the provided argument.
func (cmd *StripGateOnCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Strip.Gate.On(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate state: %t\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetOn(strip.Index.Index, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set gate state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate state set to: %s\n", strip.Index.Index, *cmd.Enable)
	return nil
}

// StripGateModeCmd defines the command for getting or setting the gate mode of a strip, allowing users to choose from different gate modes such as exp2, exp3, exp4, gate, or duck.
type StripGateModeCmd struct {
	Mode *string `arg:"" help:"The gate mode to set." optional:"" enum:"exp2,exp3,exp4,gate,duck"`
}

// Run executes the StripGateModeCmd command, either retrieving the current gate mode of the strip or setting it based on the provided argument.
func (cmd *StripGateModeCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Strip.Gate.Mode(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate mode: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate mode: %s\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetMode(strip.Index.Index, *cmd.Mode); err != nil {
		return fmt.Errorf("failed to set gate mode: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate mode set to: %s\n", strip.Index.Index, *cmd.Mode)
	return nil
}

// StripGateThresholdCmd defines the command for getting or setting the gate threshold of a strip, allowing users to specify the threshold level at which the gate will start to attenuate the signal.
type StripGateThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The gate threshold to set (in dB)." optional:""`
}

// Run executes the StripGateThresholdCmd command, either retrieving the current gate threshold of the strip or setting it based on the provided argument.
func (cmd *StripGateThresholdCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.Strip.Gate.Threshold(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate threshold: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate threshold: %.2f\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetThreshold(strip.Index.Index, *cmd.Threshold); err != nil {
		return fmt.Errorf("failed to set gate threshold: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate threshold set to: %.2f\n", strip.Index.Index, *cmd.Threshold)
	return nil
}

// StripGateRangeCmd defines the command for getting or setting the gate range of a strip, allowing users to specify the amount of attenuation applied by the gate when the signal falls below the threshold.
type StripGateRangeCmd struct {
	Range *float64 `arg:"" help:"The gate range to set (in dB)." optional:""`
}

// Run executes the StripGateRangeCmd command, either retrieving the current gate range of the strip or setting it based on the provided argument.
func (cmd *StripGateRangeCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Range == nil {
		resp, err := ctx.Client.Strip.Gate.Range(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate range: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate range: %.2f\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetRange(strip.Index.Index, *cmd.Range); err != nil {
		return fmt.Errorf("failed to set gate range: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate range set to: %.2f\n", strip.Index.Index, *cmd.Range)
	return nil
}

// StripGateAttackCmd defines the command for getting or setting the gate attack time of a strip, allowing users to specify the time it takes for the gate to fully open after the signal exceeds the threshold.
type StripGateAttackCmd struct {
	Attack *float64 `arg:"" help:"The gate attack time to set (in ms)." optional:""`
}

// Run executes the StripGateAttackCmd command, either retrieving the current gate attack time of the strip or setting it based on the provided argument.
func (cmd *StripGateAttackCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.Strip.Gate.Attack(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate attack time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate attack time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetAttack(strip.Index.Index, *cmd.Attack); err != nil {
		return fmt.Errorf("failed to set gate attack time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate attack time set to: %.2f ms\n", strip.Index.Index, *cmd.Attack)
	return nil
}

// StripGateHoldCmd defines the command for getting or setting the gate hold time of a strip, allowing users to specify the time that the gate remains open after the signal falls below the threshold before it starts to close.
type StripGateHoldCmd struct {
	Hold *float64 `arg:"" help:"The gate hold time to set (in ms)." optional:""`
}

// Run executes the StripGateHoldCmd command, either retrieving the current gate hold time of the strip or setting it based on the provided argument.
func (cmd *StripGateHoldCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.Strip.Gate.Hold(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate hold time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate hold time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetHold(strip.Index.Index, *cmd.Hold); err != nil {
		return fmt.Errorf("failed to set gate hold time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate hold time set to: %.2f ms\n", strip.Index.Index, *cmd.Hold)
	return nil
}

// StripGateReleaseCmd defines the command for getting or setting the gate release time of a strip, allowing users to specify the time it takes for the gate to fully close after the signal falls below the threshold and the hold time has elapsed.
type StripGateReleaseCmd struct {
	Release *float64 `arg:"" help:"The gate release time to set (in ms)." optional:""`
}

// Run executes the StripGateReleaseCmd command, either retrieving the current gate release time of the strip or setting it based on the provided argument.
func (cmd *StripGateReleaseCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.Strip.Gate.Release(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get gate release time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d gate release time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Gate.SetRelease(strip.Index.Index, *cmd.Release); err != nil {
		return fmt.Errorf("failed to set gate release time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d gate release time set to: %.2f ms\n", strip.Index.Index, *cmd.Release)
	return nil
}

// StripEqCmdGroup defines the command group for controlling the EQ settings of a strip, including commands for getting and setting the EQ on/off state and parameters for each EQ band such as gain, frequency, Q factor, and type.
type StripEqCmdGroup struct {
	On   StripEqOnCmd `help:"Get or set the EQ on/off state of the strip."              cmd:""`
	Band struct {
		Band int                `arg:"" help:"The EQ band number."`
		Gain StripEqBandGainCmd `help:"Get or set the gain of the EQ band." cmd:""`
		Freq StripEqBandFreqCmd `help:"Get or set the frequency of the EQ band." cmd:""`
		Q    StripEqBandQCmd    `help:"Get or set the Q factor of the EQ band." cmd:""`
		Type StripEqBandTypeCmd `help:"Get or set the type of the EQ band." cmd:""`
	} `help:"Commands for controlling a specific EQ band of the strip."        arg:""`
}

// Validate checks if the provided EQ band number is valid (between 1 and 4) and returns an error if it is not.
func (cmd *StripEqCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 4 {
		return fmt.Errorf("EQ band number must be between 1 and 4")
	}
	return nil
}

// StripEqOnCmd defines the command for getting or setting the EQ on/off state of a strip, allowing users to enable or disable the EQ effect on the strip.
type StripEqOnCmd struct {
	Enable *string `arg:"" help:"Whether to enable or disable the EQ." optional:"" enum:"true,false"`
}

// Run executes the StripEqOnCmd command, either retrieving the current EQ on/off state of the strip or setting it based on the provided argument.
func (cmd *StripEqOnCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Strip.Eq.On(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get EQ state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d EQ state: %t\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Eq.SetOn(strip.Index.Index, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set EQ state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d EQ state set to: %s\n", strip.Index.Index, *cmd.Enable)
	return nil
}

// StripEqBandGainCmd defines the command for getting or setting the gain of a specific EQ band on a strip, allowing users to adjust the level of the signal for that band in decibels (dB).
type StripEqBandGainCmd struct {
	Gain *float64 `arg:"" help:"The gain to set for the EQ band (in dB)." optional:""`
}

// Run executes the StripEqBandGainCmd command, either retrieving the current gain of the specified EQ band on the strip or setting it based on the provided argument.
func (cmd *StripEqBandGainCmd) Run(ctx *context, strip *StripCmdGroup, stripEq *StripEqCmdGroup) error {
	if cmd.Gain == nil {
		resp, err := ctx.Client.Strip.Eq.Gain(strip.Index.Index, stripEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get EQ band gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d EQ band %d gain: %.2f\n", strip.Index.Index, stripEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Strip.Eq.SetGain(strip.Index.Index, stripEq.Band.Band, *cmd.Gain); err != nil {
		return fmt.Errorf("failed to set EQ band gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d EQ band %d gain set to: %.2f\n", strip.Index.Index, stripEq.Band.Band, *cmd.Gain)
	return nil
}

// StripEqBandFreqCmd defines the command for getting or setting the frequency of a specific EQ band on a strip, allowing users to adjust the center frequency of the band in hertz (Hz).
type StripEqBandFreqCmd struct {
	Freq *float64 `arg:"" help:"The frequency to set for the EQ band (in Hz)." optional:""`
}

// Run executes the StripEqBandFreqCmd command, either retrieving the current frequency of the specified EQ band on the strip or setting it based on the provided argument.
func (cmd *StripEqBandFreqCmd) Run(ctx *context, strip *StripCmdGroup, stripEq *StripEqCmdGroup) error {
	if cmd.Freq == nil {
		resp, err := ctx.Client.Strip.Eq.Frequency(strip.Index.Index, stripEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get EQ band frequency: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d EQ band %d frequency: %.2f Hz\n", strip.Index.Index, stripEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Strip.Eq.SetFrequency(strip.Index.Index, stripEq.Band.Band, *cmd.Freq); err != nil {
		return fmt.Errorf("failed to set EQ band frequency: %w", err)
	}
	fmt.Fprintf(
		ctx.Out,
		"Strip %d EQ band %d frequency set to: %.2f Hz\n",
		strip.Index.Index,
		stripEq.Band.Band,
		*cmd.Freq,
	)
	return nil
}

// StripEqBandQCmd defines the command for getting or setting the Q factor of a specific EQ band on a strip, allowing users to adjust the bandwidth of the band, which determines how wide or narrow the affected frequency range is.
type StripEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the EQ band." optional:""`
}

// Run executes the StripEqBandQCmd command, either retrieving the current Q factor of the specified EQ band on the strip or setting it based on the provided argument.
func (cmd *StripEqBandQCmd) Run(ctx *context, strip *StripCmdGroup, stripEq *StripEqCmdGroup) error {
	if cmd.Q == nil {
		resp, err := ctx.Client.Strip.Eq.Q(strip.Index.Index, stripEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get EQ band Q factor: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d EQ band %d Q factor: %.2f\n", strip.Index.Index, stripEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Strip.Eq.SetQ(strip.Index.Index, stripEq.Band.Band, *cmd.Q); err != nil {
		return fmt.Errorf("failed to set EQ band Q factor: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d EQ band %d Q factor set to: %.2f\n", strip.Index.Index, stripEq.Band.Band, *cmd.Q)
	return nil
}

// StripEqBandTypeCmd defines the command for getting or setting the type of a specific EQ band on a strip, allowing users to choose from different EQ types such as low cut (lcut), low shelf (lshv), parametric (peq), variable Q (veq), high shelf (hshv), or high cut (hcut).
type StripEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the EQ band." optional:"" enum:"lcut,lshv,peq,veq,hshv,hcut"`
}

// Run executes the StripEqBandTypeCmd command, either retrieving the current type of the specified EQ band on the strip or setting it based on the provided argument.
func (cmd *StripEqBandTypeCmd) Run(ctx *context, strip *StripCmdGroup, stripEq *StripEqCmdGroup) error {
	if cmd.Type == nil {
		resp, err := ctx.Client.Strip.Eq.Type(strip.Index.Index, stripEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get EQ band type: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d EQ band %d type: %s\n", strip.Index.Index, stripEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Strip.Eq.SetType(strip.Index.Index, stripEq.Band.Band, *cmd.Type); err != nil {
		return fmt.Errorf("failed to set EQ band type: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d EQ band %d type set to: %s\n", strip.Index.Index, stripEq.Band.Band, *cmd.Type)
	return nil
}

// StripCompCmdGroup defines the command group for controlling the compressor settings of a strip, including commands for getting and setting the compressor on/off state, mode, threshold, ratio, mix, makeup gain, attack time, hold time, and release time.
type StripCompCmdGroup struct {
	On        StripCompOnCmd        `help:"Get or set the compressor on/off state of the strip." cmd:""`
	Mode      StripCompModeCmd      `help:"Get or set the compressor mode of the strip."         cmd:""`
	Threshold StripCompThresholdCmd `help:"Get or set the compressor threshold of the strip."    cmd:""`
	Ratio     StripCompRatioCmd     `help:"Get or set the compressor ratio of the strip."        cmd:""`
	Mix       StripCompMixCmd       `help:"Get or set the compressor mix of the strip."          cmd:""`
	Makeup    StripCompMakeupCmd    `help:"Get or set the compressor makeup gain of the strip."  cmd:""`
	Attack    StripCompAttackCmd    `help:"Get or set the compressor attack time of the strip."  cmd:""`
	Hold      StripCompHoldCmd      `help:"Get or set the compressor hold time of the strip."    cmd:""`
	Release   StripCompReleaseCmd   `help:"Get or set the compressor release time of the strip." cmd:""`
}

// StripCompOnCmd defines the command for getting or setting the compressor on/off state of a strip, allowing users to enable or disable the compressor effect on the strip.
type StripCompOnCmd struct {
	Enable *string `arg:"" help:"Whether to enable or disable the compressor." optional:"" enum:"true,false"`
}

// Run executes the StripCompOnCmd command, either retrieving the current compressor on/off state of the strip or setting it based on the provided argument.
func (cmd *StripCompOnCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Strip.Comp.On(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor state: %t\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetOn(strip.Index.Index, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set compressor state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor state set to: %s\n", strip.Index.Index, *cmd.Enable)
	return nil
}

// StripCompModeCmd defines the command for getting or setting the compressor mode of a strip, allowing users to choose from different compressor modes such as comp (standard compression) or exp (expander).
type StripCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set." optional:"" enum:"comp,exp"`
}

// Run executes the StripCompModeCmd command, either retrieving the current compressor mode of the strip or setting it based on the provided argument.
func (cmd *StripCompModeCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Strip.Comp.Mode(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor mode: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor mode: %s\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetMode(strip.Index.Index, *cmd.Mode); err != nil {
		return fmt.Errorf("failed to set compressor mode: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor mode set to: %s\n", strip.Index.Index, *cmd.Mode)
	return nil
}

// StripCompThresholdCmd defines the command for getting or setting the compressor threshold of a strip, allowing users to specify the threshold level at which the compressor will start to reduce the signal level.
type StripCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set (in dB)." optional:""`
}

// Run executes the StripCompThresholdCmd command, either retrieving the current compressor threshold of the strip or setting it based on the provided argument.
func (cmd *StripCompThresholdCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.Strip.Comp.Threshold(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor threshold: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor threshold: %.2f\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetThreshold(strip.Index.Index, *cmd.Threshold); err != nil {
		return fmt.Errorf("failed to set compressor threshold: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor threshold set to: %.2f\n", strip.Index.Index, *cmd.Threshold)
	return nil
}

// StripCompRatioCmd defines the command for getting or setting the compressor ratio of a strip, allowing users to specify the amount of gain reduction applied by the compressor once the signal exceeds the threshold.
type StripCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set." optional:""`
}

// Run executes the StripCompRatioCmd command, either retrieving the current compressor ratio of the strip or setting it based on the provided argument.
func (cmd *StripCompRatioCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Ratio == nil {
		resp, err := ctx.Client.Strip.Comp.Ratio(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor ratio: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor ratio: %.2f\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetRatio(strip.Index.Index, *cmd.Ratio); err != nil {
		return fmt.Errorf("failed to set compressor ratio: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor ratio set to: %.2f\n", strip.Index.Index, *cmd.Ratio)
	return nil
}

// StripCompMixCmd defines the command for getting or setting the compressor mix of a strip, allowing users to specify the blend between the dry (unprocessed) signal and the wet (compressed) signal, typically expressed as a percentage.
type StripCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix to set (0-100%)." optional:""`
}

// Run executes the StripCompMixCmd command, either retrieving the current compressor mix of the strip or setting it based on the provided argument.
func (cmd *StripCompMixCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Mix == nil {
		resp, err := ctx.Client.Strip.Comp.Mix(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor mix: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor mix: %.2f%%\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetMix(strip.Index.Index, *cmd.Mix); err != nil {
		return fmt.Errorf("failed to set compressor mix: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor mix set to: %.2f%%\n", strip.Index.Index, *cmd.Mix)
	return nil
}

// StripCompMakeupCmd defines the command for getting or setting the compressor makeup gain of a strip, allowing users to specify the amount of gain applied to the signal after compression to compensate for any reduction in level caused by the compressor.
type StripCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set (in dB)." optional:""`
}

// Run executes the StripCompMakeupCmd command, either retrieving the current compressor makeup gain of the strip or setting it based on the provided argument.
func (cmd *StripCompMakeupCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Makeup == nil {
		resp, err := ctx.Client.Strip.Comp.Makeup(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor makeup gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor makeup gain: %.2f\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetMakeup(strip.Index.Index, *cmd.Makeup); err != nil {
		return fmt.Errorf("failed to set compressor makeup gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor makeup gain set to: %.2f\n", strip.Index.Index, *cmd.Makeup)
	return nil
}

// StripCompAttackCmd defines the command for getting or setting the compressor attack time of a strip, allowing users to specify the time it takes for the compressor to start reducing the signal level after the signal exceeds the threshold.
type StripCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set (in ms)." optional:""`
}

// Run executes the StripCompAttackCmd command, either retrieving the current compressor attack time of the strip or setting it based on the provided argument.
func (cmd *StripCompAttackCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.Strip.Comp.Attack(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor attack time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor attack time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetAttack(strip.Index.Index, *cmd.Attack); err != nil {
		return fmt.Errorf("failed to set compressor attack time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor attack time set to: %.2f ms\n", strip.Index.Index, *cmd.Attack)
	return nil
}

// StripCompHoldCmd defines the command for getting or setting the compressor hold time of a strip, allowing users to specify the time that the compressor continues to reduce the signal level after the signal falls below the threshold before it starts to release.
type StripCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set (in ms)." optional:""`
}

// Run executes the StripCompHoldCmd command, either retrieving the current compressor hold time of the strip or setting it based on the provided argument.
func (cmd *StripCompHoldCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.Strip.Comp.Hold(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor hold time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor hold time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetHold(strip.Index.Index, *cmd.Hold); err != nil {
		return fmt.Errorf("failed to set compressor hold time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor hold time set to: %.2f ms\n", strip.Index.Index, *cmd.Hold)
	return nil
}

// StripCompReleaseCmd defines the command for getting or setting the compressor release time of a strip, allowing users to specify the time it takes for the compressor to stop reducing the signal level after the signal falls below the threshold and the hold time has elapsed.
type StripCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set (in ms)." optional:""`
}

// Run executes the StripCompReleaseCmd command, either retrieving the current compressor release time of the strip or setting it based on the provided argument.
func (cmd *StripCompReleaseCmd) Run(ctx *context, strip *StripCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.Strip.Comp.Release(strip.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get compressor release time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Strip %d compressor release time: %.2f ms\n", strip.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.Strip.Comp.SetRelease(strip.Index.Index, *cmd.Release); err != nil {
		return fmt.Errorf("failed to set compressor release time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Strip %d compressor release time set to: %.2f ms\n", strip.Index.Index, *cmd.Release)
	return nil
}
