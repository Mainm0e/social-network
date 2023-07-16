import React, { useState, useEffect } from "react";
import { getCookie } from "../../tools/cookie";
import "./UserList.css";
import { fetchData } from "../../tools/fetchData";
import { profile } from "../../tools/link";

const UserList = ({ title, id, clearBox }) => {
  const [profilename, setProfilename] = useState("");
  const [data, setData] = useState(null);
  const closeBox = () => {
    window.location.hash = "";
    clearBox();
  };
  const getProfileName = () => {
    // get element by id profile-firstName
    const firstName = document.getElementById("profile-firstName");
    setProfilename(firstName.innerHTML);
  }
  useEffect(() => {
    const method = "POST";
    const type = "profileList";
    const payload = {
      sessionId: getCookie("sessionId"),
      userId: id,
      request: title,
    };
    console.log("idont no",payload);
    fetchData(method, type, payload).then((data) => {
      setData(data);
    });
    getProfileName();
  }, []);
  if (data === null) {
    //no follower or following
    return (
        <div className="notification-container">
          <div className="notification">
            <div className="notification-content">
              <span>No {title}</span>
            </div>
          </div>
          <div className="user-list-footer">
            <button onClick={closeBox}>Close</button>
          </div>
        </div>
      );
  } else {
    return (
      <div className="user-list">
        <div className="user-list-container">
          <div className="user-list-header">
            {title === "followers" ? (
              <h2>{profilename} is being followed by:</h2>
            ) : (<></>)}
            {title === "followings" ? (
              <h2>{profilename} is following</h2>
            ) : (<></>)}
          </div>
          <div className="user-list-body">
            <ul>
              {data.map((user) => (
                <li key={user.userId}>
                  <div className="user-list-item" onClick={() => profile(user.userId)}>
                    <div className="user-list-item-left">
                      <img src={user.avatar} alt="user" />
                    </div>
                    <div className="user-list-item-right" >
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

export default UserList;
