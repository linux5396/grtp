package grtp

type GrtpError struct {
	trace string //error information
}

func (grtp *GrtpError) Error() string {
	return grtp.trace
}
func NewGrtpError(error string) *GrtpError {
	g := GrtpError{trace: error}
	return &g
}
