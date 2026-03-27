package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
	"github.com/charmbracelet/log"
	kongcompletion "github.com/jotaen/kong-completion"

	"github.com/onyx-and-iris/xair-cli/internal/xair"
)

const (
	trueStr  = "true"
	falseStr = "false"
)

var version string // Version of the CLI, set at build time.

// VersionFlag is a custom flag type that prints the version and exits.
type VersionFlag string

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }  // nolint: revive
func (v VersionFlag) IsBool() bool                       { return true } // nolint: revive
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error { // nolint: revive, unparam
	fmt.Printf("x32-cli version: %s\n", vars["version"])
	app.Exit(0)
	return nil
}

type context struct {
	Client *xair.X32Client
	Out    io.Writer
}

type Config struct {
	Host     string        `default:"mixer.local" help:"The host of the X32 device." env:"X32_CLI_HOST"     short:"H"`
	Port     int           `default:"10023"       help:"The port of the X32 device." env:"X32_CLI_PORT"     short:"P"`
	Timeout  time.Duration `default:"100ms"       help:"Timeout for OSC operations." env:"X32_CLI_TIMEOUT"  short:"T"`
	Loglevel string        `default:"warn"        help:"Log level for the CLI."      env:"X32_CLI_LOGLEVEL" short:"L" enum:"debug,info,warn,error,fatal"`
}

// CLI is the main struct for the command-line interface.
// It embeds the Config struct for global configuration and defines the available commands and flags.
type CLI struct {
	Config `embed:"" prefix:"" help:"The configuration for the CLI."`

	Man     mangokong.ManFlag `help:"Print man page."`
	Version VersionFlag       `help:"Print x32-cli version information and quit" name:"version" short:"v"`

	Completion kongcompletion.Completion `cmd:"" help:"Generate shell completion scripts."`
	Info       InfoCmd                   `cmd:"" help:"Print mixer information."`
	Raw        RawCmd                    `cmd:"" help:"Send raw OSC messages to the mixer."`

	Main     MainCmdGroup     `cmd:"" help:"Control the Main L/R output"           group:"Main"`
	Mainmono MainMonoCmdGroup `cmd:"" help:"Control the Main Mono output"          group:"MainMono"`
	Matrix   MatrixCmdGroup   `cmd:"" help:"Control the matrix outputs."           group:"Matrix"`
	Strip    StripCmdGroup    `cmd:"" help:"Control the strips."                   group:"Strip"`
	Bus      BusCmdGroup      `cmd:"" help:"Control the buses."                    group:"Bus"`
	Headamp  HeadampCmdGroup  `cmd:"" help:"Control input gain and phantom power." group:"Headamp"`
	Snapshot SnapshotCmdGroup `cmd:"" help:"Save and load mixer states."           group:"Snapshot"`
	Dca      DCACmdGroup      `cmd:"" help:"Control DCA groups."                   group:"DCA"`
}

func main() {
	var cli CLI
	kongcompletion.Register(kong.Must(&cli))
	ctx := kong.Parse(
		&cli,
		kong.Name("x32-cli"),
		kong.Description("A CLI to control Behringer X32 mixers."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": func() string {
				if version != "" {
					return version
				}

				info, ok := debug.ReadBuildInfo()
				if !ok {
					return "(unable to read version)"
				}
				return strings.Split(info.Main.Version, "-")[0]
			}(),
		},
	)

	ctx.FatalIfErrorf(run(ctx, cli.Config))
}

// run is the main entry point for the CLI.
// It connects to the X32 device, retrieves mixer info, and then runs the command.
func run(ctx *kong.Context, config Config) error {
	loglevel, err := log.ParseLevel(config.Loglevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	log.SetLevel(loglevel)

	client, err := connect(config)
	if err != nil {
		return fmt.Errorf("failed to connect to X32 device: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Errorf("failed to close client connection: %v", err)
		}
	}()

	client.StartListening()
	resp, err := client.RequestInfo()
	if err != nil {
		return err
	}
	log.Infof("Received mixer info: %+v", resp)

	ctx.Bind(&context{
		Client: client,
		Out:    os.Stdout,
	})

	return ctx.Run()
}

// connect creates a new X32 client based on the provided configuration.
func connect(config Config) (*xair.X32Client, error) {
	client, err := xair.NewX32Client(
		config.Host,
		config.Port,
		xair.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create X32 client: %w", err)
	}

	return client, nil
}
