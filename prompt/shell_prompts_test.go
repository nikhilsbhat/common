package prompt_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/nikhilsbhat/common/prompt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type failingReader struct{}

var errReadFailed = errors.New("read failed")

func (failingReader) Read(_ []byte) (int, error) {
	return 0, errReadFailed
}

func TestOption_Contains(t *testing.T) {
	options := []prompt.Options{
		{
			Name:  "testing",
			Short: "t",
		},
		{
			Name:  "coding",
			Short: "c",
		},
	}

	t.Run("should be able to find an option", func(t *testing.T) {
		actual, option := prompt.Option(options).Contains("t")
		assert.True(t, actual)
		assert.Equal(t, options[0], option)
	})
	t.Run("should be able to find an option", func(t *testing.T) {
		actual, option := prompt.Option(options).Contains("test")
		assert.False(t, actual)
		assert.Equal(t, prompt.Options{}, option)
	})
}

func TestNewReadConfig(t *testing.T) {
	options := []prompt.Options{{Name: "yes", Short: "y"}}

	actual := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New())

	assert.Equal(t, "gocd-cli", actual.ShellName)
	assert.Equal(t, "continue?", actual.ShellMessage)
	assert.Equal(t, options, actual.InputOptions)
}

func TestReadConfig_WithInput(t *testing.T) {
	options := []prompt.Options{{Name: "yes", Short: "y"}}
	config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New())

	actual := config.WithInput(strings.NewReader("y\n"))

	assert.Same(t, config, actual)
}

func TestReadConfig_Reader(t *testing.T) {
	options := []prompt.Options{
		{Name: "yes", Short: "y"},
		{Name: "no", Short: "n"},
	}

	t.Run("returns matching option by short name", func(t *testing.T) {
		config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New()).
			WithInput(strings.NewReader("y\n"))

		actual, option := config.Reader()

		assert.True(t, actual)
		assert.Equal(t, options[0], option)
	})

	t.Run("returns matching option by name", func(t *testing.T) {
		config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New()).
			WithInput(strings.NewReader("no\n"))

		actual, option := config.Reader()

		assert.True(t, actual)
		assert.Equal(t, options[1], option)
	})

	t.Run("returns false for empty input", func(t *testing.T) {
		config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New()).
			WithInput(strings.NewReader("\n"))

		actual, option := config.Reader()

		assert.False(t, actual)
		assert.Equal(t, prompt.Options{}, option)
	})

	t.Run("retries until a valid option is provided", func(t *testing.T) {
		config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New()).
			WithInput(strings.NewReader("maybe\ny\n"))

		actual, option := config.Reader()

		assert.True(t, actual)
		assert.Equal(t, options[0], option)
	})

	t.Run("returns false when input cannot be read", func(t *testing.T) {
		config := prompt.NewReadConfig("gocd-cli", "continue?", options, logrus.New()).
			WithInput(failingReader{})

		actual, option := config.Reader()

		assert.False(t, actual)
		assert.Equal(t, prompt.Options{}, option)
	})
}
