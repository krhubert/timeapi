package timeapi

// ErrJsonValue defines an error that occurs when a value
// used as a json value is invalid.
type ErrJsonValue struct {
	err error
}

func NewErrJsonValue(err error) ErrJsonValue {
	return ErrJsonValue{err: err}
}

func (e ErrJsonValue) Error() string {
	return "timeapi: " + e.err.Error()
}

func (e ErrJsonValue) Unwrap() error {
	return e.err
}
