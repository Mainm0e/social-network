import "./AlertBox.css";
import React, { useState } from 'react';

function AlertBox(props) {
  const { title, message, status } = props;
  const [showAlert, setShowAlert] = useState(status);

  const closeAlert = () => {
    setShowAlert(false);
  };

  return (
    <>
      {showAlert ? (
        <div className="alert-box">
          <span className="closebtn" onClick={closeAlert}>
            &times;
          </span>
          <strong>{title}</strong>
          <br />
          <p>{message}</p>
        </div>
      ) : null}
    </>
  );
}

export default AlertBox;
