package core

type Result[T any] struct {
	isSuccess bool
	value     T
	err       error
}

func Ok[T any](value T) *Result[T] {
	return &Result[T]{isSuccess: true, value: value, err: nil}
}

func Fail[T any](err error) *Result[T] {
	return &Result[T]{isSuccess: false, value: *new(T), err: err}
}

func (r *Result[T]) IsSuccess() bool {
	return r.isSuccess
}

func (r *Result[T]) GetValue() T {
	if !r.isSuccess {
		var zero T
		return zero
	}
	return r.value
}

func (r *Result[T]) GetError() error {
	if r.isSuccess {
		return nil
	}
	return r.err
}