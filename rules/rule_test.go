package rules

import (
	"fmt"
	"testing"

	"github.com/eternal-flame-AD/gotify-broadcast/model"
	plugin "github.com/gotify/plugin-api"

	. "github.com/smartystreets/goconvey/convey"
)

func shouldUseAction(actual interface{}, expected ...interface{}) string {
	return ShouldEqual(actual, expected...)
}

func shouldBeValidChain(actual interface{}, expected ...interface{}) string {
	actualRule := actual.(RuleChain)
	if err := actualRule.Check(); err != nil {
		return fmt.Sprintf("chain is not valid: %s", err.Error())
	}
	return ""
}

func shouldBeInvalidChain(actual interface{}, errExpected ...interface{}) string {
	actualRule := actual.(RuleChain)
	err := actualRule.Check()
	if err == nil {
		return "rule should not be valid"
	}
	for _, errItem := range errExpected {
		switch errItem := errItem.(type) {
		case int:
			found := false
			for _, err := range err.(RuleChainError).Errors {
				if err.Index == errItem {
					found = true
				}
			}
			if !found {
				return fmt.Sprintf("chain at index %d does not have error", errItem)
			}
		}
	}
	return ""
}

func TestChainCheck(t *testing.T) {
	Convey("Test Rule Chain Checking", t, func(c C) {
		testChain := RuleChain{
			Rule{
				Match: MatchSet{
					Match{
						Mode: ModeAny,
					},
				},
				Action: Accept,
			},
		}

		prependRule := func(item Rule) {
			testChain = append(RuleChain{item}, testChain...)
		}

		c.So(testChain, shouldBeValidChain)

		prependRule(Rule{
			Match: MatchSet{
				Match{
					Mode:   ModeUserID,
					UserID: 1,
				},
			},
			Action: Action("???"),
		})

		c.So(testChain, shouldBeInvalidChain, 0)

		prependRule(Rule{
			Match: MatchSet{
				Match{
					Mode: ModeIsAdmin,
				},
			},
			Action: Reject,
		})

		c.So(testChain, shouldBeInvalidChain, 0, 1)
	})
}

func TestChainMatch(t *testing.T) {
	Convey("Test Rule Chain Matching", t, func(c C) {
		testChain := RuleChain{}

		prependRule := func(item Rule) {
			testChain = append(RuleChain{item}, testChain...)
		}

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

		c.Convey("default action", func(c C) {
			c.So(testChain.Match(testMessage, Reject), shouldUseAction, Reject)
		})

		c.Convey("sender is admin", func(c C) {
			isAdmin := true
			prependRule(Rule{
				Match: MatchSet{
					Match{
						Mode:    ModeIsAdmin,
						IsAdmin: &isAdmin,
					},
				},
				Action: Accept,
			})
			c.So(testChain.Match(testMessage, Reject), shouldUseAction, Accept)
		})
		c.Convey("has extra", func(c C) {
			prependRule(Rule{
				Match: MatchSet{
					Match{
						Mode:         ModeMessageExtra,
						Regex:        true,
						MessageExtra: "test::*",
					},
				},
				Action: Reject,
			})
			c.So(testChain.Match(testMessage, Accept), shouldUseAction, Reject)
		})
		c.Convey("AND matching", func(c C) {
			testChain = RuleChain{}

			prependRule(Rule{
				Match: MatchSet{
					Match{
						Mode:         ModeMessageExtra,
						Regex:        true,
						MessageExtra: "test::*",
					},
					Match{
						Mode:         ModeMessageTitle,
						MessageTitle: "title",
					},
				},
				Action: Accept,
			})

			prependRule(Rule{
				Match: MatchSet{
					Match{
						Mode:         ModeMessageExtra,
						Regex:        true,
						MessageExtra: "test::*",
					},
					Match{
						Mode:         ModeMessageTitle,
						MessageTitle: "not_title",
					},
				},
				Action: Reject,
			})

			c.So(testChain.Match(testMessage, Reject), shouldUseAction, Accept)
		})
	})
}
