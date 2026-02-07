package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
)

// MatrixCmdGroup defines the command group for controlling the Matrix outputs, including commands for mute state, fader level, and fade-in/fade-out times.
type MatrixCmdGroup struct {
	Index struct {
		Index int           `arg:"" help:"The index of the Matrix output (1-6)."`
		Mute  MatrixMuteCmd `help:"Get or set the mute state of the Matrix output." cmd:""`

		Fader   MatrixFaderCmd   `help:"Get or set the fader level of the Matrix output."      cmd:""`
		Fadein  MatrixFadeinCmd  `help:"Fade in the Matrix output over a specified duration."  cmd:""`
		Fadeout MatrixFadeoutCmd `help:"Fade out the Matrix output over a specified duration." cmd:""`

		Eq   MatrixEqCmdGroup   `help:"Commands for controlling the equalizer settings of the Matrix output."  cmd:"eq"`
		Comp MatrixCompCmdGroup `help:"Commands for controlling the compressor settings of the Matrix output." cmd:"comp"`
	} `help:"Commands for controlling individual Matrix outputs." arg:""`
}

func (cmd *MatrixCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Index.Index < 1 || cmd.Index.Index > 6 {
		return fmt.Errorf("invalid Matrix output index: %d. Valid range is 1-6", cmd.Index.Index)
	}
	return nil
}

// MatrixMuteCmd defines the command for getting or setting the mute state of the Matrix output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MatrixMuteCmd struct {
	Mute *string `arg:"" help:"The mute state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MatrixMuteCmd command, either retrieving the current mute state of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixMuteCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Mute == nil {
		resp, err := ctx.Client.Matrix.Mute(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix mute state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix mute state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.SetMute(matrix.Index.Index, *cmd.Mute == "true"); err != nil {
		return fmt.Errorf("failed to set Matrix mute state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix mute state set to: %s\n", *cmd.Mute)
	return nil
}

// MatrixFaderCmd defines the command for getting or setting the fader level of the Matrix output, allowing users to specify the desired level in dB.
type MatrixFaderCmd struct {
	Level *float64 `arg:"" help:"The fader level to set. If not provided, the current level will be printed." optional:""`
}

// Run executes the MatrixFaderCmd command, either retrieving the current fader level of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixFaderCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Matrix.Fader(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix fader level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix fader level: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.SetFader(matrix.Index.Index, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set Matrix fader level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix fader level set to: %.2f\n", *cmd.Level)
	return nil
}

// MatrixFadeinCmd defines the command for getting or setting the fade-in time of the Matrix output, allowing users to specify the desired duration for the fade-in effect.
type MatrixFadeinCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-in. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-in. If not provided, the current target level will be printed." default:"0.0" arg:""`
}

// Run executes the MatrixFadeinCmd command, either retrieving the current fade-in time of the Matrix output or setting it based on the provided argument, with an optional target level for the fade-in effect.
func (cmd *MatrixFadeinCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	currentLevel, err := ctx.Client.Matrix.Fader(matrix.Index.Index)
	if err != nil {
		return fmt.Errorf("failed to get Matrix fader level: %w", err)
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
		if err := ctx.Client.Matrix.SetFader(matrix.Index.Index, currentLevel); err != nil {
			return fmt.Errorf("failed to set Matrix fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Matrix fade-in completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MatrixFadeoutCmd defines the command for getting or setting the fade-out time of the Matrix output, allowing users to specify the desired duration for the fade-out effect and an optional target level to fade out to.
type MatrixFadeoutCmd struct {
	Duration time.Duration `flag:"" help:"The duration of the fade-out. (in seconds.)"                                                   default:"5s"`
	Target   float64       `        help:"The target level for the fade-out. If not provided, the current target level will be printed." default:"-90.0" arg:""`
}

// Run executes the MatrixFadeoutCmd command, either retrieving the current fade-out time of the Matrix output or setting it based on the provided argument, with an optional target level for the fade-out effect.
func (cmd *MatrixFadeoutCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	currentLevel, err := ctx.Client.Matrix.Fader(matrix.Index.Index)
	if err != nil {
		return fmt.Errorf("failed to get Matrix fader level: %w", err)
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
		if err := ctx.Client.Matrix.SetFader(matrix.Index.Index, currentLevel); err != nil {
			return fmt.Errorf("failed to set Matrix fader level: %w", err)
		}
		time.Sleep(stepDuration)
	}
	fmt.Fprintf(ctx.Out, "Matrix fade-out completed. Final level: %.2f\n", currentLevel)
	return nil
}

// MatrixEqCmdGroup defines the command group for controlling the equalizer settings of the Matrix output, including commands for getting or setting the EQ parameters.
type MatrixEqCmdGroup struct {
	On   MatrixEqOnCmd `help:"Get or set the EQ on/off state of the Matrix output."               cmd:"on"`
	Band struct {
		Band int                 `arg:"" help:"The EQ band number."`
		Gain MatrixEqBandGainCmd `help:"Get or set the gain of the specified EQ band." cmd:"gain"`
		Freq MatrixEqBandFreqCmd `help:"Get or set the frequency of the specified EQ band." cmd:"freq"`
		Q    MatrixEqBandQCmd    `help:"Get or set the Q factor of the specified EQ band." cmd:"q"`
		Type MatrixEqBandTypeCmd `help:"Get or set the type of the specified EQ band." cmd:"type"`
	} `help:"Commands for controlling individual EQ bands of the Matrix output."          arg:""`
}

// Validate checks if the provided EQ band number is within the valid range (1-6) for the Matrix output.
func (cmd *MatrixEqCmdGroup) Validate(ctx kong.Context) error {
	if cmd.Band.Band < 1 || cmd.Band.Band > 6 {
		return fmt.Errorf("invalid EQ band number: %d. Valid range is 1-6", cmd.Band.Band)
	}
	return nil
}

// MatrixEqOnCmd defines the command for getting or setting the EQ on/off state of the Matrix output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MatrixEqOnCmd struct {
	Enable *string `arg:"" help:"The EQ on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MatrixEqOnCmd command, either retrieving the current EQ on/off state of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixEqOnCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Matrix.Eq.On(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix EQ on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix EQ on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Eq.SetOn(matrix.Index.Index, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Matrix EQ on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix EQ on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MatrixEqBandGainCmd defines the command for getting or setting the gain of a specific EQ band on the Matrix output, allowing users to specify the desired gain in dB.
type MatrixEqBandGainCmd struct {
	Level *float64 `arg:"" help:"The gain level to set for the specified EQ band. If not provided, the current gain will be printed." optional:""`
}

// Run executes the MatrixEqBandGainCmd command, either retrieving the current gain of a specific EQ band on the Matrix output or setting it based on the provided argument.
func (cmd *MatrixEqBandGainCmd) Run(ctx *context, matrix *MatrixCmdGroup, matrixEq *MatrixEqCmdGroup) error {
	if cmd.Level == nil {
		resp, err := ctx.Client.Matrix.Eq.Gain(matrix.Index.Index, matrixEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Matrix EQ band %d gain: %w", matrixEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Matrix EQ band %d gain: %.2f dB\n", matrixEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Matrix.Eq.SetGain(matrix.Index.Index, matrixEq.Band.Band, *cmd.Level); err != nil {
		return fmt.Errorf("failed to set Matrix EQ band %d gain: %w", matrixEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Matrix EQ band %d gain set to: %.2f dB\n", matrixEq.Band.Band, *cmd.Level)
	return nil
}

// MatrixEqBandFreqCmd defines the command for getting or setting the frequency of a specific EQ band on the Matrix output, allowing users to specify the desired frequency in Hz.
type MatrixEqBandFreqCmd struct {
	Frequency *float64 `arg:"" help:"The frequency to set for the specified EQ band. If not provided, the current frequency will be printed." optional:""`
}

// Run executes the MatrixEqBandFreqCmd command, either retrieving the current frequency of a specific EQ band on the Matrix output or setting it based on the provided argument.
func (cmd *MatrixEqBandFreqCmd) Run(ctx *context, matrix *MatrixCmdGroup, matrixEq *MatrixEqCmdGroup) error {
	if cmd.Frequency == nil {
		resp, err := ctx.Client.Matrix.Eq.Frequency(matrix.Index.Index, matrixEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Matrix EQ band %d frequency: %w", matrixEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Matrix EQ band %d frequency: %.2f Hz\n", matrixEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Matrix.Eq.SetFrequency(matrix.Index.Index, matrixEq.Band.Band, *cmd.Frequency); err != nil {
		return fmt.Errorf("failed to set Matrix EQ band %d frequency: %w", matrixEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Matrix EQ band %d frequency set to: %.2f Hz\n", matrixEq.Band.Band, *cmd.Frequency)
	return nil
}

// MatrixEqBandQCmd defines the command for getting or setting the Q factor of a specific EQ band on the Matrix output, allowing users to specify the desired Q factor.
type MatrixEqBandQCmd struct {
	Q *float64 `arg:"" help:"The Q factor to set for the specified EQ band. If not provided, the current Q factor will be printed." optional:""`
}

// Run executes the MatrixEqBandQCmd command, either retrieving the current Q factor of a specific EQ band on the Matrix output or setting it based on the provided argument.
func (cmd *MatrixEqBandQCmd) Run(ctx *context, matrix *MatrixCmdGroup, matrixEq *MatrixEqCmdGroup) error {
	if cmd.Q == nil {
		resp, err := ctx.Client.Matrix.Eq.Q(matrix.Index.Index, matrixEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Matrix EQ band %d Q factor: %w", matrixEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Matrix EQ band %d Q factor: %.2f\n", matrixEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Matrix.Eq.SetQ(matrix.Index.Index, matrixEq.Band.Band, *cmd.Q); err != nil {
		return fmt.Errorf("failed to set Matrix EQ band %d Q factor: %w", matrixEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Matrix EQ band %d Q factor set to: %.2f\n", matrixEq.Band.Band, *cmd.Q)
	return nil
}

// MatrixEqBandTypeCmd defines the command for getting or setting the type of a specific EQ band on the Matrix output, allowing users to specify the desired type as "peaking", "low_shelf", "high_shelf", "low_pass", or "high_pass".
type MatrixEqBandTypeCmd struct {
	Type *string `arg:"" help:"The type to set for the specified EQ band. If not provided, the current type will be printed." optional:"" enum:"peaking,low_shelf,high_shelf,low_pass,high_pass"`
}

// Run executes the MatrixEqBandTypeCmd command, either retrieving the current type of a specific EQ band on the Matrix output or setting it based on the provided argument.
func (cmd *MatrixEqBandTypeCmd) Run(ctx *context, matrix *MatrixCmdGroup, matrixEq *MatrixEqCmdGroup) error {
	if cmd.Type == nil {
		resp, err := ctx.Client.Matrix.Eq.Type(matrix.Index.Index, matrixEq.Band.Band)
		if err != nil {
			return fmt.Errorf("failed to get Matrix EQ band %d type: %w", matrixEq.Band.Band, err)
		}
		fmt.Fprintf(ctx.Out, "Matrix EQ band %d type: %s\n", matrixEq.Band.Band, resp)
		return nil
	}

	if err := ctx.Client.Matrix.Eq.SetType(matrix.Index.Index, matrixEq.Band.Band, *cmd.Type); err != nil {
		return fmt.Errorf("failed to set Matrix EQ band %d type: %w", matrixEq.Band.Band, err)
	}
	fmt.Fprintf(ctx.Out, "Matrix EQ band %d type set to: %s\n", matrixEq.Band.Band, *cmd.Type)
	return nil
}

// MatrixCompCmdGroup defines the command group for controlling the compressor settings of the Matrix output, including commands for getting or setting the compressor parameters.
type MatrixCompCmdGroup struct {
	On        MatrixCompOnCmd        `help:"Get or set the compressor on/off state of the Matrix output." cmd:"on"`
	Mode      MatrixCompModeCmd      `help:"Get or set the compressor mode of the Matrix output."         cmd:"mode"`
	Threshold MatrixCompThresholdCmd `help:"Get or set the compressor threshold of the Matrix output."    cmd:"threshold"`
	Ratio     MatrixCompRatioCmd     `help:"Get or set the compressor ratio of the Matrix output."        cmd:"ratio"`
	Mix       MatrixCompMixCmd       `help:"Get or set the compressor mix level of the Matrix output."    cmd:"mix"`
	Makeup    MatrixCompMakeupCmd    `help:"Get or set the compressor makeup gain of the Matrix output."  cmd:"makeup"`
	Attack    MatrixCompAttackCmd    `help:"Get or set the compressor attack time of the Matrix output."  cmd:"attack"`
	Hold      MatrixCompHoldCmd      `help:"Get or set the compressor hold time of the Matrix output."    cmd:"hold"`
	Release   MatrixCompReleaseCmd   `help:"Get or set the compressor release time of the Matrix output." cmd:"release"`
}

// MatrixCompOnCmd defines the command for getting or setting the compressor on/off state of the Matrix output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MatrixCompOnCmd struct {
	Enable *string `arg:"" help:"The compressor on/off state to set. If not provided, the current state will be printed." optional:"" enum:"true,false"`
}

// Run executes the MatrixCompOnCmd command, either retrieving the current compressor on/off state of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompOnCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Enable == nil {
		resp, err := ctx.Client.Matrix.Comp.On(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor on/off state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor on/off state: %t\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetOn(matrix.Index.Index, *cmd.Enable == "true"); err != nil {
		return fmt.Errorf("failed to set Matrix compressor on/off state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor on/off state set to: %t\n", *cmd.Enable == "true")
	return nil
}

// MatrixCompModeCmd defines the command for getting or setting the compressor mode of the Matrix output, allowing users to specify the desired mode as "comp" or "exp".
type MatrixCompModeCmd struct {
	Mode *string `arg:"" help:"The compressor mode to set. If not provided, the current mode will be printed." optional:"" enum:"comp,exp"`
}

// Run executes the MatrixCompModeCmd command, either retrieving the current compressor mode of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompModeCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Mode == nil {
		resp, err := ctx.Client.Matrix.Comp.Mode(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor mode: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor mode: %s\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetMode(matrix.Index.Index, *cmd.Mode); err != nil {
		return fmt.Errorf("failed to set Matrix compressor mode: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor mode set to: %s\n", *cmd.Mode)
	return nil
}

// MatrixCompThresholdCmd defines the command for getting or setting the compressor threshold of the Matrix output, allowing users to specify the desired threshold in dB.
type MatrixCompThresholdCmd struct {
	Threshold *float64 `arg:"" help:"The compressor threshold to set. If not provided, the current threshold will be printed." optional:""`
}

// Run executes the MatrixCompThresholdCmd command, either retrieving the current compressor threshold of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompThresholdCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Threshold == nil {
		resp, err := ctx.Client.Matrix.Comp.Threshold(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor threshold: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor threshold: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetThreshold(matrix.Index.Index, *cmd.Threshold); err != nil {
		return fmt.Errorf("failed to set Matrix compressor threshold: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor threshold set to: %.2f dB\n", *cmd.Threshold)
	return nil
}

// MatrixCompRatioCmd defines the command for getting or setting the compressor ratio of the Matrix output, allowing users to specify the desired ratio.
type MatrixCompRatioCmd struct {
	Ratio *float64 `arg:"" help:"The compressor ratio to set. If not provided, the current ratio will be printed." optional:""`
}

// Run executes the MatrixCompRatioCmd command, either retrieving the current compressor ratio of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompRatioCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Ratio == nil {
		resp, err := ctx.Client.Matrix.Comp.Ratio(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor ratio: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor ratio: %.2f\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetRatio(matrix.Index.Index, *cmd.Ratio); err != nil {
		return fmt.Errorf("failed to set Matrix compressor ratio: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor ratio set to: %.2f\n", *cmd.Ratio)
	return nil
}

// MatrixCompMixCmd defines the command for getting or setting the compressor mix level of the Matrix output, allowing users to specify the desired mix level in percentage.
type MatrixCompMixCmd struct {
	Mix *float64 `arg:"" help:"The compressor mix level to set. If not provided, the current mix level will be printed." optional:""`
}

// Run executes the MatrixCompMixCmd command, either retrieving the current compressor mix level of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompMixCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Mix == nil {
		resp, err := ctx.Client.Matrix.Comp.Mix(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor mix level: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor mix level: %.2f%%\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetMix(matrix.Index.Index, *cmd.Mix); err != nil {
		return fmt.Errorf("failed to set Matrix compressor mix level: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor mix level set to: %.2f%%\n", *cmd.Mix)
	return nil
}

// MatrixCompMakeupCmd defines the command for getting or setting the compressor makeup gain of the Matrix output, allowing users to specify the desired makeup gain in dB.
type MatrixCompMakeupCmd struct {
	Makeup *float64 `arg:"" help:"The compressor makeup gain to set. If not provided, the current makeup gain will be printed." optional:""`
}

// Run executes the MatrixCompMakeupCmd command, either retrieving the current compressor makeup gain of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompMakeupCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Makeup == nil {
		resp, err := ctx.Client.Matrix.Comp.Makeup(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor makeup gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor makeup gain: %.2f dB\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetMakeup(matrix.Index.Index, *cmd.Makeup); err != nil {
		return fmt.Errorf("failed to set Matrix compressor makeup gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor makeup gain set to: %.2f dB\n", *cmd.Makeup)
	return nil
}

// MatrixCompAttackCmd defines the command for getting or setting the compressor attack time of the Matrix output, allowing users to specify the desired attack time in milliseconds.
type MatrixCompAttackCmd struct {
	Attack *float64 `arg:"" help:"The compressor attack time to set. If not provided, the current attack time will be printed." optional:""`
}

// Run executes the MatrixCompAttackCmd command, either retrieving the current compressor attack time of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompAttackCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Attack == nil {
		resp, err := ctx.Client.Matrix.Comp.Attack(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor attack time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor attack time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetAttack(matrix.Index.Index, *cmd.Attack); err != nil {
		return fmt.Errorf("failed to set Matrix compressor attack time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor attack time set to: %.2f ms\n", *cmd.Attack)
	return nil
}

// MatrixCompHoldCmd defines the command for getting or setting the compressor hold time of the Matrix output, allowing users to specify the desired hold time in milliseconds.
type MatrixCompHoldCmd struct {
	Hold *float64 `arg:"" help:"The compressor hold time to set. If not provided, the current hold time will be printed." optional:""`
}

// Run executes the MatrixCompHoldCmd command, either retrieving the current compressor hold time of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompHoldCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Hold == nil {
		resp, err := ctx.Client.Matrix.Comp.Hold(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor hold time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor hold time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetHold(matrix.Index.Index, *cmd.Hold); err != nil {
		return fmt.Errorf("failed to set Matrix compressor hold time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor hold time set to: %.2f ms\n", *cmd.Hold)
	return nil
}

// MatrixCompReleaseCmd defines the command for getting or setting the compressor release time of the Matrix output, allowing users to specify the desired release time in milliseconds.
type MatrixCompReleaseCmd struct {
	Release *float64 `arg:"" help:"The compressor release time to set. If not provided, the current release time will be printed." optional:""`
}

// Run executes the MatrixCompReleaseCmd command, either retrieving the current compressor release time of the Matrix output or setting it based on the provided argument.
func (cmd *MatrixCompReleaseCmd) Run(ctx *context, matrix *MatrixCmdGroup) error {
	if cmd.Release == nil {
		resp, err := ctx.Client.Matrix.Comp.Release(matrix.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get Matrix compressor release time: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Matrix compressor release time: %.2f ms\n", resp)
		return nil
	}

	if err := ctx.Client.Matrix.Comp.SetRelease(matrix.Index.Index, *cmd.Release); err != nil {
		return fmt.Errorf("failed to set Matrix compressor release time: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Matrix compressor release time set to: %.2f ms\n", *cmd.Release)
	return nil
}
