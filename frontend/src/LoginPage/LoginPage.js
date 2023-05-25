import React, { useState } from 'react';
import './LoginPage.css';
import WelcomeBox from '../Common/WelcomeBox/WelcomeBox';
import AlertBox from '../Common/AlertBox/AlertBox';
function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
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

  const [loginStatus, setLoginStatus] = useState(true);

    const handleLogin = () => {
      // Perform login logic here
      if (checkUsername(username, password)) {
        setLoginStatus(true);
      } else {
        setLoginStatus(false);
      }
  };
  console.log("loginStatus",loginStatus)
  return (
    <div className='login-page'>
    <WelcomeBox />
    <div className="login-container">
      {loginStatus ? null : <AlertBox title="HI" message="im fine" status={loginStatus} />}
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
    </div>
  );
  
}

function checkUsername(username, password) {
  // Perform username validation here
  return false;
}

export default LoginPage;
