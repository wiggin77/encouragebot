// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	BOT_APP_NAME = "Mattermost Encourage Bot"
)

func main() {
	fmt.Println(BOT_APP_NAME)

	cfg := getConfig()
	botInfo, err := setup(cfg)
	if err != nil {
		printErrorAndExit(err, 2)
	}

	msg := fmt.Sprintf("_%s has **started** running_", BOT_APP_NAME)
	err = sendMsgToChannel(msg, botInfo.logChannelId, "", botInfo)
	if err != nil {
		printErrorAndExit(err, 2)
	}

	// Start listening to some channels via the websocket!
	webSocketClient, appErr := model.NewWebSocketClient4(cfg.WebSocketUrl, botInfo.client.AuthToken)
	if appErr != nil {
		printErrorAndExit(fmt.Errorf("failed to connect to websocket: %w", appErr), 2)
	}

	setupGracefulShutdown(webSocketClient, botInfo)
	webSocketClient.Listen()

	go func() {
		for resp := range webSocketClient.EventChannel {
			handleWebSocketResponse(resp, botInfo)
		}
	}()

	// You can block forever with
	select {}
}

func sendMsgToChannel(msg string, channelId, replyId string, botInfo *BotInfo) error {
	post := &model.Post{}
	post.ChannelId = channelId
	post.RootId = replyId
	post.Message = msg

	if _, resp := botInfo.client.CreatePost(post); resp.Error != nil {
		return resp.Error
	}
	return nil
}

func handleWebSocketResponse(event *model.WebSocketEvent, botInfo *BotInfo) {
	// Only respond to message posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	fmt.Println("responding to post")

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		// ignore my events
		if post.UserId == botInfo.botUserId {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			_ = sendMsgToChannel("Yes I'm running", post.ChannelId, post.Id, botInfo)
			return
		}

		// if you see any word matching 'up' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			_ = sendMsgToChannel("Yes I'm running", post.ChannelId, post.Id, botInfo)
			return
		}

		// if you see any word matching 'running' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			_ = sendMsgToChannel("Yes I'm running", post.ChannelId, post.Id, botInfo)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello(?:$|\W)`, post.Message); matched {
			_ = sendMsgToChannel("Yes I'm running", post.ChannelId, post.Id, botInfo)
			return
		}
	}

	msg := getEncouragement()
	_ = sendMsgToChannel(msg, post.ChannelId, "", botInfo)
}

func setupGracefulShutdown(wsc *model.WebSocketClient, botInfo *BotInfo) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if wsc != nil {
				wsc.Close()
			}

			msg := fmt.Sprintf("_%s has **stopped** running_", BOT_APP_NAME)
			_ = sendMsgToChannel(msg, botInfo.logChannelId, "", botInfo)
			os.Exit(0)
		}
	}()
}

func printErrorAndExit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitCode)
}
