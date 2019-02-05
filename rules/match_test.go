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
	Convey("Test Match Matching", t, func() {
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
		Convey("empty rule should not panic", func() {
			So(func() {
				rule := Match{}
				rule.Match(testMessage)
			}, ShouldNotPanic)
		})
		Convey("any matching", func() {
			So(testMessage, shouldMatchRule, Match{
				Mode: ModeAny,
			})
		})
		Convey("channel matching", func() {
			So(testMessage, shouldMatchRule, Match{
				Mode:        ModeChannelName,
				ChannelName: "test_channel",
			})
			So(testMessage, shouldNotMatchRule, Match{
				Mode:        ModeChannelName,
				ChannelName: "test.channel",
			})
			So(testMessage, shouldMatchRule, Match{
				Mode:        ModeChannelName,
				Regex:       true,
				ChannelName: "test.channel",
			})
		})
		Convey("message matching", func() {
			Convey("match title", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "title",
				})
				So(testMessage, shouldNotMatchRule, Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "t...e",
				})
				So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageTitle,
					Regex:        true,
					MessageTitle: "t...e",
				})
			})
			Convey("match body", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:        ModeMessageText,
					MessageText: "message",
				})
				So(testMessage, shouldNotMatchRule, Match{
					Mode:        ModeMessageText,
					MessageText: "m.....e",
				})
				So(testMessage, shouldMatchRule, Match{
					Mode:        ModeMessageText,
					Regex:       true,
					MessageText: "m.....e",
				})
			})
			Convey("match extra", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "test::string",
				})
				So(testMessage, shouldNotMatchRule, Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "test::*",
				})
				So(testMessage, shouldMatchRule, Match{
					Mode:         ModeMessageExtra,
					Regex:        true,
					MessageExtra: "test::*",
				})
			})
			Convey("match priority", func() {
				targetPriority := 5
				So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriority,
					MessagePriority: &targetPriority,
				})
				targetPriority = 4
				So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriorityGt,
					MessagePriority: &targetPriority,
				})
				targetPriority = 6
				So(testMessage, shouldMatchRule, Match{
					Mode:            ModePriorityLt,
					MessagePriority: &targetPriority,
				})
				So(func() {
					rule := Match{
						Mode: ModePriority,
					}
					rule.Match(testMessage)
				}, ShouldNotPanic)
			})
		})
		Convey("sender rule matching", func() {
			Convey("match user id", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:   ModeUserID,
					UserID: 1,
				})
			})
			Convey("match user name", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					UserName: "sender",
				})
				So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					Regex:    true,
					UserName: "s....r",
				})
			})
			Convey("match is admin", func() {
				isAdmin := true
				So(testMessage, shouldMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				isAdmin = false
				So(testMessage, shouldNotMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				So(func() {
					rule := Match{
						Mode: ModeIsAdmin,
					}
					rule.Match(testMessage)
				}, ShouldNotPanic)
			})
		})
		testMessage.IsSend = true
		Convey("receiver rule matching", func() {
			Convey("match user id", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:   ModeUserID,
					UserID: 2,
				})
			})
			Convey("match user name", func() {
				So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					UserName: "receiver",
				})
				So(testMessage, shouldMatchRule, Match{
					Mode:     ModeUserName,
					Regex:    true,
					UserName: "re.....r",
				})
			})
			Convey("match is admin", func() {
				isAdmin := false
				So(testMessage, shouldMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				isAdmin = true
				So(testMessage, shouldNotMatchRule, Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				})
				So(func() {
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
	Convey("Test Check Match Syntax", t, func() {
		Convey("mode not valid", func() {
			rule := Match{
				Mode: "???",
			}
			So(rule, shouldBeInvalidRule, "mode")
		})
		Convey("channel name mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeChannelName,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:        ModeChannelName,
					ChannelName: "some_channel_name",
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:        ModeUserID,
					UserID:      1,
					ChannelName: "some_channel_name",
				}
				So(rule, shouldBeInvalidRule, ErrExtraParam{})
			})
		})
		Convey("user name mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeUserName,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:     ModeUserName,
					UserName: "some_user_name",
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:     ModeUserID,
					UserID:   1,
					UserName: "some_user_name",
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("user id mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeUserID,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:   ModeUserID,
					UserID: 1,
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:     ModeUserName,
					UserID:   1,
					UserName: "some_user_name",
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("is admin mode", func() {
			isAdmin := true
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeIsAdmin,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:    ModeIsAdmin,
					IsAdmin: &isAdmin,
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:    ModeUserID,
					UserID:  1,
					IsAdmin: &isAdmin,
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("message title mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeMessageTitle,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:         ModeMessageTitle,
					MessageTitle: "title",
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:         ModeMessageTitle,
					UserID:       1,
					MessageTitle: "title",
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("message text mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeMessageText,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:        ModeMessageText,
					MessageText: "text",
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:        ModeUserID,
					UserID:      1,
					MessageText: "text",
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("message extra mode", func() {
			Convey("missing field", func() {
				rule := Match{
					Mode: ModeMessageExtra,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:         ModeMessageExtra,
					MessageExtra: "text",
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:         ModeUserID,
					UserID:       1,
					MessageExtra: "text",
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("priority mode", func() {
			priority := 2
			Convey("missing field", func() {
				rule := Match{
					Mode: ModePriority,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:            ModePriority,
					MessagePriority: &priority,
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:            ModePriority,
					UserID:          1,
					MessagePriority: &priority,
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("priority_gt mode", func() {
			priority := 2
			Convey("missing field", func() {
				rule := Match{
					Mode: ModePriorityGt,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:            ModePriorityGt,
					MessagePriority: &priority,
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:            ModePriorityGt,
					UserID:          1,
					MessagePriority: &priority,
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
		Convey("priority_lt mode", func() {
			priority := 2
			Convey("missing field", func() {
				rule := Match{
					Mode: ModePriorityLt,
				}
				So(rule, shouldBeInvalidRule, ErrMissingParam{})
			})
			Convey("valid config", func() {
				rule := Match{
					Mode:            ModePriorityLt,
					MessagePriority: &priority,
				}
				So(rule, shouldBeValidRule)
			})
			Convey("extra field", func() {
				rule := Match{
					Mode:            ModePriorityLt,
					UserID:          1,
					MessagePriority: &priority,
				}
				So(rule, shouldBeInvalidRule, "extra")
			})
		})
	})
}
