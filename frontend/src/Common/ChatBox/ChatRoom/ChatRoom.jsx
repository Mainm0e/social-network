import React, { useState, useContext, useEffect } from "react";
import "./ChatRoom.css";
import { dummyMessages } from "../DummyData";
import { WebSocketContext } from "../../../WebSocketContext/websocketcontext";
import { getCookie, getUserId } from "../../../tools/cookie";

const ChatRoom = (props) => {
  const sender = getUserId("userId");
  const { receiver, type, id, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);
  const socket = useContext(WebSocketContext);
  const [messageInput, setMessageInput] = useState("");
  const [chatHistory, setChatHistory] = useState([]);
  const [newChatContent, setNewChatContent] = useState(null); // store new chat content from server
  const [isTyping, setIsTyping] = useState(false); // store isTyping event from server
  // for chat history

  // for sending message
  const [chatType, setChatType] = useState("privateMsg");

  // todo: add functionality when user is typing

  // !! how to get chat history from server is confusing right now  in GroupChat

  const getChatContent = () => {
    const payload = {
      sessionID: getCookie("sessionId"),
      chatType: type,
      clientID: getUserId("userId"),
      targetID: receiver.userId,
    };
    if (type === "group") {
      payload.targetID = id;
    }
    const chatHistoryRequest = {
      type: "chatHistoryRequest",
      payload: payload,
    };
    socket.send(JSON.stringify(chatHistoryRequest)); // Send the message as a string
  };

  // get chat content from server
  useEffect(() => {
    const updateChatSettings = async () => {
      if (type === "group") {
        // for sending message
        setChatType("groupMsg");
      }
    };

    const getChatContentAsync = async () => {
      await updateChatSettings();
      getChatContent();
    };

    getChatContentAsync();
  }, [receiver,chatHistory]);

  const getSenderName = (senderId) => {
    if (type === "group") {
      for (let i = 0; i < receiver.members.length; i++) {
        if (receiver.members[i].userId === senderId) {
          return receiver.members[i].firstName;
        }
      }
    }
  };
  const chatContent = chatHistory.map((message, index) => {
    if (type === "private" && message.msgType === "PrivateMsg") {
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
    } else if (type === "group" && message.msgType === "GroupMsg") {
      const isSender = message.senderId === sender;
      const otherMember = message.senderId !== sender;
      if (!isSender && !otherMember) {
        return null;
      }
      return (
        <div
          className={`${
            isSender ? "sender" : otherMember ? "receiver" : ""
          }-message`}
          key={index}
        >
          <div className="chat-message">
            {message.messageContent}
            {otherMember && (
              <div className="chat-message-sender">
                {getSenderName(message.senderId)}
              </div>
            )}
          </div>
        </div>
      );
    }
  });

  // send typing event to server when user is typing
  const typingMessage = (e) => {
    const payload = {
      sessionID: getCookie("sessionId"),
      chatType: type,
      clientID: getUserId("userId"),
      targetID: receiver.userId,
    };
    if (type === "group") {
      payload.targetID = id;
    }
    const chatHistoryRequest = {
      type: "isTyping",
      payload: payload,
    };
    socket.send(JSON.stringify(chatHistoryRequest)); // Send the message as a string
    setMessageInput(e);
  };

  // send message to server when user press enter
  const sendMessage = () => {
    if (messageInput.trim() !== "") {
      // ! struct message that need to send to server is different from message that need to display on client
      // ! that why we have message and newMessage
      const message = {
        sessionID: getCookie("sessionId"),
        senderID: getUserId("userId"),
        receiverID: receiver.userId,
        message: messageInput,
        timeStamp: "",
      };
      if (type === "group") {
        message.receiverID = id;
      }
      const messageEvent = {
        type: chatType,
        payload: message,
      };
      socket.send(JSON.stringify(messageEvent)); // Send the message as a string
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
      const newMessage = {
        senderId: message.payload.senderID,
        receiverId: message.payload.receiverID,
        messageContent: message.payload.message,
        sendTime: message.payload.timeStamp,
      };

      setChatHistory((prevChatHistory) => [...prevChatHistory, newMessage]);
    } else if (message.type === "GroupMsg") {
      // ! SAME HERE
      // ! struct message that i got from server is different from getChatHistory
      const newMessage = {
        senderId: message.payload.senderID,
        receiverId: message.payload.receiverID,
        messageContent: message.payload.message,
        sendTime: message.payload.timeStamp,
      };

      setChatHistory((prevChatHistory) => [...prevChatHistory, newMessage]);
    } else if (message.type === "isTyping") {
      setIsTyping(true);
      setTimeout(() => {
        setIsTyping(false);
      }, 1000);
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
          {type === "private" && (
            <img src={receiver.avatar} alt={receiver.firstName} />
          )}
          {type === "group" && <p>{receiver.title}</p>}
        </div>
        {type === "private" && (
          <>
            <div className="chat-room-name">{receiver.firstName}
              </div>
            {isTyping && (
              <div className="typing-indicator">
                <div className="dot-1"></div>
                <div className="dot-2"></div>
                <div className="dot-3"></div>
              </div>
            )}
          </>
        )}
        <span className="close-button" onClick={handleUserClick}>
          close
        </span>
      </div>
      <div className="chat-room-content">
        <div id="chat-container" className="chat-messages">
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
