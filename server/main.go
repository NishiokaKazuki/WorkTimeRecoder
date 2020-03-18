package main

import (
	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(
		"MY TOKEN",
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			rtm.SendMessage(rtm.NewOutgoingMessage("ホリネズミです。塊茎(球根)食べたい", ev.Channel))
		}
	}
}
