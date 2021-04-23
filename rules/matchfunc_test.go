package rules

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringMatch(t *testing.T) {
	Convey("Test String Match", t, func(c C) {
		c.Convey("regexp", func(c C) {
			c.Convey("should match", func(c C) {
				c.So(stringMatch(true, "^\\[(INFO|DEBUG)\\]", "[INFO]Server started"), ShouldBeTrue)
			})
			c.Convey("should not match", func(c C) {
				c.So(stringMatch(true, "^\\[(INFO|DEBUG)\\]", "[SEVERE]Server errored"), ShouldBeFalse)
			})
		})
		c.Convey("plain", func(c C) {
			c.Convey("should match", func(c C) {
				c.So(stringMatch(false, "ok", "ok"), ShouldBeTrue)
			})
			c.Convey("should not match", func(c C) {
				c.So(stringMatch(false, "ok?", "ok"), ShouldBeFalse)
			})
		})
	})
}
