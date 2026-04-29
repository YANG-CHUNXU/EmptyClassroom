import PropTypes from "prop-types";
import { useEffect, useState } from "react";
import Modal from "./ui/Modal";
import "./CampusButtonGroup.css";

function SettingToggle({ checked, onChange, label }) {
  return (
    <label className="campus-settings__row">
      <span>{label}</span>
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        className={`campus-settings__switch ${checked ? "is-on" : ""}`}
        onClick={() => onChange(!checked)}
      >
        <span className="campus-settings__knob" />
      </button>
    </label>
  );
}

SettingToggle.propTypes = {
  checked: PropTypes.bool.isRequired,
  onChange: PropTypes.func.isRequired,
  label: PropTypes.string.isRequired,
};

function CampusButtonGroup(props) {
  const { todayData, selectedCampus, setSelectedCampus, setSelectedBuildings } =
    props;
  const [campusList, setCampusList] = useState([]);
  const [openSettingModal, setOpenSettingModal] = useState(false);

  useEffect(() => {
    if (todayData.code != 0) {
      setCampusList([]);
      return;
    }

    const campusInfoMap = todayData.data?.campus_info_map ?? {};
    const list = Object.keys(campusInfoMap);
    const order = ["西土城", "沙河"];

    list.sort((a, b) => {
      if (order.indexOf(a) == -1) {
        if (order.indexOf(b) == -1) {
          return a.localeCompare(b);
        }
        return 1;
      }

      if (order.indexOf(b) == -1) {
        return -1;
      }

      return order.indexOf(a) - order.indexOf(b);
    });

    setCampusList(list);

    if (list.length == 0) {
      setSelectedCampus("");
      setSelectedBuildings([]);
      return;
    }

    if (!list.includes(selectedCampus)) {
      setSelectedCampus(list[0]);
      setSelectedBuildings([]);
    }
  }, [
    selectedCampus,
    setSelectedBuildings,
    setSelectedCampus,
    todayData.code,
    todayData.data?.campus_info_map,
  ]);

  return (
    <div className="campus-button-group">
      <div className="campus-button-group__toolbar">
        <div className="campus-button-group__date">
          {props.selectedDate ? props.selectedDate.format("YYYY-MM-DD") : ""}
        </div>
        <button
          type="button"
          className="ui-button ui-button--ghost campus-button-group__settings-button"
          onClick={() => setOpenSettingModal(true)}
          aria-label="打开设置"
        >
          设置
        </button>
      </div>
      <div className="campus-button-group__options" role="tablist" aria-label="校区选择">
        {campusList.map((campus) => (
          <button
            type="button"
            key={campus}
            className={`campus-button-group__option ${
              selectedCampus == campus ? "is-active" : ""
            }`}
            onClick={() => {
              setSelectedCampus(campus);
              setSelectedBuildings([]);
            }}
          >
            {campus}
          </button>
        ))}
      </div>

      <Modal
        title="设置"
        open={openSettingModal}
        onClose={() => setOpenSettingModal(false)}
      >
        <div className="campus-settings">
          <SettingToggle
            checked={props.showClassTime}
            onChange={(value) => {
              localStorage.setItem("showClassTime", value ? "true" : "false");
              props.setShowClassTime(value);
            }}
            label="显示课程时间"
          />
          <SettingToggle
            checked={props.canSelectAllDay}
            onChange={(value) => {
              localStorage.setItem("canSelectAllDay", value ? "true" : "false");
              props.setCanSelectAllDay(value);
            }}
            label="全选时选全天"
          />
          <SettingToggle
            checked={props.useClassTable}
            onChange={(value) => {
              localStorage.setItem("useClassTable", value ? "true" : "false");
              props.setUseClassTable(value);
            }}
            label="非必要情况下也使用课表数据"
          />
          <div className="campus-settings__section">
            <div className="campus-settings__line">
              数据来源：
              <button
                type="button"
                className="ui-button ui-button--link"
                onClick={() => {
                  window.open(
                    "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Zip8dx2rlobE5hxW00CcHwOOnre"
                  );
                }}
              >
                了解更多
              </button>
            </div>
            <div className="campus-settings__line">
              问答 Q&amp;A：
              <button
                type="button"
                className="ui-button ui-button--link"
                onClick={() => {
                  window.open(
                    "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj"
                  );
                }}
              >
                空教室查询Q&amp;A
              </button>
            </div>
            <div className="campus-settings__line">
              当前数据刷新时间：
              {new Date(props.todayData.data?.update_at).toLocaleString()}
            </div>
            <div className="campus-settings__line">
              声明：本项目是基于原项目 vibe coding 生成的。
            </div>
            <div className="campus-settings__line">
              原项目已开源：
              <button
                type="button"
                className="ui-button ui-button--link"
                onClick={() =>
                  window.open("https://github.com/Jraaay/EmptyClassroom")
                }
              >
                Github
              </button>
            </div>
            <div className="campus-settings__line">
              本项目已开源：
              <button
                type="button"
                className="ui-button ui-button--link"
                onClick={() =>
                  window.open("https://github.com/YANG-CHUNXU/EmptyClassroom")
                }
              >
                Github
              </button>
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
}

CampusButtonGroup.propTypes = {
  todayData: PropTypes.object.isRequired,
  selectedCampus: PropTypes.string,
  setSelectedCampus: PropTypes.func,
  setSelectedBuildings: PropTypes.func,
  showClassTime: PropTypes.bool,
  setShowClassTime: PropTypes.func,
  canSelectAllDay: PropTypes.bool,
  setCanSelectAllDay: PropTypes.func,
  useClassTable: PropTypes.bool,
  setUseClassTable: PropTypes.func,
  selectedDate: PropTypes.object,
};

export default CampusButtonGroup;
