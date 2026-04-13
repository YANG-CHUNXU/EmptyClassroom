package handler

import (
	"EmptyClassroom/bootstrap"
	"EmptyClassroom/httpapi"
	"EmptyClassroom/logs"
	"EmptyClassroom/service"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpapi.MethodNotAllowed(w)
		return
	}

	bootstrap.Init(false)

	ctx := logs.GenNewContext()
	response, status := service.GetDataResponse(ctx, bootstrap.NewSnapshotStore())
	httpapi.WriteJSON(w, status, response)
}
