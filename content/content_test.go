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
	"github.com/stretchr/testify/require"
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

	t.Run("should validate ansi colored content as json", func(t *testing.T) {
		obj := content.Object("\x1b[38;2;208;208;208m{\"name\":\"testing\"}\x1b[0m")

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

	t.Run("should validate bom prefixed content as yaml", func(t *testing.T) {
		obj := content.Object("\ufeff---\nname: \"testing\"")

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate ansi colored content as yaml", func(t *testing.T) {
		obj := content.Object("\x1b[38;2;106;184;37mname\x1b[0m: \x1b[38;2;237;157;19m\"testing\"\x1b[0m")

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate anchored content as yaml", func(t *testing.T) {
		obj := content.Object(`common: &common
  material: git
  destination: destination
pipeline:
  materials:
    my-repo:
      <<: *common
      url: https://github.com/example/repo.git`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate tagged anchored content as yaml", func(t *testing.T) {
		obj := content.Object(`pipeline: !pipeline
  materials:
    defaults: &defaults
      branch: main
    repo:
      <<: *defaults
      url: https://github.com/example/repo.git`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate custom tagged content with dotted anchors as yaml", func(t *testing.T) {
		obj := content.Object(`defaults: &pipeline.defaults
  environment_variables:
    SERVICE: api
pipelines:
  app:
    group: !gocd-group "default"
    <<: *pipeline.defaults`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate content with duplicate merge keys as yaml", func(t *testing.T) {
		obj := content.Object(`common: &common
  material: git
  destination: destination
env: &env
  environment_variables:
    SERVICE: api
pipelines:
  app:
    <<: *common
    <<: *env`)

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
		require.NoError(t, err)

		obj := content.Object(fileData)
		actual := obj.CheckFileType(log)
		assert.Equal(t, "csv", actual)
	})

	t.Run("should fail while validating content as csv", func(t *testing.T) {
		fileData, err := os.ReadFile("../fixtures/sample_faulty.csv")
		require.NoError(t, err)

		obj := content.Object(fileData)
		actual := obj.CheckFileType(log)
		assert.Equal(t, "string", actual)
	})

	t.Run("should validate table content as string", func(t *testing.T) {
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
		require.NoError(t, err)

		obj := content.Object(strReader.String())
		actual := obj.CheckFileType(log)
		assert.Equal(t, "string", actual)
	})
}

func TestMarshal(t *testing.T) {
	t.Run("should marshal data into object", func(t *testing.T) {
		obj, err := content.Marshal(map[string]string{"name": "testing"})

		require.NoError(t, err)
		assert.JSONEq(t, `{"name":"testing"}`, obj.String())
	})

	t.Run("should return error for unsupported data", func(t *testing.T) {
		obj, err := content.Marshal(func() {})

		require.Error(t, err)
		assert.Empty(t, obj)
	})
}

func TestObject_String(t *testing.T) {
	obj := content.Object("testing")

	assert.Equal(t, "testing", obj.String())
}
