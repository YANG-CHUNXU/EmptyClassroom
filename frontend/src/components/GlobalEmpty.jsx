import PropTypes from "prop-types";
import EmptyState from "./ui/EmptyState";
import SurfaceCard from "./ui/SurfaceCard";
import "./GlobalEmpty.css";

function GlobalEmpty(props) {
  const campusInfoMap = props.todayData.data?.campus_info_map ?? {};
  const hasCampusData = Object.keys(campusInfoMap).length > 0;

  if (props.todayData.code == 0 && hasCampusData) {
    return null;
  }

  const fallbackReason = props.todayData.data?.fallback_reason ?? {};
  const fallbackEntries = Object.entries(fallbackReason);
  const description =
    props.todayData.code == 0
      ? props.todayData.data?.empty_reason ?? "当前暂无可用教室数据"
      : props.todayData.msg || (props.isError ? "数据获取失败，请刷新重试" : "加载中");

  return (
    <SurfaceCard className="global-empty">
      <EmptyState
        title={description}
        description={
          fallbackEntries.length > 0 ? (
            <div className="global-empty__fallback">
              {fallbackEntries.map(([campus, reason]) => (
                <div key={campus}>
                  {campus}：{reason}
                </div>
              ))}
            </div>
          ) : null
        }
      />
    </SurfaceCard>
  );
}

GlobalEmpty.propTypes = {
  todayData: PropTypes.object.isRequired,
  isError: PropTypes.bool.isRequired,
};

export default GlobalEmpty;
