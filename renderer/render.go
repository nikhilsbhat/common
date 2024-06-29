package renderer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gocarina/gocsv"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// Config implements methods to render output in JSON/YAML format.
type Config struct {
	YAML    bool `json:"yaml,omitempty" yaml:"yaml,omitempty"`
	JSON    bool `json:"json,omitempty" yaml:"json,omitempty"`
	CSV     bool `json:"csv,omitempty" yaml:"csv,omitempty"`
	Table   bool `json:"table,omitempty" yaml:"table,omitempty"`
	NoColor bool `json:"no_color,omitempty" yaml:"no_color,omitempty"`
	writer  *bufio.Writer
	logger  *logrus.Logger
}

// Renderer implements methods that Prints values in YAML,JSON,CSV and Table format.
type Renderer interface {
	ToYAML(value interface{}) error
	ToJSON(value interface{}) error
	ToCSV(value interface{}) error
	ToTable(value interface{}) error
}

// Render renders the output based on the output format selection (toYAML, toJSON).
// If none is selected it prints as the source.
func (cfg *Config) Render(value interface{}) error {
	if cfg.JSON {
		return cfg.ToJSON(value)
	}

	if cfg.YAML {
		return cfg.ToYAML(value)
	}

	if cfg.CSV {
		return cfg.ToCSV(value)
	}

	if cfg.Table {
		return cfg.ToTable(value)
	}

	cfg.logger.Debug("no format was specified for rendering output to defaults")

	_, err := cfg.writer.Write([]byte(fmt.Sprintf("%v\n", value)))
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cfg.logger.Fatalln(err)
		}
	}(cfg.writer)

	return nil
}

// ToYAML renders the value to YAML format.
func (cfg *Config) ToYAML(value interface{}) error {
	cfg.logger.Debug("rendering output in yaml format since Config.YAML is enabled")

	valueYAML, err := yaml.Marshal(value)
	if err != nil {
		return err
	}

	yamlString := strings.Join([]string{"---", string(valueYAML)}, "\n")

	if !cfg.NoColor {
		coloredYAMLString, err := cfg.Color("yaml", string(valueYAML))
		if err != nil {
			return err
		}

		yamlString = coloredYAMLString
	}

	_, err = cfg.writer.Write([]byte(yamlString))
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cfg.logger.Fatalln(err)
		}
	}(cfg.writer)

	return nil
}

// ToJSON renders the value to JSON format.
func (cfg *Config) ToJSON(value interface{}) error {
	cfg.logger.Debug("rendering output in json format since Config.JSON is enabled")

	valueJSON, err := json.MarshalIndent(value, "", "     ")
	if err != nil {
		return err
	}

	jsonString := string(valueJSON)

	if !cfg.NoColor {
		coloredJSONString, err := cfg.Color("json", jsonString)
		if err != nil {
			return err
		}

		jsonString = coloredJSONString
	}

	_, err = cfg.writer.Write([]byte(jsonString))
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cfg.logger.Fatalln(err)
		}
	}(cfg.writer)

	return nil
}

// ToCSV renders the value to CSV format.
func (cfg *Config) ToCSV(value interface{}) error {
	cfg.logger.Debug("rendering output in csv format since Config.CSV is enabled")

	csvString, err := gocsv.MarshalString(value)
	if err != nil {
		return err
	}

	_, err = cfg.writer.Write([]byte(csvString))
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cfg.logger.Fatalln(err)
		}
	}(cfg.writer)

	return nil
}

// ToTable renders the value to Table format.
func (cfg *Config) ToTable(value interface{}) error {
	cfg.logger.Debug("rendering output in table format since Config.ToTabled is enabled")

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetAlignment(tablewriter.ALIGN_CENTER) //nolint:nosnakecase
	table.SetAutoWrapText(true)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	table.AppendBulk(value.([][]string))
	table.Render()

	_, err := cfg.writer.Write([]byte(tableString.String()))
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cfg.logger.Fatalln(err)
		}
	}(cfg.writer)

	return nil
}

// GetRenderer returns the new instance of Config.
func GetRenderer(writer io.Writer, log *logrus.Logger, noColor, yaml, json, csv, table bool) Config {
	renderer := Config{
		logger:  log,
		YAML:    yaml,
		JSON:    json,
		CSV:     csv,
		Table:   table,
		NoColor: noColor,
	}

	if writer == nil {
		renderer.writer = bufio.NewWriter(os.Stdout)
	} else {
		renderer.writer = bufio.NewWriter(writer)
	}

	return renderer
}
