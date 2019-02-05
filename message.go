package main

import (
	"bytes"
	"html/template"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
	"github.com/eternal-flame-AD/gotify-broadcast/rules"
	plugin "github.com/gotify/plugin-api"
)

var msgTemplate = template.Must(template.New("message").Parse(`
{{.Msg.Message}}

==============

Sent with gotify-broadcast plugin.

Sender: {{.Sender.Name}}{{if .Sender.Admin}} (Admin){{end}}
Channel: {{.ChannelName}}
Priority: {{.Msg.Priority}}
`))

func (c *Plugin) recvMessage(msg model.Message) {
	if !c.enabled {
		return
	}
	if msg.Sender.ID == c.UserCtx.ID {
		return
	}
	if msg.Receiver.ID != c.UserCtx.ID {
		return
	}
	if action := c.config.SenderFilter.Match(msg, rules.Accept); action == rules.Accept {
		wrappedMsg := bytes.NewBuffer([]byte{})
		if err := msgTemplate.Execute(wrappedMsg, msg); err == nil {
			msg.Msg.Message = wrappedMsg.String()
		}
		c.msgHandler.SendMessage(msg.Msg)
	}
}

func (c *Plugin) sendMessage(msg plugin.Message, chanName string) int {
	sent := 0
	for _, recipient := range usersList.GetUsersList() {
		msgWrapped := model.Message{
			Sender:      c.UserCtx,
			Receiver:    recipient,
			Msg:         msg,
			ChannelName: chanName,

			IsSend: true,
		}
		if action := c.config.ReceiverFilter.Match(msgWrapped, rules.Accept); action == rules.Accept {
			msgExchanger.MsgChan <- msgWrapped
			sent++
		}
	}
	return sent
}
