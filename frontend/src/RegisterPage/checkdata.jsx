
// checkData function
// This function is used to check if the data is valid
// If the data is valid, the function will return an object with status = true
// and an empty array of message
// If the data is invalid, the function will return an object with status = false
// and an array of message
export function checkData(data){
 const alert = {
    title: "",
    message: [],
    status: true
 }
 console.log("check alert", alert,data)
    let msg = [];
    for (let key in data) {
    if (data.hasOwnProperty(key)) {
        const value = data[key];
        // check if the input is empty
        if (value === ""&& key !== "avatar"&& key !== "password") {
        msg.push(key + " is empty");
        alert.status = false;
        }

        if (key === "password" && value === "") {
        msg.push("Password is empty");
        alert.status = false;
        }
         if (key === "password" && value.length < 8) {
         msg.push("Password is too short");
         alert.status = false;
         }
         if (key === "password" && value.length > 20) {
         msg.push("Password is too long");
         alert.status = false;
         }
         if (key === "email" && checkEmail(value) === false) {
         msg.push("Email is invalid");
         alert.status = false;
         }
    }

    }
    if (msg.length === 0){
    alert.status = true;
    } else if (msg.length > 0){
    alert.status = false;
    }
    alert.title = alert.status ? "" : "Error";
    alert.message = msg;
 return (alert)
}

// checkEmail function
// This function is used to check if the email is valid
// If the email is valid, the function will return true
// If the email is invalid, the function will return false
export function checkEmail(email) {
      // regex for email
      const regex = /\S+@\S+\.\S+/;
      return regex.test(email);
}
