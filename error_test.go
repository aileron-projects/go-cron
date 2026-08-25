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
		err := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		target := &ParseError{What: "foo", Value: "bbb", Inner: io.EOF}
		tester.AssertEqual(t, true, errors.Is(err, target))
	})
	t.Run("errors not equal", func(t *testing.T) {
		err := &ParseError{What: "foo"}
		target := &ParseError{What: "bar"}
		tester.AssertEqual(t, false, errors.Is(err, target))
	})
	t.Run("wrapped error equal", func(t *testing.T) {
		inner := &ParseError{What: "foo", Value: "bbb", Inner: io.EOF}
		err := fmt.Errorf("outer error [%w]", inner)
		target := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		tester.AssertEqual(t, true, errors.Is(err, target))
	})
	t.Run("wrapped error not equal", func(t *testing.T) {
		inner := &ParseError{What: "bar"}
		err := fmt.Errorf("outer error [%w]", inner)
		target := &ParseError{What: "foo"}
		tester.AssertEqual(t, false, errors.Is(err, target))
	})
	t.Run("wrapped errors equal", func(t *testing.T) {
		inner := &ParseError{What: "foo", Value: "bbb", Inner: io.EOF}
		err := fmt.Errorf("outer error [%w] [%w]", io.EOF, inner)
		target := &ParseError{What: "foo", Value: "aaa", Inner: nil}
		tester.AssertEqual(t, true, errors.Is(err, target))
	})
	t.Run("wrapped errors not equal", func(t *testing.T) {
		inner := &ParseError{What: "bar"}
		err := fmt.Errorf("outer error [%w] [%w]", io.EOF, inner)
		target := &ParseError{What: "foo"}
		tester.AssertEqual(t, false, errors.Is(err, target))
	})
}
