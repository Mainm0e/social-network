import React, { useState } from 'react';
import './ChatBox.css';
import ChatList from './ChatList/ChatList';
import ChatRoom from './ChatRoom/ChatRoom';
const ChatBox = () => {
    const [chat_list, setChatlist] = useState(false);
    const [room, setRoom] = useState(null);
  
    const openChatList = () => {
      setChatlist(true);
    };
  
    const closeChatList = () => {
      setChatlist(false);
    };
  
    const chatButton = () => {
      return (
        <div className='chat-button' onClick={openChatList}>
          <span>Chat</span>
        </div>
      );
    };
  
    const handleUserSelection = (selectedUser) => {
      setRoom(selectedUser);
    };
  
    const handleCloseChatRoom = (isClosed) => {
      if (isClosed) {
        setRoom(null);
      }
    };
  
    return (
      <>
        <div className="chat-container">
          {room && (
            <ChatRoom receiver={room} onClose={handleCloseChatRoom} />
          )}
          {chat_list ? (
            <>
              <div className='contact-list'>
                <div className='top-bar'>
                  <span className='close-button' onClick={closeChatList}>
                    close
                  </span>
                </div>
                <ChatList onUserSelection={handleUserSelection} />
              </div>
            </>
          ) : (
            chatButton()
          )}
        </div>
      </>
    );
  };
  
  
export default ChatBox;
