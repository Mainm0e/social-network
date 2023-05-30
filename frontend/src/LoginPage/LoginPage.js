import React, { useState } from 'react';
import './LoginPage.css';
import WelcomeBox from '../Common/WelcomeBox/WelcomeBox';
import AlertBox from '../Common/AlertBox/AlertBox';
import ChatBox from '../Common/ChatBox/ChatBox';
function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [alertTitle, setAlertTitle] = useState('');
  const [alertMessage, setAlertMessage] = useState([]);
  // is part is for the color changing text animation
  /*
  // import useEffect from react befor using this
  // and add "style={{ color: textColor }}" into the h1 tag 
   const [textColor, setTextColor] = useState('red');

   useEffect(() => {
    const colors = ['red', 'blue', 'green', 'yellow']; // Add more colors if needed
    const loginText = 'Login';

    const textColors = loginText.split('').map((char, index) => ({
      char,
      color: colors[index % colors.length]
    }));

    setTextColors(textColors);
  }, []); */

  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };
  let msg = [];
  const checkUsername = (username, password) => {
    if (username === "admin" && password === "admin") {
      return true;
    } else {
      setAlertTitle("Error");
      msg.push("Username or password is incorrect");
      setAlertMessage(msg);
      return false;
    }
  }

  const [loginStatus, setLoginStatus] = useState(true);
  const handleLogin = () => {
    // Perform login logic here
    if (checkUsername(username, password)) {
      setLoginStatus(true);
      document.querySelector(".alert-box").style.display = "none";
    } else {
      setLoginStatus(false);
      document.querySelector(".alert-box").style.display = "block";
    }
  };
  return (
    <div className='login-page'>
    <WelcomeBox />
    <div className="login-container">
    <AlertBox title={alertTitle} message={alertMessage} status={true} />
      <h1 >Login Page</h1>
      <form>
        <div>
          <label>Username:</label>
          <input
            type="text"
            value={username}
            style={{ background: loginStatus === false ? "#FFEA00" : "" }}
            onChange={handleUsernameChange}
            required
          />
        </div>
        <div>
          <label>Password:</label>
          <input
            type="password"
            value={password}
            style={{ background: loginStatus === false ? "#FFEA00" : "" }}
            onChange={handlePasswordChange}
            required
          />
        </div>
        <button type="button" onClick={handleLogin}>
          Login
        </button>
      </form>
      <div className='links'>
        <a href='/register'>Register</a>
      </div>
    </div>
    <ChatBox />
    </div>
  );
  
}

export default LoginPage;
