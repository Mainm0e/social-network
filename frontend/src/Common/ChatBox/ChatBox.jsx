import React, { useState , useEffect} from 'react';
import './ChatBox.css';
import ChatList from './ChatList/ChatList';
import ChatRoom from './ChatRoom/ChatRoom';
const ChatBox = () => {

  //start socket connection
  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080.ws');

    // Connection opened
    socket.addEventListener('open', (event) => {
      console.log('WebSocket connected');
    });

    // Listen for messages
    socket.addEventListener('message', (event) => {
      console.log('Received message:', event.data);
    });

    // Connection closed
    socket.addEventListener('close', (event) => {
      console.log('WebSocket disconnected');
    });

    // Clean up the WebSocket connection on component unmount
    return () => {
      socket.close();
    };
  }, []); // Empty dependency array to run the effect only once

  

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
