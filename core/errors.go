package core

const (
	SuccessCompleted = iota + 100
	SystemError
)

type WrapError struct {
	Err  error
	Code int
}

func (w *WrapError) Error() string {
	return w.Err.Error()
}
