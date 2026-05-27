import PropTypes from "prop-types";
import { useMemo } from "react";
import EmptyState from "./ui/EmptyState";
import {
  calculateTodayIdleRate,
  calculateWeeklyIdleRate,
} from "../utils/idleRate";
import "./IdleRatePanel.css";

function formatPercent(rate) {
  return `${(rate * 100).toFixed(1).replace(".0", "")}%`;
}

function formatCount(value) {
  return Number.isInteger(value)
    ? String(value)
    : value.toFixed(1).replace(".0", "");
}

function getPeak(items, rateKey) {
  if (items.length === 0) {
    return null;
  }
  return items.reduce((best, item) =>
    item[rateKey] > best[rateKey] ? item : best
  );
}

function RateChart({ items, rateKey, labelKey, describeItem }) {
  const peak = getPeak(items, rateKey);
  const width = Math.max(360, items.length * 54 + 40);
  const height = 210;
  const top = 28;
  const bottom = 62;
  const side = 16;
  const chartHeight = height - top - bottom;
  const gap = 10;
  const barWidth = (width - side * 2 - gap * (items.length - 1)) / items.length;

  return (
    <div className="idle-rate-panel__chart-scroll">
      <svg
        className="idle-rate-panel__chart"
        role="img"
        viewBox={`0 0 ${width} ${height}`}
        style={{ minWidth: `${width}px` }}
        aria-label="闲置率条形图"
      >
        <line
          x1={side}
          y1={top + chartHeight}
          x2={width - side}
          y2={top + chartHeight}
          className="idle-rate-panel__axis"
        />
        {items.map((item, index) => {
          const rate = item[rateKey];
          const x = side + index * (barWidth + gap);
          const barHeight = rate === 0 ? 0 : Math.max(3, rate * chartHeight);
          const y = top + chartHeight - barHeight;
          const isPeak = peak === item && rate > 0;

          return (
            <g key={`${item[labelKey]}-${index}`}>
              <title>{describeItem(item)}</title>
              {isPeak ? (
                <text
                  x={x + barWidth / 2}
                  y={18}
                  textAnchor="middle"
                  className="idle-rate-panel__peak-text"
                >
                  最高
                </text>
              ) : null}
              <rect
                x={x}
                y={top}
                width={barWidth}
                height={chartHeight}
                rx="6"
                className="idle-rate-panel__bar-bg"
              />
              <rect
                x={x}
                y={y}
                width={barWidth}
                height={barHeight}
                rx="6"
                className={`idle-rate-panel__bar ${isPeak ? "is-peak" : ""}`}
              />
              <text
                x={x + barWidth / 2}
                y={top + chartHeight + 20}
                textAnchor="middle"
                className="idle-rate-panel__label"
              >
                {item[labelKey]}
              </text>
              {item.startTime && item.endTime ? (
                <text
                  x={x + barWidth / 2}
                  y={top + chartHeight + 36}
                  textAnchor="middle"
                  className="idle-rate-panel__time-label"
                >
                  <tspan x={x + barWidth / 2}>{item.startTime}</tspan>
                  <tspan x={x + barWidth / 2} dy="13">
                    {item.endTime}
                  </tspan>
                </text>
              ) : null}
            </g>
          );
        })}
      </svg>
    </div>
  );
}

RateChart.propTypes = {
  items: PropTypes.array.isRequired,
  rateKey: PropTypes.string.isRequired,
  labelKey: PropTypes.string.isRequired,
  describeItem: PropTypes.func.isRequired,
};

function IdleRatePanel({ todayData, selectedCampus, selectedBuildings }) {
  const classInfo = todayData.data;
  const campusInfo = classInfo?.campus_info_map?.[selectedCampus] ?? null;

  const todayPeriods = useMemo(
    () =>
      calculateTodayIdleRate({
        classInfo,
        selectedCampus,
        selectedBuildings,
      }),
    [classInfo, selectedCampus, selectedBuildings]
  );

  const weeklyDays = useMemo(
    () =>
      calculateWeeklyIdleRate({
        idleRateHistory: classInfo?.idle_rate_history,
        selectedCampus,
        selectedBuildings,
        campusInfo,
      }),
    [campusInfo, classInfo?.idle_rate_history, selectedBuildings, selectedCampus]
  );

  if (todayData.code != 0 || selectedCampus == "") {
    return null;
  }

  const todayPeak = getPeak(todayPeriods, "idleRate");
  const weekPeak = getPeak(weeklyDays, "averageRate");

  return (
    <div className="idle-rate-panel">
      <div className="idle-rate-panel__pills">
        {todayPeak ? (
          <span className="ui-pill idle-rate-panel__pill">
            今日最高 {todayPeak.label} {formatPercent(todayPeak.idleRate)}
          </span>
        ) : null}
        {weekPeak ? (
          <span className="ui-pill idle-rate-panel__pill">
            日均最高 {weekPeak.date.slice(5)} {formatPercent(weekPeak.averageRate)}
          </span>
        ) : null}
      </div>

      {selectedBuildings.length === 0 ? (
        <EmptyState title="请选择教学楼后查看闲置率" />
      ) : (
        <div className="idle-rate-panel__sections">
          <section className="idle-rate-panel__section">
            <h3>今日各节课</h3>
            {todayPeriods.length > 0 ? (
              <RateChart
                items={todayPeriods}
                rateKey="idleRate"
                labelKey="label"
                describeItem={(item) =>
                  `第${item.label}节 ${item.timeRange} ${formatPercent(item.idleRate)}，${item.emptyClassrooms}/${item.totalClassrooms}`
                }
              />
            ) : (
              <EmptyState title="今日暂无闲置率数据" />
            )}
          </section>

          <section className="idle-rate-panel__section">
            <h3>近7日日均</h3>
            {weeklyDays.length > 0 ? (
              <RateChart
                items={weeklyDays.map((item) => ({
                  ...item,
                  shortDate: item.date.slice(5),
                }))}
                rateKey="averageRate"
                labelKey="shortDate"
                describeItem={(item) =>
                  `${item.date} ${formatPercent(item.averageRate)}，日均空闲 ${formatCount(item.averageEmptyClassrooms)}/${formatCount(item.totalClassrooms)}`
                }
              />
            ) : (
              <EmptyState title="周统计暂无数据" />
            )}
          </section>
        </div>
      )}
    </div>
  );
}

IdleRatePanel.propTypes = {
  todayData: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
};

export default IdleRatePanel;
