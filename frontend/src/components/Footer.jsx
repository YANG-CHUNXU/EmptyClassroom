import { Typography, Button } from "antd";
import { GithubOutlined } from "@ant-design/icons";

function Footer() {
  const { Text } = Typography;
  return (
    <Text>
      © 2022-2026 YANG-CHUNXU
      <Button
        onClick={() => window.open("https://github.com/YANG-CHUNXU/EmptyClassroom")}
        type="text"
        icon={<GithubOutlined />}
      ></Button>
      <br />
      基于 <Button
        onClick={() => window.open("https://github.com/Jraaay/EmptyClassroom")}
        type="text"
        icon={<GithubOutlined />}
        style={{ padding: "0 4px" }}
      >原项目</Button> vibe coding 生成
    </Text>
  );
}

export default Footer;
