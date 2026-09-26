//go:build !urfave_cli_no_template

package cli_test

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func ExampleCommand_Suggest() {
	cmd := &cli.Command{
		Name:                          "greet",
		ErrWriter:                     os.Stdout,
		Suggest:                       true,
		HideHelp:                      false,
		HideHelpCommand:               true,
		CustomRootCommandHelpTemplate: "(this space intentionally left blank)\n",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			fmt.Printf("Hello %v\n", cmd.String("name"))
			return nil
		},
	}

	if cmd.Run(context.Background(), []string{"greet", "--nema", "chipmunk"}) == nil {
		fmt.Println("Expected error")
	}
	// Output:
	// Incorrect Usage: flag provided but not defined: -nema
	//
	// Did you mean "--name"?
	//
	// (this space intentionally left blank)
}

func ExampleCommand_Suggest_command() {
	cmd := &cli.Command{
		ErrWriter: os.Stdout,
		Name:      "greet",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			fmt.Printf("Hello %v\n", cmd.String("name"))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:               "neighbors",
				HideHelp:           true,
				HideHelpCommand:    true,
				Suggest:            true,
				CustomHelpTemplate: "(this space intentionally left blank)\n",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "smiling"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Bool("smiling") {
						fmt.Println("😀")
					}
					fmt.Println("Hello, neighbors")
					return nil
				},
			},
		},
	}

	if cmd.Run(context.Background(), []string{"greet", "neighbors", "--sliming"}) == nil {
		fmt.Println("Expected error")
	}
	// Output:
	// Incorrect Usage: flag provided but not defined: -sliming
	//
	// Did you mean "--smiling"?
}
