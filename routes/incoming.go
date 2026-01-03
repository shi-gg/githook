package routes

import (
	"fmt"
	"net/http"

	"github.com/shi-gg/githook/config"
	"github.com/shi-gg/githook/events"
	"github.com/shi-gg/githook/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func HandleIncoming(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	conf := config.Get()
	githubEvent := r.Header.Get("X-Github-Event")

	id := r.PathValue("id")

	url, err := utils.Decrypt(id, []byte(conf.Secret))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if githubEvent == "" {
		http.Error(w, "Missing X-Github-Event header", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received GitHub event: %s\n", githubEvent)

	switch githubEvent {
	case "deployment_status":
		events.DeploymentStatus(w, r, url)
	case "ping":
		events.Ping(w, r, url)
	case "push":
		events.Push(w, r, url)
	case "star":
		events.Star(w, r, url)
	case "workflow_job":
		events.WorkflowJob(w, r, client)
	case "workflow_run":
		events.WorkflowRun(w, r, url, client)
	default:
		http.Error(w, fmt.Sprintf("Unsupported event: %s", githubEvent), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
