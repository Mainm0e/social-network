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
  // const handleOptionChange = (event) => {
  //   const selectedOption = event.target.value;
  //   setSelectedOption(selectedOption);
  // };

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

  // button to navigate forward through info pages
  const handleNext = () => {
    if (selectedOption === "info1") {
      setSelectedOption("info2");
    } else if (selectedOption === "info2") {
      setSelectedOption("info3");
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
    console.log("check when send data", response.status);
    if (response.status) {
      setRegisterStatus(false);
      const method = "POST";
      const type = "register";
      const payload = data;
      fetchData(method, type, payload)
        .then((response) => response.json())
        .then((data) => {
          if (data.statusCode === 200) {
            setRegisterStatus(true);
            // Wait for the current execution cycle to finish
            setTimeout(() => {
              window.location.href = "/login";
            }, 0);
          } else {
            setAlertTitle(data.message);
            setAlertMessage(["Email already exists"]);
            setRegisterStatus(false);
          }
        })
        .catch((error) => {
          console.error("Error:", error);
        });
    }
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
  );

}

export default RegisterPage;
