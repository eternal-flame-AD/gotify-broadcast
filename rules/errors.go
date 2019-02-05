package rules

import (
	"bytes"
	"fmt"
)

// ErrMissingParam is returned when a Match lacks a paramater.
type ErrMissingParam struct {
	Tag string
}

func (c ErrMissingParam) Error() string {
	return fmt.Sprintf("missing parameter %s", c.Tag)
}

// ErrExtraParam is returned when a Match contains extra parameters.
type ErrExtraParam struct {
	ExtraParams []string
}

func (c ErrExtraParam) Error() string {
	return fmt.Sprintf("extra parameter(s): %s", c.ExtraParams)
}

// ErrMatchSetInvalid is returned when a MatchSet contains errors.
type ErrMatchSetInvalid struct {
	Errors []struct {
		Index int
		Error error
	}
}

func (c ErrMatchSetInvalid) Error() string {
	b := bytes.NewBuffer([]byte{})
	for _, err := range c.Errors {
		b.WriteString(fmt.Sprintf("in match rule index %d: %s", err.Index, err.Error.Error()))
	}
	return b.String()
}

// RuleItemError is returned when a Rule in a RuleChain contains errors.
type RuleItemError struct {
	Index int
	Err   error
}

func (c RuleItemError) Error() string {
	return fmt.Sprintf("in chain index %d: %s", c.Index, c.Err.Error())
}

// RuleChainError is the errors in a rulechain.
type RuleChainError struct {
	Errors []struct {
		Index int
		Error error
	}
}

func (c RuleChainError) Error() string {
	b := bytes.NewBuffer([]byte{})
	for _, err := range c.Errors {
		b.WriteString(fmt.Sprintf("chain index %d: %s", err.Index, err.Error.Error()))
	}
	return b.String()
}
