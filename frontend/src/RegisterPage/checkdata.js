
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
    let msg = [];
    for (let key in data) {
    if (data.hasOwnProperty(key)) {
        const value = data[key];
        if (value === ""&& key !== "avatar"&& key !== "matchPassword") {
        msg.push(key + " is empty");
        alert.status = false;
        }
        if (key === "matchPassword" && value === "") {
        msg.push("Password does not match");
        alert.status = false;
        }
    }
    }
    if (msg.length === 0){
    alert.status = true;
    } else if (msg.length > 0){
    alert.status = false;
    }
    alert.title = alert.status ? "Success" : "Error";
    alert.message = msg;
 return (alert)
}
