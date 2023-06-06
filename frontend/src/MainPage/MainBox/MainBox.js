import React from "react";
import Header from "./Header.js";
import "./MainBox.css";
const MainBox = ({ user }) => {
    return (
        <div className="main-box">
            <Header user={user} />
        </div>
    );

}

export default MainBox;