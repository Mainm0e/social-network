import React, { useState } from 'react';
import './ChatRoom.css';
import { dummyMessages } from '../DummyData';
const ChatRoom = (props) => {
  const { receiver, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);
  const dummySender = 1;
  const dummyReceiver = 2;

  const chatContent = dummyMessages.map((message) => {
    const isSender = message.sender === dummySender;
    return (
      <div
        className={`${isSender ? 'sender' : 'receiver'}-message`}
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
           {dummyMessages ?  chatContent : <div className="no-message">No messages</div>}
        </div>
        <div className="chat-room-input">
            <input type="text" placeholder="Type a message..." />
        </div>
    </div>
  );
};

export default ChatRoom;
