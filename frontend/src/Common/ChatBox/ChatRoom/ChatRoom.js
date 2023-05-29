import React, { useState } from 'react';
import './ChatRoom.css';
const ChatRoom = (props) => {
  const { receiver, onClose } = props;
  const [isClosed, setIsClosed] = useState(false);

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
        </div>
        <div className="chat-room-input">
            <input type="text" placeholder="Type a message..." />
        </div>
    </div>
  );
};

export default ChatRoom;
