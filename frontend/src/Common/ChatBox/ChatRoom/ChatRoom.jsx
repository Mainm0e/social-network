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

  const getChatContent = () => {
    const payload = {
      sessionID: getCookie("sessionId"),
      chatType: "private",
      clientID: getUserId("userId"),
      targetID: receiver.id,
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
  }, []);

  const chatContent = chatHistory.map((message) => {
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
        key={message.id}
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
        receiver: parseInt(receiver.id),
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
      const message = {
        sessionID: getCookie("sessionId"),
        senderID: getUserId("userId"),
        receiverID: parseInt(receiver.id),
        message: messageInput,
        timeStamp: "",
      };
      const privateMessageEvent = {
        type: "privateMsg",
        payload: message,
      };

      socket.send(JSON.stringify(privateMessageEvent)); // Send the message as a string
      setMessageInput(""); // Clear the input field
    }
  };


  const handleUserClick = () => {
    setIsClosed(true);
    onClose(true); // Pass the boolean value back to the parent component
  };

  if (isClosed) {
    return null; // Return null if the ChatRoom is closed
  }
  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === "chatHistory") {
      setChatHistory(message.payload.chatHistory);
    }
    if (message.type === "PrivateMsg") {
     const newMessage = addNewMessage("receiver",message.payload);
     


      document.getElementsByClassName("chat-messages")[0].scrollTop =
      document.getElementsByClassName("chat-messages")[0].scrollHeight;
    } else if (message.type !== "chatHistory") {
      console.log("message", message);
    }
  };


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
        <div className="chat-messages">
          {chatContent}

          </div>
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

const addNewMessage = (who,message) => {
    return (
      <div className={who+"-message"}>
        <div className="chat-message">{message.messageContent}</div>
      </div>
    );
}