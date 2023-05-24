import React, { useState } from 'react';
import './RegisterPage.css';
import WelcomeBox from '../Common/WelcomeBox/WelcomeBox';
function RegisterPage() {
  const [username, setUsername] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [birthdate, setBirthdate] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };

  const handleFirstNameChange = (event) => {
    setFirstName(event.target.value);
  };

  const handleLastNameChange = (event) => {
    setLastName(event.target.value);
  };

  const handleBirthdateChange = (event) => {
    setBirthdate(event.target.value);
  };

  const handlePasswordChange = (event) => {
    setPassword(event.target.value);
  };

  const handleConfirmPasswordChange = (event) => {
    setConfirmPassword(event.target.value);
  };

  const handleRegister = () => {
    // Perform registration logic here
    console.log('Username:', username);
    console.log('First Name:', firstName);
    console.log('Last Name:', lastName);
    console.log('Birthdate:', birthdate);
    console.log('Password:', password);
    console.log('Confirm Password:', confirmPassword);
  };

  return (
    <div className='register-page'>
    <WelcomeBox />
    <div className="register-container">
      <h1>Register Page</h1>
      <form>
        <div>
          <label>Username:</label>
          <input
            type="text"
            value={username}
            onChange={handleUsernameChange}
            required
          />
        </div>
        <div>
          <label>First Name:</label>
          <input
            type="text"
            value={firstName}
            onChange={handleFirstNameChange}
            required
          />
        </div>
        <div>
          <label>Last Name:</label>
          <input
            type="text"
            value={lastName}
            onChange={handleLastNameChange}
            required
          />
        </div>
        <div>
          <label>Birthdate:</label>
          <input
            type="date"
            value={birthdate}
            onChange={handleBirthdateChange}
            required
          />
        </div>
        <div>
          <label>Password:</label>
          <input
            type="password"
            value={password}
            onChange={handlePasswordChange}
            required
          />
        </div>
        <div>
          <label>Confirm Password:</label>
          <input
            type="password"
            value={confirmPassword}
            onChange={handleConfirmPasswordChange}
            required
          />
        </div>
        <button type="button" onClick={handleRegister}>
          Register
        </button>
      </form>
      <div className='links'>
        <a href='/login'>Login</a>
      </div>
    </div>
    </div>
  );
}

export default RegisterPage;