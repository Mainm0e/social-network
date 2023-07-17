import { getCookie, getUserId } from "./cookie";
import { fetchData } from "./fetchData";

export async function logout() {
  const method = "POST";
  const type = "logout"
  const payload = {
    userId: getUserId("userId"),
    sessionId: getCookie("sessionId")
   };
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
  document.cookie = "sessionId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  localStorage.removeItem("userId");
  // sent cookie to server to delete
  window.location.href = "/login";
}