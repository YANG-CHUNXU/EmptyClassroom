import assert from "node:assert/strict";
import test from "node:test";
import {
  calculateTodayIdleRate,
  calculateWeeklyIdleRate,
} from "./idleRate.js";

function makeMatrix(firstPeriod, restPeriod = [1, 1]) {
  return Array.from({ length: 14 }, (_, index) =>
    index === 0 ? firstPeriod : restPeriod
  );
}

function makePeriods(emptyClassrooms, totalClassrooms) {
  return Array.from({ length: 14 }, () => ({
    empty_classrooms: emptyClassrooms,
    total_classrooms: totalClassrooms,
    idle_rate: totalClassrooms === 0 ? 0 : emptyClassrooms / totalClassrooms,
  }));
}

const campusInfo = {
  building_id_map: {
    N: 0,
    S: 1,
  },
  building_info_map: {
    0: {
      name: "N",
      classroom_id_map: {
        101: 0,
        102: 1,
      },
      class_matrix: makeMatrix([0, 1]),
    },
    1: {
      name: "S",
      classroom_id_map: {
        201: 0,
      },
      class_matrix: makeMatrix([0], [1]),
    },
  },
};

const classInfo = {
  campus_info_map: {
    沙河: campusInfo,
  },
  classroom_catalog: {
    沙河: {
      N: {
        101: { classroom_name: "101" },
        102: { classroom_name: "102" },
        103: { classroom_name: "103" },
      },
      S: {
        201: { classroom_name: "201" },
      },
    },
  },
  idle_rate_history: [
    {
      date: "2026-05-26",
      campus: {
        沙河: {
          N: {
            known_classrooms: 3,
            periods: makePeriods(1, 3),
            average_rate: 1 / 3,
          },
          S: {
            known_classrooms: 1,
            periods: makePeriods(1, 1),
            average_rate: 1,
          },
        },
      },
    },
  ],
};

test("calculateTodayIdleRate aggregates selected buildings and uses catalog denominator", () => {
  const periods = calculateTodayIdleRate({
    classInfo,
    selectedCampus: "沙河",
    selectedBuildings: [0, 1],
  });

  assert.equal(periods.length, 14);
  assert.equal(periods[0].emptyClassrooms, 2);
  assert.equal(periods[0].totalClassrooms, 4);
  assert.equal(periods[0].idleRate, 0.5);
  assert.equal(periods[0].timeRange, "08:00-08:45");
  assert.equal(periods[13].timeRange, "20:10-20:55");
});

test("calculateWeeklyIdleRate maps selected building ids to historical building names", () => {
  const days = calculateWeeklyIdleRate({
    idleRateHistory: classInfo.idle_rate_history,
    selectedCampus: "沙河",
    selectedBuildings: [0, 1],
    campusInfo,
  });

  assert.deepEqual(days, [
    {
      date: "2026-05-26",
      averageRate: 0.5,
      averageEmptyClassrooms: 2,
      totalClassrooms: 4,
    },
  ]);
});

test("idle rate utilities return empty arrays when scope data is missing", () => {
  assert.deepEqual(
    calculateTodayIdleRate({
      classInfo,
      selectedCampus: "沙河",
      selectedBuildings: [],
    }),
    []
  );
  assert.deepEqual(
    calculateWeeklyIdleRate({
      idleRateHistory: classInfo.idle_rate_history,
      selectedCampus: "沙河",
      selectedBuildings: [0],
      campusInfo: null,
    }),
    []
  );
});
