package errors_test

import (
	"testing"

	"github.com/nikhilsbhat/common/errors"
	"github.com/stretchr/testify/assert"
)

func TestCommonError_Error(t *testing.T) {
	err := &errors.CommonError{Message: "something failed"}

	assert.Equal(t, "something failed", err.Error())
}
