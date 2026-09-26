package cli

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Modifying the default templates is not supported with urfave_cli_no_template
// tag, and should result in a compile error rather than silently ignored.
func TestNoTemplateTagAssign(t *testing.T) {
	build := func(tags string) (string, error) {
		out := filepath.Join(t.TempDir(), "notemplate-assign")
		cmd := exec.Command("go", "build", "-tags", tags, "-o", out, "./testdata/notemplate-assign")
		res, err := cmd.CombinedOutput()
		return string(res), err
	}

	// Sanity check: without the tag, it builds fine.
	res, err := build("")
	require.NoError(t, err, res)

	res, err = build("urfave_cli_no_template")
	require.Error(t, err)
	assert.Contains(t, res, "cannot assign to cli.RootCommandHelpTemplate")
}
