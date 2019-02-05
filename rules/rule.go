package rules

import (
	"fmt"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
)

// RuleChain is a set of rules that are applied in chain.
type RuleChain []Rule

// Rule is a message rule that performs an Action when a MatchSet is matched.
type Rule struct {
	Match  MatchSet `yaml:"match"`
	Action Action   `yaml:"action"`
}

// Match matches a message against a RuleChain.
// defaultAction is returned when non of the Rule matches.
func (c RuleChain) Match(msg model.Message, defaultAction Action) Action {
	for _, rule := range c {
		if matched := rule.Match.Match(msg); matched {
			return rule.Action
		}
	}
	return defaultAction
}

// Check checks a RuleChain for errors.
func (c RuleChain) Check() error {
	var errors []struct {
		Index int
		Error error
	}
	for index, rule := range c {
		if rule.Action != Accept && rule.Action != Reject {
			errors = append(errors, struct {
				Index int
				Error error
			}{index, RuleItemError{index, fmt.Errorf("unrecognized action: %s", rule.Action)}})
		}
		if err := rule.Match.Check(); err != nil {
			errors = append(errors, struct {
				Index int
				Error error
			}{index, RuleItemError{index, err}})
		}
	}
	if errors != nil {
		return RuleChainError{errors}
	}
	return nil
}
