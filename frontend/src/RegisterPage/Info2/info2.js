import React, { useState } from "react";
import "../RegisterPage.css";
import "./info2.css";

// Info2 component
// Props: selectedOption, onChange, registerStatus
// selectedOption is the state of the form
// onChange is the function to change the state of the form
// registerStatus is the state of the register status
const Info2 = ({ selectedOption, onChange, registerStatus }) => {
  const type = "info2";

  // Call the onChange prop with the input values whenever they change
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const handleInputChanges = (event) => {
    const { name, value } = event.target;
    if (name === "username") {
      setUsername(value);
    } else if (name === "password") {
      setPassword(value);
    } else if (name === "confirmPassword") {
      setConfirmPassword(value);
    }
  };

  React.useEffect(() => {
    onChange({ type, username, password, confirmPassword });
  }, [type, username, password, confirmPassword]);

  return (
    <form>
    <div className={`user-info2_${selectedOption !== 'info2' && 'hidden'}`}>
      <div className='input-container'>
        <label>Username:</label>
        <input type='text' name='username' style={{ background: registerStatus === false && username === "" ? "#FFEA00" : "" }} value={username} onChange={handleInputChanges} />
      </div>
      <div className='input-container'>
        <label>Password:</label>
        <input type='password' name='password' style={{ background: registerStatus === false && password === "" ? "#FFEA00" : "" }} value={password} onChange={handleInputChanges} />
      </div>
      <div className='input-container'>
        <label>Confirm Password:</label>
        <input type='password' name='confirmPassword' style={{ background: registerStatus === false && confirmPassword === "" ? "#FFEA00" : "" }} value={confirmPassword} onChange={handleInputChanges} />
      </div>
    </div>
    </form>
  );
};

export default Info2;
