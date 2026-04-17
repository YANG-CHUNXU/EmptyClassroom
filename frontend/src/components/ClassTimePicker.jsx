import PropTypes from "prop-types";
import dayjs from "dayjs";
import { useEffect } from "react";
import SurfaceCard from "./ui/SurfaceCard";
import "./ClassTimePicker.css";

const CLASSES = [
  "01",
  "02",
  "03",
  "04",
  "05",
  "06",
  "07",
  "08",
  "09",
  "10",
  "11",
  "12",
  "13",
  "14",
];

const CLASS_TIME = [
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

function isClassTimeDisabled(
  index,
  classTime,
  currentTime,
  canSelectAllDay,
  isSelectedToday
) {
  return (
    classTime[index].localeCompare(currentTime) < 0 &&
    classTime[classTime.length - 1].localeCompare(currentTime) >= 0 &&
    !canSelectAllDay &&
    isSelectedToday
  );
}

function ClassTimePicker(props) {
  const {
    todayData,
    selectedClassTimes,
    setSelectedClassTimes,
    selectedCampus,
    selectedDate,
    showClassTime,
    canSelectAllDay,
    isDark,
  } = props;

  const shouldRender = todayData.code == 0 && selectedCampus != "";

  function fillZero(x) {
    if (x < 10) {
      return `0${x}`;
    }
    return x;
  }

  const options = [];
  const now = new Date();
  const nowHour = fillZero(now.getHours());
  const nowMinute = fillZero(now.getMinutes());
  const currentTime = `${nowHour}:${nowMinute}`;
  const isSelectedToday = selectedDate.isSame(dayjs(), "day");

  for (let i = 0; i <= 13; i++) {
    options.push({
      label: CLASSES[i],
      value: i,
      disabled: isClassTimeDisabled(
        i,
        CLASS_TIME,
        currentTime,
        canSelectAllDay,
        isSelectedToday
      ),
    });
  }

  useEffect(() => {
    if (!shouldRender) {
      return;
    }

    const enabledSelectedClassTimes = selectedClassTimes.filter(
      (value) =>
        !isClassTimeDisabled(
          value,
          CLASS_TIME,
          currentTime,
          canSelectAllDay,
          isSelectedToday
        )
    );

    const hasChanged =
      enabledSelectedClassTimes.length !== selectedClassTimes.length ||
      enabledSelectedClassTimes.some(
        (value, index) => value !== selectedClassTimes[index]
      );

    if (hasChanged) {
      setSelectedClassTimes(enabledSelectedClassTimes);
    }
  }, [
    canSelectAllDay,
    currentTime,
    isSelectedToday,
    selectedClassTimes,
    setSelectedClassTimes,
    shouldRender,
  ]);

  function isAllChecked() {
    for (let i = 0; i <= 13; i++) {
      if (options[i].disabled) {
        continue;
      }
      if (!selectedClassTimes.includes(i)) {
        return false;
      }
    }
    return true;
  }

  function onCheckAllChange() {
    if (!isAllChecked()) {
      let newSelectedClassTimes = [];
      for (let i = 0; i <= 13; i++) {
        if (options[i].disabled) {
          continue;
        }
        newSelectedClassTimes.push(i);
      }
      setSelectedClassTimes(newSelectedClassTimes);
    } else {
      setSelectedClassTimes([]);
    }
  }

  if (!shouldRender) {
    return null;
  }

  return (
    <SurfaceCard className="class-time-picker">
      <div className="class-time-picker__grid">
        {options.map((x) => (
          <button
            type="button"
            key={x.value}
            className={`class-time-picker__button ${
              selectedClassTimes.includes(x.value) ? "primary" : "outline"
            }`}
            onClick={() => {
              if (selectedClassTimes.includes(x.value)) {
                setSelectedClassTimes(selectedClassTimes.filter((y) => y != x.value));
              } else {
                setSelectedClassTimes([...selectedClassTimes, x.value]);
              }
            }}
            style={{
              height: showClassTime ? "52px" : "36px",
              color: x.disabled ? (isDark ? "#ffffff73" : "#00000073") : undefined,
            }}
            disabled={x.disabled}
          >
            <div>
              {showClassTime ? (
                <div className="class-time-picker__tiny class-time-picker__tiny--top">
                  {CLASS_START_TIME[x.label - 1]}
                </div>
              ) : null}
              {x.label}
              {showClassTime ? (
                <div className="class-time-picker__tiny class-time-picker__tiny--bottom">
                  {CLASS_TIME[x.label - 1]}
                </div>
              ) : null}
            </div>
          </button>
        ))}
        <button
          type="button"
          className={`class-time-picker__button ${isAllChecked() ? "primary" : "outline"}`}
          onClick={onCheckAllChange}
          style={{ height: showClassTime ? "52px" : "36px" }}
        >
          {isAllChecked() ? "全不选" : "全选"}
        </button>
      </div>
    </SurfaceCard>
  );
}

ClassTimePicker.propTypes = {
  todayData: PropTypes.object,
  selectedClassTimes: PropTypes.array,
  setSelectedClassTimes: PropTypes.func,
  selectedCampus: PropTypes.string,
  selectedDate: PropTypes.object,
  showClassTime: PropTypes.bool,
  canSelectAllDay: PropTypes.bool,
  isDark: PropTypes.bool,
};

export default ClassTimePicker;
