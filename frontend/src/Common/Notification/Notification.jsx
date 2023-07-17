import { fetchData } from "../../tools/fetchData";
import React, { useState, useEffect } from "react";
import "./Notification.css";
import { getUserId, getCookie } from "../../tools/cookie";

const Notification = ({ clearBox }) => {
  const [notificationData, setNotificationData] = useState([]);
  const [showNotification, setShowNotification] = useState(true);
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
  const handleAcceptDecline = () => {
    setShowNotification(false); // Hide the notification after accepting or declining
  };
  const closeBox = () => {
    window.location.hash = "";
    clearBox();
  };


  const renderNotifications = () => {
    return notificationData.map((notification, index) => (
      <DisplayNotification
        key={index}
        notifications={notification.notifications}
        user={notification.profile}
        handleAcceptDecline={handleAcceptDecline}
      />
    ));
  };
  if (notificationData.length  > 0) {
    console.log(notificationData, "notificationData")
  return (
    <div className="notification-container">
      {showNotification && renderNotifications()}{" "}
      {/* Conditionally render the notifications */}
      <div className="user-list-footer">
        <button onClick={closeBox}>Close</button>
      </div>
    </div>
  );
}
else {
  return (
    <div className="notification-container">
      <div className="notification">
        <div className="notification-content">
          <span>No notifications</span>
        </div>
      </div>
      <div className="user-list-footer">
        <button onClick={closeBox}>Close</button>
      </div>
    </div>
  )
}
};

export default Notification;

const DisplayNotification = ({ notifications, user, handleAcceptDecline }) => {
  if (notifications.type === "follow_request" || notifications.type === "group_request" || notifications.type === "group_invitation") {
    const handleAccept = (value) => {
      console.log(notifications, "notifications in handleAccept")
      const method = "POST";
      const type = "followResponse";
      const payload = {
        sessionId: getCookie("sessionId"),
        receiverId: getUserId("userId"),
        senderId: notifications.senderId,
        groupId: notifications.groupId,
        notifId: notifications.notificationId,
        content: value, // Use the value parameter here
      };
      if (notifications.type === "group_request"){
        payload.receiverId = 0;
      }else if (notifications.type === "group_invitation"){
        payload.senderId = 0;
      }
 
      fetchData(method, type, payload).then((data) => {
        window.location.reload();
        handleAcceptDecline(); // Call the handler function to hide the notification
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
          {console.log(notifications, "notifications")} 
          {notifications.type === "follow_request" && <span>sent you a follow request</span>}
          {notifications.type === "group_request" && <span>sent you a group request</span>}
          {notifications.type === "group_invitation" && <span>invited you to a group</span>}
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
  if (notifications.type === "following") {
    const handleDeleteNotif = () => {
      const method = "POST";
      const type = "followResponse";
      const payload = {
        sessionId: getCookie("sessionId"),
        receiverId: getUserId("userId"),
        senderId: notifications.senderId,
        notifId: notifications.notificationId,
        response: "", // Use the value parameter here
      };
      fetchData(method, type, payload).then((data) => {
        window.location.reload();
        handleAcceptDecline(); // Call the handler function to hide the notification
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
          <span>started following you</span>
        </div>
        <div className="notification-btn">
          <span onClick={handleDeleteNotif}>x</span>
        </div>
      </div>
    );
  }
};
