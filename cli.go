package riex

import (
	"context"

	"github.com/alecthomas/kong"
)

type Option struct {
	Active  bool   `help:"Show active reserved instances."`
	Expired int    `help:"Show reserved instances expired in the last specified days."`
	Days    int    `arg:"" help:"Show reserved instances that will be expired within specified days."`
	Format  string `enum:"json,markdown,tsv" help:"Output format.(json, markdown, tsv)" default:"json"`
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
	return app.Run(ctx)
}
