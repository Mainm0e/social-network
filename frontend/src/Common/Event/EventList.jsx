import { useEffect } from "react";
import { fetchData } from "../../tools/fetchData";
import { getCookie, getUserId } from "../../tools/cookie";
import { useState } from "react";

const EventList = ({ clearBox }) => {
  const url = new URL(window.location.href);
  const searchParams = new URLSearchParams(url.search);
  const id = searchParams.get("id");
  const [data, setData] = useState([]);

  useEffect(() => {
    const method = "POST";
    const type = "getGroupEvents";
    const payload = {
      sessionId: getCookie("sessionId"),
      senderId : getUserId("userId"),
      groupId: parseInt(id),
    };
    fetchData(method, type, payload).then((data) => {
      console.log("data",data)
      if (data !== undefined) setData(data.events
        );
    });
  }, [id]);

  const acceptEvent = (event,option) => {
    const method = "POST";
    const type = "participate";
    const payload = {
      sessionId: getCookie("sessionId"),
      groupId: event.groupId,
      eventId: event.eventId,
      memberId: getUserId("userId"),
      option: option,
    };
    fetchData(method, type, payload).then((data) => {
      if (data !== []) setData(data.events
      );
    });
  };

  const renderNotifications = () => {
    return (
      <>
        {data.map((event, index) => (
          <div className="event-list-item" key={index}>
            <div className="event-list-item-header">
              <h2>{event.event.title}</h2>
            </div>
            <div className="event-list-item-body">
              <p>{event.event.content}</p>
              <div className="event-creater">
                <p>Created by: {event.creatorProfile.firstName}</p>
              </div>
            </div>
            <div className="event-list-item-footer">
              <p>{event.event.date}</p>
            </div>
            <div className="event-list-item-footer">
              {event.participate === "going" ? (
                  <p>Going</p>
                ) : (
                  <p>Not going</p>
                )}
              <button onClick={() => acceptEvent(event.event,"going")}>Going</button>
              <button onClick={() => acceptEvent(event.event,"not_going")}>
                Not going
              </button>
            </div>
          </div>
        ))}
      </>
    ); 
  };
 if (data === []){
    return (
      <div className="event-list">
        <div className="event-list-header">
          <h1>Events</h1>
        </div>
        <div className="event-list-body">
          <p>No events</p>
        </div>
        <div className="event-list-footer">
          <button onClick={clearBox}>Close</button>
        </div>
      </div>
    );
  } else   if (data !== undefined) {
    return (
      <div className="event-list">
        <div className="event-list-header">
          <h1>Events</h1>
        </div>
        <div className="event-list-body">
         {renderNotifications()}  
        </div>
        <div className="event-list-footer">
          <button onClick={clearBox}>Close</button>
        </div>
      </div>
    );
    }
};

export default EventList;
