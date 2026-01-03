package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shi-gg/githook/discord"
	"github.com/shi-gg/githook/utils"
	"github.com/google/go-github/v61/github"
	"github.com/redis/go-redis/v9"
)

func WorkflowRun(w http.ResponseWriter, r *http.Request, url string, client *redis.Client) {
	var body github.WorkflowRunEvent
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&body)

	if *body.Action == "requested" {
		err := client.Incr(r.Context(), "workflow:run:"+*body.WorkflowRun.HeadSHA).Err()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if *body.Action != "completed" {
		return
	}

	time.Sleep(2 * time.Second)

	num, err := client.Decr(r.Context(), "workflow:run:"+*body.WorkflowRun.HeadSHA).Result()
	if num != 0 || err != nil {
		return
	}

	keys := client.Keys(r.Context(), fmt.Sprintf("workflow:job:%s:*", *body.WorkflowRun.HeadSHA)).Val()
	defer client.Del(r.Context(), keys...)

	conclusion := "success"
	var desc strings.Builder

	for _, key := range keys {
		data, err := client.Get(r.Context(), key).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var job WorkflowJobCache
		err = json.Unmarshal([]byte(data), &job)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if job.Conclusions != "success" && job.Conclusions != "cancelled" && job.Conclusions != "skipped" {
			conclusion = job.Conclusions
		}

		fmt.Fprintf(&desc,
			"%s %s [↗︎](%s)\n",
			utils.Ternary(job.Conclusions == "success", "<:tick:1017781086102761543>", "<:cross:1017781065340964934>"),
			job.Name,
			fmt.Sprintf("https://github.com/%s/actions/runs/%s/job/%s", *body.Repo.FullName, job.RunID, job.ID),
		)
	}

	discord.SendWebhook(
		url,
		discord.WebhookPayload{
			Username:  *body.Sender.Login,
			AvatarURL: *body.Sender.AvatarURL,
			Embeds: []discord.Embed{
				{
					Title: fmt.Sprintf(
						"%s%s: Workflow %s",
						*body.Repo.FullName,
						utils.Ternary(
							*body.WorkflowRun.HeadBranch == "" || *body.WorkflowRun.HeadBranch == "master" || *body.WorkflowRun.HeadBranch == "main",
							"",
							"@"+*body.WorkflowRun.HeadBranch,
						),
						conclusion,
					),
					Description: desc.String(),
					URL:         *body.Repo.HTMLURL + "/actions",
					Color:       getColor(conclusion),
				},
			},
		},
	)
}

func getColor(conclusion string) int {
	switch conclusion {
	case "success":
		return utils.GetColors().Success
	case "failure":
		return utils.GetColors().Error
	default:
		return utils.GetColors().Default
	}
}
