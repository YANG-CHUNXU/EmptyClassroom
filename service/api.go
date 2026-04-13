package service

import (
	"EmptyClassroom/logs"
	"EmptyClassroom/service/model"
	"EmptyClassroom/snapshot"
	"context"
	"errors"
)

type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func GetDataResponse(ctx context.Context, store snapshot.Store) (APIResponse, int) {
	classInfo, err := ResolveClassInfo(ctx, store)
	if err != nil {
		logs.CtxError(ctx, "resolve class info failed: %v", err)
		return APIResponse{
			Code: 500,
			Msg:  "query failed",
			Data: nil,
		}, 500
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
	if store != nil {
		classInfo, err := store.Load(ctx)
		if err == nil {
			if classInfo.ConfigVersion == currentConfigVersion {
				return classInfo, nil
			}
			logs.CtxInfo(ctx, "snapshot config version mismatch, refreshing snapshot")
		}
		if err != nil && !errors.Is(err, snapshot.ErrSnapshotNotFound) {
			logs.CtxWarn(ctx, "load snapshot failed, falling back to refresh: %v", err)
		}
	}

	classInfo, err := QueryAll(ctx)
	if err != nil {
		return nil, err
	}

	if store != nil {
		if err := store.Save(ctx, classInfo); err != nil {
			logs.CtxWarn(ctx, "save snapshot failed after refresh: %v", err)
		}
	}

	return classInfo, nil
}

func RefreshSnapshot(ctx context.Context, store snapshot.Store) (*model.ClassInfo, error) {
	classInfo, err := QueryAll(ctx)
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
