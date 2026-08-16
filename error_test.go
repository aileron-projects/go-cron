package cron

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestParseError(t *testing.T) {
	t.Parallel()
	t.Run("unwrap", func(t *testing.T) {
		err := &ParseError{Inner: io.EOF}
		inner := err.Unwrap()
		tester.AssertEqualErr(t, io.EOF, inner)
	})
	t.Run("message without value", func(t *testing.T) {
		err := &ParseError{Inner: io.EOF, What: "xxx", Value: ""}
		msg := err.Error()
		tester.AssertEqual(t, "go-cron/cron: parse: invalid xxx [EOF]", msg)
	})
	t.Run("message with value", func(t *testing.T) {
		err := &ParseError{Inner: io.EOF, What: "xxx", Value: "yyy"}
		msg := err.Error()
		tester.AssertEqual(t, "go-cron/cron: parse: invalid xxx. `yyy` [EOF]", msg)
	})
	t.Run("errors equal", func(t *testing.T) {
		err1 := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		err2 := &ParseError{What: "foo", Value: "bbb", Inner: io.EOF}
		tester.AssertEqualErr(t, err1, err2)
	})
	t.Run("wrapped error equal", func(t *testing.T) {
		err1 := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		err2 := &ParseError{What: "foo", Value: "bbb", Inner: io.EOF}
		err3 := fmt.Errorf("outer error [%w]", err2)
		tester.AssertEqual(t, true, err1.Is(err3))
		tester.AssertEqualErr(t, err1, err3)
	})
	t.Run("errors not equal", func(t *testing.T) {
		err1 := &ParseError{What: "foo", Value: "", Inner: nil}
		err2 := &ParseError{What: "bar", Value: "", Inner: io.EOF}
		tester.AssertEqual(t, false, errors.Is(err1, err2))
	})
	t.Run("wrapped error not equal", func(t *testing.T) {
		err1 := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		err2 := fmt.Errorf("outer error [%w]", io.EOF)
		tester.AssertEqual(t, false, err1.Is(err2))
	})
}
