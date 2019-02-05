package rules

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func shouldContainErrString(actual interface{}, expected ...interface{}) string {
	actualErrString := (actual.(error)).Error()
	for _, exp := range expected {
		if err := ShouldContainSubstring(actualErrString, exp); err != "" {
			return err
		}
	}
	return ""
}

func TestErrors(t *testing.T) {
	Convey("Test Custom Errors", t, func() {
		Convey("missing param", func() {
			So(ErrMissingParam{
				Tag: "test tag",
			}, shouldContainErrString, "test tag", "missing")
		})
		Convey("extra param", func() {
			So(ErrExtraParam{
				ExtraParams: []string{"test_tag_1", "test_tag_2"},
			}, shouldContainErrString, "test_tag_1", "test_tag_2", "extra")
		})
		Convey("rule item error", func() {
			So(RuleItemError{
				Index: 12,
				Err:   errors.New("test_error"),
			}, shouldContainErrString, "test_error", "12")
		})
	})
}
