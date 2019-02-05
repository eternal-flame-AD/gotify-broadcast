package rules

import (
	"regexp"

	plugin "github.com/gotify/plugin-api"
)

func stringMatch(regex bool, matcher string, question string) bool {
	if regex {
		matched, _ := regexp.MatchString(matcher, question)
		return matched
	}
	return matcher == question
}
func containExtra(regex bool, extra string, msg plugin.Message) bool {
	if msg.Extras == nil {
		return false
	}
	for key := range msg.Extras {
		if stringMatch(regex, extra, key) {
			return true
		}
	}
	return false
}
