import { fetchData } from "../../tools/fetchData";
import "./Notification.css"
import { getUserId, getCookie } from "../../tools/cookie";

const Notification = ({data}) => {
    console.log("Notification",data)

    const renderNotifications = () => {
        return data.map((notification, index) => (
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
        </div>
    )
}

export default Notification

const DisplayNotification = ({ notifications, user }) => {

    if (notifications.type === "follow_request"){

        const handleAccept = () => {
            const method = 'POST';
            const type = "acceptFollowRequest";
            const payload = {
                userId: getUserId("userId"),
                sessionId: getCookie('sessionId'),
                followerId: notifications.followerId,
            };
            fetchData(method, type, payload).then((data) => {
                console.log(data);
            }
            );
        }
        const handleDecline = () => {
           /*  fetchData(method, type, payload).then((data) => {
                console.log(data);
            }
            ); */
        }
        
        return (
            <div className="notification">
                <div className="notification-user">
                    <img src={user.avatar} alt="avatar" />
                    <span>{user.firstName} {user.lastName}</span>
                </div>
                <div className="notification-content">
                    <span>sent you a follow request</span>
                </div>
                <div className="notification-btn">
                    <button onClick={handleAccept}>Accept</button>
                    <button onClick={handleDecline}>Decline</button>
                </div>
            </div>
        );
    }
  };