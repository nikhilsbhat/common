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
	YAML   bool `json:"yaml,omitempty" yaml:"yaml,omitempty"`
	JSON   bool `json:"json,omitempty" yaml:"json,omitempty"`
	CSV    bool `json:"csv,omitempty" yaml:"csv,omitempty"`
	Table  bool `json:"table,omitempty" yaml:"table,omitempty"`
	writer *bufio.Writer
	logger *logrus.Logger
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
func (r *Config) Render(value interface{}) error {
	if r.JSON {
		return r.ToJSON(value)
	}

	if r.YAML {
		return r.ToYAML(value)
	}

	if r.CSV {
		return r.ToCSV(value)
	}

	if r.Table {
		return r.ToTable(value)
	}

	r.logger.Debug("no format was specified for rendering output to defaults")

	_, err := r.writer.Write([]byte(fmt.Sprintf("%v\n", value)))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// ToYAML renders the value to YAML format.
func (r *Config) ToYAML(value interface{}) error {
	r.logger.Debug("rendering output in yaml format since Config.YAML is enabled")

	valueYAML, err := yaml.Marshal(value)
	if err != nil {
		return err
	}

	yamlString := strings.Join([]string{"---", string(valueYAML)}, "\n")

	_, err = r.writer.Write([]byte(yamlString))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// ToJSON renders the value to JSON format.
func (r *Config) ToJSON(value interface{}) error {
	r.logger.Debug("rendering output in json format since Config.JSON is enabled")

	valueJSON, err := json.MarshalIndent(value, "", "     ")
	if err != nil {
		return err
	}

	jsonString := strings.Join([]string{string(valueJSON), "\n"}, "")

	_, err = r.writer.Write([]byte(jsonString))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// ToCSV renders the value to CSV format.
func (r *Config) ToCSV(value interface{}) error {
	r.logger.Debug("rendering output in csv format since Config.CSV is enabled")

	csvString, err := gocsv.MarshalString(value)
	if err != nil {
		return err
	}

	_, err = r.writer.Write([]byte(csvString))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// ToTable renders the value to Table format.
func (r *Config) ToTable(value interface{}) error {
	r.logger.Debug("rendering output in table format since Config.ToTabled is enabled")

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetAlignment(tablewriter.ALIGN_CENTER) //nolint:nosnakecase
	table.SetAutoWrapText(true)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	table.AppendBulk(value.([][]string))
	table.Render()

	_, err := r.writer.Write([]byte(tableString.String()))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// GetRenderer returns the new instance of Config.
func GetRenderer(writer io.Writer, log *logrus.Logger, yaml, json, csv, table bool) Config {
	renderer := Config{
		logger: log,
		YAML:   yaml,
		JSON:   json,
		CSV:    csv,
		Table:  table,
	}

	if writer == nil {
		renderer.writer = bufio.NewWriter(os.Stdout)
	} else {
		renderer.writer = bufio.NewWriter(writer)
	}

	return renderer
}
