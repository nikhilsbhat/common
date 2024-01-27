package content_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var log *logrus.Logger

//nolint:gochecknoinits
func init() {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logger.GetLoglevel("info"))
	logrusLogger.WithField("gocd-cli", true)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	log = logrusLogger
}

func TestObject_CheckFileType(t *testing.T) {
	t.Run("should validate content as json", func(t *testing.T) {
		obj := content.Object(`{"name": "testing"}`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "json", actual)
	})

	t.Run("should validate content as unknown since malformed json passed", func(t *testing.T) {
		obj := content.Object(`{"name": "testing"`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "unknown", actual)
	})

	t.Run("should validate content as yaml", func(t *testing.T) {
		obj := content.Object(`---
name: "testing"`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate content as unknown since malformed yaml passed", func(t *testing.T) {
		obj := content.Object(`{"name": "testing"`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "unknown", actual)
	})

	t.Run("should validate content as csv", func(t *testing.T) {
		fileData, err := os.ReadFile("../fixtures/sample.csv")
		assert.NoError(t, err)

		obj := content.Object(fileData)
		actual := obj.CheckFileType(log)
		assert.Equal(t, "csv", actual)
	})

	t.Run("should fail while validating content as csv", func(t *testing.T) {
		fileData, err := os.ReadFile("../fixtures/sample_faulty.csv")
		assert.NoError(t, err)

		obj := content.Object(fileData)
		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate table content as table", func(t *testing.T) {
		data := [][]string{
			{"sn", "cat", "value"},
			{"A", "The Good", "500"},
			{"B", "The Very very Bad Man", "288"},
			{"C", "The Ugly", "120"},
			{"D", "The Gopher", "800"},
		}

		strReader := new(bytes.Buffer)

		render := renderer.GetRenderer(strReader, log, false, false, false, false, true)

		err := render.Render(data)
		assert.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(log)
		assert.Equal(t, "csv", actual)
	})
}
