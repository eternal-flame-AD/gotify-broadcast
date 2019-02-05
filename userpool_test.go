package main

import (
	"testing"

	"github.com/gotify/plugin-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserPool(t *testing.T) {
	Convey("Test User Pool", t, func() {
		pool := new(UserPool)
		So(pool.GetUsersList(), ShouldBeEmpty)
		Convey("Add users", func() {
			pool.AddUser(plugin.UserContext{
				ID:   1,
				Name: "test",
			})
			So(pool.GetUsersList(), ShouldHaveLength, 1)
		})
	})

}
