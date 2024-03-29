import React, { useState } from "react";
import Info1 from "./Info1/info1";
import Info2 from "./Info2/info2";
import Info3 from "./Info3/info3";
import "./RegisterPage.css";
import WelcomeBox from "../Common/WelcomeBox/WelcomeBox";
import AlertBox from "../Common/AlertBox/AlertBox";
import { checkData } from "./checkdata";
import { fetchData } from "../tools/fetchData";

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
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [nickName, setNickname] = useState("");
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
      setBirthdate(values.birthdate);
    } else if (values.type === "info2") {
      setEmail(values.email);
      setPassword(values.password);
      setConfirmPassword(values.confirmPassword);
      setAvatar(values.avatar);
    } else if (values.type === "info3") {
      setNickname(values.nickname);
      setAboutme(values.aboutme);
      setAvatar(values.avatar);
    }
  };

  // button to navigate back through info pages 
  const handleBack = () => {
    if (selectedOption === "info2") {
      setSelectedOption("info1");
    } else if (selectedOption === "info3") {
      setSelectedOption("info2");
    }

  };

  const [setStep] = useState(1);

  // button to navigate forward through info pages
  const handleNext = () => {
    let emptyFields = [];
    if (selectedOption === 'info1') {
      // Check if required fields in Info1 are filled
      if (!firstName) emptyFields.push('firstName');
      if (!lastName) emptyFields.push('lastName');
      if (!birthdate) emptyFields.push('birthdate');
    } else if (selectedOption === 'info2') {
      // Check if required fields in Info2 are filled
      if (!email) emptyFields.push('email');
      if (!password) emptyFields.push('password');
      if (!confirmPassword) emptyFields.push('confirmPassword');
    } else if (selectedOption === 'info3') {
      // Check if required fields in Info3 are filled
      if (!nickName) emptyFields.push('nickName');
      if (!aboutme) emptyFields.push('aboutme');
    }

    if (emptyFields.length === 0) {
      // All required fields are filled, proceed to the next step
      if (selectedOption === 'info1') setSelectedOption('info2');
      else if (selectedOption === 'info2') setSelectedOption('info3');
      else if (selectedOption === 'info3') {
        // If on info3, move to the next step (you can trigger the registration here)
        setStep((prevStep) => prevStep + 1);
      }
    } else {
      // Some required fields are empty, add the shake animation to the empty fields
      emptyFields.forEach((field) => {
        document.querySelector(`[name=${field}]`).classList.add('shake');
      });

      // Remove the shake animation after a short delay
      setTimeout(() => {
        emptyFields.forEach((field) => {
          document.querySelector(`[name=${field}]`).classList.remove('shake');
        });
      }, 500);
    }
  };
  

  // handleRegister function
  const [alertTitle, setAlertTitle] = useState("");
  const [alertMessage, setAlertMessage] = useState([]);
  const [registerStatus, setRegisterStatus] = useState(true);
  // matchPassword is used to check if the password and confirmPassword are the same
  // if they are the same, {matchPassword = password} else {matchPassword = ""}
  let data = {
    firstName,
    lastName,
    email,
    birthdate,
    password,
    nickName,
    aboutme,
    avatar,
  };
  let matchPassword = password === confirmPassword;
  if (matchPassword) {
    data = {
      firstName,
      lastName,
      email,
      birthdate,
      password,
      nickName,
      aboutme,
      avatar,
    };
  }

  // register function
  // this function is main function of the register page
  const register = () => {
    const response = checkData(data);
    setAlertTitle(response.title);
    setAlertMessage(response.message);
    setRegisterStatus(response.status);
    // check if the registerStatus is true
    if (registerStatus) {
      document.querySelector(".alert-box").style.display = "none";
    } else {
      document.querySelector(".alert-box").style.display = "block";
    }

    // if the registerStatus is true, send the data to the backend
    if (response.status) {
      setRegisterStatus(false);
      const method = "POST";
      const type = "register";
      const payload = data;
      fetchData(method, type, payload).then((data) => {
          if (data.statusCode === 200) {
            setRegisterStatus(true);
            // Wait for the current execution cycle to finish
            setTimeout(() => {
              window.location.href = "/login";
            }, 0);
          } else {
            setAlertTitle("Error");
            setAlertMessage([data.message]);
            setRegisterStatus(false);
          }
        })
    }
  };

  return (
    <div className="main-container">
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
        {selectedOption !== "info1" && (
      <button type="button" onClick={handleBack}>
        Back
      </button>
    )}
        {selectedOption === "info3" && (
      <button type="button" onClick={register}>
        Register
      </button>
    )}
        {selectedOption !== "info3" && (
      <button type="next-btn" onClick={handleNext}>
        Next
      </button>
    )}
        <div className="links">
          <a href="/login"> Login </a>{" "}
        </div>{" "}
      </div>{" "}
    </div>
    </div>
  );

}

export default RegisterPage;
