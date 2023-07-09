import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import GroupHeader from "./Group/GroupHeader";
import GroupBody from "./Group/GroupBody";
import Body from "./User/Body";
import Explore from "../../Common/explore/explore";
import RegisterGroup from "./Group/RegisterGroup";
import "./MainBox.css";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";



const MainBox = ({ profileId,groupId, type ,state}) => {
  const [refreshKey, setRefreshKey] = useState(0);
  // refreshKey is used to refresh component
  const refreshComponent = () => {
    setRefreshKey((prevKey) => prevKey + 1);
  };

  if (state === "explore"){
    if (type ==="user"){
      return (
        <Explore type={"exploreUsers"}/> 
        )
      } else if (type === "group"){
        return (
          <Explore type={"exploreGroups"}/>
        )
      } else {
        return (
          <Explore type={"exploreUsers"}/> 
          )
      }
  } else if (state === "profile"){
    if (type === "user"){
      return(
        <Profile   key={refreshKey} profileId={profileId}   refreshComponent={refreshComponent}/>
        )
      } else if (type === "group"){
        return (
          <Group key={refreshKey} groupId={groupId} refreshComponent={refreshComponent} />
        )
      }
} else if (state === "create_group"){
  if (type === "group"){
    return(
      <CreateGroup/>
    )
  }
};
}
export default MainBox;

const Profile = ({profileId,refreshComponent}) =>{
   const [data, setData] = useState(null);
  useEffect(() => {
    const method = "POST"
    const type = "profile"
    const payload = { sessionId: getCookie("sessionId"), userId: getUserId("userId"), profileId: profileId }

    fetchData(method,type, payload).then((data) => setData(data) );
  }, [ profileId]);

  const handleRefresh = () => {
    refreshComponent(); // Call the refresh function from the parent component
  };

  
  if (data === null) {
    return <div className="loading"><div>Loading...</div></div>;
  } else {
    return (
      <div className="main-box">
        <Header profile={data} handleRefresh={handleRefresh} />
        <Body id={profileId}/>
      </div>
    );
  }
}

const CreateGroup = () => {
  const [data, setData] = useState(null);
  useEffect(() => {
    const method = "POST"
    const type = "profile"
    const payload = { sessionId: getCookie("sessionId"), userId: getUserId("userId"), profileId: getUserId("userId") }

    fetchData(method,type, payload).then((data) => setData(data) );
  }, []);
  return (
    <div className="main-box">
     <RegisterGroup user={data}/>
    </div>
  );
}

const Group = ({ groupId, refreshComponent }) => {
  const [data, setData] = useState(null);
 /*  useEffect(() => {
    const method = "POST"
    const type = "group"
    const payload = { 
      sessionId: getCookie("sessionId"), 
      userId: getUserId("userId"), 
      groupId: groupId }
    fetchData(method,type, payload).then((data) => setData(data) );
  }, [groupId]); */
  useEffect(() => {
    const method = "POST"
   const type = "exploreGroups"
    const payload = { sessionId: getCookie("sessionId"), userId: getUserId("userId")}
    fetchData(method,type,payload).then((data)=>{
        setData(data)
    })
    }, []);
  const handleRefresh = () => {
    refreshComponent(); // Call the refresh function from the parent component
  };
  if (data === null) {
    return <div className="loading"><div>Loading...</div></div>;
  }
  else {
    let group = null
    for (let i = 0; i < data.length; i++) {
      if (parseInt(data[i].groupId) === parseInt(groupId)) {
        group = data[i]
      }
    }
    return (
      <div className="main-box">
        <GroupHeader group={group} handleRefresh={handleRefresh} />
        <GroupBody id={groupId}/>
      </div>
    );
  }
}