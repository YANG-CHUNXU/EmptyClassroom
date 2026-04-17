import "./Footer.css";

function Footer() {
  return (
    <div className="footer">
      © 2022-2026 YANG-CHUNXU
      <button
        type="button"
        className="ui-button ui-button--link footer__link"
        onClick={() => window.open("https://github.com/YANG-CHUNXU/EmptyClassroom")}
      >
        Github
      </button>
      <br />
      基于
      <button
        type="button"
        className="ui-button ui-button--link footer__link"
        onClick={() => window.open("https://github.com/Jraaay/EmptyClassroom")}
      >
        原项目
      </button>
      vibe coding 生成
    </div>
  );
}

export default Footer;
