package errorz

import "errors"

var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrResultNotReady = errors.New("result is not ready")
var ErrInvalidExpression = errors.New("invalid expression")
