package main

import (
	"testing"

	plugin "github.com/gotify/plugin-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPICompatibility(t *testing.T) {
	Convey("test API compatibility", t, func(c C) {
		c.Convey("should implement plugin", func(c C) {
			c.So(new(Plugin), ShouldImplement, (*plugin.Plugin)(nil))
		})
		c.Convey("should have configurer", func(c C) {
			c.So(new(Plugin), ShouldImplement, (*plugin.Configurer)(nil))
		})
		c.Convey("should have webhooker", func(c C) {
			c.So(new(Plugin), ShouldImplement, (*plugin.Webhooker)(nil))
		})
		c.Convey("should have messenger", func(c C) {
			c.So(new(Plugin), ShouldImplement, (*plugin.Messenger)(nil))
		})
		c.Convey("should have displayer", func(c C) {
			c.So(new(Plugin), ShouldImplement, (*plugin.Displayer)(nil))
		})
	})
}
