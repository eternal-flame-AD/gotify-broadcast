package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

const docStr = `

## Creating Broadcasts

1. Enable the plugin
1. Set up channel in plugin configuration
1. Optionally configure receiver filter
1. POST your message to %s

## Receiving Broadcasts

1. Enable the plugin
1. Optionally configure sender filter

For more information, go to [the project page](%s).
`

// GetDisplay implements public.Displayer
func (c *Plugin) GetDisplay(baseURL *url.URL) string {
	baseURL.Path = c.basePath
	messageURL := &url.URL{
		Path: "message",
	}
	messageURL.RawQuery = "channel=channel_name"
	messageURL = baseURL.ResolveReference(messageURL)

	docs := bytes.NewBufferString(fmt.Sprintf(docStr, messageURL, gitHubURL))

	docs.WriteString("\r\n\r\nPublic channels on this server:\r\n\r\n")
	docs.WriteString("```")
	w := tablewriter.NewWriter(docs)
	w.SetHeader([]string{"UserID", "UserName", "ChannelName"})
	for _, channel := range publicChannels.GetAllChannels() {
		w.Append([]string{strconv.Itoa(int(channel.UserContext.ID)), channel.UserContext.Name, channel.Channel.Name})
	}
	w.Render()
	docs.WriteString("```")

	return docs.String()
}
