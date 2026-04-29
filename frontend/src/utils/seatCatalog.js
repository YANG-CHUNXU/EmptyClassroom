export function searchSeatCatalog({
  catalog,
  campus,
  selectedBuildings,
  campusInfo,
  roomQuery,
  minSeats,
  observedAt,
}) {
  if (!campus || !campusInfo || selectedBuildings.length == 0) {
    return [];
  }

  const campusCatalog = catalog?.[campus] ?? {};
  const selectedBuildingNames = resolveSelectedBuildingNames(
    campusInfo,
    selectedBuildings
  );
  const normalizedRoomQuery = normalizeDigits(roomQuery);
  const minSeatCount = parseSeatCount(minSeats);
  const rowMap = new Map();

  for (const buildingName of selectedBuildingNames) {
    const rooms = campusCatalog[buildingName] ?? {};
    for (const [roomName, entry] of Object.entries(rooms)) {
      mergeRoom(rowMap, buildingName, roomName, entry);
    }
  }

  for (const buildingId of selectedBuildings) {
    const buildingInfo = campusInfo.building_info_map?.[buildingId];
    if (!buildingInfo) {
      continue;
    }
    for (const classroomInfo of Object.values(buildingInfo.classroom_info_map ?? {})) {
      if (!classroomInfo?.name) {
        continue;
      }
      mergeRoom(
        rowMap,
        buildingInfo.name,
        classroomInfo.name,
        buildCurrentSnapshotEntry(buildingInfo.name, classroomInfo, observedAt)
      );
    }
  }

  const rows = Array.from(rowMap.values()).filter((entry) => {
    if (
      normalizedRoomQuery &&
      !normalizeDigits(entry.classroom_name).includes(normalizedRoomQuery)
    ) {
      return false;
    }
    return minSeatCount == null || (entry.size ?? 0) >= minSeatCount;
  });

  rows.sort((a, b) => {
    const buildingCompare = (a.building_name ?? "").localeCompare(
      b.building_name ?? ""
    );
    if (buildingCompare != 0) {
      return buildingCompare;
    }
    return (a.classroom_name ?? "").localeCompare(b.classroom_name ?? "", undefined, {
      numeric: true,
    });
  });

  return rows;
}

function mergeRoom(rowMap, buildingName, roomName, candidate) {
  const key = `${buildingName}:${roomName}`;
  const existing = rowMap.get(key);
  if (!existing) {
    rowMap.set(key, candidate);
    return;
  }

  rowMap.set(key, {
    ...existing,
    ...candidate,
    size: candidate.size > 0 ? candidate.size : existing.size,
    type: candidate.type || existing.type,
    source:
      candidate.can_trust || existing.can_trust
        ? "教务"
        : candidate.source || existing.source,
    can_trust: candidate.can_trust || existing.can_trust,
    first_seen: existing.first_seen || candidate.first_seen,
    seen_count: existing.seen_count ?? candidate.seen_count,
  });
}

function buildCurrentSnapshotEntry(buildingName, classroomInfo, observedAt) {
  return {
    name: formatClassroomDisplayName(buildingName, classroomInfo.name),
    building_name: buildingName,
    classroom_name: classroomInfo.name,
    size: classroomInfo.size,
    can_trust: classroomInfo.can_trust,
    type: classroomInfo.type,
    source: classroomInfo.can_trust ? "教务" : "课表",
    last_seen: observedAt,
    seen_count: 1,
  };
}

function formatClassroomDisplayName(buildingName, classroomName) {
  if (buildingName == "" || classroomName == "") {
    return classroomName;
  }
  if (classroomName.startsWith(buildingName)) {
    return classroomName;
  }
  return `${buildingName}-${classroomName}`;
}

function resolveSelectedBuildingNames(campusInfo, selectedBuildings) {
  const idToName = {};
  for (const [buildingName, buildingId] of Object.entries(
    campusInfo.building_id_map ?? {}
  )) {
    idToName[buildingId] = buildingName;
  }
  return selectedBuildings.map((buildingId) => idToName[buildingId]).filter(Boolean);
}

function parseSeatCount(value) {
  const normalized = normalizeDigits(value);
  if (normalized == "") {
    return null;
  }
  return Number.parseInt(normalized, 10);
}

function normalizeDigits(value) {
  return String(value ?? "").replace(/\D/g, "");
}
