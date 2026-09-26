//go:build urfave_cli_no_template

package cli

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoTemplate_CustomHelpTemplate(t *testing.T) {
	tests := []struct {
		name string
		cmd  *Command
	}{
		{
			name: "root",
			cmd:  &Command{Name: "foo", CustomRootCommandHelpTemplate: "{{.Name}}"},
		},
		{
			name: "command",
			cmd:  &Command{Name: "foo", CustomHelpTemplate: "{{.Name}}"},
		},
		{
			name: "subcommand",
			cmd: &Command{Name: "foo", Commands: []*Command{
				{Name: "bar", Commands: []*Command{
					{Name: "baz", CustomHelpTemplate: "{{.Name}}"},
				}},
			}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			tc.cmd.Writer = &out
			tc.cmd.ErrWriter = &errOut

			// Not only help, but any run should fail.
			err := tc.cmd.Run(context.Background(), []string{"foo"})
			require.ErrorIs(t, err, errNoTemplate)
			assert.Empty(t, out.String())
		})
	}
}

func TestNoTemplate_DefaultPrintHelpCustom(t *testing.T) {
	defer func(old io.Writer) { ErrWriter = old }(ErrWriter)

	tests := []struct {
		name  string
		templ string
		data  any
	}{
		{name: "custom template", templ: "{{.Name}}", data: &Command{Name: "foo"}},
		{name: "not a command", templ: RootCommandHelpTemplate, data: "not a command"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			ErrWriter = &errOut

			DefaultPrintHelpCustom(&out, tc.templ, tc.data, nil)
			assert.Empty(t, out.String())
			assert.Equal(t, errNoTemplate.Error()+"\n", errOut.String())
		})
	}
}
