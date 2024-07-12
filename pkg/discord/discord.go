package discord

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cakramediadata2022/chs_cloud_general/config"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var mainConfig *config.Config
var token = ""
var webhookClient = webhook.New(snowflake.ID(1095914000157515836), "",
	// webhook.WithLogger(logrus.New()),
	webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
		RepliedUser: false,
	}),
)
var client bot.Client
var channelAPI = 1095912772543463426
var channelLOG = 1141241190893756456

func UpdateChannelAPIStatus(icon string, status string) {
	// webhookClient.CreateMessage(discord.WebhookMessageCreate{
	// 	Content:   "API Service: " + icon + status,
	// 	AvatarURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTZZNHW7MnfjqYKCf_N5K6D54pBsPCAMx92DQ&usqp=CAU",
	// })
	statusX := icon + "-api-status"
	_, err := client.Rest().UpdateChannel(snowflake.ID(channelAPI), discord.GuildTextChannelUpdate{
		Name: &statusX,
	})
	if err != nil {
		log.Info(err.Error())
	}
	client.Rest().CreateMessage(snowflake.ID(channelAPI), discord.NewMessageCreateBuilder().SetContent("API Service: "+icon+status).Build())
}

func SendDiscordLog(title string, msg string) {
	if mainConfig.Discord.Run {
		client.Rest().CreateMessage(snowflake.ID(channelLOG), discord.NewMessageCreateBuilder().SetContentf("**LOG**: %s\n**MESSAGE**: %s", title, msg).Build())
	}
}

func Init(ConfigX *config.Config) {
	mainConfig = ConfigX
	if !ConfigX.Discord.Run {
		return
	}
	var err error
	client, err = disgo.New((token),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuilds,
			),
		),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	UpdateChannelAPIStatus("\\🟢", " UP")
	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("discord client is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}
	var message string
	if event.Message.Content == "ping" {
		message = "pong"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}
	if message != "" {
		_, _ = event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}
