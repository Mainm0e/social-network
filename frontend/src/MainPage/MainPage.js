import "./MainPage.css";
import MainBox from "./MainBox/MainBox";
import LeftBox from "./LeftBox/LeftBox";
import RightBox from "./RightBox/RightBox";
import ChatBox from '../Common/ChatBox/ChatBox';

// dummy data
import { user1 } from "./dummyData";


function MainPage(name) {
 
   fetch('http://localhost:8080/api', {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json',     },
     body: JSON.stringify({event_type:'profile', payload: {user_id: 17}})
   }).then((response) => {
  console.log(response);
     return response.json();
     
   }).then((data) => {
     console.log(data);
   })
}
export default MainPage;