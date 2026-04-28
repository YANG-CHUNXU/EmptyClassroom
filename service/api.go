package service

import (
	"EmptyClassroom/logs"
	"EmptyClassroom/service/model"
	"EmptyClassroom/snapshot"
	"context"
)

const dataNotReadyMessage = "数据暂未准备好，请稍后刷新"

var queryAll = QueryAll

type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func GetDataResponse(ctx context.Context, store snapshot.Store) (APIResponse, int) {
	classInfo, err := ResolveClassInfo(ctx, store)
	if err != nil {
		logs.CtxWarn(ctx, "load class info snapshot failed: %v", err)
		return APIResponse{
			Code: 503,
			Msg:  dataNotReadyMessage,
			Data: nil,
		}, 503
	}

	return APIResponse{
		Code: 0,
		Data: classInfo,
	}, 200
}

func RefreshResponse(ctx context.Context, store snapshot.Store) (APIResponse, int) {
	classInfo, err := RefreshSnapshot(ctx, store)
	if err != nil {
		logs.CtxError(ctx, "refresh snapshot failed: %v", err)
		return APIResponse{
			Code: 500,
			Msg:  "refresh failed",
			Data: nil,
		}, 500
	}

	return APIResponse{
		Code: 0,
		Data: classInfo,
	}, 200
}

func ResolveClassInfo(ctx context.Context, store snapshot.Store) (*model.ClassInfo, error) {
	if store == nil {
		return nil, snapshot.ErrSnapshotNotFound
	}

	classInfo, err := store.Load(ctx)
	if err != nil {
		return nil, err
	}
	if classInfo.ConfigVersion != currentConfigVersion {
		logs.CtxWarn(ctx, "snapshot config version mismatch, serving stale snapshot")
		classInfo.IsStale = true
		classInfo.StaleReason = "当前展示的是缓存数据，正在等待后台刷新"
	}

	return classInfo, nil
}

func RefreshSnapshot(ctx context.Context, store snapshot.Store) (*model.ClassInfo, error) {
	classInfo, err := queryAll(ctx)
	if err != nil {
		return nil, err
	}

	if store != nil {
		if err := store.Save(ctx, classInfo); err != nil {
			return nil, err
		}
	}

	return classInfo, nil
}
