package riex

import (
	"context"

	"github.com/alecthomas/kong"
)

type Option struct {
	Active   bool `help:"Show active reserved instances."`
	Expired  int  `help:"Show expired reserved instances in specified days"`
	Duration int  `arg:"" help:"Show reserved instances will be expired in specified days"`
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
