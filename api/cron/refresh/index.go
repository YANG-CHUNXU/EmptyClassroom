package handler

import (
	"EmptyClassroom/bootstrap"
	"EmptyClassroom/httpapi"
	"EmptyClassroom/logs"
	"EmptyClassroom/service"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpapi.MethodNotAllowed(w)
		return
	}

	secret := os.Getenv("CRON_SECRET")
	if secret == "" || r.Header.Get("Authorization") != "Bearer "+secret {
		httpapi.Unauthorized(w, "unauthorized")
		return
	}

	bootstrap.Init(false)

	ctx := logs.WithLogID(r.Context())
	response, status := service.RefreshResponse(ctx, bootstrap.NewSnapshotStore())
	httpapi.WriteJSON(w, status, response)
}
