package bizerr

type BizError struct {
	Code       int
	HTTPStatus int
	Message    string
	Cause      error // 内部原因，可选
}

func (e *BizError) Error() string {
	return e.Message
}

func (e *BizError) WithCause(err error) *BizError {
	e.Cause = err
	return e
}

func NewBizError(code, status int, msg string) *BizError {
	return &BizError{
		Code:       code,
		HTTPStatus: status,
		Message:    msg,
	}
}
