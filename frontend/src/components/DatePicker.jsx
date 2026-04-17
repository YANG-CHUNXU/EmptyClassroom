import PropTypes from "prop-types";
import { Typography } from "antd";
import "./DatePicker.css";

function DatePicker(props) {
  if (props.todayData.code != 0) {
    return null;
  }
  return (
    <div className="date-picker">
      <Typography.Text style={{ fontSize: "16px", fontWeight: 500 }}>
        {props.selectedDate ? props.selectedDate.format("YYYY-MM-DD") : ""}
      </Typography.Text>
    </div>
  );
}

DatePicker.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  setSelectedDate: PropTypes.func,
};

export default DatePicker;
