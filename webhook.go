package main

import (
	"errors"

	"github.com/gotify/plugin-api"

	"github.com/gin-gonic/gin"
)

func (c *Plugin) hasChannel(channel string) bool {
	for _, ch := range c.config.Channels {
		if ch.Name == channel {
			return true
		}
	}
	return false
}

type message struct {
	Message  string                 `json:"message" query:"message" form:"message"`
	Title    string                 `json:"title" query:"title" form:"title"`
	Priority int                    `json:"priority" query:"priority" form:"priority"`
	Extras   map[string]interface{} `json:"extras" query:"-" form:"-"`
}

// RegisterWebhook implements plugin.Webhooker
func (c *Plugin) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	mux.POST("/message", func(ctx *gin.Context) {
		channel := ctx.Query("channel")
		if !c.hasChannel(channel) {
			ctx.AbortWithError(400, errors.New("channel not found"))
			return
		}
		msg := new(message)
		if err := ctx.Bind(msg); err == nil {
			c.sendMessage(plugin.Message{
				Message:  msg.Message,
				Title:    msg.Title,
				Priority: msg.Priority,
				Extras:   msg.Extras,
			}, channel)
			ctx.JSON(200, msg)
		}
	})
}
