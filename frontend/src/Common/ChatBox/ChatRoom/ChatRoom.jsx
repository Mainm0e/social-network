import React, { useState, useContext } from "react";
import "./ChatRoom.css";
import { dummyMessages } from "../DummyData";
import { WebSocketContext } from "../../../WebSocketContext/websocketcontext";
import { getUserId } from "../../../tools/cookie";

const ChatRoom = (props) => {
  const sender = getUserId("userId");
  const { receiver, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);
  const socket = useContext(WebSocketContext);
  const [messageInput, setMessageInput] = useState("");


  // get chat content from server
  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log("Message received: ", message);
    if (message.message === "typing") {
      console.log("typing");
    } else {
      console.log("message");
    }
  };

  const getChatContent = () => {
    const message = {
      sender: sender,
      receiver: parseInt(receiver.id),
      message: "chatHistoryRequest",
    };
    socket.send(JSON.stringify(message)); // Send the message as a string
    console.log("Message sent: ", message);
  }

  // send typing event to server when user is typing
  const typingMessage = (e) => {
    if (receiver.status === "online") {
    const message = {
      sender: getUserId("userId"),
      receiver: parseInt(receiver.id),
      message: "typing",
    };
    /* socket.send(JSON.stringify(message)); // Send the message as a string */
    console.log("Message sent: ", message);
    }
    setMessageInput(e)
  };

  // send message to server when user press enter
  const sendMessage = () => {
    if (messageInput.trim() !== "") {
      const message = {
        sender: getUserId("userId"),
        receiver: receiver.id,
        message: messageInput,
      };
      socket.send(JSON.stringify(message)); // Send the message as a string
      setMessageInput(""); // Clear the input field
    }
  };
  const chatContent = dummyMessages.map((message) => {
    const isSender = message.sender === sender;
    const isReceiver = message.sender === receiver.id;
    if (!isSender && !isReceiver) {
      return null;
    }
    return (
      <div
        className={`${
          isSender ? "sender" : isReceiver ? "receiver" : ""
        }-message`}
        key={message.id}
      >
        <div className="chat-message">{message.message}</div>
      </div>
    );
  });

  const handleUserClick = () => {
    setIsClosed(true);
    onClose(true); // Pass the boolean value back to the parent component
  };

  if (isClosed) {
    return null; // Return null if the ChatRoom is closed
  }

  return (
    <div className="chat-room">
      <div className="top-bar">
        <div className="chat-room-avatar">
          <img src={receiver.avatar} alt={receiver.name} />
        </div>
        <div className="chat-room-name">{receiver.name}</div>
        <span className="close-button" onClick={handleUserClick}>
          close
        </span>
      </div>
      <div className="chat-room-content">
        <div className="chat-messages">{chatContent}</div>
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
