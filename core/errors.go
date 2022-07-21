package core

const (
	SuccessCompleted = iota + 100
	UnknownError
)

type WrapError struct {
	Err  error
	Code int
}

func (err WrapError) Error() string {
	return err.Error()
}
