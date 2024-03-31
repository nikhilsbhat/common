package errors

func (e *CommonError) Error() string {
	return e.Message
}
