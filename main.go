package main

import (
	plugin "github.com/gotify/plugin-api"
)

const gitHubURL = "https://www.github.com/eternal-flame-AD/gotify-broadcast"

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/eternal-flame-AD/gotify-broadcast",
		Name:        "Gotify Broadcaster",
		Description: "A plugin which brings broadcasts to gotify.",
		Author:      "eternal-flame-AD",
		Website:     gitHubURL,
	}
}

func main() {
	panic("this should be built as plugin")
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	usersList.AddUser(ctx)
	p := &Plugin{
		UserCtx: ctx,
	}
	msgExchanger.OnMessage(p.recvMessage)
	return p
}
