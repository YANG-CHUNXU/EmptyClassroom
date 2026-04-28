package service

import (
	"EmptyClassroom/service/model"
	"EmptyClassroom/snapshot"
	"context"
	"errors"
	"testing"
)

type fakeSnapshotStore struct {
	loadClassInfo *model.ClassInfo
	loadErr       error
	saved         *model.ClassInfo
	saveErr       error
}

func (s *fakeSnapshotStore) Load(context.Context) (*model.ClassInfo, error) {
	return s.loadClassInfo, s.loadErr
}

func (s *fakeSnapshotStore) Save(_ context.Context, classInfo *model.ClassInfo) error {
	s.saved = classInfo
	return s.saveErr
}

func withQueryAllStub(t *testing.T, fn func(context.Context) (*model.ClassInfo, error)) {
	t.Helper()

	previous := queryAll
	queryAll = fn
	t.Cleanup(func() {
		queryAll = previous
	})
}

func TestGetDataResponseLoadsMatchingSnapshotWithoutRefresh(t *testing.T) {
	calledQueryAll := false
	withQueryAllStub(t, func(context.Context) (*model.ClassInfo, error) {
		calledQueryAll = true
		return nil, errors.New("queryAll should not be called")
	})

	want := &model.ClassInfo{ConfigVersion: currentConfigVersion}
	response, status := GetDataResponse(context.Background(), &fakeSnapshotStore{loadClassInfo: want})

	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	if response.Code != 0 {
		t.Fatalf("response.Code = %d, want 0", response.Code)
	}
	if response.Data != want {
		t.Fatalf("response.Data = %#v, want snapshot", response.Data)
	}
	if calledQueryAll {
		t.Fatal("GetDataResponse should not refresh realtime data")
	}
}

func TestGetDataResponseReturnsStaleSnapshotOnVersionMismatch(t *testing.T) {
	calledQueryAll := false
	withQueryAllStub(t, func(context.Context) (*model.ClassInfo, error) {
		calledQueryAll = true
		return nil, errors.New("queryAll should not be called")
	})

	stale := &model.ClassInfo{ConfigVersion: "old-version"}
	response, status := GetDataResponse(context.Background(), &fakeSnapshotStore{loadClassInfo: stale})

	if status != 200 {
		t.Fatalf("status = %d, want 200", status)
	}
	got, ok := response.Data.(*model.ClassInfo)
	if !ok {
		t.Fatalf("response.Data type = %T, want *model.ClassInfo", response.Data)
	}
	if !got.IsStale {
		t.Fatal("stale snapshot should be marked stale")
	}
	if got.StaleReason == "" {
		t.Fatal("stale snapshot should include a stale reason")
	}
	if calledQueryAll {
		t.Fatal("GetDataResponse should not refresh realtime data for stale snapshots")
	}
}

func TestGetDataResponseReturnsUnavailableWhenSnapshotMissing(t *testing.T) {
	calledQueryAll := false
	withQueryAllStub(t, func(context.Context) (*model.ClassInfo, error) {
		calledQueryAll = true
		return nil, errors.New("queryAll should not be called")
	})

	response, status := GetDataResponse(
		context.Background(),
		&fakeSnapshotStore{loadErr: snapshot.ErrSnapshotNotFound},
	)

	if status != 503 {
		t.Fatalf("status = %d, want 503", status)
	}
	if response.Code != 503 {
		t.Fatalf("response.Code = %d, want 503", response.Code)
	}
	if response.Msg != dataNotReadyMessage {
		t.Fatalf("response.Msg = %q, want %q", response.Msg, dataNotReadyMessage)
	}
	if response.Data != nil {
		t.Fatalf("response.Data = %#v, want nil", response.Data)
	}
	if calledQueryAll {
		t.Fatal("GetDataResponse should not refresh realtime data when snapshot is missing")
	}
}

func TestRefreshSnapshotQueriesRealtimeAndSaves(t *testing.T) {
	want := &model.ClassInfo{ConfigVersion: currentConfigVersion}
	calledQueryAll := false
	withQueryAllStub(t, func(context.Context) (*model.ClassInfo, error) {
		calledQueryAll = true
		return want, nil
	})

	store := &fakeSnapshotStore{}
	got, err := RefreshSnapshot(context.Background(), store)
	if err != nil {
		t.Fatalf("RefreshSnapshot() error = %v", err)
	}
	if !calledQueryAll {
		t.Fatal("RefreshSnapshot should query realtime data")
	}
	if got != want {
		t.Fatalf("RefreshSnapshot() = %#v, want query result", got)
	}
	if store.saved != want {
		t.Fatalf("saved snapshot = %#v, want query result", store.saved)
	}
}
