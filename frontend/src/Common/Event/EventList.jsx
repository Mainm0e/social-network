import { useEffect } from "react"
import { fetchData } from "../../tools/fetchData"
import { getCookie, getUserId } from "../../tools/cookie"
import { useState } from "react"

const EventList = ({ clearBox }) => {
    const url = new URL(window.location.href)
    const searchParams = new URLSearchParams(url.search)
    const id = searchParams.get("id")
    const [event, setEvent] = useState([])


    useEffect(() => {
        const method = "POST"
        const type = "getGroupEvents"
        const payload = { sessionId: getCookie("sessionId"), groupId: parseInt(id) }
        fetchData(method, type, payload).then((data) => {
            console.log("EventList",payload)
            setEvent(data)
        }
        )
    }, [id])
        const acceptEvent = (event) => {
            const method = "POST"
            const type = "participate"
            const payload = { sessionId: getCookie("sessionId"), eventId: event.id,memberId:getUserId("userId"),option:event }
            fetchData(method, type, payload).then((data) => {
                console.log(data)
            })
        }

        const renderNotifications = () => {
            return (
                <>
                {event.events.map((event, index) => (
                    <div className="event-list-item" key={index}>
                      <div className="event-list-item-header">
                        <h2>{event.event.title}</h2>
                      </div>
                      <div className="event-list-item-body">
                        <p>{event.event.content}</p>
                        {/* creater */}
                        <div className="event-creater">
                            <p>Created by: {event.creatorProfile.firstName}</p>
                        </div>
                      </div>
                      <div className="event-list-item-footer">
                        <p>{event.event.date}</p>
                      </div>
                        <div className="event-list-item-footer">
                            <button onClick={() => acceptEvent("going")}>Going</button>
                            <button onClick={() => acceptEvent("not_going")}>Not going</button>
                        </div>
                    </div>
                ))}
                </>
            )
        };
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
    )
}

export default EventList
