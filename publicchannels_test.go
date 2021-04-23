package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gotify/plugin-api"

	. "github.com/smartystreets/goconvey/convey"
)

func shouldAllBePublicChannel(actual interface{}, expected ...interface{}) string {
	var list []ChannelDef
	switch actual := actual.(type) {
	case []ChannelDef:
		list = actual
	case []ChannelWithUserContext:
		for _, def := range actual {
			list = append(list, def.Channel)
		}
	case ChannelDef:
		list = append(list, actual)
	case ChannelWithUserContext:
		list = append(list, actual.Channel)
	}
	for _, def := range list {
		if !def.Public {
			return fmt.Sprintf("channel %s is not public", def.Name)
		}
	}
	return ""
}

func TestPublicChannels(t *testing.T) {
	Convey("Test Public Channel Registry", t, func(c C) {
		registry := new(PublicChannelListManager)
		c.So(registry.GetAllChannels(), ShouldBeEmpty)
		c.Convey("adds single user", func(c C) {
			registry.UpdateChannelsForUser(plugin.UserContext{
				ID:    1,
				Name:  "test",
				Admin: true,
			}, []ChannelDef{
				{"test_channel", true},
				{"test_private_channel", false},
			})
			c.So(registry.GetAllChannels(), ShouldHaveLength, 1)
			c.So(registry.GetAllChannels(), shouldAllBePublicChannel)
		})
		c.Convey("updates user", func(c C) {
			registry.UpdateChannelsForUser(plugin.UserContext{
				ID:    1,
				Name:  "test",
				Admin: true,
			}, []ChannelDef{
				{"test_private_channel", true},
				{"test_public_channel", true},
			})
			c.So(registry.GetAllChannels(), ShouldHaveLength, 2)
			c.So(registry.GetAllChannels(), shouldAllBePublicChannel)
		})
		c.Convey("sync safety", func(c C) {

			registry.UpdateChannelsForUser(plugin.UserContext{
				ID:    1,
				Name:  "test_1",
				Admin: true,
			}, []ChannelDef{
				{"test_channel", true},
				{"test_private_channel", false},
			})

			generaterChan := make(chan struct{})
			go func() {
				for i := 2; i < 1000; i++ {
					registry.UpdateChannelsForUser(plugin.UserContext{
						ID:    1,
						Name:  "test_1",
						Admin: true,
					}, []ChannelDef{
						{"test_channel", true},
						{"test_private_channel", false},
					})
					registry.UpdateChannelsForUser(plugin.UserContext{
						ID:    uint(i),
						Name:  "test_" + strconv.Itoa(i),
						Admin: true,
					}, []ChannelDef{
						{"test_channel", true},
						{"test_private_channel", false},
					})
				}
				close(generaterChan)
			}()
			done := false
			for !done {

				channels := registry.GetAllChannels()
				c.So(channels, shouldAllBePublicChannel)
				uID1ChanCount := 0
				for _, c := range channels {
					if c.UserContext.ID == 1 {
						uID1ChanCount++
					}
				}
				c.So(uID1ChanCount, ShouldEqual, 1)

				select {
				case <-generaterChan:
					done = true
				default:
				}
			}
		})
	})
}
