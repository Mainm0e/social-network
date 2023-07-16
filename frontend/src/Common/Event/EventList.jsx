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
      senderId: getUserId("userId"),
      groupId: parseInt(id),
    };
    fetchData(method, type, payload).then((data) => {
      console.log("data", data);
      if (data !== undefined) setData(data.events);
    });
  }, [id]);

  const acceptEvent = (event, option, status) => {
    if (status!==option) {
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
      if (data !== []) setData(data.events);
    });
    closeBox();
  }
  };
  const getGoingAndNotGoing = (event) => {
    if (event.participants !== null) {
     if (event.participants.going !== undefined && event.participants.not_going !== undefined) {
      return (
        <div className="event-list-item-footer">
          <p>Not going: {event.participants.not_going.length}</p>
          <p>Going: {event.participants.going.length}</p>
        </div>
      );
    } else if (event.participants.going !== undefined) {  
      return (
        <div className="event-list-item-footer">
          <p>Not going: 0</p>
          <p>Going: {event.participants.going.length}</p>
        </div>
      );
    } else if (event.participants.not_going !== undefined) {
      return (
        <div className="event-list-item-footer">
          <p>Not going: {event.participants.not_going.length}</p>
          <p>Going: 0</p>
        </div>
      );
    } else {
      return (
        <div className="event-list-item-footer">
          <p>Not going: 0</p>
          <p>Going: 0</p>
        </div>
      );
    }
  } else {
    return (
      <div className="event-list-item-footer">
        <p>Not going: 0</p>
        <p>Going: 0</p>
      </div>
    );
  };
  };

  const renderNotifications = () => {
    return (
      <>
        {data.map((event, index) => (
          <div className="event-list-item" key={index}>
            <div className="event-list-item-header">
              <h1>{event.event.title}</h1>
            </div>
            <div className="event-list-item-body">
              <p>{event.event.content}</p>
              <div className="event-creater">
                <p>Created by: {event.creatorProfile.firstName}</p>
              </div>
            </div>
            <div className="event-list-item-footer">
              <p>date :{event.event.date}</p>
             {getGoingAndNotGoing(event)}
            </div>
            <div className="event-list-item-footer">
              <button
                className={event.status === "not_going" ? "highlight" : "default"}
                onClick={() => acceptEvent(event.event, "not_going",event.status)}
              >
                Not going
              </button>
              <button
                className={event.status === "going" ? "highlight" : "default"}
                onClick={() => acceptEvent(event.event, "going",event.status)}
              >
                Going
              </button>
            </div>
          </div>
        ))}
      </>
    );
  };
  const closeBox = () => {
    window.location.hash = "";
    clearBox();
  };
  if (data === []) {
    return (
      <div className="event-list">
        <div className="event-list-header">
          <h1>Events</h1>
        </div>
        <div className="event-list-body">
          <p>No events</p>
        </div>
        <div className="event-list-footer">
          <button onClick={closeBox}>Close</button>
        </div>
      </div>
    );
  } else if (data !== undefined) {
    return (
      <div className="event-list">
        <div className="event-list-header">
          <h1>Events</h1>
        </div>
        <div className="event-list-body">{renderNotifications()}</div>
        <div className="event-list-footer">
          <button onClick={closeBox}>Close</button>
        </div>
      </div>
    );
  }
};

export default EventList;
