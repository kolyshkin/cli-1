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
	var out, errOut bytes.Buffer
	cmd := &Command{
		Name:                          "foo",
		Writer:                        &out,
		ErrWriter:                     &errOut,
		CustomRootCommandHelpTemplate: "{{.Name}}",
	}

	require.NoError(t, cmd.Run(context.Background(), []string{"foo", "--help"}))
	assert.Empty(t, out.String())
	assert.Equal(t, errNoTemplate.Error()+"\n", errOut.String())
}

func TestNoTemplate_ModifiedHelpTemplate(t *testing.T) {
	defer func(old string) { RootCommandHelpTemplate = old }(RootCommandHelpTemplate)
	RootCommandHelpTemplate += "extra"

	var out, errOut bytes.Buffer
	cmd := &Command{
		Name:      "foo",
		Writer:    &out,
		ErrWriter: &errOut,
	}

	require.NoError(t, cmd.Run(context.Background(), []string{"foo", "--help"}))
	assert.Empty(t, out.String())
	assert.Equal(t, errNoTemplate.Error()+"\n", errOut.String())
}

func TestNoTemplate_NotCommand(t *testing.T) {
	defer func(old io.Writer) { ErrWriter = old }(ErrWriter)
	var out, errOut bytes.Buffer
	ErrWriter = &errOut

	DefaultPrintHelpCustom(&out, RootCommandHelpTemplate, "not a command", nil)
	assert.Empty(t, out.String())
	assert.Equal(t, errNoTemplate.Error()+"\n", errOut.String())
}

func TestNoTemplate_CustomFishCompletionTemplate(t *testing.T) {
	defer func(old string) { FishCompletionTemplate = old }(FishCompletionTemplate)
	FishCompletionTemplate = "{{.Command.Name}}"

	_, err := (&Command{Name: "foo"}).ToFishCompletion()
	assert.ErrorIs(t, err, errNoTemplate)
}
