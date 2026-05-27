package service

import (
	"EmptyClassroom/service/model"
	"math"
	"testing"
	"time"
)

func TestBuildIdleRateHistoryEntryUsesCatalogAsDenominator(t *testing.T) {
	classInfo := &model.ClassInfo{
		ClassroomCatalog: model.ClassroomCatalog{
			"沙河": {
				"N": {
					"101": {ClassroomName: "101"},
					"102": {ClassroomName: "102"},
					"103": {ClassroomName: "103"},
				},
			},
		},
		CampusInfoMap: map[string]*model.CampusInfo{
			"沙河": {
				Name:          "沙河",
				BuildingIdMap: map[string]int{"N": 0},
				BuildingInfoMap: map[int]*model.BuildingInfo{
					0: {
						Name:           "N",
						ClassroomIdMap: map[string]int{"101": 0, "102": 1},
						ClassMatrix: [][]int{
							{0, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
							{1, 1},
						},
					},
				},
			},
		},
	}

	entry := buildIdleRateHistoryEntry(classInfo, time.Date(2026, 5, 27, 8, 0, 0, 0, time.UTC))
	building := entry.Campus["沙河"]["N"]

	if building.KnownClassrooms != 3 {
		t.Fatalf("KnownClassrooms = %d, want 3", building.KnownClassrooms)
	}
	if len(building.Periods) != 14 {
		t.Fatalf("len(Periods) = %d, want 14", len(building.Periods))
	}
	period := building.Periods[0]
	if period.EmptyClassrooms != 1 {
		t.Fatalf("EmptyClassrooms = %d, want 1", period.EmptyClassrooms)
	}
	if period.TotalClassrooms != 3 {
		t.Fatalf("TotalClassrooms = %d, want 3", period.TotalClassrooms)
	}
	if !almostEqual(period.IdleRate, 1.0/3.0) {
		t.Fatalf("IdleRate = %f, want %f", period.IdleRate, 1.0/3.0)
	}
}

func TestBuildIdleRateHistoryEntryUsesShanghaiDate(t *testing.T) {
	classInfo := &model.ClassInfo{}
	observedAt := time.Date(2026, 5, 26, 16, 30, 0, 0, time.UTC)

	entry := buildIdleRateHistoryEntry(classInfo, observedAt)

	if entry.Date != "2026-05-27" {
		t.Fatalf("Date = %q, want 2026-05-27", entry.Date)
	}
}

func TestUpsertIdleRateHistoryReplacesSameDate(t *testing.T) {
	previous := []model.IdleRateHistoryEntry{
		{Date: "2026-05-26"},
		{
			Date: "2026-05-27",
			Campus: map[string]model.IdleRateCampusStats{
				"沙河": {
					"N": {KnownClassrooms: 1},
				},
			},
		},
	}
	current := model.IdleRateHistoryEntry{
		Date: "2026-05-27",
		Campus: map[string]model.IdleRateCampusStats{
			"沙河": {
				"N": {KnownClassrooms: 2},
			},
		},
	}

	got := upsertIdleRateHistory(previous, current)

	if len(got) != 2 {
		t.Fatalf("len(history) = %d, want 2", len(got))
	}
	if got[1].Campus["沙河"]["N"].KnownClassrooms != 2 {
		t.Fatalf("same date entry was not replaced: %#v", got[1])
	}
}

func TestUpsertIdleRateHistoryKeepsLatestSevenDays(t *testing.T) {
	previous := []model.IdleRateHistoryEntry{
		{Date: "2026-05-20"},
		{Date: "2026-05-21"},
		{Date: "2026-05-22"},
		{Date: "2026-05-23"},
		{Date: "2026-05-24"},
		{Date: "2026-05-25"},
		{Date: "2026-05-26"},
	}
	current := model.IdleRateHistoryEntry{Date: "2026-05-27"}

	got := upsertIdleRateHistory(previous, current)

	if len(got) != 7 {
		t.Fatalf("len(history) = %d, want 7", len(got))
	}
	if got[0].Date != "2026-05-21" {
		t.Fatalf("first Date = %q, want 2026-05-21", got[0].Date)
	}
	if got[6].Date != "2026-05-27" {
		t.Fatalf("last Date = %q, want 2026-05-27", got[6].Date)
	}
}

func almostEqual(a float64, b float64) bool {
	return math.Abs(a-b) < 0.000001
}
