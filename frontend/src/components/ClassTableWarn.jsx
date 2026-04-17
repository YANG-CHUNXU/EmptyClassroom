import PropTypes from "prop-types";
import dayjs from "dayjs";
import { useState } from "react";
import "./ClassTableWarn.css";

function ClassTableWarn(props) {
  const [confirmingClose, setConfirmingClose] = useState(false);
  let hasClassTableData = false;

  if (props.selectedCampus == "海南") {
    hasClassTableData = true;
  }
  if (props.useClassTable) {
    hasClassTableData = true;
  }
  if (!props.selectedDate.isSame(dayjs(), "day")) {
    hasClassTableData = true;
  }
  if (props.todayData.code == 0 && props.todayData.data.is_fallback != undefined) {
    if (props.todayData.data.is_fallback[props.selectedCampus]) {
      hasClassTableData = true;
    }
  }

  if (!hasClassTableData || props.dontWarnClassTable) {
    return null;
  }

  return (
    <div className="class-table-warn">
      <div className="class-table-warn__content">
        当前空教室数据包含来自课表的数据，相比教务数据，可能不准确。
      </div>
      <div className="ui-inline-actions">
        <button
          type="button"
          className="ui-button ui-button--link"
          onClick={() => {
            window.open(
              "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-UykqdD8nboWEi9xo6jTcbNAunL2"
            );
          }}
        >
          为什么
        </button>
        {!confirmingClose ? (
          <button
            type="button"
            className="ui-button ui-button--ghost ui-button--danger"
            onClick={() => setConfirmingClose(true)}
          >
            不再提醒
          </button>
        ) : (
          <>
            <span className="class-table-warn__confirm">确认永久关闭提醒？</span>
            <button
              type="button"
              className="ui-button ui-button--primary"
              onClick={() => {
                localStorage.setItem("dontWarnClassTable", "true");
                props.setDontWarnClassTable(true);
              }}
            >
              确定
            </button>
            <button
              type="button"
              className="ui-button ui-button--ghost"
              onClick={() => setConfirmingClose(false)}
            >
              取消
            </button>
          </>
        )}
      </div>
    </div>
  );
}

ClassTableWarn.propTypes = {
  selectedDate: PropTypes.object,
  selectedCampus: PropTypes.string,
  useClassTable: PropTypes.bool,
  todayData: PropTypes.object,
  dontWarnClassTable: PropTypes.bool,
  setDontWarnClassTable: PropTypes.func,
};

export default ClassTableWarn;
