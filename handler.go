//go:build localserver

package main

import (
	"EmptyClassroom/bootstrap"
	"EmptyClassroom/logs"
	"EmptyClassroom/service"
	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	ctx := logs.GetContextFromGinContext(c)
	logs.CtxInfo(ctx, "GetData")
	response, status := service.GetDataResponse(ctx, bootstrap.NewSnapshotStore())
	c.JSON(status, response)
}
