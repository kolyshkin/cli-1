// This program modifies the default help template, which should result in
// a compile error when built with urfave_cli_no_template tag.
package main

import "github.com/urfave/cli/v3"

func main() {
	cli.RootCommandHelpTemplate += "extra"
}
