// Package errors contains common error types used by this module.
package errors

func (e *CommonError) Error() string {
	return e.Message
}
