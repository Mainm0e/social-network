import React, { useState, useEffect } from "react";
import { getCookie } from "../../tools/cookie";


const UserList = ({title,id}) => {
    const [data, setData] = useState(null);
    console.log("in userlist",title,"wtf",id)
    useEffect(() => {
    const getUserList = async () => {
        const response = await fetch("http://localhost:8080/api", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ type: "profileList", payload: {sessionId:getCookie("sessionId"), userId: id, request:title}}),
        });
        const responseData = await response.json();
        setData(responseData);
        console.log("in getUserList",data)
    }
    getUserList();
    }, []);
    if (data===null) {
        return <div>Loading...</div>;
    }else{
    return (
        <div className="user-list">
            <div className="user-list-container">
                <div className="user-list-header">
                    <h2>Users</h2>
                </div>
                <div className="user-list-body">
                    <ul>
                    <p>{console.log(data)}</p>
                    </ul>
                </div>
                <div className="user-list-footer">
                    <button>Close</button>
                </div>
            </div>
        </div>
    );
    };
}

export default UserList;