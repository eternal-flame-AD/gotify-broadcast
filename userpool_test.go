package main

import (
	"testing"

	"github.com/gotify/plugin-api"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserPool(t *testing.T) {
	Convey("Test User Pool", t, func(c C) {
		pool := new(UserPool)
		c.So(pool.GetUsersList(), ShouldBeEmpty)
		c.Convey("Add users", func(c C) {
			pool.AddUser(plugin.UserContext{
				ID:   1,
				Name: "test",
			})
			c.So(pool.GetUsersList(), ShouldHaveLength, 1)
		})
	})

}
