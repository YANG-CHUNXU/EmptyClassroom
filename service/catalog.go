package service

import (
	"EmptyClassroom/service/model"
	"strings"
	"time"
)

const (
	catalogSourceRealtime   = "教务"
	catalogSourceClassTable = "课表"
)

func MergeClassroomCatalog(previous *model.ClassInfo, current *model.ClassInfo, observedAt time.Time) {
	if current == nil {
		return
	}

	catalog := model.ClassroomCatalog{}
	if previous != nil {
		catalog = cloneClassroomCatalog(previous.ClassroomCatalog)
	}

	for campusName, campusInfo := range current.CampusInfoMap {
		if campusInfo == nil {
			continue
		}
		for _, buildingInfo := range campusInfo.BuildingInfoMap {
			if buildingInfo == nil {
				continue
			}
			for _, classroomInfo := range buildingInfo.ClassroomInfoMap {
				if classroomInfo == nil || classroomInfo.Name == "" {
					continue
				}
				entry := newCatalogEntry(buildingInfo.Name, classroomInfo, observedAt)
				mergeCatalogEntry(catalog, campusName, buildingInfo.Name, classroomInfo.Name, entry)
			}
		}
	}

	current.ClassroomCatalog = catalog
}

func cloneClassroomCatalog(source model.ClassroomCatalog) model.ClassroomCatalog {
	catalog := model.ClassroomCatalog{}
	for campusName, buildings := range source {
		catalog[campusName] = model.BuildingClassroomCatalog{}
		for buildingName, rooms := range buildings {
			catalog[campusName][buildingName] = model.RoomClassroomCatalog{}
			for roomName, entry := range rooms {
				if entry == nil {
					continue
				}
				copied := *entry
				catalog[campusName][buildingName][roomName] = &copied
			}
		}
	}
	return catalog
}

func newCatalogEntry(buildingName string, classroomInfo *model.ClassroomInfo, observedAt time.Time) *model.ClassroomCatalogEntry {
	source := catalogSourceClassTable
	if classroomInfo.CanTrust {
		source = catalogSourceRealtime
	}
	return &model.ClassroomCatalogEntry{
		Name:          formatClassroomDisplayName(buildingName, classroomInfo.Name),
		BuildingName:  buildingName,
		ClassroomName: classroomInfo.Name,
		Size:          classroomInfo.Size,
		CanTrust:      classroomInfo.CanTrust,
		Type:          classroomInfo.Type,
		Source:        source,
		FirstSeen:     observedAt,
		LastSeen:      observedAt,
		SeenCount:     1,
	}
}

func mergeCatalogEntry(catalog model.ClassroomCatalog, campusName string, buildingName string, roomName string, candidate *model.ClassroomCatalogEntry) {
	if catalog[campusName] == nil {
		catalog[campusName] = model.BuildingClassroomCatalog{}
	}
	if catalog[campusName][buildingName] == nil {
		catalog[campusName][buildingName] = model.RoomClassroomCatalog{}
	}

	existing := catalog[campusName][buildingName][roomName]
	if existing == nil {
		catalog[campusName][buildingName][roomName] = candidate
		return
	}

	if existing.Name == "" {
		existing.Name = candidate.Name
	}
	if existing.BuildingName == "" {
		existing.BuildingName = candidate.BuildingName
	}
	if existing.ClassroomName == "" {
		existing.ClassroomName = candidate.ClassroomName
	}
	if existing.FirstSeen.IsZero() {
		existing.FirstSeen = candidate.FirstSeen
	}
	existing.LastSeen = candidate.LastSeen
	existing.SeenCount++

	if shouldUseCandidateSize(existing, candidate) {
		existing.Size = candidate.Size
	}
	if candidate.Type != "" {
		existing.Type = candidate.Type
	}
	if candidate.CanTrust {
		existing.CanTrust = true
		existing.Source = catalogSourceRealtime
	} else if existing.Source == "" {
		existing.Source = catalogSourceClassTable
	}
}

func shouldUseCandidateSize(existing *model.ClassroomCatalogEntry, candidate *model.ClassroomCatalogEntry) bool {
	if candidate.Size <= 0 {
		return false
	}
	if existing.Size <= 0 {
		return true
	}
	if candidate.CanTrust && !existing.CanTrust {
		return true
	}
	return candidate.CanTrust == existing.CanTrust
}

func formatClassroomDisplayName(buildingName string, classroomName string) string {
	if buildingName == "" || classroomName == "" {
		return classroomName
	}
	if strings.HasPrefix(classroomName, buildingName) {
		return classroomName
	}
	return buildingName + "-" + classroomName
}
