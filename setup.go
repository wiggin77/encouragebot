package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type BotInfo struct {
	client       *model.Client4
	botUserId    string
	teamId       string
	logChannelId string
}

func setup(cfg Config) (*BotInfo, error) {
	client := model.NewAPIv4Client(cfg.ServerUrl)
	client.SetToken(cfg.BotAccessToken)

	if err := checkServerConnection(client, cfg.ServerUrl); err != nil {
		return nil, err
	}

	botInfo := &BotInfo{
		teamId:       TEAM_ID,
		client:       client,
		botUserId:    cfg.BotUserId,
		logChannelId: LOG_CHANNEL_ID,
	}

	return botInfo, nil
}

func checkServerConnection(client *model.Client4, url string) error {
	props, resp := client.GetOldClientConfig("")
	if resp.Error != nil {
		return fmt.Errorf("cannot connect to Mattermost server %s: %w", url, resp.Error)
	}
	fmt.Println("Server detected and is running version ", props["Version"])
	return nil
}
