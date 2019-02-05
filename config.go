package main

import (
	"fmt"

	"github.com/eternal-flame-AD/gotify-broadcast/rules"
)

// ChannelDef is the definition of a channel in the configuration
type ChannelDef struct {
	Name   string `yaml:"name"`
	Public bool   `yaml:"public"`
}

// Config is user plugin configuration
type Config struct {
	Channels       []ChannelDef    `yaml:"channels"`
	SenderFilter   rules.RuleChain `yaml:"sender_filter"`
	ReceiverFilter rules.RuleChain `yaml:"receiver_filter"`
}

// DefaultConfig implements plugin.Configurer
func (c *Plugin) DefaultConfig() interface{} {
	return &Config{
		Channels: []ChannelDef{
			ChannelDef{
				Name:   "example",
				Public: false,
			},
		},
		SenderFilter: rules.RuleChain{
			rules.Rule{
				Match: rules.MatchSet{
					rules.Match{
						Mode:     rules.ModeUserName,
						UserName: "my_server",
					},
					rules.Match{
						Mode:        rules.ModeMessageText,
						Regex:       true,
						MessageText: "^\\[(INFO|DEBUG)\\]",
					},
				},
				Action: rules.Reject,
			},
			rules.Rule{
				Match: rules.MatchSet{
					rules.Match{
						Mode:     rules.ModeUserName,
						UserName: "some_one_i_dont_want_to_see_broadcast_from",
					},
				},
				Action: rules.Reject,
			},
			rules.Rule{
				Match: rules.MatchSet{
					rules.Match{
						Mode: rules.ModeAny,
					},
				},
				Action: rules.Accept,
			},
		},
		ReceiverFilter: rules.RuleChain{
			rules.Rule{
				Match: rules.MatchSet{
					rules.Match{
						Mode:     rules.ModeUserName,
						UserName: "some_one_i_dont_want_to_send_broadcast_to",
					},
				},
				Action: rules.Reject,
			},
			rules.Rule{
				Match: rules.MatchSet{
					rules.Match{
						Mode: rules.ModeAny,
					},
				},
				Action: rules.Accept,
			},
		},
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *Plugin) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)

	if err := newConfig.ReceiverFilter.Check(); err != nil {
		return err
	}

	channels := make(map[string]struct{})
	for _, ch := range newConfig.Channels {
		if _, ok := channels[ch.Name]; ok {
			return fmt.Errorf("channel name %s is duplicated", ch.Name)
		}
		channels[ch.Name] = struct{}{}
	}

	publicChannels.UpdateChannelsForUser(c.UserCtx, newConfig.Channels)
	c.config = newConfig
	return nil
}
