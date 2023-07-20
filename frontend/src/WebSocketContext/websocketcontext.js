import React, { createContext, useEffect, useState } from 'react';
import { logout } from '../tools/logout';

// Create a context that will hold the websocket connection
export const WebSocketContext = createContext(null);

// Create a provider component that will wrap other components and provide them with the websocket connection
export const WebSocketProvider = ({ children, isLoggedIn }) => {
  const [socket, setSocket] = useState(null);
  useEffect(() => {
    if (!isLoggedIn) {
      return;
    }
    // On component mount, try to connect to the websocket server
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
      console.log('connected to ws server');
      setSocket(ws);
    };

    ws.onerror = (error) => {
      console.log('failed to connect to ws server:', error);
      logout();
    };

    ws.onclose = (event) => {
      console.log('WebSocket is closed now. With code:', event.code, ' Reason:', event.reason);
    };

    // On component unmount, close the websocket connection
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [isLoggedIn]);

  return (
    // Provide the websocket connection to child components
    <WebSocketContext.Provider value={socket}>
      {children}
    </WebSocketContext.Provider>
  );
};