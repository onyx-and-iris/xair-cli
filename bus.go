package main

import (
	"fmt"
	"time"
)

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
	} `arg:"" help:"The index of the bus."`
}

type BusMuteCmd struct {
	State *string `arg:"" help:"The mute state to set (true or false). If not provided, the current mute state will be returned." optional:"" enum:"true,false"`
}

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

type BusFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set (in dB). If not provided, the current fader level will be returned."`
}

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

type BusFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in effect." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."     default:"0.0" arg:""`
}

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

type BusFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out effect." default:"5s"`
	Target   float64       `        help:"The target fader level (in dB)."      default:"-90.0" arg:""`
}

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

type BusNameCmd struct {
	Name *string `arg:"" help:"The name to set for the bus. If not provided, the current name will be returned."`
}

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

type BusEqCmdGroup struct {
	On   BusEqOnCmd   `help:"Get or set the EQ on/off state of the bus."                 cmd:"on"`
	Mode BusEqModeCmd `help:"Get or set the EQ mode of the bus (graphic or parametric)." cmd:"mode"`
	Band struct {
		Band int              `arg:"" help:"The EQ band number."`
		Gain BusEqBandGainCmd `help:"Get or set the gain of the EQ band." cmd:"gain"`
		Freq BusEqBandFreqCmd `help:"Get or set the frequency of the EQ band." cmd:"freq"`
		Q    BusEqBandQCmd    `help:"Get or set the Q factor of the EQ band." cmd:"q"`
		Type BusEqBandTypeCmd `help:"Get or set the type of the EQ band (bell, high shelf, low shelf, high pass, low pass)." cmd:"type"`
	} `help:"Commands for controlling a specific EQ band of the bus."               arg:""`
}

func (cmd *BusEqCmdGroup) Validate() error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 6 {
		return fmt.Errorf("EQ band number must be between 1 and 6")
	}
	return nil
}

type BusEqOnCmd struct {
	State *string `arg:"" help:"The EQ on/off state to set (true or false). If not provided, the current EQ state will be returned." optional:"" enum:"true,false"`
}

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

type BusEqModeCmd struct {
	Mode *string `arg:"" help:"The EQ mode to set (graphic or parametric). If not provided, the current EQ mode will be returned." optional:"" enum:"peq,geq,teq"`
}

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

type BusEqBandGainCmd struct {
	Gain *float64 `arg:"" help:"The gain to set for the EQ band (in dB). If not provided, the current gain will be returned."`
}

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

type BusEqBandFreqCmd struct {
	Freq *float64 `arg:"" help:"The frequency to set for the EQ band (in Hz). If not provided, the current frequency will be returned."`
}

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

type BusEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the EQ band. If not provided, the current Q factor will be returned."`
}

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

type BusEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the EQ band (bell, high shelf, low shelf, high pass, low pass). If not provided, the current type will be returned." optional:"" enum:"lcut,lshv,peq,veq,hshv,hcut"`
}

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

type BusCompCmdGroup struct {
	On        BusCompOnCmd        `help:"Get or set the compressor on/off state of the bus."                        cmd:"on"`
	Mode      BusCompModeCmd      `help:"Get or set the compressor mode of the bus (standard, vintage, or modern)." cmd:"mode"`
	Threshold BusCompThresholdCmd `help:"Get or set the compressor threshold of the bus (in dB)."                   cmd:"threshold"`
	Ratio     BusCompRatioCmd     `help:"Get or set the compressor ratio of the bus."                               cmd:"ratio"`
	Mix       BusCompMixCmd       `help:"Get or set the compressor mix level of the bus (in %)."                    cmd:"mix"`
	Makeup    BusCompMakeupCmd    `help:"Get or set the compressor makeup gain of the bus (in dB)."                 cmd:"makeup"`
	Attack    BusCompAttackCmd    `help:"Get or set the compressor attack time of the bus (in ms)."                 cmd:"attack"`
	Hold      BusCompHoldCmd      `help:"Get or set the compressor hold time of the bus (in ms)."                   cmd:"hold"`
	Release   BusCompReleaseCmd   `help:"Get or set the compressor release time of the bus (in ms)."                cmd:"release"`
}

type BusCompOnCmd struct {
	State *string `arg:"" help:"The compressor on/off state to set (true or false). If not provided, the current compressor state will be returned." optional:"" enum:"true,false"`
}

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

type BusCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set (standard, vintage, or modern). If not provided, the current compressor mode will be returned." optional:"" enum:"comp,exp"`
}

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

type BusCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set (in dB). If not provided, the current compressor threshold will be returned."`
}

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

type BusCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set. If not provided, the current compressor ratio will be returned."`
}

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

type BusCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix level to set (in %). If not provided, the current compressor mix level will be returned."`
}

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

type BusCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set (in dB). If not provided, the current compressor makeup gain will be returned."`
}

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

type BusCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set (in ms). If not provided, the current compressor attack time will be returned."`
}

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

type BusCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set (in ms). If not provided, the current compressor hold time will be returned."`
}

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

type BusCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set (in ms). If not provided, the current compressor release time will be returned."`
}

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
