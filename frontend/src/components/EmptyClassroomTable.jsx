import PropTypes from "prop-types";
import { Empty, Card, Table, Button, Tag, Modal, Descriptions } from "antd";
import { QuestionCircleOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";
import "./EmptyClassroomTable.css";
import CalculateEmptyClassroom from "../utils/calculte";
import dayjs from "dayjs";

function formatEmptyTimeList(emptyTimeList) {
  const classTime = [
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

  const classStartTime = [
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

  let emptyTimeListStr = "";
  if (emptyTimeList[0] == 0) {
    emptyTimeListStr += "00:00";
  } else {
    emptyTimeListStr += "00:00-08:00, " + classTime[emptyTimeList[0] - 1];
  }

  for (let i = 1; i < emptyTimeList.length; i++) {
    if (emptyTimeList[i] - emptyTimeList[i - 1] == 1) {
      continue;
    }

    emptyTimeListStr +=
      "-" +
      classStartTime[emptyTimeList[i - 1] + 1] +
      ", " +
      classTime[emptyTimeList[i] - 1];
  }

  if (emptyTimeList[emptyTimeList.length - 1] != 13) {
    emptyTimeListStr +=
      "-" + classStartTime[emptyTimeList[emptyTimeList.length - 1] + 1];
    emptyTimeListStr += ", " + classTime[classTime.length - 1] + "-24:00";
  } else {
    emptyTimeListStr += "-24:00";
  }

  return emptyTimeListStr;
}

function EmptyClassroomTable(props) {
  const {
    todayData,
    selectedDate,
    selectedCampus,
    selectedBuildings,
    selectedClassTimes,
    setIsError,
    useClassTable,
  } = props;
  const [modalTitle, setModalTitle] = useState("");
  const [modalContent, setModalContent] = useState([]);
  const [openModal, setOpenModal] = useState(false);

  const { emptyClassroom, hasCalculationError } = useMemo(() => {
    if (
      todayData.code != 0 ||
      selectedCampus == "" ||
      selectedBuildings.length == 0 ||
      selectedClassTimes.length == 0
    ) {
      return {
        emptyClassroom: [],
        hasCalculationError: false,
      };
    }

    try {
      return {
        emptyClassroom: CalculateEmptyClassroom(
          todayData.data,
          selectedCampus,
          selectedDate.toDate(),
          selectedBuildings,
          selectedClassTimes
        ),
        hasCalculationError: false,
      };
    } catch {
      return {
        emptyClassroom: [],
        hasCalculationError: true,
      };
    }
  }, [
    todayData.code,
    todayData.data,
    selectedCampus,
    selectedDate,
    selectedBuildings,
    selectedClassTimes,
  ]);

  useEffect(() => {
    if (hasCalculationError) {
      setIsError(true);
    }
  }, [hasCalculationError, setIsError]);

  if (todayData.code != 0) {
    return null;
  }

  if (selectedCampus == "") {
    return null;
  }

  if (
    selectedBuildings.length == 0 ||
    selectedClassTimes.length == 0 ||
    emptyClassroom.length == 0
  ) {
    return (
      <Card
        className="empty-classroom-table"
        style={{
          maxWidth: 400,
          width: "90%",
          boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
        }}
        bodyStyle={{
          maxWidth: "300px",
        }}
      >
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            selectedBuildings.length == 0
              ? selectedClassTimes.length == 0
                ? "请选择教学楼和上课时间"
                : "请选择教学楼"
              : selectedClassTimes.length == 0
              ? "请选择上课时间"
              : "没有空教室了😭"
          }
        />
      </Card>
    );
  }

  function ShowClassroomEmptyInfo(classroomInfo) {
    const data = [
      {
        key: "座位数",
        value: classroomInfo.size,
      },
      {
        key: "类型",
        value: classroomInfo.type,
      },
      {
        key: "空闲时间",
        value: formatEmptyTimeList(classroomInfo.empty_class_time),
      },
      {
        key: "数据来源",
        value: classroomInfo.can_trust ? "教务（可信）" : "课表（参考）",
      },
    ];
    setModalTitle(classroomInfo.name);
    setModalContent(data);
    setOpenModal(true);
  }

  const columns = [
    {
      title: "教室",
      key: "name",
      dataIndex: "name",
      align: "center",
      render: (text, record) => {
        return (
          <span style={{ display: "flex", justifyContent: "center" }}>
            <Button
              size="small"
              onClick={() => {
                ShowClassroomEmptyInfo(record);
              }}
            >
              {text}
            </Button>
          </span>
        );
      },
    },
    {
      title: "座位数",
      key: "size",
      dataIndex: "size",
      align: "center",
    },
    {
      title: "类型",
      key: "type",
      dataIndex: "type",
      align: "center",
    },
    {
      title: (
        <>
          来源
          <Button
            size="small"
            type="text"
            icon={<QuestionCircleOutlined />}
            onClick={() => {
              window.open(
                "https://jraaaaay.feishu.cn/docx/HAu9dbYF1oRb4nxFd7RcugMTnHj#part-Bwf5dnkr0o3G6ExcedZcLzaSnBe"
              );
            }}
          />
        </>
      ),
      key: "can_trust",
      dataIndex: "can_trust",
      align: "center",
      render: (text) => {
        return text ? (
          <Tag color="green" bordered={false}>
            教务
          </Tag>
        ) : (
          <Tag color="red" bordered={false}>
            课表
          </Tag>
        );
      },
    },
  ];

  return (
    <div className="empty-classroom-table">
      <Card
        style={{
          maxWidth: 400,
          width: "90%",
          boxShadow: "0 12px 32px 4px #0000000a, 0 8px 20px #00000014",
        }}
        bodyStyle={{
          padding: "0px",
        }}
      >
        <Table
          dataSource={
            !useClassTable &&
            !(selectedCampus == "海南") &&
            selectedDate.isSame(dayjs(), "day") &&
            (todayData.data.is_fallback == undefined ||
              !todayData.data.is_fallback[selectedCampus])
              ? emptyClassroom.filter((item) => item.can_trust)
              : emptyClassroom
          }
          columns={columns}
          pagination={false}
          bordered={false}
          tableLayout="auto"
          size="small"
          rowKey={(record) => record.name}
          style={{
            width: "100%",
          }}
        />
      </Card>
      <Modal
        title={modalTitle}
        open={openModal}
        footer={null}
        onCancel={() => {
          setOpenModal(false);
        }}
      >
        <div>
          <Descriptions column={1} size="small" layout="vertical">
            {modalContent.map((item, index) => {
              return (
                <Descriptions.Item key={index} label={item.key}>
                  {item.value}
                </Descriptions.Item>
              );
            })}
          </Descriptions>
        </div>
      </Modal>
    </div>
  );
}

EmptyClassroomTable.propTypes = {
  todayData: PropTypes.object,
  selectedDate: PropTypes.object,
  selectedCampus: PropTypes.string,
  selectedBuildings: PropTypes.array,
  selectedClassTimes: PropTypes.array,
  setIsError: PropTypes.func,
  useClassTable: PropTypes.bool,
};

export default EmptyClassroomTable;
