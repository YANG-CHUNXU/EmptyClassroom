import PropTypes from "prop-types";
import { useEffect, useMemo, useState } from "react";
import dayjs from "dayjs";
import Modal from "./ui/Modal";
import EmptyState from "./ui/EmptyState";
import SurfaceCard from "./ui/SurfaceCard";
import "./EmptyClassroomTable.css";
import CalculateEmptyClassroom from "../utils/calculte";

function formatEmptyTimeList(emptyTimeList) {
  const classTime = [
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

  const classStartTime = [
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

  let emptyTimeListStr = "";
  if (emptyTimeList[0] == 0) {
    emptyTimeListStr += "00:00";
  } else {
    emptyTimeListStr += `00:00-08:00, ${classTime[emptyTimeList[0] - 1]}`;
  }

  for (let i = 1; i < emptyTimeList.length; i++) {
    if (emptyTimeList[i] - emptyTimeList[i - 1] == 1) {
      continue;
    }

    emptyTimeListStr +=
      `-${classStartTime[emptyTimeList[i - 1] + 1]}, ` +
      classTime[emptyTimeList[i] - 1];
  }

  if (emptyTimeList[emptyTimeList.length - 1] != 13) {
    emptyTimeListStr +=
      `-${classStartTime[emptyTimeList[emptyTimeList.length - 1] + 1]}`;
    emptyTimeListStr += `, ${classTime[classTime.length - 1]}-24:00`;
  } else {
    emptyTimeListStr += "-24:00";
  }

  return emptyTimeListStr;
}

function EmptyClassroomTable(props) {
  const {
    todayData,
    selectedDate,
    selectedCampus,
    selectedBuildings,
    selectedClassTimes,
    setIsError,
    useClassTable,
  } = props;
  const [selectedClassroom, setSelectedClassroom] = useState(null);

  const { emptyClassroom, hasCalculationError } = useMemo(() => {
    if (
      todayData.code != 0 ||
      selectedCampus == "" ||
      selectedBuildings.length == 0 ||
      selectedClassTimes.length == 0
    ) {
      return {
        emptyClassroom: [],
        hasCalculationError: false,
      };
    }

    try {
      return {
        emptyClassroom: CalculateEmptyClassroom(
          todayData.data,
          selectedCampus,
          selectedDate.toDate(),
          selectedBuildings,
          selectedClassTimes
        ),
        hasCalculationError: false,
      };
    } catch {
      return {
        emptyClassroom: [],
        hasCalculationError: true,
      };
    }
  }, [
    todayData.code,
    todayData.data,
    selectedCampus,
    selectedDate,
    selectedBuildings,
    selectedClassTimes,
  ]);

  useEffect(() => {
    if (hasCalculationError) {
      setIsError(true);
    }
  }, [hasCalculationError, setIsError]);

  if (todayData.code != 0 || selectedCampus == "") {
    return null;
  }

  if (
    selectedBuildings.length == 0 ||
    selectedClassTimes.length == 0 ||
    emptyClassroom.length == 0
  ) {
    return (
      <SurfaceCard className="empty-classroom-table">
        <EmptyState
          title={
            selectedBuildings.length == 0
              ? selectedClassTimes.length == 0
                ? "请选择教学楼和上课时间"
                : "请选择教学楼"
              : selectedClassTimes.length == 0
                ? "请选择上课时间"
                : "没有空教室了"
          }
        />
      </SurfaceCard>
    );
  }

  const rows =
    !useClassTable &&
    selectedCampus != "海南" &&
    selectedDate.isSame(dayjs(), "day") &&
    (todayData.data.is_fallback == undefined ||
      !todayData.data.is_fallback[selectedCampus])
      ? emptyClassroom.filter((item) => item.can_trust)
      : emptyClassroom;

  return (
    <div className="empty-classroom-table">
      <SurfaceCard className="empty-classroom-table__card" padded={false}>
        <div className="empty-classroom-table__scroll">
          <table className="empty-classroom-table__table">
            <thead>
              <tr>
                <th>教室</th>
                <th>座位数</th>
                <th>类型</th>
                <th>
                  <span className="empty-classroom-table__header-inline">
                    来源
                    <button
                      type="button"
                      className="ui-button ui-button--link"
                      onClick={() => {
                        window.open(
                          "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Bwf5dnkr0o3G6ExcedZcLzaSnBe"
                        );
                      }}
                    >
                      说明
                    </button>
                  </span>
                </th>
              </tr>
            </thead>
            <tbody>
              {rows.map((record) => (
                <tr key={record.name}>
                  <td>
                    <button
                      type="button"
                      className="ui-button ui-button--link"
                      onClick={() => setSelectedClassroom(record)}
                    >
                      {record.name}
                    </button>
                  </td>
                  <td>{record.size}</td>
                  <td>{record.type}</td>
                  <td>
                    <span
                      className={`ui-pill ${
                        record.can_trust ? "ui-pill--success" : "ui-pill--danger"
                      }`}
                    >
                      {record.can_trust ? "教务" : "课表"}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </SurfaceCard>

      <Modal
        open={selectedClassroom != null}
        title={selectedClassroom?.name || ""}
        onClose={() => setSelectedClassroom(null)}
      >
        {selectedClassroom ? (
          <dl className="empty-classroom-table__details">
            <div>
              <dt>座位数</dt>
              <dd>{selectedClassroom.size}</dd>
            </div>
            <div>
              <dt>类型</dt>
              <dd>{selectedClassroom.type}</dd>
            </div>
            <div>
              <dt>空闲时间</dt>
              <dd>{formatEmptyTimeList(selectedClassroom.empty_class_time)}</dd>
            </div>
            <div>
              <dt>数据来源</dt>
              <dd>{selectedClassroom.can_trust ? "教务（可信）" : "课表（参考）"}</dd>
            </div>
          </dl>
        ) : null}
      </Modal>
    </div>
  );
}

EmptyClassroomTable.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
  selectedClassTimes: PropTypes.array,
  setIsError: PropTypes.func,
  useClassTable: PropTypes.bool,
};

export default EmptyClassroomTable;
