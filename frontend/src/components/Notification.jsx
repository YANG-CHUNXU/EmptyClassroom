import PropTypes from "prop-types";
import { useEffect, useState } from "react";
import Modal from "./ui/Modal";
import "./Notification.css";

function Notification(props) {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (
      props.todayData.code == 0 &&
      props.todayData.data?.notification != undefined &&
      props.todayData.data.notification.showNotification
    ) {
      setOpen(true);
    }
  }, [props.todayData.code, props.todayData.data?.notification]);

  if (props.todayData.code != 0) {
    return null;
  }

  if (props.todayData.code == 0 && props.todayData.data.notification) {
    return (
      <>
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="notification-button"
        >
          {props.todayData.data.notification.title}
        </button>
        <Modal
          open={open}
          onClose={() => setOpen(false)}
          title={props.todayData.data.notification.title}
        >
          <div
            className="notification-content"
            dangerouslySetInnerHTML={{
              __html: props.todayData.data.notification.content,
            }}
          />
        </Modal>
      </>
    );
  }
  return null;
}

Notification.propTypes = {
  todayData: PropTypes.object.isRequired,
};

export default Notification;
