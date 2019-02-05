package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
	"github.com/gotify/plugin-api"
)

const (
	// Accept accepts the message.
	Accept Action = "accept"
	// Reject drops the message.
	Reject Action = "reject"
)

// Action describes how the message is handled after matching a RuleSet.
type Action string

const (
	// ModeAny matches all messages.
	// No parameters are required.
	ModeAny Mode = "any"

	// ModeChannelName matches the channel name the message is sent through.
	// Use parameter channel_name to specity the channel name to match.
	// Use parameter regex: true to enable regex matching.
	ModeChannelName Mode = "channel_name"
	// ModeUserName matches the user name of the message (matches the sender on the recipient side and matches the recipient on the sender side).
	// Use parameter user_name to specify the user name to match.
	// Use parameter regex: true to enable regex matching.
	ModeUserName Mode = "user_name"
	// ModeUserID matches the user ID of the message (matches the sender on the recipient side and matches the recipient on the sender side).
	// Use parameter user_id to specify the user ID to match.
	// Use parameter regex: true to enable regex matching.
	ModeUserID Mode = "user_id"
	// ModeIsAdmin matches whether the user is an admin (matches the sender on the recipient side and matches the recipient on the sender side).
	// Use parameter is_admin to specity whether to match admins or non-admins.
	ModeIsAdmin Mode = "is_admin"

	// ModeMessageTitle matches the message title.
	// Use parameter message_title to specify the title to match.
	// Use parameter regex: true to enable regex matching.
	ModeMessageTitle Mode = "message_title"
	// ModeMessageText matches the message text.
	// Use parameter message_text to specity the text to match.
	// Use parameter regex: true to enable regex matching.
	ModeMessageText Mode = "message_text"
	// ModeMessageExtra matches whether the message possesses an extra.
	// Use parameter message_extra to specity the key of the extra to match.
	// Use parameter regex: true to enable regex matching.
	ModeMessageExtra Mode = "message_extra"

	// ModePriorityLt matches messages with priority less then a specified value.
	// Use parameter priority to specity the priority threshold.
	ModePriorityLt Mode = "message_priority_lt"
	// ModePriorityGt matches messages with priority greater then a specified value.
	// Use parameter priority to specity the priority threshold.
	ModePriorityGt Mode = "message_priority_gt"
	// ModePriority matches messages with priority at a specified value.
	// Use parameter priority to specity the priority.
	ModePriority Mode = "message_priority"
)

// Mode describes a Match matches which aspect of the message.
type Mode string

// MatchSet is a set of Matches.
// MatchSet matches message by validating every Match and only success when all Match is satisfied.
type MatchSet []Match

// Check checks a MatchSet for syntax errors.
func (c MatchSet) Check() error {
	var errors []struct {
		Index int
		Error error
	}
	for index, rule := range c {
		if err := rule.Check(); err != nil {
			errors = append(errors, struct {
				Index int
				Error error
			}{index, err})
		}
	}
	if errors != nil {
		return ErrMatchSetInvalid{errors}
	}
	return nil
}

// Match matches a Matchset with a message.
func (c MatchSet) Match(msg model.Message) bool {
	for _, rule := range c {
		if !rule.Match(msg) {
			return false
		}
	}
	return true
}

// Match describes a Match.
type Match struct {
	// Match mode
	// See Also: Mode
	Mode Mode `yaml:"mode"`

	// Match options
	Regex bool `yaml:"regex,omitempty"`

	// Match parameters
	// Only filled in as required by the Mode specified.
	// See Also: Mode
	ChannelName string `yaml:"channel_name,omitempty"`
	UserName    string `yaml:"user_name,omitempty"`
	UserID      uint   `yaml:"user_id,omitempty"`
	IsAdmin     *bool  `yaml:"is_admin,omitempty"`

	MessageTitle    string `yaml:"message_title,omitempty"`
	MessageText     string `yaml:"message_text,omitempty"`
	MessageExtra    string `yaml:"message_extra,omitempty"`
	MessagePriority *int   `yaml:"priority,omitempty"`
}

func (c Match) getYAMLTagName(fieldname string) string {
	fieldKey, _ := reflect.TypeOf(c).FieldByName(fieldname)
	yamlTag := fieldKey.Tag.Get("yaml")
	yamlTag = strings.SplitN(yamlTag, ",", 2)[0]
	return yamlTag
}

func (c Match) paramFields() []string {
	var res []string

	val := reflect.ValueOf(c)

	for _, field := range []string{
		"ChannelName",
		"UserName",
		"UserID",
		"IsAdmin",
		"MessageTitle",
		"MessageText",
		"MessageExtra",
		"MessagePriority",
	} {
		fieldVal := val.FieldByName(field)
		if !isZero(fieldVal) {
			res = append(res, c.getYAMLTagName(field))
		}
	}
	return res
}

// Check checks a match for syntax errors.
// Possible errors include: extra parameters, missing parameter.
func (c Match) Check() error {
	switch c.Mode {
	case ModeAny:
	case ModeChannelName:
		if c.ChannelName == "" {
			return ErrMissingParam{c.getYAMLTagName("ChannelName")}
		}
		c.ChannelName = ""
	case ModeUserName:
		if c.UserName == "" {
			return ErrMissingParam{c.getYAMLTagName("UserName")}
		}
		c.UserName = ""
	case ModeUserID:
		if c.UserID == 0 {
			return ErrMissingParam{c.getYAMLTagName("UserID")}
		}
		c.UserID = 0
	case ModeIsAdmin:
		if c.IsAdmin == nil {
			return ErrMissingParam{c.getYAMLTagName("IsAdmin")}
		}
		c.IsAdmin = nil
	case ModeMessageTitle:
		if c.MessageTitle == "" {
			return ErrMissingParam{c.getYAMLTagName("MessageTitle")}
		}
		c.MessageTitle = ""
	case ModeMessageText:
		if c.MessageText == "" {
			return ErrMissingParam{c.getYAMLTagName("MessageText")}
		}
		c.MessageText = ""
	case ModeMessageExtra:
		if c.MessageExtra == "" {
			return ErrMissingParam{c.getYAMLTagName("MessageExtra")}
		}
		c.MessageExtra = ""
	case ModePriority, ModePriorityGt, ModePriorityLt:
		if c.MessagePriority == nil {
			return ErrMissingParam{c.getYAMLTagName("MessagePriority")}
		}
		c.MessagePriority = nil
	default:
		return fmt.Errorf("unsupported mode: %s", c.Mode)
	}

	if extraFields := c.paramFields(); len(extraFields) > 0 {
		return ErrExtraParam{extraFields}
	}
	return nil
}

// Match matches a message against a Match.
func (c Match) Match(msg model.Message) (matched bool) {
	var userInfo plugin.UserContext
	if msg.IsSend {
		userInfo = msg.Receiver
	} else {
		userInfo = msg.Sender
	}
	switch c.Mode {
	case ModeAny:
		return true
	case ModeChannelName:
		return stringMatch(c.Regex, c.ChannelName, msg.ChannelName)
	case ModeUserName:
		return stringMatch(c.Regex, c.UserName, userInfo.Name)
	case ModeUserID:
		return c.UserID == userInfo.ID
	case ModeIsAdmin:
		if c.IsAdmin == nil {
			return false
		}
		return *c.IsAdmin == userInfo.Admin
	case ModeMessageTitle:
		return stringMatch(c.Regex, c.MessageTitle, msg.Msg.Title)
	case ModeMessageText:
		return stringMatch(c.Regex, c.MessageText, msg.Msg.Message)
	case ModeMessageExtra:
		return containExtra(c.Regex, c.MessageExtra, msg.Msg)
	case ModePriority, ModePriorityGt, ModePriorityLt:
		if c.MessagePriority == nil {
			return false
		}
		switch c.Mode {
		case ModePriority:
			return *c.MessagePriority == msg.Msg.Priority
		case ModePriorityGt:
			return *c.MessagePriority < msg.Msg.Priority
		case ModePriorityLt:
			return *c.MessagePriority > msg.Msg.Priority
		}
	}
	return false
}
