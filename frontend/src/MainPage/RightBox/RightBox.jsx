import React, { useState, useEffect } from "react";
import "./RightBox.css";
import UserList from "../../Common/UserList/UserList";

const RightBox = ({ profileId }) => {
  const [box, setBox] = useState(null);

  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.substring(1); // Remove the '#' character

      // Clear the box and create a new UserList component after a slight delay
      clearBox();
      setTimeout(() => {
        if (hash === "followers") {
          setBox(<UserList title="followers" id={profileId} clearBox={clearBox} />);
        } else if (hash === "followings") {
          setBox(<UserList title="followings" id={profileId} clearBox={clearBox} />);
        } else {
          setBox(<div className="loading">Loading...</div>);
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
  }, [profileId]);

  const clearBox = () => {
    setBox(null);
  };

  return <div className="right-box">{box}</div>;
};

export default RightBox;
