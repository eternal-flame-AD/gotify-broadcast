package main

import (
	"sync"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
)

var msgExchanger = newMessageExchange()

type messageExchange struct {
	MsgChan   chan<- model.Message
	callbacks []func(model.Message)
	mutex     sync.RWMutex
}

func (c *messageExchange) OnMessage(cb func(model.Message)) {
	c.callbacks = append(c.callbacks, cb)
}

func newMessageExchange() *messageExchange {
	messageExchanger := new(messageExchange)
	msgChan := make(chan model.Message)
	messageExchanger.MsgChan = msgChan
	go func() {
		for {
			msg := <-msgChan
			msg.IsSend = false
			func() {
				messageExchanger.mutex.RLock()
				defer messageExchanger.mutex.RUnlock()
				for _, cb := range messageExchanger.callbacks {
					cb(msg)
				}
			}()
		}
	}()
	return messageExchanger
}
