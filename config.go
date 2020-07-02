package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DEF_SERVER_DOMAIN = "localhost"
	DEF_SERVER_PORT   = 8065

	BOT_USERNAME         = "encourage-bot"
	BOT_USERID           = "fmh11pq9ei8g8fkkkxsxnnoj1y"
	DEF_BOT_ACCESS_TOKEN = "gmxeqimj8f8atjicuyabpo6f5r"

	DEF_TEAM_NAME    = "encourage-team"
	CHANNEL_LOG_NAME = "encourage-bot-log"

	TEAM_ID        = "fnrmpodiet8yzfhapdgnejh9nc"
	LOG_CHANNEL_ID = "f8z6efmy5t8wmx1dsr6goge9cw"
)

type Config struct {
	ServerUrl    string
	WebSocketUrl string

	BotUsername    string
	BotUserId      string
	BotAccessToken string

	TeamName       string
	LogChannelName string
}

func getConfig() Config {
	cfg := &Config{
		BotUsername:    BOT_USERNAME,
		BotUserId:      BOT_USERID,
		LogChannelName: CHANNEL_LOG_NAME,
	}

	var domain string
	var port int
	var https bool

	fs := flag.NewFlagSet(BOT_APP_NAME, flag.ContinueOnError)
	fs.StringVar(&domain, "domain", DEF_SERVER_DOMAIN, "IP/domain of Mattermost server")
	fs.IntVar(&port, "port", DEF_SERVER_PORT, "port of MatterMost server")
	fs.BoolVar(&https, "https", false, "use HTTPS to connect")
	fs.StringVar(&cfg.TeamName, "team", DEF_TEAM_NAME, "team name")
	fs.StringVar(&cfg.BotAccessToken, "bot_access_token", DEF_BOT_ACCESS_TOKEN, "bot access token")

	err := fs.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		fs.Usage()
		os.Exit(2)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fs.Usage()
		os.Exit(2)
	}

	cfg.ServerUrl = getServerUrl(domain, port, https)
	cfg.WebSocketUrl = getWebSocketUrl(domain, port, https)

	return *cfg
}

func getServerUrl(domain string, port int, https bool) string {
	protocol := "http"
	if https {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, domain, port)
}

func getWebSocketUrl(domain string, port int, https bool) string {
	protocol := "ws"
	if https {
		protocol = "wss"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, domain, port)
}
