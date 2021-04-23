package main

import (
	"testing"
	"time"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
	plugin "github.com/gotify/plugin-api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageExchange(t *testing.T) {
	Convey("Test Message Exchange", t, func(c C) {
		exchanger := newMessageExchange()
		c.Convey("sends messages", func(c C) {
			exchanger.MsgChan <- model.Message{}
		})
		c.Convey("callback receives mesasges", func(c C) {
			test1Received, test2Received := make(chan struct{}), make(chan struct{})
			exchanger.OnMessage(func(msg model.Message) {
				if msg.Sender.ID == 1 {
					close(test1Received)
				}
			})
			exchanger.OnMessage(func(msg model.Message) {
				if msg.Sender.ID == 1 {
					close(test2Received)
				}
			})
			exchanger.MsgChan <- model.Message{Sender: plugin.UserContext{ID: 1}}
			select {
			case <-test1Received:
			case <-time.After(1 * time.Second):
				t.Error("timeout")
			}
			select {
			case <-test2Received:
			case <-time.After(1 * time.Second):
				t.Error("timeout")
			}
		})
	})

}
