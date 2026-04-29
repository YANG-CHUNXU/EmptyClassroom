# Classroom Seat Query Design

## Goal

Add a seat query feature backed by a cumulative known-classroom catalog. Users can search seat counts by selecting a campus and one or more buildings, then entering a room number and/or minimum seat count.

## Product Behavior

- The feature is named "座位查询" in the UI.
- It uses the current campus and building selections as the search scope.
- Users may enter a room number such as `101` or `302`; the UI matches rooms within the selected buildings.
- Users may enter a minimum seat count such as `80`; the UI lists known rooms with at least that many seats.
- Room number and minimum seat filters can be combined.
- Results show full classroom name, seat count, type, data source, and last-seen time.
- Empty states must say the room is not found in the known catalog, not that the room does not exist.

## Data Model

The API response adds `classroom_catalog` to `ClassInfo`. It is a nested map:

```json
{
  "沙河": {
    "N": {
      "101": {
        "name": "N-101",
        "building_name": "N",
        "classroom_name": "101",
        "size": 80,
        "can_trust": true,
        "type": "普通教室",
        "source": "教务",
        "first_seen": "2026-04-29T00:00:00Z",
        "last_seen": "2026-04-29T00:00:00Z",
        "seen_count": 3
      }
    }
  }
}
```

## Catalog Merge

Daily refresh loads the previous snapshot before saving the new one. The new snapshot is merged into the previous catalog:

- New known rooms are inserted.
- Existing rooms update `last_seen` and increment `seen_count`.
- Non-zero seat counts are preferred over missing values.
- Trusted realtime data is preferred over class-table data.
- Empty type values never overwrite existing non-empty type values.
- `source` is `教务` if the room has ever been observed from trusted realtime data; otherwise it is `课表`.

## Architecture

- Backend keeps persistence simple by reusing the existing snapshot store. No extra database service is introduced.
- `service/catalog.go` owns catalog extraction and merge rules.
- `RefreshSnapshot` loads the previous snapshot, merges the catalog after realtime query, then saves the enriched snapshot.
- Frontend adds a focused `SeatQueryPanel` component and a small utility for catalog search.

## Testing

- Go unit tests cover catalog insertion, merge priority, and `RefreshSnapshot` preserving the previous catalog.
- Frontend verification uses lint and production build.
