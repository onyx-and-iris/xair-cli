package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type HeadampCmdGroup struct {
	Index struct {
		Index   int               `arg:"" help:"The index of the headamp."`
		Gain    HeadampGainCmd    `help:"Get or set the gain of the headamp."                cmd:""`
		Phantom HeadampPhantomCmd `help:"Get or set the phantom power state of the headamp." cmd:""`
	} `arg:"" help:"Control a specific headamp by index."`
}

type HeadampGainCmd struct {
	Duration time.Duration `help:"The duration of the fade in/out when setting the gain." default:"5s"`
	Gain     *float64      `help:"The gain of the headamp in dB."                                      arg:""`
}

func (cmd *HeadampGainCmd) Run(ctx *context, headamp *HeadampCmdGroup) error {
	if cmd.Gain == nil {
		resp, err := ctx.Client.HeadAmp.Gain(headamp.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get headamp gain: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Headamp %d gain: %.2f dB\n", headamp.Index.Index, resp)
		return nil
	}

	currentGain, err := ctx.Client.HeadAmp.Gain(headamp.Index.Index)
	if err != nil {
		return fmt.Errorf("failed to get current headamp gain: %w", err)
	}

	if err := gradualGainAdjust(ctx, headamp.Index.Index, currentGain, *cmd.Gain, cmd.Duration); err != nil {
		return fmt.Errorf("failed to set headamp gain: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Headamp %d gain set to: %.2f dB\n", headamp.Index.Index, *cmd.Gain)
	return nil
}

// gradualGainAdjust gradually adjusts gain from current to target over specified duration
func gradualGainAdjust(
	ctx *context,
	index int,
	currentGain, targetGain float64,
	duration time.Duration,
) error {
	gainDiff := targetGain - currentGain

	stepInterval := 100 * time.Millisecond
	totalSteps := int(duration / stepInterval)

	if totalSteps < 1 {
		totalSteps = 1
		stepInterval = duration
	}

	stepIncrement := gainDiff / float64(totalSteps)

	log.Debugf("Adjusting Headamp %d gain from %.2f dB to %.2f dB over %v...\n",
		index, currentGain, targetGain, duration)

	for step := 1; step <= totalSteps; step++ {
		newGain := currentGain + (stepIncrement * float64(step))

		if step == totalSteps {
			newGain = targetGain
		}

		err := ctx.Client.HeadAmp.SetGain(index, newGain)
		if err != nil {
			return err
		}

		if step%10 == 0 || step == totalSteps {
			log.Debugf("  Step %d/%d: %.2f dB\n", step, totalSteps, newGain)
		}

		if step < totalSteps {
			time.Sleep(stepInterval)
		}
	}

	return nil
}

type HeadampPhantomCmd struct {
	State *string `help:"The phantom power state of the headamp." arg:"" enum:"true,on,false,off" optional:""`
}

func (cmd *HeadampPhantomCmd) Validate() error {
	if cmd.State != nil {
		switch *cmd.State {
		case "true", "on":
			*cmd.State = "true"
		case "false", "off":
			*cmd.State = "false"
		default:
			return fmt.Errorf("invalid phantom power state: %s", *cmd.State)
		}
	}
	return nil
}

func (cmd *HeadampPhantomCmd) Run(ctx *context, headamp *HeadampCmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.HeadAmp.PhantomPower(headamp.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get headamp phantom power state: %w", err)
		}
		fmt.Fprintf(ctx.Out, "Headamp %d phantom power: %t\n", headamp.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.HeadAmp.SetPhantomPower(headamp.Index.Index, *cmd.State == "true"); err != nil {
		return fmt.Errorf("failed to set headamp phantom power state: %w", err)
	}
	fmt.Fprintf(ctx.Out, "Headamp %d phantom power set to: %s\n", headamp.Index.Index, *cmd.State)
	return nil
}
