package renderer_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/nikhilsbhat/common/prompt"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_colorJSON(t *testing.T) {
	t.Run("", func(t *testing.T) {
		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logrus.New())
		jsonOut, err := yaml.Marshal(cliShellReadConfig)
		assert.NoError(t, err)

		strReader := new(bytes.Buffer)
		config := renderer.GetRenderer(strReader, logrus.New(), false, false, true, false, false)

		colorString, err := config.Color("json", string(jsonOut))
		assert.NoError(t, err)
		fmt.Printf("%v", colorString)
	})

	t.Run("", func(t *testing.T) {
		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logrus.New())

		valueYAML, err := yaml.Marshal(cliShellReadConfig)
		assert.NoError(t, err)

		strReader := new(bytes.Buffer)
		config := renderer.GetRenderer(strReader, logrus.New(), false, true, false, false, false)

		colorString, err := config.Color("yaml", string(valueYAML))
		assert.NoError(t, err)
		fmt.Printf("%v", colorString)
	})

	t.Run("", func(t *testing.T) {
		yamlContent := `
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
    - name: myfrontend
      image: nginx`

		config := renderer.GetRenderer(nil, logrus.New(), false, true, false, false, false)
		out, err := config.Color("yaml", yamlContent)
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})
}
