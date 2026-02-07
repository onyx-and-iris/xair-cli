package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

// MainCmdGroup defines the command group for controlling the Main L/R output, including commands for mute state, fader level, and fade-in/fade-out times.
type MainCmdGroup struct {
	Mute MainMuteCmd `help:"Get or set the mute state of the Main L/R output." cmd:""`

	Fader   MainFaderCmd   `help:"Get or set the fader level of the Main L/R output."      cmd:""`
	Fadein  MainFadeinCmd  `help:"Fade in the Main L/R output over a specified duration."  cmd:""`
	Fadeout MainFadeoutCmd `help:"Fade out the Main L/R output over a specified duration." cmd:""`

	Eq   MainEqCmdGroup   `help:"Commands for controlling the equalizer settings of the Main L/R output."  cmd:"eq"`
	Comp MainCompCmdGroup `help:"Commands for controlling the compressor settings of the Main L/R output." cmd:"comp"`
}

// MainMuteCmd defines the command for getting or setting the mute state of the Main L/R output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainMuteCmd struct {
	Mute *bool `arg:"" help:"The mute state to set. If not provided, the current state will be printed." optional:""`
}

// Run executes the MainMuteCmd command, either retrieving the current mute state of the Main L/R output or setting it based on the provided argument.
func (cmd *MainMuteCmd) Run(ctx *context) error {
	if cmd.Mute == nil {
		resp, err := ctx.Client.Main.Mute()
		if err != nil {
			return fmt.Errorf("failed to get Main L/R mute state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R mute state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Main.SetMute(*cmd.Mute); err != nil {
		return fmt.Errorf("failed to set Main L/R mute state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R mute state set to: %t\n", *cmd.Mute)
	return nil
}

// MainFaderCmd defines the command for getting or setting the fader level of the Main L/R output, allowing users to specify the desired level in dB.
type MainFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set. If not provided, the current level will be printed." optional:""`
}

// Run executes the MainFaderCmd command, either retrieving the current fader level of the Main L/R output or setting it based on the provided argument.
func (cmd *MainFaderCmd) Run(ctx *context) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Main.Fader()
		if err != nil {
			return fmt.Errorf("failed to get Main L/R fader level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R fader level: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.Main.SetFader(*cmd.Level); err != nil {
		return fmt.Errorf("failed to set Main L/R fader level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R fader level set to: %.2f\n", *cmd.Level)
	return nil
}

// MainFadeinCmd defines the command for getting or setting the fade-in time of the Main L/R output, allowing users to specify the desired duration for the fade-in effect.
type MainFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-in. If not provided, the current target level will be printed." default:"0.0" arg:""`
}

// Run executes the MainFadeinCmd command, either retrieving the current fade-in time of the Main L/R output or setting it based on the provided argument, with an optional target level for the fade-in effect.
func (cmd *MainFadeinCmd) Run(ctx *context) error {
	currentLevel, err := ctx.Client.Main.Fader()
	if err != nil {
		return fmt.Errorf("failed to get Main L/R fader level: %w", err)
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
		if err := ctx.Client.Main.SetFader(currentLevel); err != nil {
			return fmt.Errorf("failed to set Main L/R fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Main L/R fade-in completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MainFadeoutCmd defines the command for getting or setting the fade-out time of the Main L/R output, allowing users to specify the desired duration for the fade-out effect and an optional target level to fade out to.
type MainFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-out. If not provided, the current target level will be printed." default:"-90.0" arg:""`
}

// Run executes the MainFadeoutCmd command, either retrieving the current fade-out time of the Main L/R output or setting it based on the provided argument, with an optional target level for the fade-out effect.
func (cmd *MainFadeoutCmd) Run(ctx *context) error {
	currentLevel, err := ctx.Client.Main.Fader()
	if err != nil {
		return fmt.Errorf("failed to get Main L/R fader level: %w", err)
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
		if err := ctx.Client.Main.SetFader(currentLevel); err != nil {
			return fmt.Errorf("failed to set Main L/R fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Main L/R fade-out completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MainEqCmdGroup defines the command group for controlling the equalizer settings of the Main L/R output, including commands for getting or setting the EQ parameters.
type MainEqCmdGroup struct {
	On   MainEqOnCmd `help:"Get or set the EQ on/off state of the Main L/R output."               cmd:"on"`
	Band struct {
		Band int               `arg:"" help:"The EQ band number."`
		Gain MainEqBandGainCmd `help:"Get or set the gain of the specified EQ band." cmd:"gain"`
		Freq MainEqBandFreqCmd `help:"Get or set the frequency of the specified EQ band." cmd:"freq"`
		Q    MainEqBandQCmd    `help:"Get or set the Q factor of the specified EQ band." cmd:"q"`
		Type MainEqBandTypeCmd `help:"Get or set the type of the specified EQ band." cmd:"type"`
	} `help:"Commands for controlling individual EQ bands of the Main L/R output."          arg:""`
}

// Validate checks if the provided EQ band number is within the valid range (1-6) for the Main L/R output.
func (cmd *MainEqCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 6 {
		return fmt.Errorf("invalid EQ band number: %d. Valid range is 1-6", cmd.Band.Band)
	}
	return nil
}

// MainEqOnCmd defines the command for getting or setting the EQ on/off state of the Main L/R output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainEqOnCmd struct {
	Enable *string `arg:"" help:"The EQ on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MainEqOnCmd command, either retrieving the current EQ on/off state of the Main L/R output or setting it based on the provided argument.
func (cmd *MainEqOnCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Main.Eq.On(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R EQ on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R EQ on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Eq.SetOn(0, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Main L/R EQ on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R EQ on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MainEqBandGainCmd defines the command for getting or setting the gain of a specific EQ band on the Main L/R output, allowing users to specify the desired gain in dB.
type MainEqBandGainCmd struct {
	Level *float64 `arg:"" help:"The gain level to set for the specified EQ band. If not provided, the current gain will be printed." optional:""`
}

// Run executes the MainEqBandGainCmd command, either retrieving the current gain of a specific EQ band on the Main L/R output or setting it based on the provided argument.
func (cmd *MainEqBandGainCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainEqCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Main.Eq.Gain(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R EQ band %d gain: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R EQ band %d gain: %.2f dB\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Main.Eq.SetGain(0, mainEq.Band.Band, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set Main L/R EQ band %d gain: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R EQ band %d gain set to: %.2f dB\n", mainEq.Band.Band, *cmd.Level)
	return nil
}

// MainEqBandFreqCmd defines the command for getting or setting the frequency of a specific EQ band on the Main L/R output, allowing users to specify the desired frequency in Hz.
type MainEqBandFreqCmd struct {
	Frequency *float64 `arg:"" help:"The frequency to set for the specified EQ band. If not provided, the current frequency will be printed." optional:""`
}

// Run executes the MainEqBandFreqCmd command, either retrieving the current frequency of a specific EQ band on the Main L/R output or setting it based on the provided argument.
func (cmd *MainEqBandFreqCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainEqCmdGroup) error {
	if cmd.Frequency == nil {
		resp, err := ctx.Client.Main.Eq.Frequency(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R EQ band %d frequency: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R EQ band %d frequency: %.2f Hz\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Main.Eq.SetFrequency(0, mainEq.Band.Band, *cmd.Frequency); err != nil {
		return fmt.Errorf("failed to set Main L/R EQ band %d frequency: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R EQ band %d frequency set to: %.2f Hz\n", mainEq.Band.Band, *cmd.Frequency)
	return nil
}

// MainEqBandQCmd defines the command for getting or setting the Q factor of a specific EQ band on the Main L/R output, allowing users to specify the desired Q factor.
type MainEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the specified EQ band. If not provided, the current Q factor will be printed." optional:""`
}

// Run executes the MainEqBandQCmd command, either retrieving the current Q factor of a specific EQ band on the Main L/R output or setting it based on the provided argument.
func (cmd *MainEqBandQCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainEqCmdGroup) error {
	if cmd.Q == nil {
		resp, err := ctx.Client.Main.Eq.Q(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R EQ band %d Q factor: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R EQ band %d Q factor: %.2f\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Main.Eq.SetQ(0, mainEq.Band.Band, *cmd.Q); err != nil {
		return fmt.Errorf("failed to set Main L/R EQ band %d Q factor: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R EQ band %d Q factor set to: %.2f\n", mainEq.Band.Band, *cmd.Q)
	return nil
}

// MainEqBandTypeCmd defines the command for getting or setting the type of a specific EQ band on the Main L/R output, allowing users to specify the desired type as "peaking", "low_shelf", "high_shelf", "low_pass", or "high_pass".
type MainEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the specified EQ band. If not provided, the current type will be printed." optional:"" enum:"peaking,low_shelf,high_shelf,low_pass,high_pass"`
}

// Run executes the MainEqBandTypeCmd command, either retrieving the current type of a specific EQ band on the Main L/R output or setting it based on the provided argument.
func (cmd *MainEqBandTypeCmd) Run(ctx *context, main *MainCmdGroup, mainEq *MainEqCmdGroup) error {
	if cmd.Type == nil {
		resp, err := ctx.Client.Main.Eq.Type(0, mainEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R EQ band %d type: %w", mainEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R EQ band %d type: %s\n", mainEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Main.Eq.SetType(0, mainEq.Band.Band, *cmd.Type); err != nil {
		return fmt.Errorf("failed to set Main L/R EQ band %d type: %w", mainEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R EQ band %d type set to: %s\n", mainEq.Band.Band, *cmd.Type)
	return nil
}

// MainCompCmdGroup defines the command group for controlling the compressor settings of the Main L/R output, including commands for getting or setting the compressor parameters.
type MainCompCmdGroup struct {
	On        MainCompOnCmd        `help:"Get or set the compressor on/off state of the Main L/R output." cmd:"on"`
	Mode      MainCompModeCmd      `help:"Get or set the compressor mode of the Main L/R output."         cmd:"mode"`
	Threshold MainCompThresholdCmd `help:"Get or set the compressor threshold of the Main L/R output."    cmd:"threshold"`
	Ratio     MainCompRatioCmd     `help:"Get or set the compressor ratio of the Main L/R output."        cmd:"ratio"`
	Mix       MainCompMixCmd       `help:"Get or set the compressor mix level of the Main L/R output."    cmd:"mix"`
	Makeup    MainCompMakeupCmd    `help:"Get or set the compressor makeup gain of the Main L/R output."  cmd:"makeup"`
	Attack    MainCompAttackCmd    `help:"Get or set the compressor attack time of the Main L/R output."  cmd:"attack"`
	Hold      MainCompHoldCmd      `help:"Get or set the compressor hold time of the Main L/R output."    cmd:"hold"`
	Release   MainCompReleaseCmd   `help:"Get or set the compressor release time of the Main L/R output." cmd:"release"`
}

// MainCompOnCmd defines the command for getting or setting the compressor on/off state of the Main L/R output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainCompOnCmd struct {
	Enable *string `arg:"" help:"The compressor on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MainCompOnCmd command, either retrieving the current compressor on/off state of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompOnCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Main.Comp.On(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetOn(0, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MainCompModeCmd defines the command for getting or setting the compressor mode of the Main L/R output, allowing users to specify the desired mode as "comp" or "exp".
type MainCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set. If not provided, the current mode will be printed." optional:"" enum:"comp,exp"`
}

// Run executes the MainCompModeCmd command, either retrieving the current compressor mode of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompModeCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Main.Comp.Mode(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor mode: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor mode: %s\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetMode(0, *cmd.Mode); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor mode: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor mode set to: %s\n", *cmd.Mode)
	return nil
}

// MainCompThresholdCmd defines the command for getting or setting the compressor threshold of the Main L/R output, allowing users to specify the desired threshold in dB.
type MainCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set. If not provided, the current threshold will be printed." optional:""`
}

// Run executes the MainCompThresholdCmd command, either retrieving the current compressor threshold of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompThresholdCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.Main.Comp.Threshold(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor threshold: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor threshold: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetThreshold(0, *cmd.Threshold); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor threshold: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor threshold set to: %.2f dB\n", *cmd.Threshold)
	return nil
}

// MainCompRatioCmd defines the command for getting or setting the compressor ratio of the Main L/R output, allowing users to specify the desired ratio.
type MainCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set. If not provided, the current ratio will be printed." optional:""`
}

// Run executes the MainCompRatioCmd command, either retrieving the current compressor ratio of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompRatioCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Ratio == nil {
		resp, err := ctx.Client.Main.Comp.Ratio(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor ratio: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor ratio: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetRatio(0, *cmd.Ratio); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor ratio: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor ratio set to: %.2f\n", *cmd.Ratio)
	return nil
}

// MainCompMixCmd defines the command for getting or setting the compressor mix level of the Main L/R output, allowing users to specify the desired mix level in percentage.
type MainCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix level to set. If not provided, the current mix level will be printed." optional:""`
}

// Run executes the MainCompMixCmd command, either retrieving the current compressor mix level of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompMixCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Mix == nil {
		resp, err := ctx.Client.Main.Comp.Mix(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor mix level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor mix level: %.2f%%\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetMix(0, *cmd.Mix); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor mix level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor mix level set to: %.2f%%\n", *cmd.Mix)
	return nil
}

// MainCompMakeupCmd defines the command for getting or setting the compressor makeup gain of the Main L/R output, allowing users to specify the desired makeup gain in dB.
type MainCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set. If not provided, the current makeup gain will be printed." optional:""`
}

// Run executes the MainCompMakeupCmd command, either retrieving the current compressor makeup gain of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompMakeupCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Makeup == nil {
		resp, err := ctx.Client.Main.Comp.Makeup(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor makeup gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor makeup gain: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetMakeup(0, *cmd.Makeup); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor makeup gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor makeup gain set to: %.2f dB\n", *cmd.Makeup)
	return nil
}

// MainCompAttackCmd defines the command for getting or setting the compressor attack time of the Main L/R output, allowing users to specify the desired attack time in milliseconds.
type MainCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set. If not provided, the current attack time will be printed." optional:""`
}

// Run executes the MainCompAttackCmd command, either retrieving the current compressor attack time of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompAttackCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.Main.Comp.Attack(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor attack time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor attack time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetAttack(0, *cmd.Attack); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor attack time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor attack time set to: %.2f ms\n", *cmd.Attack)
	return nil
}

// MainCompHoldCmd defines the command for getting or setting the compressor hold time of the Main L/R output, allowing users to specify the desired hold time in milliseconds.
type MainCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set. If not provided, the current hold time will be printed." optional:""`
}

// Run executes the MainCompHoldCmd command, either retrieving the current compressor hold time of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompHoldCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.Main.Comp.Hold(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor hold time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor hold time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetHold(0, *cmd.Hold); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor hold time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor hold time set to: %.2f ms\n", *cmd.Hold)
	return nil
}

// MainCompReleaseCmd defines the command for getting or setting the compressor release time of the Main L/R output, allowing users to specify the desired release time in milliseconds.
type MainCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set. If not provided, the current release time will be printed." optional:""`
}

// Run executes the MainCompReleaseCmd command, either retrieving the current compressor release time of the Main L/R output or setting it based on the provided argument.
func (cmd *MainCompReleaseCmd) Run(ctx *context, main *MainCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.Main.Comp.Release(0)
		if err != nil {
			return fmt.Errorf("failed to get Main L/R compressor release time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Main L/R compressor release time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Main.Comp.SetRelease(0, *cmd.Release); err != nil {
		return fmt.Errorf("failed to set Main L/R compressor release time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Main L/R compressor release time set to: %.2f ms\n", *cmd.Release)
	return nil
}
