import React, { useState } from 'react';
import './ChatBox.css';
import ChatList from './ChatList/ChatList';

const ChatBox = () => {
    const [chat_list, setChatlist] = useState(false);
    const openChatList = () => {
       setChatlist(true);
    }
    const closeChatList = () => {
        setChatlist(false);
    }
    // chat button component
    const chatButton = () => {
        return (
            <div className='chat-button' onClick={openChatList}>
            <span>Chat</span>
          </div>
        );
    }

  return (
    <div className="chat-container">
      {chat_list ? <><div className='top-bar'><span className='close-button' onClick={closeChatList}>close</span></div><ChatList/></> : chatButton() }
      
    </div>
  );
};

export default ChatBox;
