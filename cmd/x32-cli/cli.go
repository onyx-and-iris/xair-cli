package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	kongcompletion "github.com/jotaen/kong-completion"
	"github.com/onyx-and-iris/xair-cli/internal/xair"
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

	Version VersionFlag `help:"Print x32-cli version information and quit" name:"version" short:"v"`

	Completion kongcompletion.Completion `help:"Generate shell completion scripts." cmd:"" aliases:"c"`

	Raw      RawCmd           `help:"Send raw OSC messages to the mixer."   cmd:"" group:"Raw"`
	Main     MainCmdGroup     `help:"Control the Main L/R output"           cmd:"" group:"Main"`
	Mainmono MainMonoCmdGroup `help:"Control the Main Mono output"          cmd:"" group:"MainMono"`
	Matrix   MatrixCmdGroup   `help:"Control the matrix outputs."           cmd:"" group:"Matrix"`
	Strip    StripCmdGroup    `help:"Control the strips."                   cmd:"" group:"Strip"`
	Bus      BusCmdGroup      `help:"Control the buses."                    cmd:"" group:"Bus"`
	Headamp  HeadampCmdGroup  `help:"Control input gain and phantom power." cmd:"" group:"Headamp"`
	Snapshot SnapshotCmdGroup `help:"Save and load mixer states."           cmd:"" group:"Snapshot"`
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
				if version == "" {
					info, ok := debug.ReadBuildInfo()
					if !ok {
						return "(unable to read build info)"
					}
					version = strings.Split(info.Main.Version, "-")[0]
				}
				return version
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
	defer client.Close()

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
		xair.WithKind("x32"),
		xair.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
