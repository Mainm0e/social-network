import React, { useState } from 'react';
import './LoginPage.css';
import WelcomeBox from '../Common/WelcomeBox/WelcomeBox';
import AlertBox from '../Common/AlertBox/AlertBox';
import { fetchData } from '../tools/fetchData';

// LoginPage component
// This component is used to render the login page
// Props: none
function LoginPage() {  
  const [email, setEmail] = useState('');
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

  const handleEmailChange = (event) => {
    setEmail(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };
  let msg = [];
  const checkemail = (email, password) => {
    console.log("checkEmail",email, password);
    const method = "POST"
    const type = "login"
    const payload = {email: email, password: password}
    fetchData(method,type,payload).then((data) => {
      console.log("data",data);
      if (data.statusCode === 200) {
        window.location.href = "/";
      } else {
        setAlertTitle("Login Failed");
        msg.push(data.message);
        setAlertMessage(msg);
      }
    }
    );
  }

  const [loginStatus, setLoginStatus] = useState(true);

  // handleLogin function
  //  this function is main function of the login page
  const handleLogin = () => {
    // Perform login logic here
    if (checkemail(email, password)) {
      console.log("do login",email, password)
      setLoginStatus(true);
      document.querySelector(".alert-box").style.display = "none";
    } else {
      setLoginStatus(false);
      document.querySelector(".alert-box").style.display = "block";
    }
  };
  return (
    <div className="main-container">
    <div className='login-page'>
    <WelcomeBox />
    <div className="login-container">
    <AlertBox title={alertTitle} message={alertMessage} status={true} />
      <h1 >Login Page</h1>
      <form>
        <div>
          <label>email:</label>
          <input
            type="text"
            value={email}
            style={{ background: loginStatus === false ? "#FFEA00" : "" }}
            onChange={handleEmailChange}
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
    </div>
  );
  
}

export default LoginPage;
