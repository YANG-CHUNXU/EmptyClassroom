import PropTypes from "prop-types";
import SurfaceCard from "./ui/SurfaceCard";
import "./BuildingPicker.css";

function BuildingPicker(props) {
  // const [style, setStyle] = useState(true);
  if (props.todayData.code != 0) {
    return null;
  }

  if (props.selectedCampus == "") {
    return null;
  }

  const campusInfo =
    props.todayData.data?.campus_info_map?.[props.selectedCampus];
  if (!campusInfo) {
    return null;
  }

  const building_name_id_map = campusInfo.building_id_map ?? {};

  const options = [];
  for (const [key, value] of Object.entries(building_name_id_map)) {
    options.push({
      label: key,
      value: value,
    });
  }
  options.sort((a, b) => {
    if (a.label.length - b.label.length != 0) {
      return a.label.length - b.label.length;
    } else {
      return a.label.localeCompare(b.label);
    }
  });

  return (
    <SurfaceCard className="building-picker">
      {options.map((item) => (
        <button
          type="button"
          key={props.selectedCampus + item.value}
          className={`building-picker__button ${
            props.selectedBuildings.includes(item.value) ? "primary" : ""
          }`}
          onClick={() => {
            if (props.selectedBuildings.includes(item.value)) {
              props.setSelectedBuildings(
                props.selectedBuildings.filter((x) => x != item.value)
              );
            } else {
              props.setSelectedBuildings([
                ...props.selectedBuildings,
                item.value,
              ]);
            }
          }}
        >
          {item.label}
        </button>
      ))}
    </SurfaceCard>
  );
}

BuildingPicker.propTypes = {
  todayData: PropTypes.object,
  selectedBuildings: PropTypes.array,
  setSelectedBuildings: PropTypes.func,
  selectedCampus: PropTypes.string,
};

export default BuildingPicker;
