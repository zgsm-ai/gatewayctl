package error

const (
	ErrDatabase    int = 100001
	ErrCreateRoute int = 100101
	ErrRequest     int = 100102
)

type CustomError struct {
	Code int
	Msg  string
}

func NewError(code int, msg string) CustomError {
	err := CustomError{
		Code: code,
		Msg:  msg,
	}
	return err
}
func (e CustomError) Error() string {
	return e.Msg
}

var (
	// common error
	Success                = NewError(0, "ok")
	ErrBadRequest          = NewError(400, "Bad Request")
	ErrUnauthorized        = NewError(401, "Unauthorized")
	ErrNotFound            = NewError(404, "Not Found")
	ErrInternalServerError = NewError(500, "Internal Server Error")

	// custom error
)
