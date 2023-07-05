import React, {useEffect, useState} from "react";
import {logout} from "../../tools/logout";
import NavList from "./NavList";
import "./LeftBox.css";
import {home,link_notifications} from "../../tools/link";
import {getUserId, getCookie} from "../../tools/cookie";
import {fetchData} from "../../tools/fetchData";

const LeftBox = ({user,link}) => {

    //get notication number
    const [notification, setNotification] = useState(0);
    useEffect(() => {
        const method = "POST";
        const type = "requestNotif";
    
        const payload = {
          userId: getUserId("userId"),
          sessionId: getCookie("sessionId"),
        };
        fetchData(method, type, payload).then((data) => {
          console.log(data.length);
          setNotification(data.length);
        });
        }, []);

    /* get screen side */
    const [screensize, setScreensize] = useState("");
    useEffect(() => {
    let size = window.innerWidth
    if (size < 1400) {
        setScreensize("medium");
    }
    if (size > 1400) {
        setScreensize("rare");
    }
    }, []);
    window.addEventListener("resize", () => {
       let  size = window.innerWidth
        if (size < 1400) {
            setScreensize("medium");
        }
        if (size > 1400) {
            setScreensize("rare");
        }

        //
    });
    const handleLogout = () => {
        logout();
    }
    const checkImage = () => {
        if (user.avatar === ''|| user.avatar === null|| user.avatar === undefined) {
            return null;
        } else {
            return  <div className="post_image "> <img src={user.avatar} alt="content"  onClick={link_notifications}/> </div>;
        }
    };
    const [show, setShow] = useState(false);
    const showLinks = () => {
        setShow(!show);
    }
    return (
        <div className="left-box">
            <div className="user_box">
                <div className="img_box">
                    {checkImage()}
                </div>
                <div className="user_info">
                    <div className="username">
                        <span onClick={home}>{user.firstName}</span>
                    </div>
                    <div className="btns">
                    <div className="logout_btn">
                        <button className="btn" onClick={handleLogout}>logout</button>
                    </div>
                    <div className="link_btn">
                        <button className="btn" onClick={showLinks}>menu </button>
                    </div>
                    </div>
                </div>
            </div>
            {screensize === "medium" && <LinkBox link={link} type={"medium"} show={show} /> ||
            screensize === "rare" && <LinkBox link={link} type={"rare"} show={true}/> }
        </div>
    );
}

export default LeftBox;

const LinkBox = ({link,type,show}) => {
    
    if (show === true) {
    if (type ==="rare"){
        return (
            <div className="links_full">
            <NavList type={"Main"} links={link[3]} />
            <NavList type={"Post"} links={link[1]} />
            </div>
        )
    }
    if (type === "medium"){
       return ( 
       <div className="links_medium">
            <NavList type={"Main"} links={link[3]} />
            <NavList type={"Post"} links={link[1]} />
        </div>
       )
    }
    } else {
        return null;
    }

}

