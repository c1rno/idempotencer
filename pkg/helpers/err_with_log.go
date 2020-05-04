package helpers

import (
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
)

const (
	errField = "err"
)

func NewErrWithLog(l logging.Logger, c int, e error) errors.Error {
	t := errors.NewError(c, e)
	l.Error(t.Error(), map[string]interface{}{
		errField: t,
	})
	return t
}
