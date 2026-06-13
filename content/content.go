// Package content provides helpers for identifying and converting serialized content.
package content

import (
	"encoding/csv"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/sirupsen/logrus"
	yamlv3 "gopkg.in/yaml.v3"
)

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

// Object implements method that check for file content type.
type Object string

const (
	// FileTypeYAML identifies YAML content.
	FileTypeYAML = "yaml"
	// FileTypeJSON identifies JSON content.
	FileTypeJSON = "json"
	// FileTypeCSV identifies CSV content.
	FileTypeCSV = "csv"
	// FileTypeString identifies string content.
	FileTypeString = "string"
	// FileTypeUnknown identifies content that could not be classified.
	FileTypeUnknown = "unknown"
)

func handlePanic(log *logrus.Logger) {
	if r := recover(); r != nil {
		log.Error("recovered from panic")
	}
}

// IsJSON checks if the passed content of JSON.
func IsJSON(content string) bool {
	var js any

	return json.Unmarshal([]byte(content), &js) == nil
}

// IsJSONString checks if the passed content of JSON string.
func IsJSONString(content string) bool {
	var js string

	return json.Unmarshal([]byte(content), &js) == nil
}

// IsYAML checks if the passed content of YAML.
func IsYAML(log *logrus.Logger, content string) bool {
	// github.com/goccy/go-yaml can produce panics
	defer handlePanic(log)

	file, err := parser.ParseBytes([]byte(content), 0, parser.AllowDuplicateMapKey())
	if err == nil {
		for _, doc := range file.Docs {
			if doc != nil && isStructuredYAMLNode(doc.Body) {
				return true
			}
		}
	}

	var node yamlv3.Node

	err = yamlv3.Unmarshal([]byte(content), &node)
	if err != nil {
		return false
	}

	return isStructuredYAMLV3Node(&node)
}

func isStructuredYAMLNode(node ast.Node) bool {
	if _, ok := node.(ast.MapNode); ok {
		return true
	}

	if _, ok := node.(ast.ArrayNode); ok {
		return true
	}

	switch astNode := node.(type) {
	case *ast.DocumentNode:
		return isStructuredYAMLNode(astNode.Body)
	case *ast.AnchorNode:
		return isStructuredYAMLNode(astNode.Value)
	case *ast.TagNode:
		return isStructuredYAMLNode(astNode.Value)
	default:
		return false
	}
}

func isStructuredYAMLV3Node(node *yamlv3.Node) bool {
	if node == nil {
		return false
	}

	switch node.Kind {
	case yamlv3.DocumentNode:
		return len(node.Content) > 0 && isStructuredYAMLV3Node(node.Content[0])
	case yamlv3.MappingNode, yamlv3.SequenceNode:
		return true
	case yamlv3.AliasNode:
		return isStructuredYAMLV3Node(node.Alias)
	case yamlv3.ScalarNode:
		return false
	default:
		return false
	}
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

	records, err := csvReader.ReadAll()
	if err != nil || len(records) < 2 {
		return false
	}

	return len(records[0]) > 1
}

func normalizeContent(content string) string {
	normalized := ansiEscapePattern.ReplaceAllString(content, "")
	normalized = strings.TrimPrefix(normalized, "\ufeff")

	return strings.Trim(normalized, "\r\n")
}

// CheckFileType checks the file type of the content passed, it validates for YAML/JSON.
func (obj Object) CheckFileType(log *logrus.Logger) string {
	log.Debug("identifying the input file type, allowed types are YAML/JSON/CSV")

	content := normalizeContent(string(obj))

	if IsJSON(content) {
		log.Debug("input file type identified as JSON")

		return FileTypeJSON
	}

	if IsCSV(content) {
		log.Debug("input file type identified as CSV")

		return FileTypeCSV
	}

	if IsYAML(log, content) {
		log.Debug("input file type identified as YAML")

		return FileTypeYAML
	}

	if IsJSONString(content) || IsYAMLString(log, content) {
		log.Debug("input file type identified as string")

		return FileTypeString
	}

	log.Debug("input file type identified as UNKNOWN")

	return FileTypeUnknown
}

// Marshal converts data into an Object by JSON marshaling it.
func Marshal(data any) (Object, error) {
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
