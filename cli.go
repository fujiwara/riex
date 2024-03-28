package riex

import (
	"context"
	"time"

	"github.com/alecthomas/kong"
)

type Option struct {
	Active       bool              `help:"Show active reserved instances."`
	Pending      bool              `help:"Show payment-pending reserved instances."`
	Expired      int               `help:"Show reserved instances expired in the last specified days."`
	Days         int               `arg:"" help:"Show reserved instances that will be expired within specified days."`
	Format       string            `enum:"json,markdown,tsv" help:"Output format.(json, markdown, tsv)" default:"json"`
	DummyOutput  bool              `help:"Dummy output for testing."`
	DummyEndTime time.Time         `help:"Endtime for testing. works only with --dummy-output."`
	IgnoreTags   map[string]string `help:"Resource tag for ignore RI."`
	LocalTime    bool              `help:"Use local time for output."`
}

func RunCLI(ctx context.Context, args []string) error {
	var cli Option
	parser, err := kong.New(&cli, kong.Vars{"version": Version})
	if err != nil {
		return err
	}
	if _, err := parser.Parse(args); err != nil {
		return err
	}
	app, err := New(ctx, &cli)
	if err != nil {
		return err
	}
	if cli.DummyOutput {
		return app.RunForDummy(ctx, cli.DummyEndTime)
	} else {
		return app.Run(ctx)
	}
}
