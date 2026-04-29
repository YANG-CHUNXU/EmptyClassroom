# Classroom Seat Query Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a known-classroom seat query feature that accumulates classroom metadata from daily refresh snapshots.

**Architecture:** Add a catalog field to the existing `ClassInfo` snapshot payload, merge previous and current classroom observations during refresh, then expose the catalog to a new React query panel. Persistence stays on the existing Vercel Blob/local JSON snapshot store.

**Tech Stack:** Go service/model tests, React 18, Vite, existing CSS/UI card components.

---

## File Structure

- Modify `service/model/realtime_data.go`: add catalog model types and `ClassInfo.ClassroomCatalog`.
- Create `service/catalog.go`: extract and merge classroom catalog entries.
- Create `service/catalog_test.go`: verify merge behavior.
- Modify `service/api.go`: load previous snapshot during refresh and merge catalogs before save.
- Modify `service/api_test.go`: verify refresh preserves and enriches previous catalog.
- Create `frontend/src/utils/seatCatalog.js`: search helper for selected campus/buildings and filters.
- Create `frontend/src/components/SeatQueryPanel.jsx`: UI card for room number and minimum seat query.
- Create `frontend/src/components/SeatQueryPanel.css`: component styling.
- Modify `frontend/src/App.jsx`: render the new card after building selection.

## Tasks

- [ ] Add backend catalog model types.
- [ ] Add catalog merge unit tests.
- [ ] Implement catalog extraction and merge.
- [ ] Wire merge into snapshot refresh.
- [ ] Add frontend search helper.
- [ ] Add seat query panel UI.
- [ ] Run `GOCACHE=/tmp/emptyclassroom-go-cache go test ./...`.
- [ ] Run `pnpm lint` and `pnpm build` in `frontend`.
