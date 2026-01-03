package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shi-gg/githook/discord"
	"github.com/shi-gg/githook/utils"
	"github.com/google/go-github/v61/github"
)

func Ping(w http.ResponseWriter, r *http.Request, url string) {
	var body github.PingEvent
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	discord.SendWebhook(
		url,
		discord.WebhookPayload{
			Username:  *body.Sender.Login,
			AvatarURL: *body.Sender.AvatarURL,
			Embeds: []discord.Embed{
				{
					Title:       fmt.Sprintf("%s: Ping", *body.Repo.FullName),
					URL:         *body.Repo.HTMLURL,
					Description: "🏓 Ping! Pong!",
					Color:       utils.GetColors().Default,
				},
			},
		},
	)
}
