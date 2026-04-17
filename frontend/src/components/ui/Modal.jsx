import PropTypes from "prop-types";
import "./ui.css";

function Modal({ open, title, onClose, children, width = "min(92vw, 32rem)" }) {
  if (!open) {
    return null;
  }

  return (
    <div
      className="ui-modal-backdrop"
      onClick={onClose}
      role="presentation"
    >
      <div
        className="ui-modal"
        onClick={(event) => event.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-label={title}
        style={{ width }}
      >
        <div className="ui-modal__header">
          <h2 className="ui-modal__title">{title}</h2>
          <button
            type="button"
            className="ui-button ui-button--ghost ui-button--icon"
            onClick={onClose}
            aria-label="关闭"
          >
            ×
          </button>
        </div>
        <div className="ui-modal__body">{children}</div>
      </div>
    </div>
  );
}

Modal.propTypes = {
  open: PropTypes.bool.isRequired,
  title: PropTypes.string,
  onClose: PropTypes.func.isRequired,
  children: PropTypes.node,
  width: PropTypes.string,
};

Modal.defaultProps = {
  title: "",
  children: null,
  width: "min(92vw, 32rem)",
};

export default Modal;
