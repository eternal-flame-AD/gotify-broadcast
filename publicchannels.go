package main

import (
	"sync"

	plugin "github.com/gotify/plugin-api"
)

var publicChannels = new(PublicChannelListManager)

// PublicChannelListManager holds a registry of public channels at a server scope
type PublicChannelListManager struct {
	mutex    sync.RWMutex
	channels []ChannelWithUserContext
}

// UpdateChannelsForUser replaces all public channels belonging to a user context with a new slice of channels
// the publicity of the channel is checked here so it is not necessary to check for it again
func (c *PublicChannelListManager) UpdateChannelsForUser(userCtx plugin.UserContext, channels []ChannelDef) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	res := make([]ChannelWithUserContext, 0)
	for _, ch := range c.channels {
		if ch.UserContext.ID == userCtx.ID {
			continue
		}
		res = append(res, ch)
	}

	for _, def := range channels {
		if def.Public {
			res = append(res, ChannelWithUserContext{def, userCtx})
		}
	}
	c.channels = res

}

// GetAllChannels gets all public channels in the manager
func (c *PublicChannelListManager) GetAllChannels() []ChannelWithUserContext {
	res := make([]ChannelWithUserContext, 0)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, ch := range c.channels {
		res = append(res, ch)
	}
	return res
}

// ChannelWithUserContext wraps a ChannelDef with the user context that possesses it
type ChannelWithUserContext struct {
	Channel     ChannelDef
	UserContext plugin.UserContext
}
