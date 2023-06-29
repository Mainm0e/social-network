import { fetchData } from "../../tools/fetchData";
import React, { useState, useEffect } from "react";
import "./Notification.css";
import { getUserId, getCookie } from "../../tools/cookie";

const Notification = ({clearBox}) => {
  const [notificationData, setNotificationData] = useState([]);
  useEffect(() => {
    const method = "POST";
    const type = "requestNotif";

    const payload = {
      userId: getUserId("userId"),
      sessionId: getCookie("sessionId"),
    };
    fetchData(method, type, payload).then((data) => {
      console.log(data);
      setNotificationData(data);
    });
  }, []);

  if (
    notificationData === undefined ||
    notificationData.lenght < 1 ||
    notificationData === null
  )
    return null;

  const renderNotifications = () => {
    return notificationData.map((notification, index) => (
      <DisplayNotification
        key={index}
        notifications={notification.notifications}
        user={notification.profile}
      />
    ));
  };
  return (
    <div className="notification-container">
      {renderNotifications()}
      <div className="user-list-footer">
        <button onClick={clearBox}>Close</button>
      </div>
    </div>
  );
};

export default Notification;

const DisplayNotification = ({ notifications, user }) => {
  if (notifications.type === "follow_request") {
    const handleAccept = (value) => {
      const method = "POST";
      const type = "followResponse";
      const payload = {
        sessionId: getCookie("sessionId"),
        followeeId: getUserId("userId"),
        followerId: notifications.senderId,
        notifId: notifications.notificationId,
        response: value, // Use the value parameter here
      };
      console.log("request response", payload);
      fetchData(method, type, payload).then((data) => {
        console.log(data);
      });
    };

    return (
      <div className="notification">
        <div className="notification-user">
          <img src={user.avatar} alt="avatar" />
          <span>
            {user.firstName} {user.lastName}
          </span>
        </div>
        <div className="notification-content">
          <span>sent you a follow request</span>
        </div>
        <div className="notification-btn">
          <button value="accept" onClick={() => handleAccept("accept")}>
            Accept
          </button>
          <button value="reject" onClick={() => handleAccept("reject")}>
            Decline
          </button>
        </div>
      </div>
    );
  }
};
