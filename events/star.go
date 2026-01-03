package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shi-gg/githook/discord"
	"github.com/shi-gg/githook/utils"
	"github.com/google/go-github/v61/github"
)

func Star(w http.ResponseWriter, r *http.Request, url string) {
	var body github.StarEvent
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if *body.Action != "created" {
		return
	}

	discord.SendWebhook(
		url,
		discord.WebhookPayload{
			Username:  *body.Sender.Login,
			AvatarURL: *body.Sender.AvatarURL,
			Embeds: []discord.Embed{
				{
					Title: fmt.Sprintf("%s: Star added", *body.Repo.FullName),
					URL:   *body.Repo.HTMLURL,
					Description: fmt.Sprintf(
						"⭐ %s now has **%d star%s**",
						*body.Repo.Name,
						*body.Repo.StargazersCount,
						utils.Ternary(*body.Repo.StargazersCount > 1, "s", ""),
					),
					Color: utils.GetColors().Default,
				},
			},
		},
	)
}
