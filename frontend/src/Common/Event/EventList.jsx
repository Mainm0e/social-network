import { useEffect } from "react"
import { fetchData } from "../../tools/fetchData"
import { getCookie } from "../../tools/cookie"
import { useState } from "react"

const EventList = ({ clearBox }) => {
    const url = new URL(window.location.href)
    const searchParams = new URLSearchParams(url.search)
    const id = searchParams.get("id")
    const [data, setData] = useState()


    useEffect(() => {
        const method = "POST"
        const type = "getGroupEvents"
        const payload = { sessionId: getCookie("sessionId"), groupId: parseInt(id) }
        fetchData(method, type, payload).then((data) => {
            console.log("EventList",payload)
            setData(data)
        }
        )
    }, [id])
    /* if (data.length === 0) {
        return (
            <div className="event-list">
                <div className="event-list-header">
                    <h1>Events</h1>
                </div>
                <div className="event-list-body">
                    <div className="event-list-item">
                        <div className="event-list-item-header">
                            <h2>No events</h2>
                        </div>
                        <div className="event-list-item-body">
                            <p>There are no events in this group</p>
                        </div>
                    </div>
                </div>
                <div className="event-list-footer">
                    <button onClick={clearBox}>Close</button>
                </div>
            </div>
        )
    } else {
        console.log(data) */
    return (
        <div className="event-list">
        <div className="event-list-header">
          <h1>Events</h1>
        </div>
        <div className="event-list-body">
            {console.log(data,"data")}
          {data.map((event, index) => (
            <div className="event-list-item" key={index}>
              <div className="event-list-item-header">
                <h2>{event.title}</h2>
              </div>
              <div className="event-list-item-body">
                <p>{event.content}</p>
              </div>
              <div className="event-list-item-footer">
                <p>{event.date}</p>
              </div>
            </div>
          ))}
        </div>
        <div className="event-list-footer">
          <button onClick={clearBox}>Close</button>
        </div>
      </div>
    )
}

export default EventList