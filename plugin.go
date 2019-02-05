package main

import (
	"github.com/gotify/plugin-api"
)

// Plugin is plugin instance
type Plugin struct {
	config     *Config
	enabled    bool
	msgHandler plugin.MessageHandler
	basePath   string

	UserCtx plugin.UserContext
}

// Enable implements plugin.Plugin
func (c *Plugin) Enable() error {
	c.enabled = true
	return nil
}

// Disable implements plugin.Disable
func (c *Plugin) Disable() error {
	c.enabled = false
	return nil
}

// SetMessageHandler implements plugin.Messenger
func (c *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}
