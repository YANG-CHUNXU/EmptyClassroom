import PropTypes from "prop-types";
import "./DatePicker.css";

function DatePicker(props) {
  if (props.todayData.code != 0) {
    return null;
  }

  return (
    <div className="date-picker">
      {props.selectedDate ? props.selectedDate.format("YYYY-MM-DD") : ""}
    </div>
  );
}

DatePicker.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  setSelectedDate: PropTypes.func,
};

export default DatePicker;
