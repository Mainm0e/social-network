import React, { useState, useContext, useEffect } from "react";
import "./ChatRoom.css";
import { dummyMessages } from "../DummyData";
import { WebSocketContext } from "../../../WebSocketContext/websocketcontext";
import { getCookie, getUserId } from "../../../tools/cookie";

const ChatRoom = (props) => {
  const sender = getUserId("userId");
  const { receiver, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);
  const socket = useContext(WebSocketContext);
  const [messageInput, setMessageInput] = useState("");
  const [chatHistory, setChatHistory] = useState([]);
  const [newChatContent , setNewChatContent] = useState(null); // store new chat content from server


 // todo: add functionality when user is typing
 // todo: add functionality when user is sending message successfully sroll to bottom
 // todo: notication when user is offline

  const getChatContent = () => {
    const payload = {
      sessionID: getCookie("sessionId"),
      chatType: "private",
      clientID: getUserId("userId"),
      targetID: receiver.userId,
    };
    const chatHistoryRequest = {
      type: "chatHistoryRequest",
      payload: payload,
    };
    socket.send(JSON.stringify(chatHistoryRequest)); // Send the message as a string
  };

  // get chat content from server
  useEffect(() => {
    getChatContent();
  }, [receiver]);

  const chatContent = chatHistory.map((message,index) => {
    const isSender = message.senderId === sender;
    const isReceiver = message.receiverId === sender;
    if (!isSender && !isReceiver) {
      return null;
    }
    return (
      <div
        className={`${
          isSender ? "sender" : isReceiver ? "receiver" : ""
        }-message`}
        key={index}
      >
        <div className="chat-message">{message.messageContent}</div>
      </div>
    );
  });

  // send typing event to server when user is typing
  const typingMessage = (e) => {
    if (receiver.status === "online") {
      const message = {
        sender: getUserId("userId"),
        receiver: parseInt(receiver.userId),
        message: "typing",
      };
      /* socket.send(JSON.stringify(message)); // Send the message as a string */
      /*   console.log("Message sent: ", message); */
    }
    setMessageInput(e);
  };

  // send message to server when user press enter
  const sendMessage = () => {
    if (messageInput.trim() !== "") {
      // ! struct message that need to send to server is different from message that need to display on client
      // ! that why we have message and newMessage
      const newMessage = {
        senderId: getUserId("userId"),
        receiverId: parseInt(receiver.userId),
        messageContent: messageInput,
      };
      const message = {
        sessionID: getCookie("sessionId"),
        senderID: getUserId("userId"),
        receiverID: parseInt(receiver.userId),
        message: messageInput,
        timeStamp: "",
      };
      setChatHistory((prevChatHistory) => [...prevChatHistory, newMessage]);
      const privateMessageEvent = {
        type: "privateMsg",
        payload: message,
      };
      socket.send(JSON.stringify(privateMessageEvent)); // Send the message as a string
      setMessageInput(""); // Clear the input field


      // scroll to bottom when user send message
      setTimeout(() => {
        const chatMessages = document.getElementById("chat-container");
        chatMessages.scrollTop = chatMessages.scrollHeight;
      }, 100);
    
    }
  };
  

 
  const handleUserClick = () => {
    setIsClosed(true);
    onClose(true); // Pass the boolean value back to the parent component
  };

  // onmessage event listener
  // for catching event from server

  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === "chatHistory") {
      setChatHistory(message.payload.chatHistory);
    } else if (message.type === "PrivateMsg") {
      // ! SAME HERE
      // ! struct message that i got from server is different from getChatHistory 
      console.log("onmessage", message)
      const newMessage = {
        senderId: message.payload.senderID,
        receiverId: message.payload.receiverID,
        messageContent: message.payload.message,
        sendTime: message.payload.timeStamp,
      };

      setChatHistory((prevChatHistory) => [...prevChatHistory, newMessage]);
    } else {
      console.log("message", message);
    }
  };
  
  if (isClosed) {
    return null; // Return null if the ChatRoom is closed
  }
  return (
    <div className="chat-room">
      <div className="top-bar">
        <div className="chat-room-avatar">
          <img src={receiver.avatar} alt={receiver.firstName} />
        </div>
        <div className="chat-room-name">{receiver.firstName}</div>
        <span className="close-button" onClick={handleUserClick}>
          close
        </span>
      </div>
      <div className="chat-room-content">
      <div id="chat-container" className="chat-messages">{chatContent}</div>
      </div>
      <div className="chat-room-input">
        <input
          type="text"
          placeholder="Type a message..."
          value={messageInput}
          onChange={(e) => typingMessage(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === "Enter") {
              sendMessage();
            }
          }}
        />
      </div>
    </div>
  );
};

export default ChatRoom;
