package service

import (
	"EmptyClassroom/service/model"
	"testing"
	"time"
)

func TestMergeClassroomCatalogBuildsCatalogFromCurrentSnapshot(t *testing.T) {
	observedAt := time.Date(2026, 4, 29, 8, 0, 0, 0, time.UTC)
	current := &model.ClassInfo{
		CampusInfoMap: map[string]*model.CampusInfo{
			"沙河": {
				Name: "沙河",
				BuildingInfoMap: map[int]*model.BuildingInfo{
					0: {
						Name: "N",
						ClassroomInfoMap: map[int]*model.ClassroomInfo{
							0: {
								Name:     "101",
								Size:     80,
								CanTrust: true,
								Type:     "普通教室",
							},
						},
					},
				},
			},
		},
	}

	MergeClassroomCatalog(nil, current, observedAt)

	entry := current.ClassroomCatalog["沙河"]["N"]["101"]
	if entry == nil {
		t.Fatal("catalog entry should be created")
	}
	if entry.Name != "N-101" {
		t.Fatalf("entry.Name = %q, want N-101", entry.Name)
	}
	if entry.Size != 80 {
		t.Fatalf("entry.Size = %d, want 80", entry.Size)
	}
	if entry.Source != "教务" {
		t.Fatalf("entry.Source = %q, want 教务", entry.Source)
	}
	if !entry.CanTrust {
		t.Fatal("entry.CanTrust should be true for realtime observations")
	}
	if !entry.FirstSeen.Equal(observedAt) || !entry.LastSeen.Equal(observedAt) {
		t.Fatalf("seen times = %v/%v, want %v", entry.FirstSeen, entry.LastSeen, observedAt)
	}
	if entry.SeenCount != 1 {
		t.Fatalf("entry.SeenCount = %d, want 1", entry.SeenCount)
	}
}

func TestMergeClassroomCatalogAvoidsDuplicateBuildingPrefixInDisplayName(t *testing.T) {
	observedAt := time.Date(2026, 4, 29, 8, 0, 0, 0, time.UTC)
	current := &model.ClassInfo{
		CampusInfoMap: map[string]*model.CampusInfo{
			"沙河": {
				Name: "沙河",
				BuildingInfoMap: map[int]*model.BuildingInfo{
					0: {
						Name: "N",
						ClassroomInfoMap: map[int]*model.ClassroomInfo{
							0: {
								Name:     "N104",
								Size:     80,
								CanTrust: true,
							},
						},
					},
				},
			},
		},
	}

	MergeClassroomCatalog(nil, current, observedAt)

	entry := current.ClassroomCatalog["沙河"]["N"]["N104"]
	if entry == nil {
		t.Fatal("catalog entry should be created")
	}
	if entry.Name != "N104" {
		t.Fatalf("entry.Name = %q, want N104", entry.Name)
	}
}

func TestMergeClassroomCatalogPreservesAndUpgradesPreviousCatalog(t *testing.T) {
	firstSeen := time.Date(2026, 4, 28, 8, 0, 0, 0, time.UTC)
	observedAt := time.Date(2026, 4, 29, 8, 0, 0, 0, time.UTC)
	previous := &model.ClassInfo{
		ClassroomCatalog: model.ClassroomCatalog{
			"沙河": {
				"N": {
					"101": {
						Name:          "N-101",
						BuildingName:  "N",
						ClassroomName: "101",
						Size:          60,
						CanTrust:      false,
						Type:          "旧类型",
						Source:        "课表",
						FirstSeen:     firstSeen,
						LastSeen:      firstSeen,
						SeenCount:     2,
					},
				},
			},
		},
	}
	current := &model.ClassInfo{
		CampusInfoMap: map[string]*model.CampusInfo{
			"沙河": {
				Name: "沙河",
				BuildingInfoMap: map[int]*model.BuildingInfo{
					0: {
						Name: "N",
						ClassroomInfoMap: map[int]*model.ClassroomInfo{
							0: {
								Name:     "101",
								Size:     80,
								CanTrust: true,
							},
						},
					},
				},
			},
		},
	}

	MergeClassroomCatalog(previous, current, observedAt)

	entry := current.ClassroomCatalog["沙河"]["N"]["101"]
	if entry.Size != 80 {
		t.Fatalf("entry.Size = %d, want trusted realtime size 80", entry.Size)
	}
	if entry.Source != "教务" {
		t.Fatalf("entry.Source = %q, want 教务", entry.Source)
	}
	if !entry.CanTrust {
		t.Fatal("entry.CanTrust should stay true after realtime observation")
	}
	if entry.Type != "旧类型" {
		t.Fatalf("entry.Type = %q, want previous non-empty type", entry.Type)
	}
	if !entry.FirstSeen.Equal(firstSeen) {
		t.Fatalf("entry.FirstSeen = %v, want %v", entry.FirstSeen, firstSeen)
	}
	if !entry.LastSeen.Equal(observedAt) {
		t.Fatalf("entry.LastSeen = %v, want %v", entry.LastSeen, observedAt)
	}
	if entry.SeenCount != 3 {
		t.Fatalf("entry.SeenCount = %d, want 3", entry.SeenCount)
	}
}
