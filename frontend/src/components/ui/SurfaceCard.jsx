import PropTypes from "prop-types";
import "./ui.css";

function SurfaceCard({ className, children, padded }) {
  const classes = ["ui-card", className];
  if (!padded) {
    classes.push("ui-card--flush");
  }

  return <section className={classes.filter(Boolean).join(" ")}>{children}</section>;
}

SurfaceCard.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node,
  padded: PropTypes.bool,
};

SurfaceCard.defaultProps = {
  className: "",
  children: null,
  padded: true,
};

export default SurfaceCard;
