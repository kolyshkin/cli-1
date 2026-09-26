//go:build !urfave_cli_no_template

package cli

import (
	"io"
	"strings"
	"text/tabwriter"
	"text/template"
)

// templatesSupported is true unless built with the urfave_cli_no_template tag.
const templatesSupported = true

// RootCommandHelpTemplate is the text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var RootCommandHelpTemplate = rootCommandHelpTemplate

// CommandHelpTemplate is the text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = commandHelpTemplate

// SubcommandHelpTemplate is the text template for the subcommand help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var SubcommandHelpTemplate = subcommandHelpTemplate

// FishCompletionTemplate is the text template for fish shell completion.
var FishCompletionTemplate = fishCompletionTemplate

// DefaultPrintHelpCustom is the default implementation of HelpPrinterCustom.
//
// The customFuncs map will be combined with a default template.FuncMap to
// allow using arbitrary functions in template rendering.
func DefaultPrintHelpCustom(out io.Writer, templ string, data any, customFuncs map[string]any) {
	const maxLineLength = 10000

	tracef("building default funcMap")
	funcMap := template.FuncMap{
		"join":           strings.Join,
		"subtract":       subtract,
		"indent":         indent,
		"nindent":        nindent,
		"trim":           strings.TrimSpace,
		"wrap":           func(input string, offset int) string { return wrap(input, offset, maxLineLength) },
		"offset":         offset,
		"offsetCommands": offsetCommands,
	}

	if wa, ok := customFuncs["wrapAt"]; ok {
		if wrapAtFunc, ok := wa.(func() int); ok {
			wrapAt := wrapAtFunc()
			customFuncs["wrap"] = func(input string, offset int) string {
				return wrap(input, offset, wrapAt)
			}
		}
	}

	for key, value := range customFuncs {
		funcMap[key] = value
	}

	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))

	if _, err := t.New("helpNameTemplate").Parse(helpNameTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("argsTemplate").Parse(argsTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("usageTemplate").Parse(usageTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("descriptionTemplate").Parse(descriptionTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visibleCommandTemplate").Parse(visibleCommandTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("copyrightTemplate").Parse(copyrightTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("versionTemplate").Parse(versionTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visibleFlagCategoryTemplate").Parse(visibleFlagCategoryTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visibleFlagTemplate").Parse(visibleFlagTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visiblePersistentFlagTemplate").Parse(visiblePersistentFlagTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visibleGlobalFlagCategoryTemplate").Parse(strings.ReplaceAll(visibleFlagCategoryTemplate, "OPTIONS", "GLOBAL OPTIONS")); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("authorsTemplate").Parse(authorsTemplate); err != nil {
		handleTemplateError(err)
	}

	if _, err := t.New("visibleCommandCategoryTemplate").Parse(visibleCommandCategoryTemplate); err != nil {
		handleTemplateError(err)
	}

	tracef("executing template")
	handleTemplateError(t.Execute(w, data))

	_ = w.Flush()
}

func checkNoCustomTemplates(*Command) error {
	return nil
}

func renderFishCompletion(w io.Writer, data *fishCommandCompletionTemplate) error {
	const name = "cli"
	t, err := template.New(name).Parse(FishCompletionTemplate)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(w, name, data)
}
