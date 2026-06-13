package renderer_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/nikhilsbhat/common/prompt"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_colorJSON(t *testing.T) {
	t.Run("", func(t *testing.T) {
		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logrus.New())
		jsonOut, err := yaml.Marshal(cliShellReadConfig)
		require.NoError(t, err)

		strReader := new(bytes.Buffer)
		config := renderer.GetRenderer(strReader, logrus.New(), false, false, true, false, false)

		colorString, err := config.Color("json", string(jsonOut))
		require.NoError(t, err)
		fmt.Printf("%v", colorString)
	})

	t.Run("", func(t *testing.T) {
		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logrus.New())

		valueYAML, err := yaml.Marshal(cliShellReadConfig)
		require.NoError(t, err)

		strReader := new(bytes.Buffer)
		config := renderer.GetRenderer(strReader, logrus.New(), false, true, false, false, false)

		colorString, err := config.Color("yaml", string(valueYAML))
		require.NoError(t, err)
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
		require.NoError(t, err)
		assert.NotNil(t, out)
	})
}
