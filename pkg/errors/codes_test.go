package errors

import (
	"errors"
	"testing"
)

func TestErrorStringer(t *testing.T) {
	expected := "|base 'error' (without Unwrap)| => ||0, Unknown error|| => |-1, | => |42, |"
	err1 := errors.New(`base "error" (without Unwrap)`)
	err2 := NewError(UnknownError, err1)
	err3 := NewError(-1, err2)
	err4 := NewError(42, err3)

	if err4.String() != expected {
		t.Fatalf(
			"Incorrect err representation:\n%s\nnot eq\n%s\n",
			err4.String(),
			expected,
		)
	}
}
