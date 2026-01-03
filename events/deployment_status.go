package events

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shi-gg/githook/discord"
	"github.com/shi-gg/githook/utils"
	"github.com/google/go-github/v61/github"
)

func DeploymentStatus(w http.ResponseWriter, r *http.Request, url string) {
	var body github.DeploymentStatusEvent
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	discord.SendWebhook(
		url,
		discord.WebhookPayload{
			Username:  *body.Sender.Login,
			AvatarURL: *body.Sender.AvatarURL,
			Embeds: []discord.Embed{
				{
					Title: fmt.Sprintf(
						"%s: %s deployment %s",
						*body.Repo.FullName,
						*body.Deployment.Environment,
						*body.DeploymentStatus.State,
					),
					URL: *body.Repo.HTMLURL,
					Color: utils.Ternary(
						*body.DeploymentStatus.State == "success",
						utils.GetColors().Success,
						utils.GetColors().Error,
					).(int),
				},
			},
			Components: []discord.ActionRowComponent{
				{
					Type: discord.ComponentType(1),
					Componments: []discord.ButtonComponent{
						{
							Type:     discord.ComponentType(2),
							Style:    discord.ButtonStyle(5),
							Label:    "View Deployment",
							URL:      *body.DeploymentStatus.EnvironmentURL,
							Disabled: false,
						},
					},
				},
			},
		},
	)
}
