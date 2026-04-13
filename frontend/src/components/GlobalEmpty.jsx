import PropTypes from "prop-types";
import { Card, Empty } from "antd";
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
      : props.todayData.msg ||
        (props.isError
          ? "数据获取失败，请刷新重试"
          : "加载中");

  return (
    <Card
      className="global-empty"
      style={{
        maxWidth: 400,
        width: "90%",
        boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
      }}
      bodyStyle={{
        maxWidth: "340px",
      }}
    >
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description={
          <>
            <div>{description}</div>
            {fallbackEntries.length > 0 && (
              <div style={{ marginTop: 12, textAlign: "left" }}>
                {fallbackEntries.map(([campus, reason]) => (
                  <div key={campus}>
                    {campus}：{reason}
                  </div>
                ))}
              </div>
            )}
          </>
        }
      />
    </Card>
  );
}

GlobalEmpty.propTypes = {
  todayData: PropTypes.object.isRequired,
  isError: PropTypes.bool.isRequired,
};

export default GlobalEmpty;
