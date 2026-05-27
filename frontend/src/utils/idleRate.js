const PERIOD_COUNT = 14;

const CLASS_START_TIME = [
  "08:00",
  "08:50",
  "09:50",
  "10:40",
  "11:30",
  "13:00",
  "13:50",
  "14:45",
  "15:40",
  "16:35",
  "17:25",
  "18:30",
  "19:20",
  "20:10",
];

const CLASS_END_TIME = [
  "08:45",
  "09:35",
  "10:35",
  "11:25",
  "12:15",
  "13:45",
  "14:35",
  "15:30",
  "16:25",
  "17:20",
  "18:10",
  "19:15",
  "20:05",
  "20:55",
];

export function getSelectedBuildingNames(campusInfo, selectedBuildings) {
  if (!campusInfo || !Array.isArray(selectedBuildings) || selectedBuildings.length === 0) {
    return [];
  }

  const idToName = new Map();
  for (const [buildingName, buildingId] of Object.entries(campusInfo.building_id_map ?? {})) {
    idToName.set(String(buildingId), buildingName);
  }

  return selectedBuildings
    .map((buildingId) => idToName.get(String(buildingId)))
    .filter(Boolean);
}

export function calculateTodayIdleRate({
  classInfo,
  selectedCampus,
  selectedBuildings,
}) {
  const campusInfo = classInfo?.campus_info_map?.[selectedCampus];
  const catalog = classInfo?.classroom_catalog?.[selectedCampus];
  const buildingNames = getSelectedBuildingNames(campusInfo, selectedBuildings);
  if (!campusInfo || !catalog || buildingNames.length === 0) {
    return [];
  }

  const totals = Array.from({ length: PERIOD_COUNT }, (_, periodIndex) => ({
    periodIndex,
    label: String(periodIndex + 1).padStart(2, "0"),
    startTime: CLASS_START_TIME[periodIndex],
    endTime: CLASS_END_TIME[periodIndex],
    timeRange: `${CLASS_START_TIME[periodIndex]}-${CLASS_END_TIME[periodIndex]}`,
    emptyClassrooms: 0,
    totalClassrooms: 0,
    idleRate: 0,
  }));

  for (const buildingName of buildingNames) {
    const rooms = catalog[buildingName];
    if (!rooms) {
      continue;
    }

    const roomNames = Object.entries(rooms)
      .map(([roomName, entry]) => entry?.classroom_name || roomName)
      .filter(Boolean);
    if (roomNames.length === 0) {
      continue;
    }

    const buildingInfo = findBuildingInfo(campusInfo, buildingName);
    for (let periodIndex = 0; periodIndex < PERIOD_COUNT; periodIndex += 1) {
      let emptyClassrooms = 0;
      for (const roomName of roomNames) {
        if (isClassroomEmptyAtPeriod(buildingInfo, roomName, periodIndex)) {
          emptyClassrooms += 1;
        }
      }
      totals[periodIndex].emptyClassrooms += emptyClassrooms;
      totals[periodIndex].totalClassrooms += roomNames.length;
    }
  }

  const hasKnownClassrooms = totals.some((period) => period.totalClassrooms > 0);
  if (!hasKnownClassrooms) {
    return [];
  }

  return totals.map((period) => ({
    ...period,
    idleRate:
      period.totalClassrooms === 0
        ? 0
        : period.emptyClassrooms / period.totalClassrooms,
  }));
}

export function calculateWeeklyIdleRate({
  idleRateHistory,
  selectedCampus,
  selectedBuildings,
  campusInfo,
}) {
  const buildingNames = getSelectedBuildingNames(campusInfo, selectedBuildings);
  if (!Array.isArray(idleRateHistory) || buildingNames.length === 0) {
    return [];
  }

  return idleRateHistory
    .slice()
    .sort((a, b) => String(a.date).localeCompare(String(b.date)))
    .map((entry) => {
      const campusStats = entry?.campus?.[selectedCampus];
      if (!campusStats) {
        return null;
      }

      let emptyClassroomPeriods = 0;
      let totalClassroomPeriods = 0;
      for (const buildingName of buildingNames) {
        const buildingStats = campusStats[buildingName];
        if (!buildingStats) {
          continue;
        }

        for (const period of buildingStats.periods ?? []) {
          emptyClassroomPeriods += period.empty_classrooms ?? 0;
          totalClassroomPeriods +=
            period.total_classrooms ?? buildingStats.known_classrooms ?? 0;
        }
      }

      if (totalClassroomPeriods === 0) {
        return null;
      }

      return {
        date: entry.date,
        averageRate: emptyClassroomPeriods / totalClassroomPeriods,
        averageEmptyClassrooms: emptyClassroomPeriods / PERIOD_COUNT,
        totalClassrooms: totalClassroomPeriods / PERIOD_COUNT,
      };
    })
    .filter(Boolean);
}

function findBuildingInfo(campusInfo, buildingName) {
  const buildingId = campusInfo?.building_id_map?.[buildingName];
  if (buildingId === undefined || buildingId === null) {
    return null;
  }
  return campusInfo?.building_info_map?.[buildingId] ?? null;
}

function isClassroomEmptyAtPeriod(buildingInfo, roomName, periodIndex) {
  if (!buildingInfo) {
    return false;
  }
  const classroomId = buildingInfo.classroom_id_map?.[roomName];
  if (classroomId === undefined || classroomId === null) {
    return false;
  }
  return buildingInfo.class_matrix?.[periodIndex]?.[classroomId] === 0;
}
