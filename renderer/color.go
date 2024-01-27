package renderer

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-colorable"
	"github.com/neilotoole/jsoncolor"
)

// ColorJSON add colors to your JSON string.
func (cfg *Config) ColorJSON(value interface{}) (string, error) {
	strReader := new(bytes.Buffer)

	jsonEnc := jsoncolor.NewEncoder(strReader)
	jsonColors := jsoncolor.DefaultColors()
	jsonEnc.SetColors(jsonColors)
	jsonEnc.SetIndent("", "   ")

	if err := jsonEnc.Encode(value); err != nil {
		return "", err
	}

	return strReader.String(), nil
}

// Shamelessly borrowed ycat functionality from https://github.com/goccy/go-yaml/blob/master/cmd/ycat/ycat.go with customized changes.

// ColorYAML add colors to your YAML string.
func (cfg *Config) ColorYAML(text string) (string, error) {
	tokens := lexer.Tokenize(text)

	var prnt printer.Printer

	prnt.Bool = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	prnt.Number = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	prnt.MapKey = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiCyan),
			Suffix: format(color.Reset),
		}
	}
	prnt.Anchor = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	prnt.Alias = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	prnt.String = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiGreen),
			Suffix: format(color.Reset),
		}
	}

	outputFileCache := "highlighted_output.yaml"

	file, err := os.Create(outputFileCache)
	if err != nil {
		if err = os.RemoveAll(outputFileCache); err != nil {
			if !os.IsNotExist(err) {
				return "", err
			}
		}

		return "", err
	}

	defer func() {
		if err = os.RemoveAll(outputFileCache); err != nil {
			if !os.IsNotExist(err) {
				cfg.logger.Errorf("removing file %s as part of clean up errored with '%v'", outputFileCache, err)
			}
		}
	}()

	writer := colorable.NewColorable(file)

	if _, err = writer.Write([]byte(prnt.PrintTokens(tokens) + "\n")); err != nil {
		return "", err
	}

	highlightedOutputFile, err := os.ReadFile(outputFileCache)
	if err != nil {
		return "", err
	}

	return string(highlightedOutputFile), nil
}

func format(attr color.Attribute) string {
	const escape = "\x1b"

	return fmt.Sprintf("%s[%dm", escape, attr)
}
