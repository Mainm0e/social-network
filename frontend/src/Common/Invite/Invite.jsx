import { getCookie, getUserId } from "../../tools/cookie";
import { fetchData } from "../../tools/fetchData";

const Invite = ({ clearBox }) => {

    const test = () => {
        const url = new URL(window.location.href);
        const searchParams = new URLSearchParams(url.search);
        const id = searchParams.get("id");
        const method = "POST";
        const type = "getNonMembers";
        /* 
        type NonMembers struct {
    SessionId  string         `json:"sessionId"`
    UserId     int            `json:"userId"`
    GroupId    int            `json:"groupId"`
    NonMembers []SmallProfile `json:"nonMembers"`
}
event type : "getNonMembers" */
const payload = {
    sessionId: getCookie("sessionId"),
    userId: getUserId("userId"),
    groupId: parseInt(id),
};
console.log("test case", payload)

        fetchData (method, type, payload).then((data) => {
            console.log(data);
        });
    };

    return (
        <div className="invite">
            <h1>hello</h1>
            <button onClick={test}>test</button>
            <button onClick={clearBox}>close</button>
        </div>
    );
};

export default Invite;