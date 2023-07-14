import React, { useState, useEffect } from "react";
import "./RightBox.css";
import UserList from "../../Common/UserList/UserList";
import Notification from "../../Common/Notification/Notification";
import Invite from "../../Common/Invite/Invite";
import EventList from "../../Common/Event/EventList"; 
import Memberlist from "../../Common/UserList/MemberList";

const RightBox = () => {
  const [box, setBox] = useState(null);
  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.substring(1); // Remove the '#' character
      // Get the URL
       // Get the query parameters from the URL
    const searchParams = new URLSearchParams(window.location.search);

    // Get the value of the 'id' parameter
    const id = searchParams.get('id');
      // Clear the box and create a new UserList component after a slight delay
      clearBox();
      setTimeout(() => {
        if (hash === "followers") {
          setBox(
            <UserList title="followers" id={parseInt(id)} clearBox={clearBox} />
          );
        } else if (hash === "followings") {
          setBox(
            <UserList title="followings" id={parseInt(id)} clearBox={clearBox} />
          );
        } else if (hash === "notifications") {
          setBox(<Notification clearBox={clearBox} />);
        } else if (hash === "invite_to_group") {
          setBox(<Invite clearBox={clearBox} />);
        } else if (hash === "eventlist") {
          setBox(<EventList clearBox={clearBox} />);
        } else if (hash === "memberlist"){
          setBox(<Memberlist id={parseInt(id)} clearBox={clearBox} />);
        } else {
          setBox(null);
        }
      }, 10);
    };

    handleHashChange(); // Call the function initially to handle the current URL hash

    // Add event listener to handle hash changes
    window.addEventListener("hashchange", handleHashChange);

    return () => {
      // Cleanup the event listener on component unmount
      window.removeEventListener("hashchange", handleHashChange);
    };
  }, []);

  const clearBox = () => {
    // delete # from url
    setBox(null);
  };
  if (box === null) {
    return null;
  } else {
    return <div className="right-box">{box}</div>;
  }

};

export default RightBox;
