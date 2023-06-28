import { fetchData } from "../../tools/fetchData";
import "./Notification.css";
import { getUserId, getCookie } from "../../tools/cookie";

const Notification = ({ data }) => {
  console.log("Notification", data);

  const renderNotifications = () => {
    return data.map((notification, index) => (
      <DisplayNotification
        key={index}
        notifications={notification.notifications}
        user={notification.profile}
      />
    ));
  };
  return <div className="notification-container">{renderNotifications()}</div>;
};

export default Notification;

const DisplayNotification = ({ notifications, user }) => {
  if (notifications.type === "follow_request") {
    const handleAccept = (value) => {
        const method = 'POST';
        const type = "followResponse";
        const payload = {
          sessionId: getCookie('sessionId'),
          userId: getUserId("userId"),
          followId: notifications.followerId,
          response: value // Use the value parameter here
        };
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
