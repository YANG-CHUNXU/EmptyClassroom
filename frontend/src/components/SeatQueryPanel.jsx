import PropTypes from "prop-types";
import { useMemo, useState } from "react";
import SurfaceCard from "./ui/SurfaceCard";
import EmptyState from "./ui/EmptyState";
import { searchSeatCatalog } from "../utils/seatCatalog";
import "./SeatQueryPanel.css";

function formatSeatCount(size) {
  return size > 0 ? size : "无数据";
}

function formatLastSeen(value) {
  if (!value) {
    return "无记录";
  }
  return new Date(value).toLocaleDateString();
}

function SeatQueryPanel({
  todayData,
  selectedCampus,
  selectedBuildings,
}) {
  const [roomQuery, setRoomQuery] = useState("");
  const [minSeats, setMinSeats] = useState("");
  const [isOpen, setIsOpen] = useState(false);

  const campusInfo =
    todayData.data?.campus_info_map?.[selectedCampus] ?? null;
  const rows = useMemo(
    () =>
      searchSeatCatalog({
        catalog: todayData.data?.classroom_catalog,
        campus: selectedCampus,
        selectedBuildings,
        campusInfo,
        roomQuery,
        minSeats,
        observedAt: todayData.data?.update_at,
      }),
    [
      campusInfo,
      minSeats,
      roomQuery,
      selectedBuildings,
      selectedCampus,
      todayData.data?.update_at,
      todayData.data?.classroom_catalog,
    ]
  );

  if (todayData.code != 0 || selectedCampus == "") {
    return null;
  }

  const hasFilter = roomQuery.trim() != "" || minSeats.trim() != "";
  const canQuery = selectedBuildings.length > 0;

  return (
    <SurfaceCard className={`seat-query-panel ${!isOpen ? "is-collapsed" : ""}`}>
      <div className="seat-query-panel__header">
        <div>
          <h2 className="seat-query-panel__title">座位查询</h2>
          <p className="seat-query-panel__subtitle">
            基于已知教室库，按当前校区和教学楼查询
          </p>
        </div>
        <button
          type="button"
          className="ui-button ui-button--ghost seat-query-panel__toggle"
          onClick={() => setIsOpen(!isOpen)}
          aria-expanded={isOpen}
        >
          {isOpen ? "收起" : "展开"}
        </button>
      </div>

      {isOpen ? (
        <>
          <div className="seat-query-panel__controls">
            <label className="seat-query-panel__field">
              <span>教室号</span>
              <input
                inputMode="numeric"
                pattern="[0-9]*"
                placeholder="如 101"
                value={roomQuery}
                onChange={(event) => setRoomQuery(event.target.value)}
              />
            </label>
            <label className="seat-query-panel__field">
              <span>最少座位</span>
              <input
                inputMode="numeric"
                pattern="[0-9]*"
                placeholder="如 80"
                value={minSeats}
                onChange={(event) => setMinSeats(event.target.value)}
              />
            </label>
          </div>

          {!canQuery ? (
            <EmptyState title="请选择教学楼后查询座位" />
          ) : !hasFilter ? (
            <EmptyState title="输入教室号或最少座位数开始查询" />
          ) : rows.length == 0 ? (
            <EmptyState
              title="已知教室库中未找到"
              description="可能是教室不存在，也可能是当前数据源尚未覆盖。"
            />
          ) : (
            <div className="seat-query-panel__results">
              {rows.map((room) => (
                <div className="seat-query-panel__result" key={room.name}>
                  <div>
                    <div className="seat-query-panel__room">{room.name}</div>
                    <div className="seat-query-panel__meta">
                      最近见到 {formatLastSeen(room.last_seen)}
                    </div>
                  </div>
                  <div className="seat-query-panel__facts">
                    <strong>{formatSeatCount(room.size)}</strong>
                    <span className={`ui-pill ${room.can_trust ? "ui-pill--success" : "ui-pill--danger"}`}>
                      {room.source || (room.can_trust ? "教务" : "课表")}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      ) : null}
    </SurfaceCard>
  );
}

SeatQueryPanel.propTypes = {
  todayData: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
};

export default SeatQueryPanel;
