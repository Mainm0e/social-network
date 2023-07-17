import React, { useState, useEffect } from "react";
import LoginPage from "./LoginPage/LoginPage";
import RegisterPage from "./RegisterPage/RegisterPage";
import MainPage from "./MainPage/MainPage";
// import ErrorPage from "./ErrorPage/ErrorPage";
import { WebSocketProvider } from "./WebSocketContext/websocketcontext"; // import WebSocketProvider
import "./App.css";
import { getCookie } from "./tools/cookie";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const sessionId = getCookie("sessionId");
    setIsLoggedIn(sessionId !== null);
  }, []);

  const getPage = () => {
    const page = window.location.pathname;
    if (!isLoggedIn){
      if (page === "/"||page ==="/login"){
        return <LoginPage />
      } else if (page === "/register"){
        return <RegisterPage />
      }
    }
    if ((isLoggedIn && page === "/") || (isLoggedIn && page === "/register") || (isLoggedIn && page === "/login")){
      return <MainPage />
    } else {
      return  <MainPage />
    }
  };

  return (
    <div className="App">
      <WebSocketProvider isLoggedIn={isLoggedIn}> {/* Pass isLoggedIn prop to WebSocketProvider */}
        {getPage()}
      </WebSocketProvider>
    </div>
  );
};

export default App;
