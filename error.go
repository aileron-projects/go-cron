package cron

import "errors"

// ParseError reports cron parse error.
type ParseError struct {
	Inner error  // Inner is the inner error.
	What  string // What is the invalid item.
	Value string // Value is the given invalid value.
}

// Unwrap returns the inner error if any.
func (e *ParseError) Unwrap() error {
	return e.Inner
}

func (e *ParseError) Error() string {
	s := "go-cron/cron: parse: invalid " + e.What
	if e.Value != "" {
		s += ". `" + e.Value + "`"
	}
	if e.Inner != nil {
		s = s + " [" + e.Inner.Error() + "]"
	}
	return s
}

func (e *ParseError) Is(target error) bool {
	for target != nil {
		ee, ok := target.(*ParseError)
		if ok {
			return e.What == ee.What
		}
		target = errors.Unwrap(target)
	}
	return false
}
