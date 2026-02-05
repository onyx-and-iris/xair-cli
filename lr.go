package main

import (
	"fmt"
	"time"
)

// MainCmdGroup defines the command group for controlling the Main L/R output, including commands for mute state, fader level, and fade-in/fade-out times.
type MainCmdGroup struct {
	Mute MainMuteCmd `help:"Get or set the mute state of the Main L/R output." cmd:""`

	Fader   MainFaderCmd   `help:"Get or set the fader level of the Main L/R output."      cmd:""`
	Fadein  MainFadeinCmd  `help:"Fade in the Main L/R output over a specified duration."  cmd:""`
	Fadeout MainFadeoutCmd `help:"Fade out the Main L/R output over a specified duration." cmd:""`
}

// MainMuteCmd defines the command for getting or setting the mute state of the Main L/R output, allowing users to specify the desired state as "true"/"on" or "false"/"off".
type MainMuteCmd struct {
	Mute *bool `arg:"" help:"The mute state to set. If not provided, the current state will be printed."`
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
	Level *float64 `arg:"" help:"The fader level to set. If not provided, the current level will be printed."`
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
