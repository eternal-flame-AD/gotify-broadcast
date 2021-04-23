package rules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gotify/plugin-api"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
	. "github.com/smartystreets/goconvey/convey"
)

func shouldBeValidRule(actual interface{}, expected ...interface{}) string {
	actualRule := actual.(Match)
	if err := actualRule.Check(); err != nil {
		return err.Error()
	}
	return ""
}

func shouldBeInvalidRule(actual interface{}, errExpected ...interface{}) string {
	actualRule := actual.(Match)
	err := actualRule.Check()
	if err == nil {
		return "rule should not be valid"
	}
	for _, expectedErr := range errExpected {
		switch expectedErr := expectedErr.(type) {
		case string:
			if !strings.Contains(err.Error(), expectedErr) {
				return fmt.Sprintf("error does not contain %s", expectedErr)
			}
		case error:
			if err := ShouldHaveSameTypeAs(err, expectedErr); err != "" {
				return err
			}
		}
	}
	return ""
}

func shouldNotMatchRule(actual interface{}, rules ...interface{}) string {
	msg := actual.(model.Message)
	for _, rule := range rules {
		if rule.(Match).Match(msg) {
			return fmt.Sprintf("Message %+v should not match %+v", msg, rule)
		}
	}
	return ""
}

func shouldMatchRule(actual interface{}, rules ...interface{}) string {
	msg := actual.(model.Message)
	for _, rule := range rules {
		if !rule.(Match).Match(msg) {
			return fmt.Sprintf("Message %+v should match %+v", msg, rule)
		}
	}
	return ""
}

func TestMatchMatch(t *testing.T) {
	Convey("Test Match Matching", t, func(c C) {
		testMessage := model.Message{
			Sender: plugin.UserContext{
				ID:    1,
				Name:  "sender",
				Admin: true,
			},
			Receiver: plugin.UserContext{
				ID:    2,
				Name:  "receiver",
				Admin: false,
			},
			Msg: plugin.Message{
				Title:   "title",
				Message: "message",
				Extras: map[string]interface{}{
					"test::string": "string",
				},
				Priority: 5,
			},
			IsSend:      false,
			ChannelName: "test_channel",
		}
		c.Convey("empty rule should not panic", func(c C) {
			c.So(func() {
				rule := Match{}
				rule.Match(testMessage)
			}, ShouldNotPanic)
		})
		c.Convey("any matching", func(c C) {
			c.So(testMessage, shouldMatchRule, Match{
				Mode: ModeAny,
			})
		})
		c.Convey("channel matching", func(c C) {
			c.So(testMessage, shouldMatchRule, Match{
				Mode:        ModeChannelName,
				ChannelName: "test_channel",
			})
			c.So(testMessage, shouldNotMatchRule, Match{
				Mode:        ModeChannelName,
				ChannelName: "test.channel",
			})
			c.So(testMessage, shouldMatchRule, Match{
				Mode:        ModeChannelName,
				Regex:       true,
				ChannelName: "test.channel",
			})
		})
		c.Convey("message matching", func(c C) {
			c.Convey("match title", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "title",
				})
				c.So(testMessage, shouldNotMatchRule, Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "t...e",
				})
				c.So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageTitle,
					Regex:        true,
					MessageTitle: "t...e",
				})
			})
			c.Convey("match body", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:        ModeMessageText,
					MessageText: "message",
				})
				c.So(testMessage, shouldNotMatchRule, Match{
					Mode:        ModeMessageText,
					MessageText: "m.....e",
				})
				c.So(testMessage, shouldMatchRule, Match{
					Mode:        ModeMessageText,
					Regex:       true,
					MessageText: "m.....e",
				})
			})
			c.Convey("match extra", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "test::string",
				})
				c.So(testMessage, shouldNotMatchRule, Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "test::*",
				})
				c.So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageExtra,
					Regex:        true,
					MessageExtra: "test::*",
				})
			})
			c.Convey("match priority", func(c C) {
				targetPriority := 5
				c.So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriority,
					MessagePriority: &targetPriority,
				})
				targetPriority = 4
				c.So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriorityGt,
					MessagePriority: &targetPriority,
				})
				targetPriority = 6
				c.So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriorityLt,
					MessagePriority: &targetPriority,
				})
				c.So(func() {
					rule := Match{
						Mode: ModePriority,
					}
					rule.Match(testMessage)
				}, ShouldNotPanic)
			})
		})
		c.Convey("sender rule matching", func(c C) {
			c.Convey("match user id", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:   ModeUserID,
					UserID: 1,
				})
			})
			c.Convey("match user name", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					UserName: "sender",
				})
				c.So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					Regex:    true,
					UserName: "s....r",
				})
			})
			c.Convey("match is admin", func(c C) {
				isAdmin := true
				c.So(testMessage, shouldMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				isAdmin = false
				c.So(testMessage, shouldNotMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				c.So(func() {
					rule := Match{
						Mode: ModeIsAdmin,
					}
					rule.Match(testMessage)
				}, ShouldNotPanic)
			})
		})
		testMessage.IsSend = true
		c.Convey("receiver rule matching", func(c C) {
			c.Convey("match user id", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:   ModeUserID,
					UserID: 2,
				})
			})
			c.Convey("match user name", func(c C) {
				c.So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					UserName: "receiver",
				})
				c.So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					Regex:    true,
					UserName: "re.....r",
				})
			})
			c.Convey("match is admin", func(c C) {
				isAdmin := false
				c.So(testMessage, shouldMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				isAdmin = true
				c.So(testMessage, shouldNotMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				c.So(func() {
					rule := Match{
						Mode: ModeIsAdmin,
					}
					rule.Match(testMessage)
				}, ShouldNotPanic)
			})
		})
	})
}

func TestMatchCheck(t *testing.T) {
	Convey("Test Check Match Syntax", t, func(c C) {
		c.Convey("mode not valid", func(c C) {
			rule := Match{
				Mode: "???",
			}
			c.So(rule, shouldBeInvalidRule, "mode")
		})
		c.Convey("channel name mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeChannelName,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:        ModeChannelName,
					ChannelName: "some_channel_name",
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:        ModeUserID,
					UserID:      1,
					ChannelName: "some_channel_name",
				}
				c.So(rule, shouldBeInvalidRule, ErrExtraParam{})
			})
		})
		c.Convey("user name mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeUserName,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:     ModeUserName,
					UserName: "some_user_name",
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:     ModeUserID,
					UserID:   1,
					UserName: "some_user_name",
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("user id mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeUserID,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:   ModeUserID,
					UserID: 1,
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:     ModeUserName,
					UserID:   1,
					UserName: "some_user_name",
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("is admin mode", func(c C) {
			isAdmin := true
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeIsAdmin,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:    ModeUserID,
					UserID:  1,
					IsAdmin: &isAdmin,
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("message title mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeMessageTitle,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "title",
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:         ModeMessageTitle,
					UserID:       1,
					MessageTitle: "title",
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("message text mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeMessageText,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:        ModeMessageText,
					MessageText: "text",
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:        ModeUserID,
					UserID:      1,
					MessageText: "text",
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("message extra mode", func(c C) {
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModeMessageExtra,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "text",
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:         ModeUserID,
					UserID:       1,
					MessageExtra: "text",
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("priority mode", func(c C) {
			priority := 2
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModePriority,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:            ModePriority,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:            ModePriority,
					UserID:          1,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("priority_gt mode", func(c C) {
			priority := 2
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModePriorityGt,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:            ModePriorityGt,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:            ModePriorityGt,
					UserID:          1,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
		c.Convey("priority_lt mode", func(c C) {
			priority := 2
			c.Convey("missing field", func(c C) {
				rule := Match{
					Mode: ModePriorityLt,
				}
				c.So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			c.Convey("valid config", func(c C) {
				rule := Match{
					Mode:            ModePriorityLt,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeValidRule)
			})
			c.Convey("extra field", func(c C) {
				rule := Match{
					Mode:            ModePriorityLt,
					UserID:          1,
					MessagePriority: &priority,
				}
				c.So(rule, shouldBeInvalidRule, "extra")
			})
		})
	})
}
