import React, { useState } from 'react';
import '../RegisterPage.css';

const Info1 = ({ selectedOption, onChange ,registerStatus}) => {
    const type = 'info1';
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [birthdate, setBirthdate] = useState('');

  const handleInputChange = (event) => {
    const { name, value } = event.target;

    // Update the respective state based on the input name
    switch (name) {
      case 'firstName':
        setFirstName(value);
        break;
      case 'lastName':
        setLastName(value);
        break;
      case 'email':
        setEmail(value);
        break;
      case 'birthdate':
        setBirthdate(value);
        break;
      default:
        break;
    }
  };

  // Call the onChange prop with the input values whenever they change
  React.useEffect(() => {
    onChange({
        type,
      firstName,
      lastName,
      email,
      birthdate,
    });
  }, [type,firstName, lastName, email, birthdate, onChange]);

  return (
    <form>
    <div className={`user-info1 ${selectedOption !== 'info1' && 'hidden'}`}>
      <div className='input-container'>
        <label>First Name:</label>
        <input type='text' style={{ background: registerStatus === false && firstName === "" ? "#FFEA00" : "" }} name='firstName' value={firstName} onChange={handleInputChange} />
      </div>
      <div className='input-container'>
        <label>Last Name:</label>
        <input type='text' style={{ background: registerStatus === false && lastName === "" ? "#FFEA00" : "" }} name='lastName' value={lastName} onChange={handleInputChange} />
      </div>
      <div className='input-container'>
        <label>Email:</label>
        <input type='text' style={{ background: registerStatus === false && email === "" ? "#FFEA00" : "" }} name='email' value={email} onChange={handleInputChange} />
      </div>
      <div className='input-container'>
        <label>Birthdate:</label>
        <input type='date' style={{ background: registerStatus === false && birthdate ==="" ? "#FFEA00" : "" }} name='birthdate' value={birthdate} onChange={handleInputChange} />
      </div>
    </div>
    </form>
  );
};

export default Info1;
