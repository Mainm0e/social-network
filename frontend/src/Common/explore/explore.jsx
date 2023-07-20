import React, { useState, useEffect } from "react";
import { getCookie,getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";
import "./explore.css"
;

const Explore = ({type}) => {
    //asdf
    const [data, setData] = useState(null);
    useEffect(() => {
        const method = "POST"
        const payload = { sessionId: getCookie("sessionId"), userId: getUserId("userId")}
        fetchData(method,type,payload).then((data)=>{
            setData(data)
        })
        }, []);  
        const followRequest = async (id) => {
        const method = "POST"
        const type = "followRequest"
        const payload ={ sessionId: getCookie("sessionId"), senderId: getUserId("userId"), groupId:id}
        fetchData(method,type,payload).then((data) => {})
        /* refect */
        window.location.reload();
      };
        
        const generateExploreList = () => {
            if (type === "exploreUsers"){
            return data.map((user) => {
                return (
                    <div className="explore_list_item" key={user.userId}  userid={user.userId}   onClick={() => navigateToProfile(type,user.userId)}>
                        <div className="explore_list_item_left">
                            <img src={user.avatar} alt="profile" />
                        </div>
                        <div className="explore_list_item_right">
                            <div className="explore_list_item_right_top">
                                <h3>{user.firstName} {user.lastName}</h3>
                            </div>
                        </div>
                    </div>
                )
            })
        } else if (type === "exploreGroups"){
            return data.map((group) => {
                return (
                    <div className="explore_list_item" key={group.groupId}  userid={group.groupId}   /* onClick={() => navigateToProfile(type,group.groupId)} */>
                        <div className="explore_list_item_title">
                            <h3>{group.title}</h3>
                        </div>
                        <div className="explore_list_item_create_info">
                            <div className="explore_list_item_create_info_left">
                                <img src={group.creatorProfile.avatar} alt="profile" />
                            </div>
                            <div className="explore_list_item_create_info_right">
                                    <h3>{group.creatorProfile.firstName}</h3>
                            </div>
                        </div>
                        <div className="explore_list_item_create_time">
                            <p>{group.date}</p>
                        </div>
                        {group.status !== "member" ? 
                        <div className="explore_list_item_follow_btn">
                        <button onClick={() => followRequest(group.groupId)}>{group.status}</button>
                    </div>: <div className="explore_list_item_go_to_page">
                            <button onClick={() => navigateToProfile(type,group.groupId)}>Go to page</button>
                        </div>  }
                    </div>
                )
            })
        }
    }
   
    if (!data) {
        return <div>Loading...</div>;
    } else {
    return (
    <div className="explore">
        <div className="explore_top">
            <h1>Explore</h1>
        </div>
        <div className="explore_body">
            <div className="explore_list">
              {generateExploreList()}
            </div>
        </div>
    </div>
    )
    }
}
export default Explore;

const navigateToProfile = (type,userId) => {
    let linkType ="user"
    if (type === "exploreUsers"){
        linkType = "user"
    } else if (type ==="exploreGroups"){
        linkType = "group"
    }
    return window.location.href = `/${linkType}?id=${userId}#postlist`;
  };