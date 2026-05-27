package service

import (
	"EmptyClassroom/service/model"
	"sort"
	"time"
)

const idleRateHistoryLimit = 7
const classPeriodsPerDay = 14

func MergeIdleRateHistory(previous *model.ClassInfo, current *model.ClassInfo, observedAt time.Time) {
	if current == nil {
		return
	}

	var previousHistory []model.IdleRateHistoryEntry
	if previous != nil {
		previousHistory = previous.IdleRateHistory
	}
	currentEntry := buildIdleRateHistoryEntry(current, observedAt)
	current.IdleRateHistory = upsertIdleRateHistory(previousHistory, currentEntry)
}

func buildIdleRateHistoryEntry(classInfo *model.ClassInfo, observedAt time.Time) model.IdleRateHistoryEntry {
	entry := model.IdleRateHistoryEntry{
		Date:   formatIdleRateDate(observedAt),
		Campus: map[string]model.IdleRateCampusStats{},
	}
	if classInfo == nil {
		return entry
	}

	for campusName, buildings := range classInfo.ClassroomCatalog {
		campusStats := model.IdleRateCampusStats{}
		for buildingName, rooms := range buildings {
			buildingStats := buildIdleRateBuildingStats(classInfo, campusName, buildingName, rooms)
			if buildingStats.KnownClassrooms == 0 {
				continue
			}
			campusStats[buildingName] = buildingStats
		}
		if len(campusStats) > 0 {
			entry.Campus[campusName] = campusStats
		}
	}

	return entry
}

func buildIdleRateBuildingStats(classInfo *model.ClassInfo, campusName string, buildingName string, rooms model.RoomClassroomCatalog) model.IdleRateBuildingStats {
	stats := model.IdleRateBuildingStats{
		Periods: make([]model.IdleRatePeriodStat, classPeriodsPerDay),
	}

	roomNames := make([]string, 0, len(rooms))
	for roomName, room := range rooms {
		if room == nil {
			continue
		}
		catalogRoomName := room.ClassroomName
		if catalogRoomName == "" {
			catalogRoomName = roomName
		}
		if catalogRoomName == "" {
			continue
		}
		roomNames = append(roomNames, catalogRoomName)
	}
	stats.KnownClassrooms = len(roomNames)
	if stats.KnownClassrooms == 0 {
		return stats
	}

	buildingInfo := findBuildingInfo(classInfo, campusName, buildingName)
	rateSum := 0.0
	for periodIndex := 0; periodIndex < classPeriodsPerDay; periodIndex++ {
		emptyClassrooms := 0
		for _, roomName := range roomNames {
			if isClassroomEmptyAtPeriod(buildingInfo, roomName, periodIndex) {
				emptyClassrooms++
			}
		}
		idleRate := float64(emptyClassrooms) / float64(stats.KnownClassrooms)
		stats.Periods[periodIndex] = model.IdleRatePeriodStat{
			EmptyClassrooms: emptyClassrooms,
			TotalClassrooms: stats.KnownClassrooms,
			IdleRate:        idleRate,
		}
		rateSum += idleRate
	}
	stats.AverageRate = rateSum / classPeriodsPerDay
	return stats
}

func findBuildingInfo(classInfo *model.ClassInfo, campusName string, buildingName string) *model.BuildingInfo {
	if classInfo == nil || classInfo.CampusInfoMap == nil {
		return nil
	}
	campusInfo := classInfo.CampusInfoMap[campusName]
	if campusInfo == nil {
		return nil
	}
	buildingID, ok := campusInfo.BuildingIdMap[buildingName]
	if !ok {
		return nil
	}
	return campusInfo.BuildingInfoMap[buildingID]
}

func isClassroomEmptyAtPeriod(buildingInfo *model.BuildingInfo, roomName string, periodIndex int) bool {
	if buildingInfo == nil || periodIndex < 0 || periodIndex >= len(buildingInfo.ClassMatrix) {
		return false
	}
	classroomID, ok := buildingInfo.ClassroomIdMap[roomName]
	if !ok {
		return false
	}
	period := buildingInfo.ClassMatrix[periodIndex]
	if classroomID < 0 || classroomID >= len(period) {
		return false
	}
	return period[classroomID] == 0
}

func upsertIdleRateHistory(previous []model.IdleRateHistoryEntry, current model.IdleRateHistoryEntry) []model.IdleRateHistoryEntry {
	byDate := map[string]model.IdleRateHistoryEntry{}
	for _, entry := range previous {
		if entry.Date == "" {
			continue
		}
		byDate[entry.Date] = entry
	}
	if current.Date != "" {
		byDate[current.Date] = current
	}

	dates := make([]string, 0, len(byDate))
	for date := range byDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	if len(dates) > idleRateHistoryLimit {
		dates = dates[len(dates)-idleRateHistoryLimit:]
	}

	history := make([]model.IdleRateHistoryEntry, 0, len(dates))
	for _, date := range dates {
		history = append(history, byDate[date])
	}
	return history
}

func formatIdleRateDate(observedAt time.Time) string {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("Asia/Shanghai", 8*60*60)
	}
	return observedAt.In(location).Format("2006-01-02")
}
