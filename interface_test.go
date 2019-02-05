package main

import (
	"testing"

	plugin "github.com/gotify/plugin-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPICompatibility(t *testing.T) {
	Convey("test API compatibility", t, func() {
		Convey("should implement plugin", func() {
			So(new(Plugin), ShouldImplement, (*plugin.Plugin)(nil))
		})
		Convey("should have configurer", func() {
			So(new(Plugin), ShouldImplement, (*plugin.Configurer)(nil))
		})
		Convey("should have webhooker", func() {
			So(new(Plugin), ShouldImplement, (*plugin.Webhooker)(nil))
		})
		Convey("should have messenger", func() {
			So(new(Plugin), ShouldImplement, (*plugin.Messenger)(nil))
		})
		Convey("should have displayer", func() {
			So(new(Plugin), ShouldImplement, (*plugin.Displayer)(nil))
		})
	})
}
