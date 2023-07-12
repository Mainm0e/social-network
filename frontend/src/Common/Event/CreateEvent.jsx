import React, { useState } from "react";
import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";
import "./Event.css";
const CreateEvent = ({profileId,groupId}) => {
  // Create Event Form

  const [eventTitle, setEventTitle] = useState("");
    const [eventDescription, setEventDescription] = useState("");
    const [eventDate, setEventDate] = useState("");
    const [eventTime, setEventTime] = useState("");

    const handleEventTitle = (e) => {
        setEventTitle(e.target.value);
    };
    const handleEventDescription = (e) => {
        setEventDescription(e.target.value);
    };
    const handleEventDate = (e) => {
        setEventDate(e.target.value);
    };
    const handleEventTime = (e) => {
        setEventTime(e.target.value);
    };
    console.log("groupId",groupId);
    const handleCreateEvent = (e) => {
        e.preventDefault();
        const method = "POST";
        const type = "createEvent";
         
        const payload = {
            sessionId: getCookie("sessionId"),
            event:{
            creatorId: getUserId("userId"),
            groupId: parseInt(groupId),
            title: eventTitle,
            content: eventDescription,
            date: eventDate,
            },

        };
        console.log("test case",type,"payload",payload);
      fetchData(method, type, payload).then((data) => {
            console.log("data", data);
            //window.location.reload();
        }); 
    };

  return (
    <div className="container">
      <div className="row">
        <div className="header">
          <h1>Create Event</h1>
          <form>
            <div className="form-group">
              <label htmlFor="eventTitle">Event Title</label>
              <input
                type="text"
                className="form-control"
                id="eventTitle"
                placeholder="Enter event title"
                onChange={handleEventTitle}
              />
            </div>
            <div className="form-group">
              <label htmlFor="eventDescription">Event Description</label>
              <textarea
                className="form-control"
                id="eventDescription"
                rows="3"
                onChange={handleEventDescription}
              ></textarea>
            </div>
            <div className="form-group">
              <label htmlFor="eventDate">Event Date</label>
              <input type="date" className="form-control" id="eventDate" onChange={handleEventDate} />
            </div>
            <div className="form-group">
              <label htmlFor="eventTime">Event Time</label>
              <input type="time" className="form-control" id="eventTime" onChange={handleEventTime}/>
            </div>
          </form>
          <button type="submit" className="btn btn-primary" onClick={handleCreateEvent}>
            Create Event
          </button>
        </div>
      </div>
    </div>
  );
};

export default CreateEvent;