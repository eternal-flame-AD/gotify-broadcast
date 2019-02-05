package model

import (
	"github.com/gotify/plugin-api"
)

// Message is a message wrapper with the channel, sender and recipient.
type Message struct {
	Sender      plugin.UserContext
	Receiver    plugin.UserContext
	Msg         plugin.Message
	ChannelName string

	IsSend bool
}
