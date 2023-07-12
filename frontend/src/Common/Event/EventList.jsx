import { useEffect } from "react"
import { fetchData } from "../../tools/fetchData"
import { getCookie } from "../../tools/cookie"

const EventList = ({ clearBox }) => {
    const url = new URL(window.location.href)
    const searchParams = new URLSearchParams(url.search)
    const id = searchParams.get("id")


    useEffect(() => {
        const method = "POST"
        const type = "getGroupEvents"
        const payload = { sessionId: getCookie("sessionId"), groupId: parseInt(id) }
        fetchData(method, type, payload).then((res) => {
            console.log("EventList",payload)
            console.log(res)
        }
        )
    }, [id])

    return (
        <>
        </>
    )
}

export default EventList