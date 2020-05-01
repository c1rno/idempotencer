package errors

import (
	"encoding"
	"errors"
	"fmt"
	"strings"
)

func NewError(code int, prev error) Error {
	return baseError{
		code: code,
		prev: prev,
	}
}

type Error interface {
	error
	fmt.Stringer
	encoding.TextMarshaler
	Code() int
	IsFatal() bool
	Unwrap() error
	Is(error) bool
	As(error, interface{}) bool
}

var (
	// compile time validation
	_ error             = baseError{}
	_ Error             = baseError{}
	_ fmt.Stringer      = baseError{}
	r *strings.Replacer = strings.NewReplacer(`"`, `'`) // stop wasting json in logs
)

type baseError struct {
	code int
	prev error
}

func (e baseError) Error() string {
	return errorsMap[e.code].msg
}

// errors.Unwrap
func (e baseError) Unwrap() error {
	return e.prev
}

// errors.Is
// TODO: test it
func (e baseError) Is(err error) bool {
	x, ok := err.(Error)
	return ok && x.Code() == e.Code()
}

// errors.As
// TODO: test it
func (e baseError) As(err error, target interface{}) bool {
	if e.Is(err) {
		target = &e
		return true
	}
	return false
}

// fatal here don't means instant panic
// fatal it's flag, that allows retry or not
// if fatal == true, it's tell us that there is not any reason to retry
func (e baseError) IsFatal() bool {
	return errorsMap[e.code].fatal
}

// needs to "value"-independent comparison
func (e baseError) Code() int {
	return e.code
}

func (e baseError) String() string {
	s := printErr(e)
	for x := errors.Unwrap(e); x != nil; x = errors.Unwrap(x) {
		s = fmt.Sprintf(`%s => %s`, printErr(x), s)
	}
	return s
}

func (e baseError) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func printErr(err error) string {
	if err == nil {
		return ""
	}
	if x, ok := err.(Error); ok {
		if x.IsFatal() {
			return fmt.Sprintf(`||%d, %s||`, x.Code(), x.Error())
		} else {
			return fmt.Sprintf(`|%d, %s|`, x.Code(), x.Error())
		}
	}
	return r.Replace(fmt.Sprintf(`|%v|`, err))
}
