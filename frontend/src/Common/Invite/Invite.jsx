import { useEffect, useState } from "react";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";

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
    console.log("test case", payload);

    fetchData(method, type, payload).then((data) => {
      setData(data);
    });
  }, []);

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
      </div>
    </div>
  );
};

export default Invite;

const DisplayInvite = ({ groupId, user, status }) => {
  const sentInv = (e) => {
    console.log(e);
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
      <div className="notification-btn">
        <button value="Invite" onClick={() => sentInv(9)}>
          Invite
        </button>
      </div>
    </div>
  );
};
