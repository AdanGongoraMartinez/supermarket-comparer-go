package errors

type StatusCoder interface {
	GetStatusCode() int
}
