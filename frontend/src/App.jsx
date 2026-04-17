import "./App.css";
import logo from "./assets/logo.png";
import { useState, useEffect } from "react";
import dayjs from "dayjs";
import Notification from "./components/Notification";
import CampusButtonGroup from "./components/CampusButtonGroup";
import BuildingPicker from "./components/BuildingPicker";
import ClassTimePicker from "./components/ClassTimePicker";
import EmptyClassroomTable from "./components/EmptyClassroomTable";
import GlobalEmpty from "./components/GlobalEmpty";
import Footer from "./components/Footer";
import ClassTableWarn from "./components/ClassTableWarn";
import "./components/ui/ui.css";

function App() {
  const [spining, setSpining] = useState(true);
  const [isError, setIsError] = useState(false);
  const [resp, setResp] = useState({ code: 1 });
  const [selectedCampus, setSelectedCampus] = useState("");
  const [selectedDate, setSelectedDate] = useState(dayjs());
  const [selectedBuildings, setSelectedBuildings] = useState([]);
  const [selectedClassTimes, setSelectedClassTimes] = useState([]);
  const [showClassTime, setShowClassTime] = useState(false);
  const [canSelectAllDay, setCanSelectAllDay] = useState(false);
  const [useClassTable, setUseClassTable] = useState(false);
  const [dontWarnClassTable, setDontWarnClassTable] = useState(false);
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const mql = window.matchMedia("(prefers-color-scheme: dark)");

    function matchMode(e) {
      const body = document.body;
      if (e.matches) {
        body.classList.add("dark");
        setIsDark(true);
        localStorage.setItem("darkMode", "true");
      } else {
        body.classList.remove("dark");
        setIsDark(false);
        localStorage.setItem("darkMode", "false");
      }
    }

    mql.addEventListener("change", matchMode);

    matchMode(mql);
    fetch("/api/get_data")
      .then((resp) => resp.json())
      .then((resp) => {
        setResp(resp);
        setIsError(false);
        setSpining(false);
      })
      .catch(() => {
        setIsError(true);
        setSpining(false);
      });
    setShowClassTime(localStorage.getItem("showClassTime") !== "false");
    setCanSelectAllDay(localStorage.getItem("canSelectAllDay") === "true");
    setUseClassTable(localStorage.getItem("useClassTable") === "true");
    setDontWarnClassTable(
      localStorage.getItem("dontWarnClassTable") === "true"
    );
  }, []);

  return (
    <div className="App">
      <img src={logo} className="logo" />
      <h1 className="app-title">BUPT 空教室查询</h1>
      {spining ? <div className="app-loading">正在加载数据...</div> : null}
      <Notification todayData={resp} />
      <CampusButtonGroup
        todayData={resp}
        selectedCampus={selectedCampus}
        setSelectedCampus={setSelectedCampus}
        setSelectedBuildings={setSelectedBuildings}
        showClassTime={showClassTime}
        setShowClassTime={setShowClassTime}
        canSelectAllDay={canSelectAllDay}
        setCanSelectAllDay={setCanSelectAllDay}
        useClassTable={useClassTable}
        setUseClassTable={setUseClassTable}
        selectedDate={selectedDate}
      />
      <BuildingPicker
        todayData={resp}
        selectedBuildings={selectedBuildings}
        setSelectedBuildings={setSelectedBuildings}
        selectedCampus={selectedCampus}
      />
      <ClassTimePicker
        todayData={resp}
        selectedClassTimes={selectedClassTimes}
        setSelectedClassTimes={setSelectedClassTimes}
        selectedCampus={selectedCampus}
        selectedDate={selectedDate}
        showClassTime={showClassTime}
        canSelectAllDay={canSelectAllDay}
        isDark={isDark}
      />
      <ClassTableWarn
        todayData={resp}
        selectedDate={selectedDate}
        selectedCampus={selectedCampus}
        useClassTable={useClassTable}
        dontWarnClassTable={dontWarnClassTable}
        setDontWarnClassTable={setDontWarnClassTable}
      />
      <EmptyClassroomTable
        todayData={resp}
        selectedDate={selectedDate}
        selectedCampus={selectedCampus}
        selectedBuildings={selectedBuildings}
        selectedClassTimes={selectedClassTimes}
        setIsError={setIsError}
        useClassTable={useClassTable}
      />
      <GlobalEmpty todayData={resp} isError={isError} />
      <Footer />
    </div>
  );
}

export default App;
