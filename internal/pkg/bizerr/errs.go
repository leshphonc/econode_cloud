package bizerr

import "net/http"

var (
	ErrParamInvalid = NewBizError(10000, http.StatusBadRequest, "非法参数")
)
