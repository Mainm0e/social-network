import { useEffect, useState } from "react";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";
import "./invite.css";

const Invite = ({ clearBox }) => {
  const [data, setData] = useState(null);
  const url = new URL(window.location.href);
  const searchParams = new URLSearchParams(url.search);
  const id = searchParams.get("id");
  useEffect(() => {
    const method = "POST";
    const type = "getNonMembers";
    const payload = {
      sessionId: getCookie("sessionId"),
      userId: getUserId("userId"),
      groupId: parseInt(id),
    };
    fetchData(method, type, payload).then((data) => {
      setData(data);
    });
  }, []);

  const closeBox = () => {
    window.location.hash = "";
    clearBox();
  };

  return (
    <div className="invite-container">
      <div className="invite-header">
        <h1>Invite</h1>
      </div>
      <div className="invite-body">
        {data &&
          data.map((u, index) => (
            <DisplayInvite groupId={id} key={index} user={u} />
          ))}
          {data === null && <div>No users to invite</div>}
      </div>
      <div className="invite-footer">
        <button onClick={closeBox}>Close</button>
        </div>
    </div>
  );
};

export default Invite;

const DisplayInvite = ({ groupId, user }) => {
  const sentInv = () => {
    const method = "POST";
    const type = "followRequest";
    const payload = {
      sessionId: getCookie("sessionId"),
      senderId: getUserId("userId"),
      receiverId: parseInt(user.userId),
      groupId: parseInt(groupId),
    };
    console.log("test case", payload);

    fetchData(method, type, payload).then((data) => {
       // refecth page
      });
      window.location.reload()
  };
  return (
    <div className="notification">
      <div className="notification-user">
        <img src={user.avatar} alt="avatar" />
        <span>
          {user.firstName} {user.lastName}
        </span>
      </div>
      <div className="notification-btn">
        <button value="Invite" onClick={() => sentInv()}>
          Invite
        </button>
      </div>
    </div>
  );
};
