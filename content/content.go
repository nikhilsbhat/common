package content

import (
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
)

// Object implements method that check for file content type.
type Object string

const (
	FileTypeYAML    = "yaml"
	FileTypeJSON    = "json"
	FileTypeCSV     = "csv"
	FileTypeString  = "string"
	FileTypeUnknown = "unknown"
)

func handlePanic(log *logrus.Logger) {
	if r := recover(); r != nil {
		log.Error("recovered from panic")
	}
}

// IsJSON checks if the passed content of JSON.
func IsJSON(content string) bool {
	var js interface{}

	return json.Unmarshal([]byte(content), &js) == nil
}

// IsJSONString checks if the passed content of JSON string.
func IsJSONString(content string) bool {
	var js string

	return json.Unmarshal([]byte(content), &js) == nil
}

// IsYAML checks if the passed content of YAML.
func IsYAML(log *logrus.Logger, content string) bool {
	var yml interface{}

	// github.com/goccy/go-yaml can produce panics
	defer handlePanic(log)

	return yaml.Unmarshal([]byte(content), &yml) == nil
}

// IsYAMLString checks if the passed content of YAML string.
func IsYAMLString(log *logrus.Logger, content string) bool {
	var yml string

	// github.com/goccy/go-yaml can produce panics
	defer handlePanic(log)

	return yaml.Unmarshal([]byte(content), &yml) == nil
}

// IsCSV checks if the passed content of CSV.
func IsCSV(content string) bool {
	csvReader := csv.NewReader(strings.NewReader(content))
	_, err := csvReader.ReadAll()

	return err == nil
}

// CheckFileType checks the file type of the content passed, it validates for YAML/JSON.
func (obj Object) CheckFileType(log *logrus.Logger) string {
	log.Debug("identifying the input file type, allowed types are YAML/JSON/CSV")

	// if IsCSV(string(obj)) {
	//	log.Debug("input file type identified as CSV")
	//
	//	return FileTypeCSV
	// }

	if IsJSONString(string(obj)) || IsYAMLString(log, string(obj)) {
		log.Debug("input file type identified as string")

		return FileTypeString
	}

	if IsJSON(string(obj)) {
		log.Debug("input file type identified as JSON")

		return FileTypeJSON
	}

	if IsYAML(log, string(obj)) {
		log.Debug("input file type identified as YAML")

		return FileTypeYAML
	}

	log.Debug("input file type identified as UNKNOWN")

	return FileTypeUnknown
}

func Marshal(data interface{}) (Object, error) {
	out, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return Object(out), nil
}

// String should return the string equivalent of Object.
func (obj Object) String() string {
	return string(obj)
}
