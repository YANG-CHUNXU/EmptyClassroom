import PropTypes from "prop-types";
import {
  Radio,
  Button,
  Modal,
  Switch,
  Typography,
  Divider,
} from "antd";
import { useEffect, useState } from "react";
import {
  SettingOutlined,
  GithubOutlined,
  HeartFilled,
} from "@ant-design/icons";
import "./CampusButtonGroup.css";

function CampusButtonGroup(props) {
  const { todayData, selectedCampus, setSelectedCampus, setSelectedBuildings } =
    props;
  const [campusList, setCampusList] = useState([]);

  useEffect(() => {
    if (todayData.code != 0) {
      setCampusList([]);
      return;
    }

    const campusInfoMap = todayData.data?.campus_info_map ?? {};
    const list = Object.keys(campusInfoMap);

    // 排序，西土城在第一，沙河在第二，其他按照字典序
    const order = ["西土城", "沙河"];
    list.sort((a, b) => {
      if (order.indexOf(a) == -1) {
        if (order.indexOf(b) == -1) {
          return a.localeCompare(b);
        } else {
          return 1;
        }
      } else {
        if (order.indexOf(b) == -1) {
          return -1;
        }
        return order.indexOf(a) - order.indexOf(b);
      }
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

  const [openSettingModal, setOpenSettingModal] = useState(false);

  function OpenSettingModal() {
    setOpenSettingModal(true);
  }

  return (
    <div className="campus-button-group">
      <Radio.Group
        value={props.selectedCampus}
        onChange={(e) => {
          props.setSelectedCampus(e.target.value);
          props.setSelectedBuildings([]);
        }}
        buttonStyle="solid"
        size="middle"
      >
        {campusList.map((campus) => {
          return (
            <Radio.Button value={campus} key={campus}>
              {campus}
            </Radio.Button>
          );
        })}
      </Radio.Group>
      <Button
        style={{
          marginLeft: "10px",
        }}
        icon={<SettingOutlined />}
        onClick={OpenSettingModal}
      />
      <Modal
        title="设置"
        open={openSettingModal}
        closable={true}
        footer={null}
        onCancel={() => {
          setOpenSettingModal(false);
        }}
      >
        <div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.showClassTime}
              onChange={(v) => {
                localStorage.setItem("showClassTime", v ? "true" : "false");
                props.setShowClassTime(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              显示课程时间
            </Typography.Title>
          </div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.canSelectAllDay}
              onChange={(v) => {
                localStorage.setItem("canSelectAllDay", v ? "true" : "false");
                props.setCanSelectAllDay(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              全选时选全天
            </Typography.Title>
          </div>
          <div style={{ display: "flex", alignItems: "center" }}>
            <Switch
              defaultChecked={props.useClassTable}
              onChange={(v) => {
                localStorage.setItem("useClassTable", v ? "true" : "false");
                props.setUseClassTable(v);
              }}
              size="small"
            />
            <Typography.Title level={5} style={{ margin: 8 }}>
              非必要情况下也使用课表数据
            </Typography.Title>
          </div>
          <Divider plain>
            <HeartFilled />
          </Divider>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            数据来源：
            <Button
              size="small"
              onClick={() => {
                window.open(
                  "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Zip8dx2rlobE5hxW00CcHwOOnre"
                );
              }}
            >
              了解更多
            </Button>
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            问答Q&A：
            <Button
              size="small"
              onClick={() => {
                window.open(
                  "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj"
                );
              }}
            >
              空教室查询Q&A
            </Button>
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            当前数据刷新时间：
            {new Date(props.todayData.data?.update_at).toLocaleString()}
          </div>
          <div
            style={{
              lineHeight: "2em",
            }}
          >
            项目已开源：
            <Button
              onClick={() =>
                window.open("https://github.com/Jraaay/EmptyClassroom")
              }
              icon={<GithubOutlined />}
              size="small"
            >
              Github
            </Button>
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
};

export default CampusButtonGroup;
