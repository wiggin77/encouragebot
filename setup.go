package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type BotInfo struct {
	client *model.Client4

	botUser *model.User

	team       *model.Team
	logChannel *model.Channel
}

func setup(cfg Config) (*BotInfo, error) {
	clientAdmin := model.NewAPIv4Client(cfg.ServerUrl)

	if err := checkServerConnection(clientAdmin, cfg.ServerUrl); err != nil {
		return nil, err
	}

	admin, err := login(clientAdmin, cfg.AdminUsername, cfg.AdminPassword)
	if err != nil {
		return nil, err
	}

	team, err := createTeamIfNeeded(clientAdmin, cfg.TeamName)
	if err != nil {
		return nil, err
	}

	err = addUserToTeamIfNeeded(clientAdmin, admin.Id, team.Id)
	if err != nil {
		return nil, err
	}

	bot, err := createBotIfNeeded(clientAdmin, team, cfg)
	if err != nil {
		return nil, err
	}

	err = addUserToTeamIfNeeded(clientAdmin, bot.UserId, team.Id)
	if err != nil {
		return nil, err
	}

	channel, err := createChannelIfNeeded(clientAdmin, cfg.LogChannelName, team.Id)
	if err != nil {
		return nil, err
	}

	client := model.NewAPIv4Client(cfg.ServerUrl)
	user, err := login(client, bot.Username, cfg.BotAccessToken)
	if err != nil {
		return nil, err
	}

	botInfo := &BotInfo{
		team:       team,
		botUser:    user,
		client:     client,
		logChannel: channel,
	}

	return botInfo, nil
}

func checkServerConnection(clientAdmin *model.Client4, url string) error {
	props, resp := clientAdmin.GetOldClientConfig("")
	if resp.Error != nil {
		return fmt.Errorf("cannot connect to Mattermost server %s: %w", url, resp.Error)
	}
	fmt.Println("Server detected and is running version ", props["Version"])
	return nil
}

func login(clientAdmin *model.Client4, username string, password string) (*model.User, error) {
	admin, resp := clientAdmin.Login(username, password)
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to login to Mattermost server as %s: %w", username, resp.Error)
	}
	return admin, nil
}

func createTeamIfNeeded(clientAdmin *model.Client4, teamName string) (*model.Team, error) {
	team, resp := clientAdmin.GetTeamByName(teamName, "")
	if resp.Error != nil {
		team = &model.Team{
			Name:            teamName,
			DisplayName:     "Encourage Bot team",
			Description:     "A very encouraging team",
			AllowOpenInvite: true,
			CompanyName:     "Mattermost",
			Type:            model.TEAM_OPEN,
		}
		team, resp = clientAdmin.CreateTeam(team)
		if resp.Error != nil {
			return nil, fmt.Errorf("failed to create team: %w", resp.Error)
		}
	}

	return team, nil
}

func createBotIfNeeded(clientAdmin *model.Client4, team *model.Team, cfg Config) (*model.Bot, error) {
	var bot *model.Bot
	var resp *model.Response

	bot, err := getBotByUsername(clientAdmin, cfg.BotUsername)
	if err != nil {
		bot = &model.Bot{
			Description: cfg.BotDescription,
			DisplayName: cfg.BotDisplayname,
			Username:    cfg.BotUsername,
		}

		bot, resp = clientAdmin.CreateBot(bot)
		if resp.Error != nil {
			return nil, fmt.Errorf("failed to create bot: %w", resp.Error)
		}
	}
	return bot, nil
}

func addUserToTeamIfNeeded(clientAdmin *model.Client4, userId string, teamId string) error {
	_, resp := clientAdmin.GetTeamMember(teamId, userId, "")
	if resp.Error != nil {
		_, resp = clientAdmin.AddTeamMember(teamId, userId)
		if resp.Error != nil {
			return fmt.Errorf("failed to add user %s to team %s: %w", userId, teamId, resp.Error)
		}
	}
	return nil
}

func createChannelIfNeeded(clientAdmin *model.Client4, name string, teamId string) (*model.Channel, error) {
	channel, resp := clientAdmin.GetChannelByName(name, teamId, "")
	if resp.Error != nil {
		channel = &model.Channel{
			Name:        name,
			DisplayName: name,
			Type:        model.CHANNEL_OPEN,
		}
		channel, resp = clientAdmin.CreateChannel(channel)
		if resp.Error != nil {
			return nil, fmt.Errorf("failed to create channel %s: %w", name, resp.Error)
		}
	}
	return channel, nil
}

func getBotByUsername(clientAdmin *model.Client4, username string) (*model.Bot, error) {
	bots, resp := clientAdmin.GetBots(0, 10000, "")
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to get list of bots: %w", resp.Error)
	}

	for _, bot := range bots {
		if bot.Username == username {
			return bot, nil
		}
	}
	return nil, fmt.Errorf("failed to get bot %s", username)
}
