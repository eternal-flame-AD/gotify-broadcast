package rules

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringMatch(t *testing.T) {
	Convey("Test String Match", t, func() {
		Convey("regexp", func() {
			Convey("should match", func() {
				So(stringMatch(true, "^\\[(INFO|DEBUG)\\]", "[INFO]Server started"), ShouldBeTrue)
			})
			Convey("should not match", func() {
				So(stringMatch(true, "^\\[(INFO|DEBUG)\\]", "[SEVERE]Server errored"), ShouldBeFalse)
			})
		})
		Convey("plain", func() {
			Convey("should match", func() {
				So(stringMatch(false, "ok", "ok"), ShouldBeTrue)
			})
			Convey("should not match", func() {
				So(stringMatch(false, "ok?", "ok"), ShouldBeFalse)
			})
		})
	})
}
