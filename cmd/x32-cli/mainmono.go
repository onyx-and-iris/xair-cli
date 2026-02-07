package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

// MainMonoCmdGroup defines the command group for controlling the Main Mono output, including commands for mute state, fader level, and fade-in/fade-out times.
type MainMonoCmdGroup struct {
	Mute MainMonoMuteCmd `help:"Get or set the mute state of the Main Mono output." cmd:""`

	Fader   MainMonoFaderCmd   `help:"Get or set the fader level of the Main Mono output."      cmd:""`
	Fadein  MainMonoFadeinCmd  `help:"Fade in the Main Mono output over a specified duration."  cmd:""`
	Fadeout MainMonoFadeoutCmd `help:"Fade out the Main Mono output over a specified duration." cmd:""`

	Eq   MainMonoEqCmdGroup   `help:"Commands for controlling the equalizer settings of the Main Mono output."  cmd:"eq"`
	Comp MainMonoCompCmdGroup `help:"Commands for controlling the compressor settings of the Main Mono output." cmd:"comp"`
}

// MainMonoMuteCmd defines the command for getting or setting the mute state of the Main Mono output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainMonoMuteCmd struct {
	Mute *string `arg:"" help:"The mute state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MainMonoMuteCmd command, either retrieving the current mute state of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoMuteCmd) Run(ctx *context) error {
	if cmd.Mute == nil {
		resp, err := ctx.Client.MainMono.Mute()
		if err != nil {
			return fmt.Errorf("failed to get Main Mono mute state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono mute state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.SetMute(*cmd.Mute == "true"); err != nil {
		return fmt.Errorf("failed to set Main Mono mute state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono mute state set to: %s\n", *cmd.Mute)
	return nil
}

// MainMonoFaderCmd defines the command for getting or setting the fader level of the Main Mono output, allowing users to specify the desired level in dB.
type MainMonoFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set. If not provided, the current level will be printed." optional:""`
}

// Run executes the MainMonoFaderCmd command, either retrieving the current fader level of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoFaderCmd) Run(ctx *context) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.MainMono.Fader()
		if err != nil {
			return fmt.Errorf("failed to get Main Mono fader level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono fader level: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.SetFader(*cmd.Level); err != nil {
		return fmt.Errorf("failed to set Main Mono fader level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono fader level set to: %.2f\n", *cmd.Level)
	return nil
}

// MainMonoFadeinCmd defines the command for getting or setting the fade-in time of the Main Mono output, allowing users to specify the desired duration for the fade-in effect.
type MainMonoFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-in. If not provided, the current target level will be printed." default:"0.0" arg:""`
}

// Run executes the MainMonoFadeinCmd command, either retrieving the current fade-in time of the Main Mono output or setting it based on the provided argument, with an optional target level for the fade-in effect.
func (cmd *MainMonoFadeinCmd) Run(ctx *context) error {
	currentLevel, err := ctx.Client.MainMono.Fader()
	if err != nil {
		return fmt.Errorf("failed to get Main Mono fader level: %w", err)
	}

	if currentLevel >= cmd.Target {
		return fmt.Errorf(
			"current fader level (%.2f) is already at or above the target level (%.2f)",
			currentLevel,
			cmd.Target,
		)
	}

	totalSteps := float64(cmd.Target - currentLevel)
	stepDuration := time.Duration(cmd.Duration.Seconds()*1000/totalSteps) * time.Millisecond
	for currentLevel < cmd.Target {
		currentLevel++
		if err := ctx.Client.MainMono.SetFader(currentLevel); err != nil {
			return fmt.Errorf("failed to set Main Mono fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Main Mono fade-in completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MainMonoFadeoutCmd defines the command for getting or setting the fade-out time of the Main Mono output, allowing users to specify the desired duration for the fade-out effect and an optional target level to fade out to.
type MainMonoFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-out. If not provided, the current target level will be printed." default:"-90.0" arg:""`
}

// Run executes the MainMonoFadeoutCmd command, either retrieving the current fade-out time of the Main Mono output or setting it based on the provided argument, with an optional target level for the fade-out effect.
func (cmd *MainMonoFadeoutCmd) Run(ctx *context) error {
	currentLevel, err := ctx.Client.MainMono.Fader()
	if err != nil {
		return fmt.Errorf("failed to get Main Mono fader level: %w", err)
	}

	if currentLevel <= cmd.Target {
		return fmt.Errorf(
			"current fader level (%.2f) is already at or below the target level (%.2f)",
			currentLevel,
			cmd.Target,
		)
	}

	totalSteps := float64(currentLevel - cmd.Target)
	stepDuration := time.Duration(cmd.Duration.Seconds()*1000/totalSteps) * time.Millisecond
	for currentLevel > cmd.Target {
		currentLevel--
		if err := ctx.Client.MainMono.SetFader(currentLevel); err != nil {
			return fmt.Errorf("failed to set Main Mono fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Main Mono fade-out completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MainMonoEqCmdGroup defines the command group for controlling the equalizer settings of the Main Mono output, including commands for getting or setting the EQ parameters.
type MainMonoEqCmdGroup struct {
	On   MainMonoEqOnCmd `help:"Get or set the EQ on/off state of the Main Mono output."               cmd:"on"`
	Band struct {
		Band int                   `arg:"" help:"The EQ band number."`
		Gain MainMonoEqBandGainCmd `help:"Get or set the gain of the specified EQ band." cmd:"gain"`
		Freq MainMonoEqBandFreqCmd `help:"Get or set the frequency of the specified EQ band." cmd:"freq"`
		Q    MainMonoEqBandQCmd    `help:"Get or set the Q factor of the specified EQ band." cmd:"q"`
		Type MainMonoEqBandTypeCmd `help:"Get or set the type of the specified EQ band." cmd:"type"`
	} `help:"Commands for controlling individual EQ bands of the Main Mono output."          arg:""`
}

// Validate checks if the provided EQ band number is within the valid range (1-6) for the Main Mono output.
func (cmd *MainMonoEqCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 6 {
		return fmt.Errorf("invalid EQ band number: %d. Valid range is 1-6", cmd.Band.Band)
	}
	return nil
}

// MainMonoEqOnCmd defines the command for getting or setting the EQ on/off state of the Main Mono output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainMonoEqOnCmd struct {
	Enable *string `arg:"" help:"The EQ on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MainMonoEqOnCmd command, either retrieving the current EQ on/off state of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoEqOnCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.MainMono.Eq.On(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono EQ on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono EQ on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Eq.SetOn(0, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Main Mono EQ on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono EQ on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MainMonoEqBandGainCmd defines the command for getting or setting the gain of a specific EQ band on the Main Mono output, allowing users to specify the desired gain in dB.
type MainMonoEqBandGainCmd struct {
	Level *float64 `arg:"" help:"The gain level to set for the specified EQ band. If not provided, the current gain will be printed." optional:""`
}

// Run executes the MainMonoEqBandGainCmd command, either retrieving the current gain of a specific EQ band on the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoEqBandGainCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainMonoEqCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.MainMono.Eq.Gain(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono EQ band %d gain: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono EQ band %d gain: %.2f dB\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.MainMono.Eq.SetGain(0, mainEq.Band.Band, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set Main Mono EQ band %d gain: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono EQ band %d gain set to: %.2f dB\n", mainEq.Band.Band, *cmd.Level)
	return nil
}

// MainMonoEqBandFreqCmd defines the command for getting or setting the frequency of a specific EQ band on the Main Mono output, allowing users to specify the desired frequency in Hz.
type MainMonoEqBandFreqCmd struct {
	Frequency *float64 `arg:"" help:"The frequency to set for the specified EQ band. If not provided, the current frequency will be printed." optional:""`
}

// Run executes the MainMonoEqBandFreqCmd command, either retrieving the current frequency of a specific EQ band on the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoEqBandFreqCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainMonoEqCmdGroup) error {
	if cmd.Frequency == nil {
		resp, err := ctx.Client.MainMono.Eq.Frequency(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono EQ band %d frequency: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono EQ band %d frequency: %.2f Hz\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.MainMono.Eq.SetFrequency(0, mainEq.Band.Band, *cmd.Frequency); err != nil {
		return fmt.Errorf("failed to set Main Mono EQ band %d frequency: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono EQ band %d frequency set to: %.2f Hz\n", mainEq.Band.Band, *cmd.Frequency)
	return nil
}

// MainMonoEqBandQCmd defines the command for getting or setting the Q factor of a specific EQ band on the Main Mono output, allowing users to specify the desired Q factor.
type MainMonoEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the specified EQ band. If not provided, the current Q factor will be printed." optional:""`
}

// Run executes the MainMonoEqBandQCmd command, either retrieving the current Q factor of a specific EQ band on the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoEqBandQCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainMonoEqCmdGroup) error {
	if cmd.Q == nil {
		resp, err := ctx.Client.MainMono.Eq.Q(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono EQ band %d Q factor: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono EQ band %d Q factor: %.2f\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.MainMono.Eq.SetQ(0, mainEq.Band.Band, *cmd.Q); err != nil {
		return fmt.Errorf("failed to set Main Mono EQ band %d Q factor: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono EQ band %d Q factor set to: %.2f\n", mainEq.Band.Band, *cmd.Q)
	return nil
}

// MainMonoEqBandTypeCmd defines the command for getting or setting the type of a specific EQ band on the Main Mono output, allowing users to specify the desired type as "peaking", "low_shelf", "high_shelf", "low_pass", or "high_pass".
type MainMonoEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the specified EQ band. If not provided, the current type will be printed." optional:"" enum:"peaking,low_shelf,high_shelf,low_pass,high_pass"`
}

// Run executes the MainMonoEqBandTypeCmd command, either retrieving the current type of a specific EQ band on the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoEqBandTypeCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainMonoEqCmdGroup) error {
	if cmd.Type == nil {
		resp, err := ctx.Client.MainMono.Eq.Type(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono EQ band %d type: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono EQ band %d type: %s\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.MainMono.Eq.SetType(0, mainEq.Band.Band, *cmd.Type); err != nil {
		return fmt.Errorf("failed to set Main Mono EQ band %d type: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono EQ band %d type set to: %s\n", mainEq.Band.Band, *cmd.Type)
	return nil
}

// MainMonoCompCmdGroup defines the command group for controlling the compressor settings of the Main Mono output, including commands for getting or setting the compressor parameters.
type MainMonoCompCmdGroup struct {
	On        MainMonoCompOnCmd        `help:"Get or set the compressor on/off state of the Main Mono output." cmd:"on"`
	Mode      MainMonoCompModeCmd      `help:"Get or set the compressor mode of the Main Mono output."         cmd:"mode"`
	Threshold MainMonoCompThresholdCmd `help:"Get or set the compressor threshold of the Main Mono output."    cmd:"threshold"`
	Ratio     MainMonoCompRatioCmd     `help:"Get or set the compressor ratio of the Main Mono output."        cmd:"ratio"`
	Mix       MainMonoCompMixCmd       `help:"Get or set the compressor mix level of the Main Mono output."    cmd:"mix"`
	Makeup    MainMonoCompMakeupCmd    `help:"Get or set the compressor makeup gain of the Main Mono output."  cmd:"makeup"`
	Attack    MainMonoCompAttackCmd    `help:"Get or set the compressor attack time of the Main Mono output."  cmd:"attack"`
	Hold      MainMonoCompHoldCmd      `help:"Get or set the compressor hold time of the Main Mono output."    cmd:"hold"`
	Release   MainMonoCompReleaseCmd   `help:"Get or set the compressor release time of the Main Mono output." cmd:"release"`
}

// MainMonoCompOnCmd defines the command for getting or setting the compressor on/off state of the Main Mono output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainMonoCompOnCmd struct {
	Enable *string `arg:"" help:"The compressor on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MainMonoCompOnCmd command, either retrieving the current compressor on/off state of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompOnCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.MainMono.Comp.On(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetOn(0, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MainMonoCompModeCmd defines the command for getting or setting the compressor mode of the Main Mono output, allowing users to specify the desired mode as "comp" or "exp".
type MainMonoCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set. If not provided, the current mode will be printed." optional:"" enum:"comp,exp"`
}

// Run executes the MainMonoCompModeCmd command, either retrieving the current compressor mode of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompModeCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.MainMono.Comp.Mode(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor mode: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor mode: %s\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetMode(0, *cmd.Mode); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor mode: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor mode set to: %s\n", *cmd.Mode)
	return nil
}

// MainMonoCompThresholdCmd defines the command for getting or setting the compressor threshold of the Main Mono output, allowing users to specify the desired threshold in dB.
type MainMonoCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set. If not provided, the current threshold will be printed." optional:""`
}

// Run executes the MainMonoCompThresholdCmd command, either retrieving the current compressor threshold of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompThresholdCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.MainMono.Comp.Threshold(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor threshold: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor threshold: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetThreshold(0, *cmd.Threshold); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor threshold: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor threshold set to: %.2f dB\n", *cmd.Threshold)
	return nil
}

// MainMonoCompRatioCmd defines the command for getting or setting the compressor ratio of the Main Mono output, allowing users to specify the desired ratio.
type MainMonoCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set. If not provided, the current ratio will be printed." optional:""`
}

// Run executes the MainMonoCompRatioCmd command, either retrieving the current compressor ratio of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompRatioCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Ratio == nil {
		resp, err := ctx.Client.MainMono.Comp.Ratio(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor ratio: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor ratio: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetRatio(0, *cmd.Ratio); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor ratio: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor ratio set to: %.2f\n", *cmd.Ratio)
	return nil
}

// MainMonoCompMixCmd defines the command for getting or setting the compressor mix level of the Main Mono output, allowing users to specify the desired mix level in percentage.
type MainMonoCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix level to set. If not provided, the current mix level will be printed." optional:""`
}

// Run executes the MainMonoCompMixCmd command, either retrieving the current compressor mix level of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompMixCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Mix == nil {
		resp, err := ctx.Client.MainMono.Comp.Mix(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor mix level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor mix level: %.2f%%\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetMix(0, *cmd.Mix); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor mix level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor mix level set to: %.2f%%\n", *cmd.Mix)
	return nil
}

// MainMonoCompMakeupCmd defines the command for getting or setting the compressor makeup gain of the Main Mono output, allowing users to specify the desired makeup gain in dB.
type MainMonoCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set. If not provided, the current makeup gain will be printed." optional:""`
}

// Run executes the MainMonoCompMakeupCmd command, either retrieving the current compressor makeup gain of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompMakeupCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Makeup == nil {
		resp, err := ctx.Client.MainMono.Comp.Makeup(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor makeup gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor makeup gain: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetMakeup(0, *cmd.Makeup); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor makeup gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor makeup gain set to: %.2f dB\n", *cmd.Makeup)
	return nil
}

// MainMonoCompAttackCmd defines the command for getting or setting the compressor attack time of the Main Mono output, allowing users to specify the desired attack time in milliseconds.
type MainMonoCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set. If not provided, the current attack time will be printed." optional:""`
}

// Run executes the MainMonoCompAttackCmd command, either retrieving the current compressor attack time of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompAttackCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.MainMono.Comp.Attack(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor attack time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor attack time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetAttack(0, *cmd.Attack); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor attack time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor attack time set to: %.2f ms\n", *cmd.Attack)
	return nil
}

// MainMonoCompHoldCmd defines the command for getting or setting the compressor hold time of the Main Mono output, allowing users to specify the desired hold time in milliseconds.
type MainMonoCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set. If not provided, the current hold time will be printed." optional:""`
}

// Run executes the MainMonoCompHoldCmd command, either retrieving the current compressor hold time of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompHoldCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.MainMono.Comp.Hold(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor hold time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor hold time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetHold(0, *cmd.Hold); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor hold time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor hold time set to: %.2f ms\n", *cmd.Hold)
	return nil
}

// MainMonoCompReleaseCmd defines the command for getting or setting the compressor release time of the Main Mono output, allowing users to specify the desired release time in milliseconds.
type MainMonoCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set. If not provided, the current release time will be printed." optional:""`
}

// Run executes the MainMonoCompReleaseCmd command, either retrieving the current compressor release time of the Main Mono output or setting it based on the provided argument.
func (cmd *MainMonoCompReleaseCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.MainMono.Comp.Release(0)
		if err != nil {
			return fmt.Errorf("failed to get Main Mono compressor release time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main Mono compressor release time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.MainMono.Comp.SetRelease(0, *cmd.Release); err != nil {
		return fmt.Errorf("failed to set Main Mono compressor release time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main Mono compressor release time set to: %.2f ms\n", *cmd.Release)
	return nil
}
