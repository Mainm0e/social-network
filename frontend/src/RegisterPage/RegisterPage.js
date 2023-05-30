import React, { useState } from 'react';
import Info1 from './Info1/info1';
import Info2 from './Info2/info2';
import './RegisterPage.css';
import WelcomeBox from '../Common/WelcomeBox/WelcomeBox';
import AlertBox from '../Common/AlertBox/AlertBox';


function RegisterPage() {
  const [selectedOption, setSelectedOption] = useState('info1');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [birthdate, setBirthdate] = useState('');
  const [avatar, setAvatar] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
   // [username, password, confirmPassword, firstName, lastName, email, birthdate, avatar
  const handleOptionChange = (event) => {
    const selectedOption = event.target.value;
    setSelectedOption(selectedOption); 
  };
  
  const handleInfoChange = (values) => {
    if (values.type === "info1"){
      setFirstName(values.firstName);
      setLastName(values.lastName);
      setEmail(values.email);
      setBirthdate(values.birthdate);
    }
    else if (values.type === "info2"){
      setUsername(values.username);
      setPassword(values.password);
      setConfirmPassword(values.confirmPassword);
      setAvatar(values.avatar);
    }
  };
  const [alertTitle, setAlertTitle] = useState('');
  const [alertMessage, setAlertMessage] = useState([]);
  const [registerStatus, setRegisterStatus] = useState(true);

  const handleRegister = () => {
  let msg = [];
    // Perform registration logic here
    let matchPassword = password === confirmPassword ? password : "";
    const data = {firstName, lastName, email, birthdate, avatar,username, matchPassword};
  for (let key in data) {
  if (data.hasOwnProperty(key)) {
    const value = data[key];
    if (value == ""&& key !== "avatar"&& key !== "matchPassword") {
      msg.push(key + " is empty");
      setRegisterStatus(false);
    }
    if (key === "matchPassword" && value === "") {
      msg.push("Password does not match");
      console.log(msg);
      setRegisterStatus(false);
    }
  }
  }


/*    if (username !== "" && password !== "" && firstName !== "" && lastName !== "" && email !== "" && birthdate !== "" && avatar !== ""){
    fetch('http://localhost:5000/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
      })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
        if (data.success){
          alert('Registration successful');
          window.location.href = '/login';
        }
        else {
          alert('Registration failed');
        }
      })
      .catch((error) => {
        console.error('Error:', error);
      });
    } */
    if (msg.length === 0){
      setRegisterStatus(true);
    } else if (msg.length > 0){
      setRegisterStatus(false);
    }

    if (registerStatus){
      document.querySelector(".alert-box").style.display = "none";
    } else {
      setAlertTitle('Error');
      setAlertMessage(msg);
      document.querySelector(".alert-box").style.display = "block";
    }
  };

  return (
    <div className='register-page'>
    <WelcomeBox />
    <div className="register-container">
    <AlertBox title={alertTitle} message={alertMessage} status={true} />
      <h1>Register Page</h1>
      <Info1 selectedOption={selectedOption} onChange={handleInfoChange} registerStatus={registerStatus}/>
      <Info2 selectedOption={selectedOption} onChange={handleInfoChange} registerStatus={registerStatus}/>
        <div className='select-container'>
        <input
          type="radio"
          value="info1"
          checked={selectedOption === 'info1'}
          onChange={handleOptionChange}
        />
        <input
          type="radio"
          value="info2"
          checked={selectedOption === 'info2'}
          onChange={handleOptionChange}
        />
      </div>
        <button type="button" onClick={handleRegister}>
          Register
        </button>
      
      <div className='links'>
        <a href='/login'>Login</a>
      </div>
    </div>
    </div>
  );
}

export default RegisterPage; 