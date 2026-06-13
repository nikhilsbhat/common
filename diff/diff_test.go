package diff_test

import (
	"testing"

	"github.com/nikhilsbhat/common/diff"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_String(t *testing.T) {
	t.Run("renders yaml", func(t *testing.T) {
		cfg := diff.NewDiff("yaml", true, logrus.New())

		actual, err := cfg.String(map[string]string{"name": "testing"})

		require.NoError(t, err)
		assert.Contains(t, actual, "---\n")
		assert.Contains(t, actual, "name: testing")
	})

	t.Run("renders json", func(t *testing.T) {
		cfg := diff.NewDiff("json", true, logrus.New())

		actual, err := cfg.String(map[string]string{"name": "testing"})

		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"testing"}`, actual)
	})

	t.Run("returns error for unsupported format", func(t *testing.T) {
		cfg := diff.NewDiff("toml", true, logrus.New())

		actual, err := cfg.String(map[string]string{"name": "testing"})

		require.Error(t, err)
		assert.Empty(t, actual)
		assert.Contains(t, err.Error(), "type 'toml' is not supported")
	})
}

func TestConfig_Diff(t *testing.T) {
	t.Run("returns no diff for matching content", func(t *testing.T) {
		cfg := diff.NewDiff("json", true, logrus.New())

		found, actual, err := cfg.Diff(`{"name":"testing"}`, `{"name":"testing"}`)

		require.NoError(t, err)
		assert.False(t, found)
		assert.Empty(t, actual)
	})

	t.Run("returns diff for changed content", func(t *testing.T) {
		cfg := diff.NewDiff("yaml", true, logrus.New())

		found, actual, err := cfg.Diff("name: old\n", "name: new\n")

		require.NoError(t, err)
		assert.True(t, found)
		assert.Contains(t, actual, "-name: old")
		assert.Contains(t, actual, "+name: new")
	})

	t.Run("colors diff when enabled", func(t *testing.T) {
		cfg := diff.NewDiff("json", false, logrus.New())

		found, actual, err := cfg.Diff("old\n", "new\n")

		require.NoError(t, err)
		assert.True(t, found)
		assert.Contains(t, actual, "old")
		assert.Contains(t, actual, "new")
	})

	t.Run("returns error for unsupported diff format", func(t *testing.T) {
		cfg := diff.NewDiff("toml", true, logrus.New())

		found, actual, err := cfg.Diff("old", "new")

		require.Error(t, err)
		assert.False(t, found)
		assert.Empty(t, actual)
		assert.Contains(t, err.Error(), "unknown format")
	})
}
