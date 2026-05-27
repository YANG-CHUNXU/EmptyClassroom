import PropTypes from "prop-types";
import { useState } from "react";
import SurfaceCard from "./ui/SurfaceCard";
import IdleRatePanel from "./IdleRatePanel";
import SeatQueryPanel from "./SeatQueryPanel";
import {
  CLASSROOM_TOOL_TABS,
  DEFAULT_CLASSROOM_TOOLS_OPEN,
  DEFAULT_CLASSROOM_TOOL_TAB,
} from "../utils/classroomTools";
import "./ClassroomToolsPanel.css";

function ClassroomToolsPanel({ todayData, selectedCampus, selectedBuildings }) {
  const [isOpen, setIsOpen] = useState(DEFAULT_CLASSROOM_TOOLS_OPEN);
  const [activeTab, setActiveTab] = useState(DEFAULT_CLASSROOM_TOOL_TAB);

  if (todayData.code != 0 || selectedCampus == "") {
    return null;
  }

  return (
    <SurfaceCard className={`classroom-tools-panel ${!isOpen ? "is-collapsed" : ""}`}>
      <div className="classroom-tools-panel__header">
        <div>
          <h2 className="classroom-tools-panel__title">查询工具</h2>
          <p className="classroom-tools-panel__subtitle">
            闲置率统计和座位查询共用当前校区与教学楼
          </p>
        </div>
        <button
          type="button"
          className="ui-button ui-button--ghost classroom-tools-panel__toggle"
          onClick={() => setIsOpen(!isOpen)}
          aria-expanded={isOpen}
        >
          {isOpen ? "收起" : "展开"}
        </button>
      </div>

      <div className="classroom-tools-panel__body" hidden={!isOpen}>
        <div
          className="classroom-tools-panel__tabs"
          role="tablist"
          aria-label="查询工具切换"
        >
          {CLASSROOM_TOOL_TABS.map((tab) => (
            <button
              type="button"
              key={tab.id}
              role="tab"
              aria-selected={activeTab === tab.id}
              className={`classroom-tools-panel__tab ${
                activeTab === tab.id ? "is-active" : ""
              }`}
              onClick={() => setActiveTab(tab.id)}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="classroom-tools-panel__content">
          <div hidden={activeTab !== "idle-rate"}>
            <IdleRatePanel
              todayData={todayData}
              selectedCampus={selectedCampus}
              selectedBuildings={selectedBuildings}
            />
          </div>
          <div hidden={activeTab !== "seat-query"}>
            <SeatQueryPanel
              todayData={todayData}
              selectedCampus={selectedCampus}
              selectedBuildings={selectedBuildings}
            />
          </div>
        </div>
      </div>
    </SurfaceCard>
  );
}

ClassroomToolsPanel.propTypes = {
  todayData: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
};

export default ClassroomToolsPanel;
