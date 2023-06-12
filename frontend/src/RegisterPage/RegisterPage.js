import React, { useState } from "react";
import axios from 'axios';
import Info1 from "./Info1/info1";
import Info2 from "./Info2/info2";
import Info3 from "./Info3/info3";
import "./RegisterPage.css";
import WelcomeBox from "../Common/WelcomeBox/WelcomeBox";
import AlertBox from "../Common/AlertBox/AlertBox";
import { checkData } from "./checkdata";


// RegisterPage component
// Props: none
function RegisterPage() {
  // !!! Junk useSate !!!
  //
  const [selectedOption, setSelectedOption] = useState("info1");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [birthdate, setBirthdate] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [nickname, setNickname] = useState("");
  const [aboutme, setAboutme] = useState("");
  const [avatar, setAvatar] = useState(null);

  // handleOptionChange function
  // This function is used to change the form state
  const handleOptionChange = (event) => {
    const selectedOption = event.target.value;
    setSelectedOption(selectedOption);
  };

  // handleInfoChange function
  // this function get the values from Info1 and Info2 components
  // and set the state of the form
  const handleInfoChange = (values) => {
    if (values.type === "info1") {
      setFirstName(values.firstName);
      setLastName(values.lastName);
      setEmail(values.email);
      setBirthdate(values.birthdate);
    } else if (values.type === "info2") {
      setUsername(values.username);
      setPassword(values.password);
      setConfirmPassword(values.confirmPassword);
      setAvatar(values.avatar);
    } else if (values.type === "info3") {
      setNickname(values.nickname);
      setAboutme(values.aboutme);
      setAvatar(values.avatar);
    }
  };

  // handleRegister function
  const [alertTitle, setAlertTitle] = useState("");
  const [alertMessage, setAlertMessage] = useState([]);
  const [registerStatus, setRegisterStatus] = useState(true);
  // matchPassword is used to check if the password and confirmPassword are the same
  // if they are the same, {matchPassword = password} else {matchPassword = ""}
  let matchPassword = password === confirmPassword ? password : "";
  const data = {
    firstName,
    lastName,
    email,
    birthdate,
    username,
    matchPassword,
    nickname,
    aboutme,
    avatar,
  };

  // register function
  // this function is main function of the register page
  const register = () => {
    console.log(data);
    const response = checkData(data);
    setAlertTitle(response.title);
    setAlertMessage(response.message);
    setRegisterStatus(response.status);
    console.log(response.status);
    if (registerStatus) {
      document.querySelector(".alert-box").style.display = "none";
    } else {
      document.querySelector(".alert-box").style.display = "block";
    }

    if (username !== "" && password !== "" && firstName !== "" && lastName !== "" && email !== "" && birthdate !== "" && avatar !== ""){
    const FormData = require("form-data");
    const formData = new FormData();
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
    formData.append("birthdate", birthdate);
    formData.append("username", username);
    formData.append("password", password);
    formData.append("nickname", nickname);
    formData.append("aboutme", aboutme);
    formData.append("avatar", avatar);

    console.log(formData);
   // Make a POST request to the Golang endpoint
   axios.post('http://localhost:8080/upload', formData)
   .then((response) => {
     console.log(response);
   })
   .catch((error) => {
     console.log(error);
   });
 
  };
  };

  return (
    <div className="register-page">
      <WelcomeBox />
      <div className="register-container">
        <AlertBox
          title={alertTitle}
          message={alertMessage}
          status={registerStatus}
        />{" "}
        <h1> Register Page </h1>{" "}
        <Info1
          selectedOption={selectedOption}
          onChange={handleInfoChange}
          registerStatus={registerStatus}
        />{" "}
        <Info2
          selectedOption={selectedOption}
          onChange={handleInfoChange}
          registerStatus={registerStatus}
        />{" "}
        <Info3
          selectedOption={selectedOption}
          onChange={handleInfoChange}
          registerStatus={registerStatus}
        />{" "}
        <div className="select-container">
          <input
            type="radio"
            value="info1"
            checked={selectedOption === "info1"}
            onChange={handleOptionChange}
          />{" "}
          <input
            type="radio"
            value="info2"
            checked={selectedOption === "info2"}
            onChange={handleOptionChange}
          />{" "}
          <input
            type="radio"
            value="info3"
            checked={selectedOption === "info3"}
            onChange={handleOptionChange}
          />{" "}
        </div>{" "}
        <button type="button" onClick={register}>
          Register{" "}
        </button>{" "}
        <div className="links">
          <a href="/login"> Login </a>{" "}
        </div>{" "}
      </div>{" "}
    </div>
  );
}

export default RegisterPage;
