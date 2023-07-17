import React, { useState, useEffect } from "react";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";
import { profile } from "../../tools/link";
import "./UserList.css";
const Memberlist = ({ id, clearBox }) => {
  const [memberlist, setMemberlist] = useState(null);
  const closeBox = () => {
    window.location.hash = "";
    clearBox();
  };
  useEffect(() => {
    const method = "POST";
    const type = "exploreGroups";
    const payload = {
      sessionId: getCookie("sessionId"),
      userId: getUserId("userId"),
    };
    fetchData(method, type, payload).then((data) => {
      for (let i = 0; i < data.length; i++) {
        if (data[i].groupId === id) {
          setMemberlist(data[i].members);
        }
      }
    });
  }, []);
  if (memberlist === null) {
    return (
      <div className="loading">
        <div>Loading...</div>
      </div>
    );
  } else {
    return (
      <div className="user-list">
        <div className="user-list-container">
          <div className="user-list-header">
            <h2>Members</h2>
          </div>
          <div className="user-list-body">
            <ul>
              {memberlist.map((user) => (
                <li key={user.userId}>
                  <div className="user-list-item" onClick={() => profile(user.userId)}>
                    <div className="user-list-item-left">
                      <img src={user.avatar} alt="user" />
                    </div>
                    <div className="user-list-item-right">
                      <span>{user.firstName}</span>
                      <span> </span>
                      <span>{user.lastName}</span>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
          <div className="user-list-footer">
            <button onClick={closeBox}>Close</button>
          </div>
        </div>
      </div>
    );
  }
};

export default Memberlist;
