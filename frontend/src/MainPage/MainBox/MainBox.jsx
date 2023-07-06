import React, { useState, useEffect } from "react";
import Header from "./User/Header";
import GroupHeader from "./Group/GroupHeader";
import Body from "./User/Body";
import Explore from "../../Common/explore/explore";
import RegisterGroup from "./Group/RegisterGroup";
import "./MainBox.css";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";



const MainBox = ({ profileId, type ,state}) => {
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
        <Group key={refreshKey} groupId={profileId} refreshComponent={refreshComponent} />
        return(
          <>
          <p>hello</p>
          </>
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
        <Body user={profileId}/>
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
  useEffect(() => {
    const method = "POST"
    const type = "group"
    const payload = { 
      sessionId: getCookie("sessionId"), 
      userId: getUserId("userId"), 
      groupId: groupId }
    fetchData(method,type, payload).then((data) => setData(data) );
  }, [groupId]);
  const handleRefresh = () => {
    refreshComponent(); // Call the refresh function from the parent component
  };
  if (data === null) {
    return <div className="loading"><div>Loading...</div></div>;
  }
  else {
    return (
      <div className="main-box">
        <GroupHeader group={data} handleRefresh={handleRefresh} />
      </div>
    );
  }
}