import PropTypes from "prop-types";
import "./ui.css";

function EmptyState({ title, description }) {
  return (
    <div className="ui-empty-state">
      <div className="ui-empty-state__icon" aria-hidden="true">
        □
      </div>
      <div className="ui-empty-state__title">{title}</div>
      {description ? (
        <div className="ui-empty-state__description">{description}</div>
      ) : null}
    </div>
  );
}

EmptyState.propTypes = {
  title: PropTypes.node.isRequired,
  description: PropTypes.node,
};

EmptyState.defaultProps = {
  description: null,
};

export default EmptyState;
