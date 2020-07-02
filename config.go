package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DEF_SERVER_DOMAIN = "localhost"
	DEF_SERVER_PORT   = 8065

	DEF_ADMIN_USERNAME = "dlauder"
	DEF_ADMIN_PASSWORD = "hello"

	//BOT_USER_EMAIL    = "bot@example.com"
	//BOT_USER_PASSWORD = "password1"
	//BOT_USER_USERNAME = "samplebot"
	//BOT_USER_FIRST    = "Sample"
	//BOT_USER_LAST     = "Bot"

	BOT_USERNAME         = "encouragebot"
	BOT_DISPLAYNAME      = "Encourage Bot"
	BOT_DESCRIPTION      = "Encourage Bot is encouraging."
	DEF_BOT_ACCESS_TOKEN = "mo1bx4u6k3rx3bbmc5as1y8zsc"

	DEF_TEAM_NAME    = "encourage-team"
	CHANNEL_LOG_NAME = "debugging-for-encourage-bot"
)

type Config struct {
	ServerUrl    string
	WebSocketUrl string

	AdminUsername string
	AdminPassword string

	BotUsername    string
	BotDisplayname string
	BotDescription string
	BotAccessToken string

	TeamName       string
	LogChannelName string
}

func getConfig() Config {
	cfg := &Config{
		BotUsername:    BOT_USERNAME,
		BotDisplayname: BOT_DISPLAYNAME,
		BotDescription: BOT_DESCRIPTION,
		BotAccessToken: DEF_BOT_ACCESS_TOKEN,
		LogChannelName: CHANNEL_LOG_NAME,
	}

	var domain string
	var port int
	var https bool

	fs := flag.NewFlagSet(BOT_APP_NAME, flag.ContinueOnError)
	fs.StringVar(&domain, "domain", DEF_SERVER_DOMAIN, "IP/domain of Mattermost server")
	fs.IntVar(&port, "port", DEF_SERVER_PORT, "port of MatterMost server")
	fs.BoolVar(&https, "https", false, "use HTTPS to connect")
	fs.StringVar(&cfg.AdminUsername, "admin_username", DEF_ADMIN_USERNAME, "username of admin user")
	fs.StringVar(&cfg.AdminPassword, "admin_password", DEF_ADMIN_PASSWORD, "password of admin user")
	fs.StringVar(&cfg.TeamName, "team", DEF_TEAM_NAME, "team name")
	fs.StringVar(&cfg.BotAccessToken, "bot_access_token", "", "bot access token")

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
