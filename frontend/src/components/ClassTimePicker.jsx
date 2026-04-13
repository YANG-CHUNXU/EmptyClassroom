import PropTypes from "prop-types";
import { Card, Button } from "antd";
import dayjs from "dayjs";
import { useEffect } from "react";
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
      return "0" + x;
    }
    return x;
  }

  const options = [];
  const now = new Date();
  const now_hour = fillZero(now.getHours());
  const now_minute = fillZero(now.getMinutes());
  const currentTime = `${now_hour}:${now_minute}`;
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
    <Card
      className="class-time-picker"
      style={{
        maxWidth: 400,
        width: "90%",
        boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
      }}
      bodyStyle={{
        maxWidth: "300px",
      }}
    >
      <div
        style={{
          display: "flex",
          flexWrap: "wrap",
          justifyContent: "center",
        }}
      >
        {options.map((x) => (
          <Button
            key={x.value}
            type={
              selectedClassTimes.includes(x.value) ? "primary" : "outline"
            }
            onClick={() => {
              if (selectedClassTimes.includes(x.value)) {
                setSelectedClassTimes(selectedClassTimes.filter((y) => y != x.value));
              } else {
                setSelectedClassTimes([...selectedClassTimes, x.value]);
              }
            }}
            style={{
              borderRadius: "0px",
              width: "45px",
              margin: "2px",
              height: showClassTime ? "45px" : "30px",
              padding: "0px",
              color: x.disabled
                ? isDark
                  ? "#ffffff73"
                  : "#00000073"
                : null,
            }}
            disabled={x.disabled}
          >
            <div>
              {showClassTime ? (
                <div
                  style={{
                    fontSize: "0.7em",
                    marginBottom: "-0.5em",
                  }}
                >
                  {CLASS_START_TIME[x.label - 1]}
                </div>
              ) : null}
              {x.label}
              {showClassTime ? (
                <div
                  style={{
                    fontSize: "0.7em",
                    marginTop: "-0.5em",
                  }}
                >
                  {CLASS_TIME[x.label - 1]}
                </div>
              ) : null}
            </div>
          </Button>
        ))}
        <Button
          type={isAllChecked() ? "primary" : "outline"}
          onClick={onCheckAllChange}
          style={{
            borderRadius: "0px",
            width: "45px",
            margin: "2px",
            height: showClassTime ? "45px" : "30px",
            padding: "0px",
          }}
        >
          {isAllChecked() ? "全不选" : "全选"}
        </Button>
      </div>
    </Card>
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
