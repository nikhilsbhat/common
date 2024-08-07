package renderer

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/nikhilsbhat/common/errors"
)

const (
	TypeYAML = "yaml"
	TypeJSON = "json"
)

// Color add colors to your YAML, JSON or any specified string.
func (cfg *Config) Color(contentType, yamlContent string) (string, error) {
	lexer := lexers.Get(contentType)
	if lexer == nil {
		return "", &errors.CommonError{Message: fmt.Sprintf("no lexer found for '%s'", contentType)}
	}

	lexer = chroma.Coalesce(lexer)

	style := styles.Get("native")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		return "", &errors.CommonError{Message: "no terminal formatter found"}
	}

	iterator, err := lexer.Tokenise(nil, yamlContent)
	if err != nil {
		return "", &errors.CommonError{Message: fmt.Sprintf("tokenise errored with '%v'", err)}
	}

	var stringBuffer bytes.Buffer

	err = formatter.Format(&stringBuffer, style, iterator)
	if err != nil {
		return "", &errors.CommonError{Message: fmt.Sprintf("formatting yaml errored with '%v'", err)}
	}

	return stringBuffer.String(), err
}
