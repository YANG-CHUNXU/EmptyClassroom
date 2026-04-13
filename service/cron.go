package service

import (
	"EmptyClassroom/logs"
	"EmptyClassroom/snapshot"
	"context"
)

func RunRefresh(ctx context.Context, store snapshot.Store) error {
	_, err := RefreshSnapshot(ctx, store)
	if err != nil {
		logs.CtxError(ctx, "QueryAll error: %v", err)
		return err
	}
	logs.CtxInfo(ctx, "QueryAll success")
	return nil
}
