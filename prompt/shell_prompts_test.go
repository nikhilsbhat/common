package prompt_test

import (
	"testing"

	"github.com/nikhilsbhat/common/prompt"
	"github.com/stretchr/testify/assert"
)

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
