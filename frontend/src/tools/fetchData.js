import { logout } from "./logout"
export async function fetchData(method, type, payload) {
  const response = await fetch("http://localhost:8080/api", {
    method: method,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      type: type,
      payload: payload,
    }),
  }).catch((error) => {
    console.log("error", error);
  });
  
if (!response.ok) {
  // Handle the error based on the response status
  if (response.status === 401) {
    // Unauthorized error
    console.log("Unauthorized error");
    logout();
  } else {
    // Other error
    console.log("HTTP error:", response.status);
  }
} else {
  const responseData = await response.json();
  if (type === "login") {
    if (responseData.statusCode === 200) {
      document.cookie = "sessionId=" + responseData.event.payload.sessionId;
      localStorage.setItem("userId", responseData.event.payload.userId);
      window.location.href = "/";
      return null
    } else {
      return responseData;
    }
  }
  if (type === "register") {
    if (responseData.statusCode === 200) {
      window.location.href = "/login";
      return responseData.event.payload;
    } else if (responseData.statusCode === 400) {
      return responseData;
    }
  }
  if (responseData.statusCode === 200) {
    return responseData.event.payload;
  }
  if (responseData.statusCode !== 200) {
    if (responseData.message === "Error handling event:Error fetchingUser:user not found"){
      logout();
    }
  }
}
}
