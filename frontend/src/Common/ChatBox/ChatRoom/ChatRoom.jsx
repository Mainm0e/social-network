import React, { useState, useContext, useEffect } from "react";
import "./ChatRoom.css";
import { WebSocketContext } from "../../../WebSocketContext/websocketcontext";
import { getCookie, getUserId } from "../../../tools/cookie";

const ChatRoom = (props) => {
  const sender = getUserId("userId");
  const { receiver, type, id, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);
  const socket = useContext(WebSocketContext);
  const [messageInput, setMessageInput] = useState("");
  const [chatHistory, setChatHistory] = useState([]);
  const [isTyping, setIsTyping] = useState(false); // store isTyping event from server
  const [currentReceiver, setCurrentReceiver] = useState(receiver);
  const [startChat, setStartChat] = useState(false);


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

  // start chat function for make sure chat history is up to date
  // when user change chat
  // chatbox will scroll to bottom when start chat
  const chatState = () => {
    if (currentReceiver !== receiver) {
      console.log("change receiver")
      setCurrentReceiver(receiver);
      setStartChat(false);
      getChatContent();
      setTimeout(() => {
        const chatMessages = document.getElementById("chat-container");
        chatMessages.scrollTop = chatMessages.scrollHeight;
      }, 100);
      setStartChat(true);
    }
    if (startChat === false) {
      getChatContent();
      setTimeout(() => {
        const chatMessages = document.getElementById("chat-container");
        chatMessages.scrollTop = chatMessages.scrollHeight;
      }, 100);
      setStartChat(true);
    }
  };
  // for chat history

  // for sending message
  const [chatType, setChatType] = useState("privateMsg");

  // todo: add functionality when user is typing

  // !! how to get chat history from server is confusing right now  in GroupChat

  // get chat content from server
  useEffect(() => {
    const updateChatSettings = async () => {
      if (type === "group") {
        // for sending message
        setChatType("groupMsg");
      }
      if (type === "private") {
        // for sending message
        setChatType("privateMsg");
      }
    };

    const getChatContentAsync = async () => {
      await updateChatSettings();
    };
    chatState();
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


  const [test, setTest] = useState([]);
  const addNewMessage = (message) => {
    const isSender = message.senderId === sender;
    const isReceiver = message.receiverId === sender;
    let index = test.length;
    let newElement = <></>
   
    if (type === "private" && message.msgType === "PrivateMsg" && (isSender || isReceiver)) {
      newElement = (
        <div
          className={`${isSender ? "sender" : "receiver"}-message`}
          key={index}
        >
          <div className="chat-message">{message.messageContent}</div>
        </div>
      );
    } else if (type === "group" && message.msgType === "GroupMsg" && (isSender || message.senderId !== sender)) {
      newElement = (
        <div
          className={`${isSender ? "sender" : "receiver"}-message`}
          key={index}
        >
          <div className="chat-message">
            {message.messageContent}
            {!isSender && (
              <div className="chat-message-sender">
                {getSenderName(message.senderId)}
              </div>
            )}
          </div>
        </div>
      );
    } else {
      newElement = null;
    }

    setTest((prevTest) => [...prevTest, newElement]);
  };

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
      console.log(message);
      const newMessage = {
        senderId: message.payload.senderID,
        receiverId: message.payload.receiverID,
        messageContent: message.payload.message,
        sendTime: message.payload.timeStamp,
        msgType: message.type,
      };

      setChatHistory((prevChatHistory) => [...prevChatHistory, newMessage]);
      if (newMessage.senderId === receiver.userId)
       addNewMessage(newMessage);
    } else if (message.type === "GroupMsg") {
      const newMessage = {
        senderId: message.payload.senderID,
        receiverId: message.payload.receiverID,
        messageContent: message.payload.message,
        sendTime: message.payload.timeStamp,
        msgType: message.type,
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
            <div className="chat-room-name">{receiver.firstName}</div>
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
          <ChatContent
            chatHistory={chatHistory}
            type={type}
            sender={sender}
            receiver={receiver}
          />
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

const ChatContent = ({ chatHistory, type, sender, receiver }) => {
  const getSenderName = (senderId) => {
    if (type === "group") {
      for (let i = 0; i < receiver.members.length; i++) {
        if (receiver.members[i].userId === senderId) {
          return receiver.members[i].firstName;
        }
      }
    }
  };

  return (
    <>
      {chatHistory.map((message, index) => {
        const isSender = message.senderId === sender;
        const isReceiver = message.receiverId === sender;

        if (type === "private" && message.msgType === "PrivateMsg" && (isSender || isReceiver)) {
          return (
            <div
              className={`${isSender ? "sender" : "receiver"}-message`}
              key={index}
            >
              <div className="chat-message">{message.messageContent}</div>
            </div>
          );
        } else if (type === "group" && message.msgType === "GroupMsg" && (isSender || message.senderId !== sender)) {
          return (
            <div
              className={`${isSender ? "sender" : "receiver"}-message`}
              key={index}
            >
              <div className="chat-message">
                {message.messageContent}
                {!isSender && (
                  <div className="chat-message-sender">
                    {getSenderName(message.senderId)}
                  </div>
                )}
              </div>
            </div>
          );
        } else {
          return null;
        }
      })}
    </>
  );
};




