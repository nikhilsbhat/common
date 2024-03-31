package diff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/nikhilsbhat/common/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sirupsen/logrus"
)

const (
	defaultContextLines = 2000000
)

// Config holds necessary information of diff.
type Config struct {
	NoColor      bool   `json:"no_color,omitempty" yaml:"no_color,omitempty"`
	Format       string `json:"format,omitempty" yaml:"format,omitempty"`
	ContextLines int    `json:"context_lines,omitempty" yaml:"context_lines,omitempty"`
	log          *logrus.Logger
}

// Diff identifies the discrepancies between two provided objects, which can be in formats such as YAML or JSON.
func (cfg *Config) Diff(oldData, newData string) (bool, string, error) {
	switch cfg.Format {
	case "yaml":
		cfg.log.Debug("loading diff in yaml format")
	case "json":
		cfg.log.Debug("loading diff in json format")
	default:
		return false, "", &errors.CommonError{Message: fmt.Sprintf("unknown format, cannot calculate diff for the format '%s'", cfg.Format)}
	}

	diffIdentified, err := cfg.diff(oldData, newData)
	if err != nil {
		return false, "", err
	}

	if len(diffIdentified) == 0 {
		return false, "", nil
	}

	return true, strings.Join(diffIdentified, "\n"), nil
}

// String returns the string representation of the DataStructure in the specified format.
func (cfg *Config) String(input interface{}) (string, error) {
	switch strings.ToLower(cfg.Format) {
	case "yaml":
		out, err := yaml.Marshal(input)
		if err != nil {
			return "", err
		}

		yamlString := strings.Join([]string{"---", string(out)}, "\n")

		return yamlString, nil
	case "json":
		out, err := json.MarshalIndent(input, "", "     ")
		if err != nil {
			return "", err
		}

		return string(out), nil
	default:
		return "", &errors.CommonError{Message: fmt.Sprintf("type '%s' is not supported for loading diff", cfg.Format)}
	}
}

// NewDiff returns a new instance of Config.
func NewDiff(format string, noColor bool, log *logrus.Logger) *Config {
	return &Config{
		NoColor: noColor,
		Format:  format,
		log:     log,
	}
}

func (cfg *Config) diff(content1, content2 string) ([]string, error) {
	contextLines := cfg.ContextLines
	if cfg.ContextLines == 0 {
		contextLines = defaultContextLines
	}

	diffVal := difflib.UnifiedDiff{
		A:        difflib.SplitLines(content1),
		B:        difflib.SplitLines(content2),
		FromFile: "old",
		ToFile:   "new",
		Context:  contextLines,
	}

	text, err := difflib.GetUnifiedDiffString(diffVal)
	if err != nil {
		return nil, err
	}

	if len(text) == 0 {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	if !cfg.NoColor {
		for index, line := range lines {
			switch {
			case strings.HasPrefix(line, "-"):
				lines[index] = color.RedString(line)
			case strings.HasPrefix(line, "+"):
				lines[index] = color.GreenString(line)
			}
		}
	}

	return lines, nil
}
