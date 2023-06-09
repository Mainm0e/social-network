import React, { useState } from 'react';
import '../RegisterPage.css';
import './info1.css';

// Info1 component
// Props: selectedOption, onChange, registerStatus
// selectedOption is the state of the form
// onChange is the function to change the state of the form
// registerStatus is the state of the register status
const Info1 = ({ selectedOption, onChange, registerStatus }) => {
  const type = "info1";
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [birthdate, setBirthdate] = useState("");

  const handleInputChange = (event) => {
    const { name, value } = event.target;

    // Update the respective state based on the input name
    switch (name) {
      case "firstName":
        setFirstName(value);
        break;
      case "lastName":
        setLastName(value);
        break;
      case "birthdate":
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
      birthdate,
    });
  }, [type, firstName, lastName, birthdate, onChange]);

  return (
    <form>
    <div className={`user-info1_${selectedOption !== 'info1' && 'hidden'}`}>
      <div className='input-container'>
        <label>First Name:</label>
        <input type='text' style={{ background: registerStatus === false && firstName === "" ? "#FFEA00" : "" }} name='firstName' value={firstName} onChange={handleInputChange} />
      </div>
      <div className='input-container'>
        <label>Last Name:</label>
        <input type='text' style={{ background: registerStatus === false && lastName === "" ? "#FFEA00" : "" }} name='lastName' value={lastName} onChange={handleInputChange} />
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
