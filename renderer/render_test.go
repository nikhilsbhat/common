package renderer_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/common/prompt"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRenderer(t *testing.T) {
	t.Run("should be able to render the value to csv successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, false, false, true, false)

		type csvRow struct {
			Name string
			Date string
		}

		err := render.Render([]csvRow{{Name: "nikhil", Date: "01-01-2024"}})

		require.NoError(t, err)
		assert.Contains(t, strReader.String(), "Name,Date")
		assert.Contains(t, strReader.String(), "nikhil,01-01-2024")
	})
	t.Run("should be able to render the value to json successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		render := renderer.GetRenderer(strReader, logger, false, false, true, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		require.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "json", actual)
	})

	t.Run("should be able to render the value to yaml successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, true, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		require.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(logger)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should be able to render the value to table successfully", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, false, false, false, true)

		data := [][]string{
			{"sn", "cat", "value"},
			{"A", "The Good", "500"},
		}

		err := render.Render(data)

		require.NoError(t, err)
		assert.Contains(t, strReader.String(), "The Good")
		assert.Contains(t, strReader.String(), "500")
	})

	t.Run("should render in defaults since no render type was selected", func(t *testing.T) {
		strReader := new(bytes.Buffer)

		logger := logrus.New()
		render := renderer.GetRenderer(strReader, logger, false, false, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		require.NoError(t, err)
	})

	t.Run("should render in defaults to stdout since no writer or render type specified", func(t *testing.T) {
		logger := logrus.New()
		render := renderer.GetRenderer(nil, logger, false, false, false, false, false)

		inputOptions := []prompt.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
		cliShellReadConfig := prompt.NewReadConfig("gocd-cli", "this is test message", inputOptions, logger)

		err := render.Render([]prompt.ReadConfig{*cliShellReadConfig})
		require.NoError(t, err)
	})

	t.Run("", func(_ *testing.T) {
		type Object struct {
			Name string
			Date string
		}

		newObject := []Object{
			{Name: "nikhil", Date: "01-01-2024"},
			{Name: "jon", Date: "01-02-2024"},
		}

		logger := logrus.New()
		render := renderer.GetRenderer(os.Stdout, logger, false, true, false, false, false)

		err := render.Render(newObject)
		require.NoError(t, err)
	})
}
