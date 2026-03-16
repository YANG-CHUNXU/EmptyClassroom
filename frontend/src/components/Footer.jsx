import { Typography, Button } from "antd";
import { GithubOutlined } from "@ant-design/icons";

function Footer() {
  const { Text } = Typography;
  return (
    <Text>
      © 2022-2026 Jray
      <Button
        onClick={() => window.open("https://github.com/Jraaay/EmptyClassroom")}
        type="text"
        icon={<GithubOutlined />}
      ></Button>
    </Text>
  );
}

export default Footer;
