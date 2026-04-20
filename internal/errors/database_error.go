package errors

type DatabaseError struct {
	Message string
	Cause   error
}

func (e *DatabaseError) Error() string {
	return e.Message
}

func NewDatabaseError(message string, cause error) *DatabaseError {
	return &DatabaseError{Message: message, Cause: cause}
}