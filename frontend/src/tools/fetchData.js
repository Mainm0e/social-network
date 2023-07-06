import { logout } from "./logout"
export async function fetchData(method, type, payload) {
  const response = await fetch("http://localhost:8080/api", {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      type: type,
      payload: payload,
    }),
  });
  const responseData = await response.json();

  //this is for debugging
  if (responseData.statusCode !== 200) {
    console.log("type", type, "error", responseData);
  }

  if (type === "login") {
    if (responseData.statusCode === 200) {
      console.log("set cookie", responseData);
      document.cookie = "sessionId=" + responseData.event.payload.sessionId;
      localStorage.setItem("userId", responseData.event.payload.userId);
      window.location.href = "/";
      return responseData.event.payload;
    } else {
      return responseData;
    }
  }
  if (type === "register") {
    if (responseData.statusCode === 200) {
      console.log("set cookie", responseData);
      window.location.href = "/login";
      return responseData.event.payload;
    } else {
      return responseData;
    }
  }
  if (responseData.statusCode === 200) {
    console.log("responseData", responseData)
    return responseData.event.payload;
  }
  if (responseData.statusCode !== 200) {
    if (responseData.message === "Error handling event:Error fetchingUser:user not found"){
      logout();
    }
  }
}
