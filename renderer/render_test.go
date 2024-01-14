package renderer_test

import (
	"bytes"
	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/common/prompt"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRenderer(t *testing.T) {
	t.Run("should be able to render the value to csv successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, false, true, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		assert.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "csv", actual)
	})

	t.Run("should be able to render the value to json successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, true, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		assert.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "json", actual)
	})

	t.Run("should be able to render the value to yaml successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, true, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		assert.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should be able to render the value to table successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		render := renderer.GetRenderer(strReader, logger, false, false, false, true)

		data := [][]string{
			{"sn", "cat", "value"},
			{"A", "The Good", "500"},
			{"B", "The Very very Bad Man", "288"},
			{"C", "The Ugly", "120"},
			{"D", "The Gopher", "800"},
		}

		err := render.Render(data)
		assert.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "csv", actual)
	})

	t.Run("should render in defaults since no render type was selected", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		assert.NoError(t, err)
	})

	t.Run("should render in defaults to stdout since no writer or render type specified", func(t *testing.T) {
		logger := logrus.New()
		render := renderer.GetRenderer(nil, logger, false, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		assert.NoError(t, err)
	})
}
